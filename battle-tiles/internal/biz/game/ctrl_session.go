// internal/biz/game/ctrl_session.go
package game

import (
	"battle-tiles/internal/consts"
	repo "battle-tiles/internal/dal/repo/game"
	"battle-tiles/internal/infra/plaza"
	plazaUtils "battle-tiles/internal/utils/plaza"
	pdb "battle-tiles/pkg/plugin/dbx"
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

type CtrlSessionUseCase struct {
	ctrlRepo repo.GameCtrlAccountRepo      // 中控账号
	linkRepo repo.GameCtrlAccountHouseRepo // 中控-店铺 绑定
	// 可选：houseRepo 若需要“先连再落库创建店铺”，放开注入后在回调里 Ensure
	// houseRepo repo.GameHouseRepo

	sessRepo      repo.SessionRepo
	mgr           plaza.Manager
	syncMgr       *BattleSyncManager      // 战绩同步管理器
	creditHandler *RoomCreditEventHandler // 房间额度事件处理器（可选）
	log           *log.Helper
}

func NewCtrlSessionUseCase(
	ctrl repo.GameCtrlAccountRepo,
	link repo.GameCtrlAccountHouseRepo,
	sess repo.SessionRepo,
	mgr plaza.Manager,
	syncMgr *BattleSyncManager,
	logger log.Logger,
	creditHandler *RoomCreditEventHandler,
) *CtrlSessionUseCase {
	uc := &CtrlSessionUseCase{
		ctrlRepo:      ctrl,
		linkRepo:      link,
		sessRepo:      sess,
		mgr:           mgr,
		syncMgr:       syncMgr,
		creditHandler: creditHandler,
		log:           log.NewHelper(log.With(logger, "module", "usecase/ctrl_session")),
	}

	// 设置重连失败回调 - 自动停用中控账号
	mgr.SetReconnectFailedCallback(func(houseGID int, retryCount int) int32 {
		// 根据 houseGID 查找对应的中控账号ID
		ctx := context.Background()
		ctrls, err := uc.linkRepo.ListByHouse(ctx, int32(houseGID))
		if err != nil || len(ctrls) == 0 {
			uc.log.Errorf("查找中控账号失败 house=%d err=%v", houseGID, err)
			return 0
		}

		ctrlID := ctrls[0].Id
		// 停用中控账号
		if err := uc.ctrlRepo.UpdateStatus(ctx, ctrlID, 0); err != nil {
			uc.log.Errorf("停用中控账号失败 ctrlID=%d err=%v", ctrlID, err)
			return 0
		}

		uc.log.Infof("已自动停用中控账号 ctrlID=%d house=%d 重试次数=%d", ctrlID, houseGID, retryCount)
		return ctrlID
	})

	return uc
}

// WithCreditHandler 注入房间额度事件处理器（可选）
func (uc *CtrlSessionUseCase) WithCreditHandler(h *RoomCreditEventHandler) *CtrlSessionUseCase {
	uc.creditHandler = h
	return uc
}

// 先连 -> 成功收到登录/房间回调后再做落库（如需）
func (uc *CtrlSessionUseCase) StartSession(ctx context.Context, userID int32, ctrlAccID int32, houseGID int32) error {
	// 1) 取中控凭证
	ctrl, err := uc.ctrlRepo.Get(ctx, ctrlAccID)
	if err != nil {
		return errors.Wrap(err, "get ctrl account")
	}

	// 2) 映射登录方式
	mode := consts.GameLoginModeAccount
	if ctrl.LoginMode != int32(consts.GameLoginModeAccount) {
		mode = consts.GameLoginModeMobile
	}

	// 3) 若该用户在该店铺已有会话（无论是否已在线/登录中），直接返回（防重复连接/重启）
	if _, ok := uc.mgr.Get(int(userID), int(houseGID)); ok {
		return nil
	}

	// 4) 从 ctx 中获取 platform
	platform := ""
	if v := ctx.Value(pdb.CtxDBKey); v != nil {
		if s, ok := v.(string); ok {
			platform = s
		}
	}

	// 5) 包一个 bootstrap handler：连接成功/收到房间列表时，做你想做的落库动作（可选）
	h := uc.newBootstrapHandler(userID, ctrl.Id, houseGID, platform)

	// 5) 不再强制关闭旧会话，改为“存在则更新、否则插入”

	// 6) 启动
	// 若未存 game_player_id，拒绝启动，提示先验证账号
	if strings.TrimSpace(ctrl.GameUserID) == "" {
		return errors.New("game_player_id not set, please verify account first")
	}
	gameUID := 0
	if v, err := strconv.Atoi(ctrl.GameUserID); err == nil {
		gameUID = v
	}
	if err := uc.mgr.StartUser(ctx, int(userID), int(houseGID), mode, ctrl.Identifier, strings.ToUpper(ctrl.PwdMD5), gameUID, h); err != nil {
		// 启动失败：更新现有记录为 error 状态，而不是插入新记录
		_ = uc.sessRepo.UpsertErrorByHouse(ctx, ctrl.Id, userID, houseGID, err.Error())
		return errors.Wrap(err, "start session")
	}

	// 启动战绩同步，传入带 platform 的 context
	uc.syncMgr.StartSync(ctx, int(userID), int(houseGID))

	// 7) 成功：若存在该店铺记录则更新最新一条为 online，否则插入
	return uc.sessRepo.UpsertOnlineByHouse(ctx, ctrl.Id, userID, houseGID)
}

func (uc *CtrlSessionUseCase) StopSession(ctx context.Context, userID int32, ctrlAccID int32, houseGID int32) error {
	if _, err := uc.ctrlRepo.Get(ctx, (ctrlAccID)); err != nil {
		return errors.Wrap(err, "get ctrl account")
	}

	uc.mgr.StopUser(int(userID), int(houseGID))

	// 停止战绩同步
	uc.syncMgr.StopSync(int(userID), int(houseGID))

	// 更新现有记录为 offline，而不是插入新记录
	return uc.sessRepo.SetOfflineByHouse(ctx, houseGID)
}

// ===== 可选：登录/房间回调后，再做落库 =====

type noopHandler struct{}

func (*noopHandler) OnSessionRestarted(*plazaUtils.Session)        {}
func (*noopHandler) OnMemberListUpdated([]*plazaUtils.GroupMember) {}
func (*noopHandler) OnMemberInserted(*plazaUtils.MemberInserted)   {}
func (*noopHandler) OnMemberDeleted(*plazaUtils.MemberDeleted)     {}
func (*noopHandler) OnMemberRightUpdated(string, int, bool)        {}
func (*noopHandler) OnLoginDone(bool)                              {}
func (*noopHandler) OnRoomListUpdated([]*plazaUtils.TableInfo)     {}
func (*noopHandler) OnUserSitDown(*plazaUtils.UserSitDown)         {}
func (*noopHandler) OnUserStandUp(*plazaUtils.UserStandUp)         {}
func (*noopHandler) OnTableRenew(*plazaUtils.TableRenew)           {}
func (*noopHandler) OnDismissTable(int)                            {}
func (*noopHandler) OnAppliesForHouse([]*plazaUtils.ApplyInfo)     {}
func (*noopHandler) OnReconnectFailed(int, int)                    {}

// compositeHandler 组合多个 handler
type compositeHandler struct {
	bootstrap *bootstrapHandler
	credit    *sessionCreditHandler
	log       *log.Helper
}

func (h *compositeHandler) OnSessionRestarted(session *plazaUtils.Session) {
	h.bootstrap.OnSessionRestarted(session)
	h.credit.OnSessionRestarted(session)
}

func (h *compositeHandler) OnMemberListUpdated(members []*plazaUtils.GroupMember) {
	h.bootstrap.OnMemberListUpdated(members)
	h.credit.OnMemberListUpdated(members)
}

func (h *compositeHandler) OnMemberInserted(member *plazaUtils.MemberInserted) {
	h.bootstrap.OnMemberInserted(member)
	h.credit.OnMemberInserted(member)
}

func (h *compositeHandler) OnMemberDeleted(member *plazaUtils.MemberDeleted) {
	h.bootstrap.OnMemberDeleted(member)
	h.credit.OnMemberDeleted(member)
}

func (h *compositeHandler) OnMemberRightUpdated(key string, memberID int, success bool) {
	h.bootstrap.OnMemberRightUpdated(key, memberID, success)
	h.credit.OnMemberRightUpdated(key, memberID, success)
}

func (h *compositeHandler) OnLoginDone(success bool) {
	h.log.Infof("[compositeHandler] OnLoginDone: success=%v", success)
	h.bootstrap.OnLoginDone(success)
	h.credit.OnLoginDone(success)
}

func (h *compositeHandler) OnRoomListUpdated(tables []*plazaUtils.TableInfo) {
	h.log.Infof("[compositeHandler] OnRoomListUpdated: 桌子数=%d", len(tables))
	h.bootstrap.OnRoomListUpdated(tables)
	h.credit.OnRoomListUpdated(tables)
}

func (h *compositeHandler) OnUserSitDown(sitdown *plazaUtils.UserSitDown) {
	h.log.Infof("[compositeHandler] OnUserSitDown: userID=%d, gameID=%d, mappedNum=%d", sitdown.UserID, sitdown.GameID, sitdown.MappedNum)
	h.bootstrap.OnUserSitDown(sitdown)
	h.credit.OnUserSitDown(sitdown)
}

func (h *compositeHandler) OnUserStandUp(standup *plazaUtils.UserStandUp) {
	h.bootstrap.OnUserStandUp(standup)
	h.credit.OnUserStandUp(standup)
}

func (h *compositeHandler) OnTableRenew(item *plazaUtils.TableRenew) {
	h.bootstrap.OnTableRenew(item)
	h.credit.OnTableRenew(item)
}

func (h *compositeHandler) OnDismissTable(table int) {
	h.bootstrap.OnDismissTable(table)
	h.credit.OnDismissTable(table)
}

func (h *compositeHandler) OnAppliesForHouse(applyInfos []*plazaUtils.ApplyInfo) {
	h.log.Infof("[compositeHandler] OnAppliesForHouse: 申请数=%d", len(applyInfos))
	h.bootstrap.OnAppliesForHouse(applyInfos)
	h.credit.OnAppliesForHouse(applyInfos)
}

func (h *compositeHandler) OnReconnectFailed(houseGID int, retryCount int) {
	h.bootstrap.OnReconnectFailed(houseGID, retryCount)
	h.credit.OnReconnectFailed(houseGID, retryCount)
}

type bootstrapHandler struct {
	noopHandler
	once      sync.Once
	ctrlID    int32
	houseGID  int32
	bootstrap func()
}

func (h *bootstrapHandler) OnLoginDone(ok bool) {
	if ok {
		h.once.Do(h.bootstrap)
	}
	h.noopHandler.OnLoginDone(ok)
}
func (h *bootstrapHandler) OnRoomListUpdated(ts []*plazaUtils.TableInfo) {
	h.once.Do(h.bootstrap)
	h.noopHandler.OnRoomListUpdated(ts)
}

func (uc *CtrlSessionUseCase) newBootstrapHandler(userID int32, ctrlID int32, houseGID int32, platform string) plaza.Handler {
	h := &bootstrapHandler{
		ctrlID:   ctrlID,
		houseGID: houseGID,
		bootstrap: func() {
			// 按你的需求：连接成功/房间有了 → 再确保店铺落库、绑定关系等
			// 示例（伪代码，按你的仓储接口替换）：
			// ctx := context.Background()
			// _ = uc.houseRepo.Ensure(ctx, int(houseGID), "")
			// _ = uc.linkRepo.Bind(ctx, &model.GameCtrlAccountHouse{
			//     CtrlAccountID: ctrlID,
			//     HouseGID:      houseGID,
			//     Status:        1,
			//     Alias:         "",
			// })
		},
	}

	// 如果有费用检查处理器，创建组合handler
	if uc.creditHandler != nil {
		uc.log.Infof("[费用检查] 创建组合handler: userID=%d, houseGID=%d, platform=%s", userID, houseGID, platform)
		creditH := uc.creditHandler.CreateHandler(int(userID), houseGID, platform)
		return &compositeHandler{
			bootstrap: h,
			credit:    creditH,
			log:       uc.log,
		}
	}

	uc.log.Infof("[警告] creditHandler为空，未启用费用检查: userID=%d, houseGID=%d", userID, houseGID)
	return h
}
