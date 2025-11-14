// internal/biz/game/ctrl_session.go
package game

import (
	"battle-tiles/internal/consts"
	model "battle-tiles/internal/dal/model/game"
	repo "battle-tiles/internal/dal/repo/game"
	"battle-tiles/internal/infra/plaza"
	plazaUtils "battle-tiles/internal/utils/plaza"
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

type CtrlSessionUseCase struct {
	ctrlRepo repo.GameCtrlAccountRepo      // 中控账号
	linkRepo repo.GameCtrlAccountHouseRepo // 中控-店铺 绑定
	// 可选：houseRepo 若需要“先连再落库创建店铺”，放开注入后在回调里 Ensure
	// houseRepo repo.GameHouseRepo

	sessRepo  repo.SessionRepo
	mgr       plaza.Manager
	syncMgr   *BattleSyncManager // 战绩同步管理器
	log       *log.Helper
}

func NewCtrlSessionUseCase(
	ctrl repo.GameCtrlAccountRepo,
	link repo.GameCtrlAccountHouseRepo,
	sess repo.SessionRepo,
	mgr plaza.Manager,
	syncMgr *BattleSyncManager,
	logger log.Logger,
) *CtrlSessionUseCase {
	return &CtrlSessionUseCase{
		ctrlRepo: ctrl,
		linkRepo: link,
		sessRepo: sess,
		mgr:      mgr,
		syncMgr:  syncMgr,
		log:      log.NewHelper(log.With(logger, "module", "usecase/ctrl_session")),
	}
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

	// 4) 包一个 bootstrap handler：连接成功/收到房间列表时，做你想做的落库动作（可选）
	h := uc.newBootstrapHandler(ctrl.Id, houseGID)

	// 5) 不再强制关闭旧会话，改为“存在则更新、否则插入”

	// 6) 启动
	// 若未存 game_user_id，拒绝启动，提示先验证账号
	if strings.TrimSpace(ctrl.GameUserID) == "" {
		return errors.New("game_user_id not set, please verify account first")
	}
	gameUID := 0
	if v, err := strconv.Atoi(ctrl.GameUserID); err == nil {
		gameUID = v
	}
	if err := uc.mgr.StartUser(ctx, int(userID), int(houseGID), mode, ctrl.Identifier, strings.ToUpper(ctrl.PwdMD5), gameUID, h); err != nil {
		// 失败记录
		_ = uc.sessRepo.Insert(ctx, &model.GameSession{
			GameCtrlAccountID: ctrl.Id,
			UserID:            userID,
			HouseGID:          houseGID,
			State:             "error",
			ErrorMsg:          err.Error(),
		})
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

	now := time.Now()
	return uc.sessRepo.Insert(ctx, &model.GameSession{
		GameCtrlAccountID: (ctrlAccID),
		UserID:            userID,
		HouseGID:          houseGID,
		State:             "offline",
		EndAt:             &now,
	})
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

func (uc *CtrlSessionUseCase) newBootstrapHandler(ctrlID int32, houseGID int32) plaza.Handler {
	return &bootstrapHandler{
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
}
