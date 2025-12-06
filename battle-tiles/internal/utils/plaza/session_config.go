// pkg/plaza/session.go
package plaza

import (
	"battle-tiles/internal/consts"
	"battle-tiles/internal/dal/vo/game"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/patrickmn/go-cache"
)

/* =========================
   可配置的 SessionConfig
   ========================= */

type SessionConfig struct {
	// 网络
	Server82     string        // 例: "androidsc.foxuc.com:8200"
	Server87Host string        // 例: "newbgp.foxuc.com"
	KeepAlive    time.Duration // 例: 30 * time.Second

	// 行为
	AutoReconnect bool

	// 身份
	Identifier string // 账号或手机号
	UserPwdMD5 string // 登录密码（MD5，大写）
	UserID     int
	HouseGID   int
	LoginMode  consts.GameLoginMode

	// 事件回调
	Handler IPlazaHandler
}

/* =========================
   回调接口（保持你原来的）
   ========================= */

type IPlazaHandler interface {
	OnSessionRestarted(session *Session)
	OnMemberListUpdated(members []*GroupMember)
	OnMemberInserted(member *MemberInserted)
	OnMemberDeleted(member *MemberDeleted)
	OnMemberRightUpdated(key string, memberID int, success bool)
	OnLoginDone(success bool)
	OnRoomListUpdated(tables []*TableInfo)
	OnUserSitDown(sitdown *UserSitDown)
	OnUserStandUp(standup *UserStandUp)
	OnTableRenew(item *TableRenew)
	OnDismissTable(table int)
	OnAppliesForHouse(applyInfos []*ApplyInfo)
	OnReconnectFailed(houseGID int, retryCount int) // 新增：重连失败回调
}

/* =========================
   内部结构（与你原来一致）
   ========================= */

type ForbidMemberTag struct {
	Key      string
	MemberID int
}

type Session struct {
	cfg SessionConfig

	autoReconnect         bool
	userName              string
	userID                int
	userPwd               string
	houseGID              int
	handler               IPlazaHandler
	lastForbidCmdKey      string
	lastGroupMemKey       string
	lastCmdType           atomic.Int64
	dontReportApplicatons bool

	shutdown   atomic.Bool
	restarting atomic.Bool
	restarted  atomic.Bool

	// 被踢下线计数器
	kickedOfflineCount int
	lastKickedTime     time.Time

	_87connection            net.Conn
	_87encoder               *Encoder
	_87buffer                *RecvBuf
	_87portChan              chan int
	_87recvChan              chan []byte
	_87quitChan              chan bool
	_87connReady             atomic.Bool
	_87forbidTagStack        ForbidTagStack
	_87cmdQueue              GameCmdQueue
	_87waitingForCmdResponse atomic.Bool

	_82connection net.Conn
	_82encoder    *Encoder
	_82buffer     *RecvBuf
	_82recvChan   chan []byte
	_82quitChan   chan bool
	_82connReady  atomic.Bool

	tables       *cache.Cache
	members      *cache.Cache
	applications *cache.Cache

	// 已处理的申请ID（防止被游戏服务器推送重新添加）
	processedApplications *cache.Cache

	// 禁用成员缓存（内存中，key=gameID，value=昵称）
	forbiddenMembers *cache.Cache

	// houses: cache latest discovered group/house ids
	houses *cache.Cache
}

/* =========================
   构造 & 生命周期
   ========================= */

// 新的构造器：完全由 cfg 决定拨号与行为
func NewSessionWithConfig(cfg SessionConfig) (*Session, error) {
	s := new(Session)
	s.cfg = cfg

	s.autoReconnect = cfg.AutoReconnect
	s.userName = cfg.Identifier
	s.userID = cfg.UserID
	s.userPwd = cfg.UserPwdMD5
	s.houseGID = cfg.HouseGID
	s.handler = cfg.Handler

	s._82quitChan = make(chan bool)
	s._87quitChan = make(chan bool)

	s.tables = cache.New(10*time.Minute, 10*time.Minute)
	s.members = cache.New(10*time.Minute, 10*time.Minute)
	s.applications = cache.New(24*time.Hour, 1*time.Hour)          // 申请数据保留24小时
	s.processedApplications = cache.New(24*time.Hour, 1*time.Hour) // 已处理的申请ID保留24小时
	s.forbiddenMembers = cache.New(24*time.Hour, 1*time.Hour)      // 禁用成员保留24小时
	s.houses = cache.New(10*time.Minute, 10*time.Minute)
	if err := s.doLogonServer82(); err != nil {
		return nil, err
	}
	if err := s.doLogonServer87(); err != nil {
		return nil, err
	}
	go s.doSendCommand()
	return s, nil

}

