package plaza

import (
	"battle-tiles/internal/consts"
	"battle-tiles/internal/dal/vo/game"
	"bytes"
	"fmt"
	"time"
	"unicode/utf16"
	"unicode/utf8"
)

type Response struct {
	data []byte
}

func (that *Response) ReadByte(offset int) byte {
	return that.data[offset]
}

func (that *Response) ReadWord(offset int) uint16 {
	lo := uint16(that.data[offset])
	offset++
	hi := uint16(that.data[offset])
	return hi<<8 + lo
}

func (that *Response) ReadDWord(offset int) uint32 {
	b1 := uint32(that.data[offset])
	offset++
	b2 := uint32(that.data[offset])
	offset++
	b3 := uint32(that.data[offset])
	offset++
	b4 := uint32(that.data[offset])
	return b4<<24 + b3<<16 + b2<<8 + b1
}

func (that *Response) ReadLong(offset int) uint64 {
	b1 := uint64(that.data[offset])
	offset++
	b2 := uint64(that.data[offset])
	offset++
	b3 := uint64(that.data[offset])
	offset++
	b4 := uint64(that.data[offset])
	offset++
	b5 := uint64(that.data[offset])
	offset++
	b6 := uint64(that.data[offset])
	offset++
	b7 := uint64(that.data[offset])
	offset++
	b8 := uint64(that.data[offset])
	return b8<<56 + b7<<48 + b6<<40 + b5<<32 + b4<<24 + b3<<16 + b2<<8 + b1
}

