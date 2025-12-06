package game

import (
	"battle-tiles/internal/infra/plaza"
	utilsplaza "battle-tiles/internal/utils/plaza"
	pdb "battle-tiles/pkg/plugin/dbx"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// RoomCreditEventHandler 房间额度事件处理器
// 这是一个包装器，用于为特定的 (userID, houseGID) 会话创建 plaza.IPlazaHandler
type RoomCreditEventHandler struct {
	creditUC *RoomCreditLimitUseCase
	plazaMgr plaza.Manager
	log      *log.Helper
}

func NewRoomCreditEventHandler(
	creditUC *RoomCreditLimitUseCase,
	plazaMgr plaza.Manager,
	logger log.Logger,
) *RoomCreditEventHandler {
	return &RoomCreditEventHandler{
		creditUC: creditUC,
		plazaMgr: plazaMgr,
		log:      log.NewHelper(log.With(logger, "module", "handler/room_credit_event")),
	}
}

// CreateHandler 为特定的会话创建 Handler
// userID: 平台用户ID（中控账号对应的用户）
// houseGID: 店铺GID
// platform: 数据库平台标识
func (h *RoomCreditEventHandler) CreateHandler(userID int, houseGID int32, platform string) *sessionCreditHandler {
	return &sessionCreditHandler{
		userID:            userID,
		houseGID:          houseGID,
		platform:          platform,
		creditUC:          h.creditUC,
		plazaMgr:          h.plazaMgr,
		log:               h.log,
		roomsCache:        make(map[int]*utilsplaza.TableInfo),
		uncheckedSitdowns: make(map[int]int32),
	}
}

// sessionCreditHandler 实现 plaza.IPlazaHandler 接口
// 每个会话有自己的 Handler 实例，携带 userID 和 houseGID 信息
type sessionCreditHandler struct {
	userID   int
	houseGID int32
	platform string // 数据库平台标识
	creditUC *RoomCreditLimitUseCase
	plazaMgr plaza.Manager
	log      *log.Helper

	// 缓存：桌子号 -> 桌子信息
	roomsCache map[int]*utilsplaza.TableInfo
	// 缓存：桌子号 -> gameID（未检查的坐下事件）
	uncheckedSitdowns map[int]int32
}

// OnUserSitDown 处理玩家坐下事件
// 当收到游戏服务器的坐下事件时，检查玩家余额是否满足额度要求
func (h *sessionCreditHandler) OnUserSitDown(sitdown *utilsplaza.UserSitDown) {
	h.log.Infof("[费用检查] 收到玩家坐下事件: userID=%d, gameID=%d, mappedNum=%d, houseGID=%d",
		sitdown.UserID, sitdown.GameID, sitdown.MappedNum, h.houseGID)

	// 1. 检查桌子信息是否已缓存（类似 passing-dragonfly）
	tableInfo, ok := h.roomsCache[int(sitdown.MappedNum)]
	if !ok {
		// 桌子信息还没有收到，尝试从session快照获取（ensureTableByMappedNum已确保桌子存在）
		h.log.Warnf("[费用检查] 桌子信息未缓存，尝试从session快照获取: mappedNum=%d, gameID=%d",
			sitdown.MappedNum, sitdown.GameID)

		// 尝试从plaza manager获取桌子信息
		tables, err := h.plazaMgr.ListTables(h.userID, int(h.houseGID))
		if err == nil && tables != nil {
			for _, table := range tables {
				if table.MappedNum == int(sitdown.MappedNum) {
					// 检查桌子信息是否完整（有KindID和BaseScore）
					if table.KindID > 0 && table.BaseScore > 0 {
						// 找到完整的桌子信息，更新缓存并处理
						h.roomsCache[table.MappedNum] = table
						h.log.Infof("[费用检查] 从session快照获取完整桌子信息: mappedNum=%d, kindID=%d, baseScore=%d",
							table.MappedNum, table.KindID, table.BaseScore)

						// 立即处理坐下事件
						h.handlePlayerSitDown(int32(sitdown.GameID), table.KindID, table.BaseScore, int32(sitdown.MappedNum))
						return
					} else {
						// 只有占位项，信息不完整，需要等待OnRoomListUpdated
						h.log.Warnf("[费用检查] 桌子信息不完整（占位项）: mappedNum=%d, kindID=%d, baseScore=%d",
							table.MappedNum, table.KindID, table.BaseScore)
					}
				}
			}
		}

		// 如果仍然找不到或信息不完整，主动请求获取成员列表（可能会触发桌子列表更新）
		h.log.Warnf("[费用检查] 桌子信息未找到或不完整，主动请求获取成员列表: mappedNum=%d, gameID=%d",
			sitdown.MappedNum, sitdown.GameID)
		if err := h.plazaMgr.GetGroupMembers(h.userID, int(h.houseGID)); err != nil {
			h.log.Errorf("[费用检查] 请求获取成员列表失败: %v", err)
		}

		// 将坐下事件暂存，等待 OnRoomListUpdated
		h.uncheckedSitdowns[int(sitdown.MappedNum)] = int32(sitdown.GameID)
		return
	}

	// 检查缓存的桌子信息是否完整
	if tableInfo.KindID <= 0 || tableInfo.BaseScore <= 0 {
		h.log.Warnf("[费用检查] 缓存的桌子信息不完整: mappedNum=%d, kindID=%d, baseScore=%d, 等待OnRoomListUpdated",
			tableInfo.MappedNum, tableInfo.KindID, tableInfo.BaseScore)
		// 主动请求获取成员列表
		if err := h.plazaMgr.GetGroupMembers(h.userID, int(h.houseGID)); err != nil {
			h.log.Errorf("[费用检查] 请求获取成员列表失败: %v", err)
		}
		// 将坐下事件暂存，等待 OnRoomListUpdated
		h.uncheckedSitdowns[int(sitdown.MappedNum)] = int32(sitdown.GameID)
		return
	}

	h.log.Infof("[费用检查] 桌子信息: gameID=%d, kindID=%d, baseScore=%d, mappedNum=%d",
		sitdown.GameID, tableInfo.KindID, tableInfo.BaseScore, sitdown.MappedNum)

	// 2. 调用额度检查逻辑（传 gameID 而不是 userID）
	h.handlePlayerSitDown(int32(sitdown.GameID), tableInfo.KindID, tableInfo.BaseScore, int32(sitdown.MappedNum))
}

// 以下是 IPlazaHandler 接口的其他方法实现（空实现）

func (h *sessionCreditHandler) OnSessionRestarted(session *utilsplaza.Session) {}

func (h *sessionCreditHandler) OnMemberListUpdated(members []*utilsplaza.GroupMember) {}

func (h *sessionCreditHandler) OnMemberInserted(member *utilsplaza.MemberInserted) {}

func (h *sessionCreditHandler) OnMemberDeleted(member *utilsplaza.MemberDeleted) {}

func (h *sessionCreditHandler) OnMemberRightUpdated(key string, memberID int, success bool) {}

func (h *sessionCreditHandler) OnLoginDone(success bool) {}

func (h *sessionCreditHandler) OnRoomListUpdated(tables []*utilsplaza.TableInfo) {
	h.log.Infof("[费用检查] 收到桌子列表更新: 桌子数=%d, 未检查坐下事件数=%d", len(tables), len(h.uncheckedSitdowns))

	// 1. 构建当前桌子列表的映射，用于清理已解散的桌子
	currentTableMap := make(map[int]bool)
	for _, table := range tables {
		currentTableMap[table.MappedNum] = true
	}

	// 2. 清理已解散的桌子缓存（不在当前列表中的桌子）
	for cachedMappedNum := range h.roomsCache {
		if !currentTableMap[cachedMappedNum] {
			h.log.Debugf("[费用检查] 清理已解散的桌子缓存: mappedNum=%d", cachedMappedNum)
			delete(h.roomsCache, cachedMappedNum)
			// 同时清理对应的未检查坐下事件
			delete(h.uncheckedSitdowns, cachedMappedNum)
		}
	}

	// 3. 更新桌子缓存并处理未检查的坐下事件
	for _, table := range tables {
		h.roomsCache[table.MappedNum] = table
		h.log.Debugf("[费用检查] 缓存桌子: mappedNum=%d, kindID=%d, baseScore=%d",
			table.MappedNum, table.KindID, table.BaseScore)

		// 4. 检查是否有未处理的坐下事件（类似 passing-dragonfly 的 uncheckedSitdownCache）
		if gameID, ok := h.uncheckedSitdowns[table.MappedNum]; ok {
			h.log.Infof("[费用检查] 发现未检查的坐下事件，立即检查: mappedNum=%d, gameID=%d, kindID=%d, baseScore=%d",
				table.MappedNum, gameID, table.KindID, table.BaseScore)
			delete(h.uncheckedSitdowns, table.MappedNum)

			// 立即处理之前未检查的坐下事件
			h.handlePlayerSitDown(gameID, table.KindID, table.BaseScore, int32(table.MappedNum))
		}
	}

	// 5. 记录仍未处理的坐下事件（用于调试）
	if len(h.uncheckedSitdowns) > 0 {
		h.log.Warnf("[费用检查] 仍有未处理的坐下事件: %d个, mappedNums=%v",
			len(h.uncheckedSitdowns), h.uncheckedSitdowns)
	}
}

func (h *sessionCreditHandler) OnUserStandUp(standup *utilsplaza.UserStandUp) {}

func (h *sessionCreditHandler) OnTableRenew(item *utilsplaza.TableRenew) {}

func (h *sessionCreditHandler) OnDismissTable(table int) {
	h.log.Infof("[费用检查] 桌子被解散: mappedNum=%d", table)
	// 清理桌子缓存和未检查的坐下事件
	delete(h.roomsCache, table)
	delete(h.uncheckedSitdowns, table)
}

func (h *sessionCreditHandler) OnAppliesForHouse(applyInfos []*utilsplaza.ApplyInfo) {}

func (h *sessionCreditHandler) OnReconnectFailed(houseGID int, retryCount int) {}

// handlePlayerSitDown 处理玩家坐下的核心逻辑（类似 passing-dragonfly 的 _handlePlayerSitDown）
func (h *sessionCreditHandler) handlePlayerSitDown(gameID int32, kindID int, baseScore int, mappedNum int32) {
	// 创建带有数据库 key 的 context
	ctx := context.WithValue(context.Background(), pdb.CtxDBKey, h.platform)

	if err := h.creditUC.HandlePlayerSitDown(
		ctx,
		h.userID,
		h.houseGID,
		gameID,
		int32(kindID),
		int32(baseScore),
		mappedNum,
	); err != nil {
		h.log.Errorf("[费用检查] HandlePlayerSitDown失败: %v", err)
	}
}