func (that *Session) close() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	// 先关连接
	if that._82connection != nil {
		_ = that._82connection.Close()
	}
	if that._87connection != nil {
		_ = that._87connection.Close()
	}

	time.Sleep(10 * time.Millisecond)

	// 再停 goroutine
	if that._87connReady.Load() {
		that._87quitChan <- true
	}
	if that._82connReady.Load() {
		that._82quitChan <- true
	}

	that._87forbidTagStack.Clear()

	// 清理 cache：先 Flush 清空数据，再设置为 nil 让 GC 回收（触发 finalizer 停止 janitor）
	if that.tables != nil {
		that.tables.Flush()
		that.tables = nil
	}
	if that.members != nil {
		that.members.Flush()
		that.members = nil
	}
	if that.applications != nil {
		that.applications.Flush()
		that.applications = nil
	}
	if that.houses != nil {
		that.houses.Flush()
		that.houses = nil
	}
}

/* =========================
   命令发送主循环（与你原来一致）
   ========================= */

func (that *Session) prepareForbidCmd(key string, member int, forbid bool) {
	cmd := that._87cmdQueue.Top()
	if cmd == nil || cmd.Type != CmdTypeForbid || cmd.Key != key {
		that._87cmdQueue.AddHead(&GameCommand{
			Pack:   CmdForbidMember(uint32(that.userID), that.userPwd, uint32(that.houseGID), uint32(member), forbid),
			Type:   CmdTypeForbid,
			Key:    key,
			Member: member,
		})
	}
}

