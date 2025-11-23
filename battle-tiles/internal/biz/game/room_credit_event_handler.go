package game

import (
	"battle-tiles/internal/infra/plaza"
	utilsplaza "battle-tiles/internal/utils/plaza"
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
func (h *RoomCreditEventHandler) CreateHandler(userID int, houseGID int32) *sessionCreditHandler {
	return &sessionCreditHandler{
		userID:   userID,
		houseGID: houseGID,
		creditUC: h.creditUC,
		plazaMgr: h.plazaMgr,
		log:      h.log,
	}
}

// sessionCreditHandler 实现 plaza.IPlazaHandler 接口
// 每个会话有自己的 Handler 实例，携带 userID 和 houseGID 信息
type sessionCreditHandler struct {
	userID   int
	houseGID int32
	creditUC *RoomCreditLimitUseCase
	plazaMgr plaza.Manager
	log      *log.Helper
}

// OnUserSitDown 处理玩家坐下事件
// 当收到游戏服务器的坐下事件时，检查玩家余额是否满足额度要求
func (h *sessionCreditHandler) OnUserSitDown(sitdown *utilsplaza.UserSitDown) {
	h.log.Debugf("OnUserSitDown: gameID=%d, mappedNum=%d",
		sitdown.UserID, sitdown.MappedNum)

	ctx := context.Background()

	// 1. 从 plaza.Manager 获取会话
	session, ok := h.plazaMgr.Get(h.userID, int(h.houseGID))
	if !ok {
		h.log.Errorf("Session not found: userID=%d, houseGID=%d", h.userID, h.houseGID)
		return
	}

	// 2. 从会话的桌子列表中查找对应的桌子信息
	tables := session.ListTables()
	var kindID, baseScore int
	found := false
	for _, table := range tables {
		if table.MappedNum == int(sitdown.MappedNum) {
			kindID = table.KindID
			baseScore = table.BaseScore
			found = true
			break
		}
	}

	if !found {
		h.log.Warnf("Table not found in session: mappedNum=%d", sitdown.MappedNum)
		return
	}

	h.log.Debugf("OnUserSitDown: gameID=%d, kindID=%d, baseScore=%d, mappedNum=%d",
		sitdown.UserID, kindID, baseScore, sitdown.MappedNum)

	// 3. 调用额度检查逻辑
	if err := h.creditUC.HandlePlayerSitDown(
		ctx,
		h.userID,
		h.houseGID,
		int32(sitdown.UserID),    // gameID
		int32(kindID),            // kindID
		int32(baseScore),         // baseScore
		int32(sitdown.MappedNum), // tableNum
	); err != nil {
		h.log.Errorf("HandlePlayerSitDown failed: %v", err)
	}
}

// 以下是 IPlazaHandler 接口的其他方法实现（空实现）

func (h *sessionCreditHandler) OnSessionRestarted(session *utilsplaza.Session) {}

func (h *sessionCreditHandler) OnMemberListUpdated(members []*utilsplaza.GroupMember) {}

func (h *sessionCreditHandler) OnMemberInserted(member *utilsplaza.MemberInserted) {}

func (h *sessionCreditHandler) OnMemberDeleted(member *utilsplaza.MemberDeleted) {}

func (h *sessionCreditHandler) OnMemberRightUpdated(key string, memberID int, success bool) {}

func (h *sessionCreditHandler) OnLoginDone(success bool) {}

func (h *sessionCreditHandler) OnRoomListUpdated(tables []*utilsplaza.TableInfo) {}

func (h *sessionCreditHandler) OnUserStandUp(standup *utilsplaza.UserStandUp) {}

func (h *sessionCreditHandler) OnTableRenew(item *utilsplaza.TableRenew) {}

func (h *sessionCreditHandler) OnDismissTable(table int) {}

func (h *sessionCreditHandler) OnAppliesForHouse(applyInfos []*utilsplaza.ApplyInfo) {}

func (h *sessionCreditHandler) OnReconnectFailed(houseGID int, retryCount int) {}
