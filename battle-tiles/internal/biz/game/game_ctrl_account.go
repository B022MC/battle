// internal/biz/game/game_ctrl_account.go
package game

import (
	"battle-tiles/internal/consts"
	model "battle-tiles/internal/dal/model/game"
	repo "battle-tiles/internal/dal/repo/game"
	"battle-tiles/internal/dal/req"
	"battle-tiles/internal/dal/resp"
	"battle-tiles/internal/dal/vo/game"
	"battle-tiles/internal/infra/plaza"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CtrlAccountUseCase struct {
	ctrlRepo repo.GameCtrlAccountRepo
	linkRepo repo.GameCtrlAccountHouseRepo
	accRepo  repo.GameAccountRepo
	mgr      plaza.Manager
	log      *log.Helper
}

func NewCtrlAccountUseCase(
	c repo.GameCtrlAccountRepo,
	link repo.GameCtrlAccountHouseRepo,
	acc repo.GameAccountRepo,
	mgr plaza.Manager,
	logger log.Logger,
) *CtrlAccountUseCase {
	return &CtrlAccountUseCase{
		ctrlRepo: c,
		linkRepo: link,
		accRepo:  acc,
		mgr:      mgr,
		log:      log.NewHelper(log.With(logger, "module", "usecase/ctrl_account")),
	}
}

// 仅创建/更新中控（不绑定店铺）—— 先探测再落库
func (uc *CtrlAccountUseCase) CreateOrUpdateCtrl(ctx context.Context, mode consts.GameLoginMode, identifier, pwdMD5 string, status int32, operatorUserID int32) (*model.GameCtrlAccount, error) {
	id := strings.TrimSpace(identifier)
	md5 := strings.ToUpper(strings.TrimSpace(pwdMD5))
	if id == "" || len(md5) != 32 {
		return nil, errors.New("invalid identifier or md5")
	}
	// 探测登录（拿到游戏端 UserID）
	info, err := uc.mgr.ProbeLoginWithInfo(ctx, mode, id, md5)
	if err != nil {
		return nil, errors.Wrap(err, "probe login failed")
	}
	now := time.Now()
	m := &model.GameCtrlAccount{
		LoginMode:    int32(mode),
		Identifier:   id,
		PwdMD5:       md5,
		GameUserID:   fmt.Sprintf("%d", info.UserID),
		GameID:       fmt.Sprintf("%d", info.GameID),
		Status:       status,
		LastVerifyAt: &now,
	}
	m, err = uc.ctrlRepo.Upsert(ctx, m)
	if err != nil {
		return nil, err
	}
	// 确保存在一个与中控关联的 game_account，便于后续绑定店铺（如果还没有）
	if uc.accRepo != nil {
		if cnt, _ := uc.accRepo.CountByCtrl(ctx, m.Id); cnt == 0 {
			loginMode := "account"
			if mode == consts.GameLoginModeMobile {
				loginMode = "mobile"
			}
			_ = uc.accRepo.Create(ctx, &model.GameAccount{
				UserID:        &operatorUserID, // 指针类型
				Account:       id,
				PwdMD5:        md5,
				Nickname:      "",
				IsDefault:     true,
				Status:        1,
				LoginMode:     loginMode,
				CtrlAccountID: &m.Id,
				GamePlayerID:  fmt.Sprintf("%d", info.GameID), // 设置游戏玩家ID
			})
		}
	}
	return m, nil
}

// 绑定/解绑
func (uc *CtrlAccountUseCase) BindCtrlToHouse(ctx context.Context, ctrlID int32, houseGID int32, status int32, _ string, operatorUserID int32) error {
	// 确保该 ctrl 账号下存在可绑定的 game_account；若不存在则自动补一条默认
	if uc.accRepo != nil {
		if cnt, _ := uc.accRepo.CountByCtrl(ctx, ctrlID); cnt == 0 {
			if ctrl, err := uc.ctrlRepo.Get(ctx, ctrlID); err == nil && ctrl != nil {
				loginMode := "account"
				if ctrl.LoginMode == int32(consts.GameLoginModeMobile) {
					loginMode = "mobile"
				}
				_ = uc.accRepo.Create(ctx, &model.GameAccount{
					UserID:        &operatorUserID, // 指针类型
					Account:       ctrl.Identifier,
					PwdMD5:        ctrl.PwdMD5,
					Nickname:      "",
					IsDefault:     true,
					Status:        1,
					LoginMode:     loginMode,
					CtrlAccountID: &ctrlID,
				})
			}
		}
	}
	// 通过 ctrlID 选择其一个 game_account 进行绑定（优先默认账号）
	err := uc.linkRepo.BindByCtrl(ctx, ctrlID, houseGID, status)
	if err != nil {
		// 如果是因为已绑定其他店铺导致的错误，返回更友好的错误信息
		if errors.Is(err, gorm.ErrInvalidData) {
			return errors.New("该中控账号已绑定其他店铺，每个中控账号只能绑定一个店铺")
		}
		return err
	}
	return nil
}

func (uc *CtrlAccountUseCase) UnbindCtrlFromHouse(ctx context.Context, ctrlID int32, houseGID int32) error {
	return uc.linkRepo.UnbindByCtrl(ctx, ctrlID, houseGID)
}

// UpdateStatus 更新中控账号状态
func (uc *CtrlAccountUseCase) UpdateStatus(ctx context.Context, ctrlID int32, status int32) error {
	// 验证状态值
	if status != 0 && status != 1 {
		return errors.New("invalid status, must be 0 or 1")
	}

	// 更新数据库状态
	if err := uc.ctrlRepo.UpdateStatus(ctx, ctrlID, status); err != nil {
		return errors.Wrap(err, "update ctrl account status failed")
	}

	uc.log.Infof("中控账号 %d 状态已更新为 %d", ctrlID, status)
	return nil
}

// 查询：按店铺列出所有中控
func (uc *CtrlAccountUseCase) ListCtrlByHouse(ctx context.Context, houseGID int32) ([]*model.GameCtrlAccount, error) {
	return uc.linkRepo.ListByHouse(ctx, houseGID)
}

/*** 兼容旧接口：创建+绑定（/shops/ctrlAccounts [POST]） ***/
func (uc *CtrlAccountUseCase) CreateAndBind(ctx context.Context, houseGID int32, mode consts.GameLoginMode, identifier, pwdMD5 string, status int32, operatorUserID int32) (*model.GameCtrlAccount, error) {
	m, err := uc.CreateOrUpdateCtrl(ctx, mode, identifier, pwdMD5, status, operatorUserID)
	if err != nil {
		return nil, err
	}
	if err := uc.BindCtrlToHouse(ctx, m.Id, houseGID, status, "", operatorUserID); err != nil {
		return nil, err
	}
	return m, nil
}
func (uc *CtrlAccountUseCase) ListAll(ctx context.Context, f req.CtrlListFilter, page, size int) ([]*resp.CtrlAccountWithHouses, int64, error) {
	items, total, err := uc.ctrlRepo.List(ctx, game.CtrlAccountListCond{
		LoginMode: f.LoginMode,
		Status:    f.Status,
		Keyword:   strings.TrimSpace(f.Keyword),
		Page:      page,
		Size:      size,
	})
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return []*resp.CtrlAccountWithHouses{}, total, nil
	}
	ids := make([]int32, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.Id)
	}

	houseMap, _ := uc.linkRepo.ListHouseMapByCtrlIDs(ctx, ids)
	out := make([]*resp.CtrlAccountWithHouses, 0, len(items))
	for _, it := range items {
		out = append(out, &resp.CtrlAccountWithHouses{
			GameCtrlAccount: it,
			Houses:          houseMap[it.Id],
		})
	}
	return out, total, nil
}