func (that *Session) doSendCommand() {
	that.lastForbidCmdKey = ""
	that.lastGroupMemKey = ""
	ticker := time.NewTicker(5 * time.Millisecond)

	for {
		if that.shutdown.Load() {
			return
		}
		if !that._87connReady.Load() {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if that.restarting.Load() {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		ticker.Reset(5 * time.Millisecond)
		for range ticker.C {
			if that.shutdown.Load() {
				return
			}
			if !that._87waitingForCmdResponse.Load() {
				goto SEND
			}
		}

	SEND:
		cmd := that._87cmdQueue.Pop()
		if cmd == nil {
			continue
		}
		if cmd.Type == CmdTypeForbid {
			if cmd.Key == that.lastForbidCmdKey {
				logger.Infof("跳过解禁命令,与上次key重复:%s", cmd.Key)
				goto SEND
			}
			that._87forbidTagStack.Push(&ForbidMemberTag{
				Key:      cmd.Key,
				MemberID: cmd.Member,
			})
		} else if cmd.Type == CmdTypeGetGroupMember {
			if cmd.Key == that.lastGroupMemKey {
				logger.Infof("跳过获取成员命令,与上次key重复:%s", cmd.Key)
				goto SEND
			}
		} else if cmd.Type == CmdTypeDismissTable {
			logger.Info("发送解散桌子指令")
		}

		data := that._87encoder.Encrypt(cmd.Pack)
		if _, err := that._87connection.Write(data); err != nil {
			if that.shutdown.Load() {
				return
			}
			if cmd.Type == CmdTypeForbid {
				that._87forbidTagStack.Pop()
			}
			if !strings.Contains(err.Error(), "time") {
				logger.Errorf("87服务器发送命令失败:%v", err)
				that.Restart()
				return
			}
		} else {
			that.lastCmdType.Store(int64(cmd.Type))
			that._87waitingForCmdResponse.Store(true)
			if cmd.Type == CmdTypeForbid {
				that.lastForbidCmdKey = cmd.Key
			} else if cmd.Type == CmdTypeGetGroupMember {
				that.lastGroupMemKey = cmd.Key
			}
		}
	}
}

/* =========================
   82/87 连接（改为使用 cfg）
   ========================= */

func (that *Session) doLogonServer82() error {
	con, err := net.DialTimeout("tcp", that.cfg.Server82, 5*time.Second)
	if err != nil {
		logger.Errorf("连接82服务器失败:%v", err)
		// 不在这里调用 Restart(),避免在 createNewSession 中造成指数级增长
		// that.Restart()
		return err
	}

	if tcp, ok := con.(*net.TCPConn); ok {
		_ = tcp.SetKeepAlive(true)
		if that.cfg.KeepAlive > 0 {
			_ = tcp.SetKeepAlivePeriod(that.cfg.KeepAlive)
		}
	}

	that._82connection = con
	that._82recvChan = make(chan []byte)
	that._87portChan = make(chan int)
	that._82buffer = NewRecvBuf(that._82recvChan)
	that._82connReady.Store(true)

	var cmd *game.Packer
	if that.cfg.LoginMode == consts.GameLoginModeMobile {
		cmd = CmdMobileLogon(that.cfg.Identifier, that.userPwd)
	} else {
		cmd = CmdAccountLogon(that.cfg.Identifier, that.userPwd)
	}
	that._82encoder = &Encoder{}
	ret := that._82encoder.Encrypt(cmd)
	if _, err = con.Write(ret); err != nil {
		return err
	}

	go that._82serverWaitForData()
	go that._82ServerHandleData()
	return nil
}

func (that *Session) doLogonServer87() error {
	var port int
	select {
	case port = <-that._87portChan:
	case <-time.After(8 * time.Second):
		return fmt.Errorf("timeout waiting 87 port from 82")
	}
	if port <= 0 {
		return fmt.Errorf("invalid port:%d", port)
	}
	addr := fmt.Sprintf("%s:%d", that.cfg.Server87Host, port)

	con, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	if tcp, ok := con.(*net.TCPConn); ok {
		_ = tcp.SetKeepAlive(true)
		if that.cfg.KeepAlive > 0 {
			_ = tcp.SetKeepAlivePeriod(that.cfg.KeepAlive)
		}
	}

	that._87connection = con
	that._87recvChan = make(chan []byte)
	that._87buffer = NewRecvBuf(that._87recvChan)
	that._87connReady.Store(true)

	that._87encoder = &Encoder{}
	cmd := CmdLogonServer(uint32(that.userID), that.userPwd)
	ret := that._87encoder.Encrypt(cmd)
	if _, err = con.Write(ret); err != nil {
		return err
	}

	go that._87serverWaitForData()
	go that._87ServerHandleData()

	// 87服务器连接成功后，发送进入房间命令（参考 battle-bot 实现）
	go func() {
		time.Sleep(2 * time.Second) // 等待连接稳定

		// 发送进入房间命令
		enterCmd := CmdGroupService(uint32(that.userID), uint32(that.houseGID))
		data := that._87encoder.Encrypt(enterCmd)
		if _, err := that._87connection.Write(data); err != nil {
			logger.Errorf("发送进入房间命令失败: %v", err)
			return
		}
		logger.Infof("已发送进入房间命令: UserID=%d, HouseGID=%d", that.userID, that.houseGID)

		// 等待进入房间响应
		time.Sleep(2 * time.Second)

		if that.handler != nil {
			that.handler.OnLoginDone(true)
			// 主动拉取成员列表，促使服务端推送房间列表快照
			that.GetGroupMembers()
		}
	}()

	return nil
}

/* =========================
   读循环 & 分发（与你原来一致）
   ========================= */

func (that *Session) _87serverWaitForData() {
	time.Sleep(1 * time.Second)
	buf := make([]byte, 1024*10)
	for {
		n, err := that._87connection.Read(buf)
		if err != nil {
			if !strings.Contains(err.Error(), "time") {
				if !that.shutdown.Load() {
					logger.Errorf("87服务器读取数据失败:%v", err)
					that.Restart()
				}
				return
			}
		}
		if n > 0 {
			that._87buffer.Add(buf[:n])
		}
	}
}

func (that *Session) _82serverWaitForData() {
	buf := make([]byte, 1024*10)
	for {
		n, err := that._82connection.Read(buf)
		if err != nil {
			logger.Error(err)
			if !strings.Contains(err.Error(), "time") {
				logger.Errorf("82服务器读取数据失败:%v", err)
				if !that._87connReady.Load() && !that.shutdown.Load() {
					that.Restart()
				}
				return
			}
		}
		if n > 0 {
			that._82buffer.Add(buf[:n])
		}
	}
}

func (that *Session) _87ServerHandleData() {
	for {
		select {
		case <-that._87quitChan:
			return
		case packet := <-that._87recvChan:
			if len(packet) == 0 {
				return
			}
			that._87handlePacket(packet)
		}
	}
}

func (that *Session) _82ServerHandleData() {
	for {
		select {
		case <-that._82quitChan:
			return
		case packet := <-that._82recvChan:
			if len(packet) == 0 {
				return
			}
			that._82handlePacket(packet)
		}
	}
}

/* =========================
   Protocol 分发（与原来一致）
   ========================= */

func (that *Session) _87handlePacket(data []byte) {
	packer, err := that._87encoder.Decrypt(data)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Infof("87收到数据包: MainCmdID=%d, SubCmdID=%d, DataLen=%d",
		packer.Head.Cmd.MainCmdID, packer.Head.Cmd.SubCmdID, len(packer.Data()))

	switch packer.Head.Cmd.MainCmdID {
	case consts.MDM_GA_BATTLE_SERVICE: // 1
		switch packer.Head.Cmd.SubCmdID {
		case consts.SUB_GA_TABLE_LIST:
			tables := ParseTableList(packer.Data())
			that.setTables(tables.Tables)
			that.handler.OnRoomListUpdated(tables.Tables)
			logger.Infof("tables push count=%d", len(tables.Tables))
		case consts.SUB_GA_DISMISS_RESULT:
			that._87waitingForCmdResponse.Store(false)
			result := ParseDismissTableResult(packer.Data())
			logger.Infof("解散桌子结果:%v\n", result)
		case consts.SUB_GA_USER_SITDOWN:
			user := ParseUserSitDown(packer.Data())
			logger.Infof("玩家坐下: UserID=%d, MappedNum=%d, ChairID=%d", user.UserID, user.MappedNum, user.ChairID)
			that.handler.OnUserSitDown(user)
			// 增量保障：根据坐下事件确保桌子存在于快照
			that.ensureTableByMappedNum(int(user.MappedNum))
		case consts.SUB_GA_USER_STANDUP:
			user := ParseUserStandUp(packer.Data())
			that.handler.OnUserStandUp(user)
			// 站起事件不移除桌，避免频繁抖动；如需可在无玩家时清理
		case consts.SUB_GA_TABLE_DISMISS:
			table := ParseTableDismissed(packer.Data())
			that.handler.OnDismissTable(table.MappedNum)
			that.removeTableByMappedNum(table.MappedNum)
		case consts.SUB_GA_TABLE_RENEW:
			item := ParseTableRenew(packer.Data())
			that.handler.OnTableRenew(item)
			that.renewTableMappedNum(item.MappedNum, item.NewMappedNum)
		}

	case consts.MDM_GA_LOGIC_SERVICE: // 2
		switch packer.Head.Cmd.SubCmdID {
		case consts.SUB_GA_APPLY_MESSAGE:
			if that.dontReportApplicatons {
				that.dontReportApplicatons = false
			} else {
				applyInfos := ParseApplyList(packer.Data())
				if len(applyInfos) != 0 {
					that.saveApplications(applyInfos)
					that.handler.OnAppliesForHouse(applyInfos)
				}
			}
		case consts.SUB_GA_OPERATE_SUCCESS:
			tag := that._87forbidTagStack.Pop()
			if tag != nil {
				that._87cmdQueue.Remove(tag.Key)
				go that.handler.OnMemberRightUpdated(strings.Split(tag.Key, ":")[0], tag.MemberID, true)
			}
			that._87waitingForCmdResponse.Store(false)
		case consts.SUB_GA_OPERATE_FAILURE:
			msg := ParseSystemMessage(packer.Data())
			logger.Infof("[%d]接收到操作失败:%s", that.houseGID, msg.Text)
			tag := that._87forbidTagStack.Pop()
			if tag != nil {
				logger.Infof("解禁用户失败:%v", tag)
				that._87cmdQueue.Remove(tag.Key)
				go that.handler.OnMemberRightUpdated(strings.Split(tag.Key, ":")[0], tag.MemberID, false)
			}
			that._87waitingForCmdResponse.Store(false)
		case consts.SUB_GA_SYSTEM_MESSAGE:
			msg := ParseSystemMessage(packer.Data())
			logger.Infof("[%d]接收到系统消息:%s", that.houseGID, msg.Text)
			if strings.Contains(msg.Text, "您的账号在其他地方登录，您被迫下线") {
				// 检查是否频繁被踢下线
				now := time.Now()
				if now.Sub(that.lastKickedTime) < 2*time.Minute {
					that.kickedOfflineCount++
				} else {
					// 超过2分钟,重置计数器
					that.kickedOfflineCount = 1
				}
				that.lastKickedTime = now

				// 如果2分钟内被踢下线超过5次,停止自动重连
				if that.kickedOfflineCount > 5 {
					logger.Errorf("[%d]账号频繁被踢下线(2分钟内%d次),停止自动重连", that.houseGID, that.kickedOfflineCount)
					that.autoReconnect = false
					that.Shutdown()
					// 通知上层需要停用中控账号
					if that.handler != nil {
						that.handler.OnReconnectFailed(that.houseGID, that.kickedOfflineCount)
					}
					return
				}

				if that.autoReconnect {
					that.dontReportApplicatons = true
					// 添加延迟,避免疯狂重试导致内存占满
					// 延迟时间随着重试次数增加: 5秒, 10秒, 15秒, 20秒, 25秒
					delaySeconds := that.kickedOfflineCount * 5
					logger.Infof("[%d]账号被踢下线,将在%d秒后重连(第%d次)", that.houseGID, delaySeconds, that.kickedOfflineCount)
					go func() {
						time.Sleep(time.Duration(delaySeconds) * time.Second)
						that.Restart()
					}()
				}
			} else if strings.Contains(msg.Text, "茶馆服务不可用。请稍后再次重试") {
				if that.autoReconnect {
					that.dontReportApplicatons = true
					that.Restart()
				}
			}
		}

	case consts.MDM_GA_GROUP_SERVICE: // 3
		switch packer.Head.Cmd.SubCmdID {
		case consts.SUB_GA_GROUP_MEMBER:
			that._87waitingForCmdResponse.Store(false)
			mems := ParseGroupMember(packer.Data())
			that.setMembers(mems)
			that.handler.OnMemberListUpdated(mems)
			logger.Infof("members push count=%d", len(mems))
		case consts.SUB_GA_MEMBER_INSERT:
			val := ParseMemberInserted(packer.Data())
			that.handler.OnMemberInserted(val)
		case consts.SUB_GA_MEMBER_DELETE:
			if that.lastCmdType.Load() == CmdTypeDeleteMember {
				logger.Info("踢出成员成功")
				that.lastCmdType.Store(-1)
				that._87waitingForCmdResponse.Store(false)
			} else {
				val := ParseMemberDeleted(packer.Data())
				that.handler.OnMemberDeleted(val)
			}
		}
	}
}

func (that *Session) _82handlePacket(data []byte) {
	packer, err := that._82encoder.Decrypt(data)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if packer.Head.Cmd.MainCmdID == 101 && packer.Head.Cmd.SubCmdID == 108 {
		accs := ParseListAccess(packer.Data())
		if len(accs) > 0 {
			that._87portChan <- int(accs[0].Port)
		}
	} else if packer.Head.Cmd.MainCmdID == consts.MDM_MB_SERVER_LIST {
		switch packer.Head.Cmd.SubCmdID {
		case consts.SUB_MB_LIST_SERVER:
			// Debug: dump first bytes for locating offsets (extended to 256 bytes)
			data := packer.Data()
			head := hexHead(data, 256)
			logger.Infof("MDM_MB_SERVER_LIST/SUB_MB_LIST_SERVER len=%d head=%s", len(data), head)
			// Extract candidate house ids
			if ids := ExtractHouseIDsFromServerList(data); len(ids) > 0 {
				n := len(ids)
				if n > 20 {
					n = 20
				}
				logger.Infof("server_list candidates=%d sample=%v", len(ids), ids[:n])
				for _, id := range ids {
					that.appendHouse(id)
				}
			}
		case consts.SUB_MB_LIST_FINISH:
			logger.Info("MDM_MB_SERVER_LIST, SUB_MB_LIST_FINISH got")
		}
	}
}

/* =========================
   控制（对外 API，与你原来一致）
   ========================= */

func (that *Session) SetAutoConnect(auto bool) {
	that.autoReconnect = auto
}

func (that *Session) ForbidMembers(key string, members []int, forbid bool) {
	for _, m := range members {
		if m == that.userID {
			continue
		}
		that.prepareForbidCmd(fmt.Sprintf("%s:%d", key, m), m, forbid)

		// 更新禁用成员缓存
		gameIDKey := fmt.Sprintf("%d", m)
		if forbid {
			// 尝试从在线成员列表获取昵称
			nickName := ""
			for _, mem := range that.ListMembers() {
				if int(mem.GameID) == m || int(mem.MemberID) == m {
					nickName = mem.NickName
					break
				}
			}
			that.forbiddenMembers.Set(gameIDKey, nickName, cache.DefaultExpiration)
		} else {
			// 解禁时从缓存移除
			that.forbiddenMembers.Delete(gameIDKey)
		}
	}
}

// IsForbidden 检查成员是否被禁用（基于内存缓存）
func (that *Session) IsForbidden(gameID int) bool {
	_, found := that.forbiddenMembers.Get(fmt.Sprintf("%d", gameID))
	return found
}

// ListForbiddenGameIDs 获取所有被禁用的 gameID 列表
func (that *Session) ListForbiddenGameIDs() []int {
	items := that.forbiddenMembers.Items()
	result := make([]int, 0, len(items))
	for key := range items {
		var gameID int
		if _, err := fmt.Sscanf(key, "%d", &gameID); err == nil {
			result = append(result, gameID)
		}
	}
	return result
}

// SetForbiddenStatus 设置成员禁用状态（仅更新缓存，不推送到游戏服务器）
func (that *Session) SetForbiddenStatus(gameID int, forbid bool) {
	gameIDKey := fmt.Sprintf("%d", gameID)
	if forbid {
		that.forbiddenMembers.Set(gameIDKey, "", cache.DefaultExpiration)
	} else {
		that.forbiddenMembers.Delete(gameIDKey)
	}
}

func (that *Session) GetGroupMembers() {
	cmd := that._87cmdQueue.Last()
	if cmd == nil || cmd.Type != CmdTypeGetGroupMember {
		gc := &GameCommand{
			Pack:   CmdGroupService(uint32(that.userID), uint32(that.houseGID)),
			Type:   CmdTypeGetGroupMember,
			Key:    fmt.Sprintf("%d", time.Now().UnixNano()),
			Member: that.houseGID,
		}
		that._87cmdQueue.Push(gc)
		logger.Infof("enqueue GetGroupMembers key=%s house=%d", gc.Key, that.houseGID)
	}
}

func (that *Session) RespondApplication(applyInfo *ApplyInfo, agree bool) {
	that._87cmdQueue.Push(&GameCommand{
		Pack: CmdRespondApplication(that.userID, that.userPwd, applyInfo.MessageId, applyInfo.HouseGid, applyInfo.AplierId, agree),
		Type: CmdTypeRespondApply,
		Key:  fmt.Sprintf("respond_app-%d", time.Now().UnixNano()),
	})
}

func (that *Session) GetDiamond() {
	that._87cmdQueue.Push(&GameCommand{
		Pack: CmdQueryDiamond(that.userID),
		Type: CmdTypeQueryDiamond,
		Key:  fmt.Sprintf("query_diamond-%d", time.Now().UnixNano()),
	})
}

func (that *Session) DismissTable(kindId int, mappedNum int) {
	that._87cmdQueue.AddHead(&GameCommand{
		Pack: CmdDismissRoom(that.userID, that.userPwd, kindId, mappedNum),
		Type: CmdTypeDismissTable,
		Key:  fmt.Sprintf("dismiss_table-%d", time.Now().UnixNano()),
	})
}

func (that *Session) KickOffMember(houseGid int, memberId int) {
	that._87cmdQueue.Push(&GameCommand{
		Pack: CmdDeleteMember(that.userID, that.userPwd, houseGid, memberId),
		Type: CmdTypeDeleteMember,
		Key:  fmt.Sprintf("delete_member-%d-%d", time.Now().UnixNano(), memberId),
	})
}

func (that *Session) QueryTable(tabMappedNum int) {
	if tabMappedNum <= 0 {
		// 防御：无效桌号不下发，避免远端强制断开
		return
	}
	that._87cmdQueue.AddHead(&GameCommand{
		Pack: CmdQueryTable(tabMappedNum),
		Type: CmdTypeQueryTable,
		Key:  fmt.Sprintf("query_table-%d", time.Now().UnixNano()),
	})
}

func (that *Session) Shutdown() {
	that.shutdown.Store(true)
	that.close() // close() 中已经清理了 cache

	safeCloseBoolChan(that._87quitChan)
	safeCloseBoolChan(that._82quitChan)
	safeCloseBytesChan(that._87recvChan)
	safeCloseBytesChan(that._82recvChan)
	safeCloseIntChan(that._87portChan)
}

/* =========================
   重启 & 安全关闭（与你原来一致）
   ========================= */

func (that *Session) Restart() {
	// 使用 CompareAndSwap 确保只有一个重启过程在运行
	// 如果已经在重启中,直接返回,避免启动多个 goroutine
	if that.shutdown.Load() {
		return
	}
	if !that.restarting.CompareAndSwap(false, true) {
		// 已经有重启过程在运行,直接返回
		logger.Warnf("[%d]重启已在进行中,跳过本次重启请求", that.houseGID)
		return
	}
	// 重置 restarting 标志将在 createNewSession 的 defer 中完成
	go that.createNewSession()
}
func (that *Session) setTables(arr []*TableInfo) {
	that.tables.Set("tables", cloneTables(arr), cache.DefaultExpiration)
}

// ListTables：读取房间列表快照（可能为 nil）
func (that *Session) ListTables() []*TableInfo {
	if v, ok := that.tables.Get("tables"); ok {
		if arr, ok2 := v.([]*TableInfo); ok2 {
			return cloneTables(arr)
		}
	}
	return nil
}

// ensureTableByMappedNum 若快照中无该 mapped_num，则加入一个最小占位项，避免列表始终为空
func (that *Session) ensureTableByMappedNum(mapped int) {
	if mapped <= 0 {
		return
	}
	cur := that.ListTables()
	for _, t := range cur {
		if t.MappedNum == mapped {
			return
		}
	}
	// 追加占位（KindID/BaseScore 等未知）；由后续 TABLE_LIST 或 RENEW 更新
	nt := &TableInfo{MappedNum: mapped}
	next := append(cur, nt)
	that.setTables(next)
}

// removeTableByMappedNum 从快照中删除指定桌
func (that *Session) removeTableByMappedNum(mapped int) {
	if mapped <= 0 {
		return
	}
	cur := that.ListTables()
	if len(cur) == 0 {
		return
	}
	next := make([]*TableInfo, 0, len(cur))
	for _, t := range cur {
		if t.MappedNum != mapped {
			next = append(next, t)
		}
	}
	that.setTables(next)
}

// renewTableMappedNum 在快照中将旧桌号变更为新桌号
func (that *Session) renewTableMappedNum(oldMapped, newMapped int) {
	if oldMapped <= 0 || newMapped <= 0 {
		return
	}
	cur := that.ListTables()
	updated := false
	for _, t := range cur {
		if t.MappedNum == oldMapped {
			t.MappedNum = newMapped
			updated = true
			break
		}
	}
	if !updated {
		cur = append(cur, &TableInfo{MappedNum: newMapped})
	}
	that.setTables(cur)
}

func (that *Session) setMembers(arr []*GroupMember) {
	if arr == nil {
		that.members.Delete("members")
		return
	}
	that.members.Set("members", cloneMembers(arr), cache.DefaultExpiration)
}

func (that *Session) appendHouse(houseGID int) {
	if houseGID <= 0 {
		return
	}
	key := fmt.Sprintf("%d", houseGID)
	that.houses.SetDefault(key, true)
}

func (that *Session) ListHouses() []int {
	items := that.houses.Items()
	out := make([]int, 0, len(items))
	for k := range items {
		var id int
		if _, err := fmt.Sscanf(k, "%d", &id); err == nil && id > 0 {
			out = append(out, id)
		}
	}
	sort.Ints(out)
	return out
}

func hexHead(b []byte, n int) string {
	if n > len(b) {
		n = len(b)
	}
	out := make([]byte, 0, n*3)
	for i := 0; i < n; i++ {
		out = append(out, fmt.Sprintf("%02X ", b[i])...)
	}
	return strings.TrimSpace(string(out))
}

// ListMembers：读取成员列表快照（可能为 nil）
func (that *Session) ListMembers() []*GroupMember {
	if v, ok := that.members.Get("members"); ok {
		if arr, ok2 := v.([]*GroupMember); ok2 {
			return cloneMembers(arr)
		}
	}
	return nil
}
func (that *Session) saveApplications(list []*ApplyInfo) {
	for _, a := range list {
		key := fmt.Sprintf("%d", a.MessageId)
		// 跳过已处理的申请
		if _, processed := that.processedApplications.Get(key); processed {
			continue
		}
		that.applications.Set(key, a, cache.DefaultExpiration)
	}
}

func (that *Session) ListApplications(houseGID int) []*ApplyInfo {
	items := that.applications.Items()
	out := make([]*ApplyInfo, 0, len(items))
	for _, it := range items {
		if ai, ok := it.Object.(*ApplyInfo); ok {
			if houseGID == 0 || ai.HouseGid == houseGID {
				out = append(out, ai)
			}
		}
	}
	// 可按创建时间降序
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt > out[j].CreatedAt })
	return out
}

