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
		// 桌子信息还没有收到，将坐下事件暂存，等待 OnRoomListUpdated
		h.log.Warnf("[费用检查] 桌子信息未缓存，等待OnRoomListUpdated: mappedNum=%d, gameID=%d",
			sitdown.MappedNum, sitdown.GameID)
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
	h.log.Infof("[费用检查] 收到桌子列表更新: 桌子数=%d", len(tables))

	// 1. 更新桌子缓存
	for _, table := range tables {
		h.roomsCache[table.MappedNum] = table
		h.log.Debugf("[费用检查] 缓存桌子: mappedNum=%d, kindID=%d, baseScore=%d",
			table.MappedNum, table.KindID, table.BaseScore)

		// 2. 检查是否有未处理的坐下事件（类似 passing-dragonfly 的 uncheckedSitdownCache）
		if gameID, ok := h.uncheckedSitdowns[table.MappedNum]; ok {
			h.log.Infof("[费用检查] 发现未检查的坐下事件，立即检查: mappedNum=%d, gameID=%d, kindID=%d, baseScore=%d",
				table.MappedNum, gameID, table.KindID, table.BaseScore)
			delete(h.uncheckedSitdowns, table.MappedNum)

			// 立即处理之前未检查的坐下事件
			h.handlePlayerSitDown(gameID, table.KindID, table.BaseScore, int32(table.MappedNum))
		}
	}
}

func (h *sessionCreditHandler) OnUserStandUp(standup *utilsplaza.UserStandUp) {}

func (h *sessionCreditHandler) OnTableRenew(item *utilsplaza.TableRenew) {}

func (h *sessionCreditHandler) OnDismissTable(table int) {}

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
