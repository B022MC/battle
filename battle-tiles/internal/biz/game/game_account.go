// internal/biz/game/account_usecase.go
package game

import (
	"battle-tiles/internal/consts"
	model "battle-tiles/internal/dal/model/game"
	repo "battle-tiles/internal/dal/repo/game"
	"battle-tiles/internal/infra/plaza"
	"context"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type GameAccountUseCase struct {
	accRepo            repo.GameAccountRepo
	accCtrlAccountRepo repo.GameCtrlAccountRepo
	accHouseRepo       repo.GameAccountHouseRepo
	sessRepo           repo.SessionRepo
	mgr                plaza.Manager
	log                *log.Helper
}

func NewGameAccountUseCase(
	acc repo.GameAccountRepo,
	ctrlAcc repo.GameCtrlAccountRepo,
	accHouse repo.GameAccountHouseRepo,
	sess repo.SessionRepo,
	mgr plaza.Manager,
	logger log.Logger,
) *GameAccountUseCase {
	return &GameAccountUseCase{
		accRepo:            acc,
		accCtrlAccountRepo: ctrlAcc,
		accHouseRepo:       accHouse,
		sessRepo:           sess,
		mgr:                mgr,
		log:                log.NewHelper(log.With(logger, "module", "usecase/game_account")),
	}
}

// 只绑定“我的”账号 普通用户才使用这个方法
func (uc *GameAccountUseCase) BindSingle(ctx context.Context, userID int32, mode consts.GameLoginMode, identifier, pwdMD5, nickname string) (*model.GameAccount, error) {
	// 探活并获取游戏用户信息
	info, err := uc.mgr.ProbeLoginWithInfo(ctx, mode, identifier, pwdMD5)
	if err != nil {
		return nil, err
	}
	// 普通用户仅允许1条（DB 触发器也兜底）
	if _, err := uc.accRepo.GetOneByUser(ctx, userID); err == nil {
		return nil, errors.New("you have already bound a game account")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	loginMode := "account"
	if mode == consts.GameLoginModeMobile {
		loginMode = "mobile"
	}
	a := &model.GameAccount{
		UserID:     userID,
		Account:    strings.TrimSpace(identifier),
		PwdMD5:     strings.ToUpper(strings.TrimSpace(pwdMD5)),
		Nickname:   nickname,
		IsDefault:  true,
		Status:     1,
		LoginMode:  loginMode,
		GameUserID: fmt.Sprintf("%d", info.UserID), // 保存游戏用户 ID
	}
	if err := uc.accRepo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (uc *GameAccountUseCase) GetMine(ctx context.Context, userID int32) (*model.GameAccount, error) {
	return uc.accRepo.GetOneByUser(ctx, userID)
}

func (uc *GameAccountUseCase) DeleteMine(ctx context.Context, userID int32) error {
	return uc.accRepo.DeleteByUser(ctx, userID)
}

func (uc *GameAccountUseCase) List(ctx context.Context, userID int32) ([]*model.GameAccount, error) {
	return uc.accRepo.ListByUser(ctx, userID)
}

func (uc *GameAccountUseCase) GetMyHouses(ctx context.Context, userID int32) ([]*model.GameAccountHouse, error) {
	// 先获取用户的游戏账号
	acc, err := uc.accRepo.GetOneByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	// 获取该账号绑定的所有店铺
	return uc.accHouseRepo.ListHousesByAccount(ctx, acc.Id)
}

func (uc *GameAccountUseCase) Verify(ctx context.Context, mode consts.GameLoginMode, identifier, pwdMD5 string) error {
	return uc.mgr.ProbeLogin(ctx, mode, identifier, pwdMD5)
}

// FixEmptyGameUserID 修复空的 game_user_id 字段
// 这个方法用于修复在修复代码之前注册的用户的 game_user_id
func (uc *GameAccountUseCase) FixEmptyGameUserID(ctx context.Context) (fixed int64, failed int64, err error) {
	// 查询所有 game_user_id 为空的记录
	accounts, err := uc.accRepo.ListByGameUserIDEmpty(ctx)
	if err != nil {
		return 0, 0, err
	}

	uc.log.Infof("Found %d accounts with empty game_user_id", len(accounts))

	// 逐个修复
	for _, acc := range accounts {
		// 调用 ProbeLoginWithInfo 获取游戏用户信息
		mode := consts.GameLoginModeAccount
		if acc.LoginMode == "mobile" {
			mode = consts.GameLoginModeMobile
		}

		info, err := uc.mgr.ProbeLoginWithInfo(ctx, mode, acc.Account, acc.PwdMD5)
		if err != nil {
			uc.log.Warnf("Failed to get game user info for account %s: %v", acc.Account, err)
			failed++
			continue
		}

		// 更新 game_user_id
		acc.GameUserID = fmt.Sprintf("%d", info.UserID)
		if err := uc.accRepo.Update(ctx, acc); err != nil {
			uc.log.Warnf("Failed to update account %s: %v", acc.Account, err)
			failed++
			continue
		}

		fixed++
		uc.log.Infof("Fixed account %s with game_user_id %s", acc.Account, acc.GameUserID)
	}

	return fixed, failed, nil
}