func (that *Session) FindApplicationByID(msgID int) (*ApplyInfo, bool) {
	if v, ok := that.applications.Get(fmt.Sprintf("%d", msgID)); ok {
		if ai, ok2 := v.(*ApplyInfo); ok2 {
			return ai, true
		}
	}
	return nil, false
}

// RemoveApplication 从缓存中删除已处理的申请，并记录到已处理列表
func (that *Session) RemoveApplication(msgID int) {
	key := fmt.Sprintf("%d", msgID)
	that.applications.Delete(key)
	// 记录到已处理列表，防止游戏服务器推送时重新添加
	that.processedApplications.Set(key, true, cache.DefaultExpiration)
}

func (that *Session) createNewSession() bool {
	// restarting 标志已经在 Restart() 中设置为 true
	defer that.restarting.Store(false)

	logger.Infof("========== 重启茶馆 %d", that.houseGID)
	start := time.Now()

	// 关闭旧连接(只关闭连接,不设置 shutdown 标志)
	that.close()
	time.Sleep(50 * time.Millisecond)

	retry := 1
LOOP:
	if retry > 30 {
		logger.Errorf("重复连接house%d失败超过30次,不再尝试", that.houseGID)
		// 通知上层重连失败,需要自动停用中控账号
		if that.handler != nil {
			that.handler.OnReconnectFailed(that.houseGID, retry-1)
		}
		return false
	}

	retry++

	// 重新用 cfg 构建
	s, err := NewSessionWithConfig(that.cfg)
	if err != nil {
		logger.Errorf(">>>重启失败 (第%d次尝试)", retry-1)
		time.Sleep(5 * time.Second) // 改为5秒间隔
		goto LOOP
	} else {
		logger.Infof("========== 重启耗时：%d", time.Since(start).Milliseconds())
		s.restarted.Store(true)
		for {
			cmd := that._87cmdQueue.Pop()
			if cmd == nil {
				break
			}
			if cmd.Type == CmdTypeForbid || cmd.Type == CmdTypeDismissTable {
				if !that.restarted.Load() {
					s._87cmdQueue.Push(cmd)
				}
			}
		}
	}

	that.handler.OnSessionRestarted(s)
	return true
}

func safeCloseBytesChan(ch chan []byte) (justClosed bool) {
	defer func() {
		if recover() != nil {
			justClosed = false
		}
	}()
	close(ch)
	return true
}
func safeCloseIntChan(ch chan int) (justClosed bool) {
	defer func() {
		if recover() != nil {
			justClosed = false
		}
	}()
	close(ch)
	return true
}
func safeCloseBoolChan(ch chan bool) (justClosed bool) {
	defer func() {
		if recover() != nil {
			justClosed = false
		}
	}()
	close(ch)
	return true
}
func cloneTables(in []*TableInfo) []*TableInfo {
	if in == nil {
		return nil
	}
	out := make([]*TableInfo, len(in))
	copy(out, in)
	return out
}