// 列出所有出现过的店铺号（去重），基于 game_account_house
func (uc *CtrlAccountUseCase) ListDistinctHouses(ctx context.Context) ([]int32, error) {
	return uc.linkRepo.ListDistinctHouses(ctx)
}

// Delete 删除中控账号
func (uc *CtrlAccountUseCase) Delete(ctx context.Context, ctrlID int32) error {
	// 先检查是否存在
	ctrl, err := uc.ctrlRepo.Get(ctx, ctrlID)
	if err != nil {
		return errors.Wrap(err, "get ctrl account failed")
	}

	// 获取所有绑定的店铺
	houses, err := uc.linkRepo.ListByCtrlID(ctx, ctrlID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "list houses failed")
	}

	// 停止所有会话
	for _, house := range houses {
		if _, ok := uc.mgr.GetAnyByHouse(int(house.HouseGID)); ok {
			uc.log.Infof("停止中控账号 %d 的会话 house=%d", ctrlID, house.HouseGID)
			uc.mgr.StopUser(1, int(house.HouseGID))
		}
	}

	// 删除所有绑定关系
	if err := uc.linkRepo.DeleteByCtrlID(ctx, ctrlID); err != nil {
		return errors.Wrap(err, "delete bindings failed")
	}

	// 删除中控账号
	if err := uc.ctrlRepo.Delete(ctx, ctrlID); err != nil {
		return errors.Wrap(err, "delete ctrl account failed")
	}

	uc.log.Infof("已删除中控账号 %d (identifier=%s)", ctrl.Id, ctrl.Identifier)
	return nil
}