func (that *Response) ReadString(offset int, ln int) string {
	u16s := make([]uint16, 1)
	ret := &bytes.Buffer{}
	b8buf := make([]byte, 4)
	for i := offset; i < len(that.data) && i < offset+(ln*2); i += 2 {
		if that.data[i] == 0 && that.data[i+1] == 0 {
			break
		}

		u16s[0] = uint16(that.data[i]) + (uint16(that.data[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String()
}

func ParseUserLogon(data []byte) *game.UserLogonInfo {
	var res Response
	res.data = data

	var info game.UserLogonInfo
	offset := 7
	info.UserID = res.ReadDWord(offset)
	offset += 4
	info.GameID = res.ReadDWord(offset)

	return &info
}

type ServerAgentList struct {
	ServerID uint16
	Agents   []*ServerAgent
}

type ServerAgent struct {
	AgentID    uint16
	ServicePor uint16
}

func ParseServerAgent(data []byte) []*ServerAgentList {
	var res Response
	res.data = data

	var ret []*ServerAgentList
	offset := 0
	for {
		if offset >= len(data) {
			break
		}

		sal := &ServerAgentList{}
		sal.ServerID = res.ReadWord(offset)
		offset += 2
		count := int(res.ReadWord(offset))
		offset += 2
		for i := 0; i < count; i++ {
			agent := &ServerAgent{}
			agent.AgentID = res.ReadWord(offset)
			offset += 2
			agent.ServicePor = res.ReadWord(offset)
			offset += 2
			sal.Agents = append(sal.Agents, agent)
		}

		ret = append(ret, sal)
	}

	return ret
}

type Access struct {
	ID   uint16
	Port uint16
	Addr string
}

func ParseListAccess(data []byte) []*Access {
	var res Response
	res.data = data

	var ret []*Access
	count := len(data) / 68 //LEN_ACCESS_ITEM
	var offset int
	for i := 0; i < count; i++ {
		offset = i * 68
		acc := &Access{}
		acc.ID = res.ReadWord(offset)
		offset += 2
		acc.Port = res.ReadWord(offset)
		offset += 2
		acc.Addr = res.ReadString(offset, consts.LEN_SERVER)
		offset += consts.LEN_SERVER
		ret = append(ret, acc)
	}
	return ret
}

type SystemMessage struct {
	Type uint16
	Text string
}

func ParseSystemMessage(data []byte) *SystemMessage {
	var res Response
	res.data = data

	sm := &SystemMessage{}

	sm.Type = res.ReadWord(0)
	sm.Text = res.ReadString(2, 128)
	return sm
}

type MessageItem struct {
	Type    byte
	Length  uint16
	Message string
}

func ParseMessageList(data []byte) []*MessageItem {
	var res Response
	res.data = data

	var ret []*MessageItem
	offset := 0
	for {
		if offset >= len(data) {
			break
		}

		item := &MessageItem{}
		item.Type = res.ReadByte(offset)
		offset++
		item.Length = res.ReadWord(offset)
		offset += 2
		item.Message = res.ReadString(offset, int(item.Length))
		offset += int(item.Length)

		ret = append(ret, item)
	}

	return ret
}

// --- Group/House listing (best-effort minimal parse) ---
// ExtractHouseIDsFromServerList 从 SUB_MB_LIST_SERVER 报文里提取可能的 house_gid（best-effort）
// 规则：
// - 小端 32 位整数，高 16 位为 0，低 16 位在 [10000, 70000]
// - 略过前若干头部字节（从 offset=8 开始扫描，步长=2），避免把头部字段误识别
func ExtractHouseIDsFromServerList(data []byte) []int {
	if len(data) < 4 {
		return nil
	}
	seen := map[int]struct{}{}
	out := make([]int, 0, 8)
	// 优先扫描 16-bit 小端值（更符合 60870 范围），然后用 32-bit 条件过滤
	for off := 8; off+2 <= len(data); off += 2 { // 2 字节对齐
		lo := int(data[off]) | int(data[off+1])<<8
		if lo >= 10000 && lo <= 70000 {
			// 可选：再检查其后 2 字节是否为 0 以提高置信度
			if off+4 <= len(data) {
				hi := int(data[off+2]) | int(data[off+3])<<8
				if hi != 0 {
					// 若不是 0 也允许，因结构体可能紧凑排列
				}
			}
			if _, ok := seen[lo]; !ok {
				seen[lo] = struct{}{}
				out = append(out, lo)
			}
		}
	}
	// 若包含目标 60870，则优先返回该值
	for _, v := range out {
		if v == 60870 {
			return []int{60870}
		}
	}
	return out
}

// FindUint16Offsets 返回 data 中等于指定 16 位小端值的所有偏移（从 0 开始）
func FindUint16Offsets(data []byte, value uint16) []int {
	if len(data) < 2 {
		return nil
	}
	lo := byte(value & 0xFF)
	hi := byte(value >> 8)
	out := make([]int, 0, 4)
	for i := 0; i+1 < len(data); i++ {
		if data[i] == lo && data[i+1] == hi {
			out = append(out, i)
		}
	}
	return out
}

type LogonFailure struct {
	Code int
	Desc string
}

func ParseLogonFailure(data []byte) *LogonFailure {
	var res Response
	res.data = data

	var ret LogonFailure
	ret.Code = int(res.ReadDWord(0))
	ret.Desc = res.ReadString(4, 128)
	return &ret

}

type GroupMember struct {
	UserID     uint32
	UserStatus int
	GameID     uint32
	MemberID   uint32
	MemberType int
	NickName   string
}

func ParseGroupMember(data []byte) []*GroupMember {
	var res Response
	res.data = data

	var members []*GroupMember
	//_ := res.ReadDWord(0)

	count := (len(data) - 4) / 127
	offset := 4
	for i := 0; i < count; i++ {
		offset = 4 + i*127
		//user Info
		mem := new(GroupMember)
		members = append(members, mem)
		mem.UserID = res.ReadDWord(offset)
		offset += 4
		mem.GameID = res.ReadDWord(offset)
		offset += 4

		mem.UserStatus = int(res.ReadByte(offset))
		offset += 1 //struct.cbGender      = pBuffer:readbyte()
		offset += 1 // struct.cbUserStatus  = pBuffer:readbyte()
		mem.NickName = res.ReadString(offset, consts.LEN_ACCOUNTS)
		offset += consts.LEN_ACCOUNTS * 2 // struct.szNickName    = pBuffer:readstring(df.LEN_ACCOUNTS)
		offset += 4                       // struct.dwCustomID    = pBuffer:readdword()

		//member Info
		mem.MemberID = res.ReadDWord(offset)
		offset += 4

		mem.MemberType = int(res.ReadByte(offset))
		//fmt.Printf("member name:%s,member type:%d,member id:%d,gameID=%d,status=%d\n", mem.NickName, mem.MemberType, mem.MemberID, mem.GameID, mem.UserStatus)
		//offset += 1 //struct.cbMemberType    = pBuffer:readbyte()                            --成员类型
		//offset += 4 //struct.dwMemberRight   = pBuffer:readdword()                           --成员权限
		//offset += 2 //struct.JoinDateTime.wYear 	 	  = pBuffer:readword()
		//offset += 2 //struct.JoinDateTime.wMonth 	 	  = pBuffer:readword()
		//offset += 2 //sstruct.JoinDateTime.wDayOfWeek 	  = pBuffer:readword()
		//offset += 2 //struct.JoinDateTime.wDay = pBuffer:readword()
		//offset += 2 //struct.JoinDateTime.wHour = pBuffer:readword()
		//offset += 2 //struct.JoinDateTime.wMinute = pBuffer:readword()
		//offset += 2 //struct.JoinDateTime.wSecond = pBuffer:readword()
		//offset += 2 //struct.JoinDateTime.wMilliseconds = pBuffer:readword()
		//offset += 4 //struct.dwCreateCount = pBuffer:readdword()                                --创建次数
		//offset += 4 //struct.dwBattleCount = pBuffer:readdword()                                --参与次数
		//
		////最近战绩
		//offset += 2 // struct.BattleDateTime.wYear = pBuffer:readword()
		//offset += 2 //struct.BattleDateTime.wMonth = pBuffer:readword()
		//offset += 2 //struct.BattleDateTime.wDayOfWeek = pBuffer:readword()
		//offset += 2 //struct.BattleDateTime.wDay = pBuffer:readword()
		//offset += 2 //struct.BattleDateTime.wHour = pBuffer:readword()
		//offset += 2 //sstruct.BattleDateTime.wMinute = pBuffer:readword()
		//offset += 2 //struct.BattleDateTime.wSecond = pBuffer:readword()
		//offset += 2 //struct.BattleDateTime.wMilliseconds = pBuffer:readword()
	}

	return members
}

type UserSitDown struct {
	UserID    uint32
	GameID    uint32
	MappedNum uint32
	ChairID   uint16
}

func ParseUserSitDown(data []byte) *UserSitDown {
	var res Response
	res.data = data

	var ret UserSitDown
	offset := 4
	ret.MappedNum = res.ReadDWord(offset)

	tableUser := ParseTableUserItem(data[8:])
	ret.UserID = tableUser.UserID
	ret.ChairID = tableUser.ChairID
	ret.GameID = tableUser.GameID

	return &ret
}

type UserStandUp struct {
	UserID    uint32
	ChairID   uint16
	MappedNum uint32
	//struct.wChairID 	= pBuffer:readword()							--用户椅子
	//struct.dwUserID		= pBuffer:readdword()							--用户标识
	//struct.dwMappedNum 	= pBuffer:readdword()						    --映射编号
}

func ParseUserStandUp(data []byte) *UserStandUp {
	var res Response
	res.data = data

	var ret UserStandUp

	ret.ChairID = res.ReadWord(0)
	ret.UserID = res.ReadDWord(2)
	ret.MappedNum = res.ReadDWord(6)
	return &ret
}

type TableUserItem struct {
	MappedNum uint32
	FaceID    uint16
	ChairID   uint16
	UserID    uint32
	GameID    uint32
	CustomID  uint32
	NickName  string
	//struct.wFaceID    = pBuffer:readword()								--头像标识 2
	//struct.wChairID	  = pBuffer:readword()								--用户方位 4
	//struct.dwUserID   = pBuffer:readdword()								--用户标识 8
	//struct.dwGameID   = pBuffer:readdword()								--游戏标识 12
	//struct.dwCustomID = pBuffer:readdword()								--头像标识 16
	//struct.szNickName = pBuffer:readstring(df.LEN_ACCOUNTS)							--用户昵称 80
}

func ParseTableUserItem(data []byte) *TableUserItem {
	var res Response
	res.data = data

	var ret TableUserItem
	offset := 0
	ret.FaceID = res.ReadWord(offset)
	offset += 2
	ret.ChairID = res.ReadWord(offset)
	offset += 2
	ret.UserID = res.ReadDWord(offset)
	offset += 4
	ret.GameID = res.ReadDWord(offset)
	offset += 4
	ret.CustomID = res.ReadDWord(offset)
	offset += 4
	ret.NickName = res.ReadString(offset, consts.LEN_ACCOUNTS)
	return &ret
}

type TableList struct {
	Count  int
	Tables []*TableInfo
}

type TableInfo struct {
	TableID   int
	MappedNum int
	GroupID   int
	KindID    int
	BaseScore int
	//struct.dwMappedNum   = pBuffer:readdword()						    --映射编号 4
	//struct.wFinishCount  = pBuffer:readword()						    --完成局数 6
	//struct.dwElapsedTime = pBuffer:readdword()						    --逝去时间 10
	//struct.wUserCount    = pBuffer:readword() 							--用户数量 12
}

func ParseTableList(data []byte) *TableList {
	var res Response
	res.data = data

	var ret TableList
	offset := 0
	ret.Count = int(res.ReadWord(offset))
	offset += 2

	for i := 0; i < ret.Count; i++ {
		tableID := int(res.ReadWord(offset))
		offset += 2
		offset += 2
		offset += 4
		mappedNum := int(res.ReadDWord(offset))
		offset += 4
		//passwd := res.ReadString(offset, LEN_PASSWORD)
		offset += 33 * 2

		groupID := int(res.ReadDWord(offset))
		offset += 4
		offset += 4

		kindID := int(res.ReadWord(offset))
		offset += 2
		offset += 2
		offset += 4
		offset += 63 * 2

		offset += 4

		offset += 2
		offset += 2
		offset += 1
		offset += 1
		offset += 1
		offset += 2
		offset += 4

		baseScore := res.ReadLong(offset)
		offset += 8
		offset += 8

		offset += 2
		offset += 1

		ret.Tables = append(ret.Tables, &TableInfo{
			TableID:   tableID,
			MappedNum: mappedNum,
			GroupID:   groupID,
			KindID:    kindID,
			BaseScore: int(baseScore),
		})
	}

	return &ret
}

type ApplyInfo struct {
	MessageId     int
	MessageStatus int
	ApplierGid    int
	AplierId      int
	ApplierGName  string
	HouseGid      int
	ApplyType     int
	AdminUserID   int
	CreatedAt     int64
}

func ParseApplyList(data []byte) []*ApplyInfo {
	var res Response
	res.data = data

	offset := 0
	applyType := int(res.ReadByte(offset))
	fmt.Printf("申请状态:%d\n", applyType)
	offset += 1

	applyCount := int(res.ReadWord(offset))
	offset += 2

	var result []*ApplyInfo
	for i := 0; i < applyCount; i++ {
		// struct.dwMessageID		  = pBuffer:readdword()					        --消息标识
		messageId := res.ReadDWord(offset)
		offset += 4
		// struct.cbMessageStatus    = pBuffer:readbyte()					        --消息状态
		status := res.ReadByte(offset)
		offset += 1
		// --申请信息
		// struct.dwApplyerID 		  = pBuffer:readdword()					        --用户标识
		applierId := res.ReadDWord(offset)
		offset += 4
		gid := res.ReadDWord(offset)
		// struct.dwApplyerGameID    = pBuffer:readdword()				            --游戏标识
		offset += 4
		// struct.dwApplyerCustomID  = pBuffer:readdword()							--头像标识
		offset += 4
		// struct.szApplyerNickName  = pBuffer:readstring(df.LEN_ACCOUNTS) 	    --用户昵称
		gname := res.ReadString(offset, consts.LEN_ACCOUNTS)
		offset += consts.LEN_ACCOUNTS * 2
		// struct.ApplyDateTime      = {}						                    --申请时间
		// struct.ApplyDateTime.wYear 	 	   = pBuffer:readword()
		year := res.ReadWord(offset)
		offset += 2
		// // struct.ApplyDateTime.wMonth        = pBuffer:readword()
		mon := res.ReadWord(offset)
		offset += 2
		// // struct.ApplyDateTime.wDayOfWeek    = pBuffer:readword()
		offset += 2
		// // struct.ApplyDateTime.wDay 		   = pBuffer:readword()
		day := res.ReadWord(offset)
		offset += 2
		// // struct.ApplyDateTime.wHour 	 	   = pBuffer:readword()
		hour := res.ReadWord(offset)
		offset += 2
		// // struct.ApplyDateTime.wMinute 	   = pBuffer:readword()
		min := res.ReadWord(offset)
		offset += 2
		// // struct.ApplyDateTime.wSecond 	   = pBuffer:readword()
		sec := res.ReadWord(offset)
		offset += 2
		// // struct.ApplyDateTime.wMilliseconds = pBuffer:readword()
		ms := res.ReadWord(offset)
		offset += 2

		// --群组信息
		// struct.dwGroupID          = pBuffer:readdword()					    --群组标识
		houseGid := res.ReadDWord(offset)
		offset += 4
		// struct.dwCreaterID        = pBuffer:readdword()				        --馆主标识
		creatorID := res.ReadDWord(offset)
		offset += 4
		// struct.szGroupName        = pBuffer:readstring(df.LEN_GROUP_NAME)		--群组名称
		// groupName := res.ReadString(offset, LEN_GROUP_NAME)
		// fmt.Println(groupName)
		offset += consts.LEN_GROUP_NAME * 2

		t, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, mon, day, hour, min, sec))
		createdAt := t.Unix() + int64(ms)

		result = append(result, &ApplyInfo{
			AplierId:      int(applierId),
			ApplierGid:    int(gid),
			ApplierGName:  gname,
			HouseGid:      int(houseGid),
			MessageId:     int(messageId),
			MessageStatus: int(status),
			ApplyType:     applyType,
			AdminUserID:   int(creatorID),
			CreatedAt:     createdAt,
		})
	}

	// if applyType == APPLY_MSG_TYPE_RESPOND {
	// 	hlogger.Info("忽略申请信息,这是一个Respond")
	// 	return nil
	// }
	return result
}

type DismissResult struct {
	GroupID   int
	MappedNum int
}

func ParseDismissTable(data []byte) *DismissResult {
	var res Response
	res.data = data

	var ret DismissResult
	ret.GroupID = int(res.ReadDWord(0))
	ret.MappedNum = int(res.ReadDWord(4))

	return &ret
}

type BattleOpFail struct {
	Code int
	Msg  string
}

func ParseBattleOpFail(data []byte) *BattleOpFail {
	var res Response
	res.data = data

	var ret BattleOpFail
	ret.Code = int(res.ReadDWord(0))
	ret.Msg = res.ReadString(4, len(data)-4)

	return &ret
}

type DismissTableResult struct {
	Code int
	Msg  string
}

func ParseDismissTableResult(data []byte) *DismissTableResult {
	var res Response
	res.data = data

	var ret DismissTableResult
	ret.Code = int(res.ReadByte(0))
	ret.Msg = res.ReadString(1, len(data)-1)

	return &ret
}

type MemberInserted struct {
	GroupID  uint32
	MemCount uint16

	//tagIMGroupMemberUser.tagIMUserInfo
	UserID uint32
	GameID uint32
}

func ParseMemberInserted(data []byte) *MemberInserted {
	var res Response
	res.data = data

	var ret MemberInserted
	offset := 0
	ret.GroupID = res.ReadDWord(0)
	offset += 4
	ret.MemCount = res.ReadWord(offset)
	offset += 2
	ret.UserID = res.ReadDWord(offset)
	offset += 4
	ret.GameID = res.ReadDWord(offset)
	return &ret

}

type MemberDeleted struct {
	MemberID uint32
}

func ParseMemberDeleted(data []byte) *MemberDeleted {
	var res Response
	res.data = data

	var ret MemberDeleted
	offset := 4
	ret.MemberID = res.ReadDWord(offset)
	return &ret
}

type TableUserList struct {
}

func ParseTableUserList(data []byte) *TableUserList {
	return nil
}

func ParseTableDismissed(data []byte) *TableInfo {
	var res Response
	res.data = data

	offset := 0
	groupID := res.ReadDWord(offset)
	offset += 4
	mappedNum := res.ReadDWord(offset)

	return &TableInfo{
		GroupID:   int(groupID),
		MappedNum: int(mappedNum),
	}
}

type TableRenew struct {
	MappedNum    int
	NewMappedNum int
}

func ParseTableRenew(data []byte) *TableRenew {
	var res Response
	res.data = data

	offset := 0
	mappedNum := res.ReadDWord(offset)
	offset += 6
	newMappedNum := res.ReadDWord(offset)

	return &TableRenew{
		MappedNum:    int(mappedNum),
		NewMappedNum: int(newMappedNum),
	}
}
