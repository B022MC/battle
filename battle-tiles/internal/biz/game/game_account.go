// internal/biz/game/account_usecase.go
package game

import (
	"battle-tiles/internal/consts"
	model "battle-tiles/internal/dal/model/game"
	repo "battle-tiles/internal/dal/repo/game"
	"battle-tiles/internal/infra/plaza"
	"context"
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
	// 探活
	if err := uc.Verify(ctx, mode, identifier, pwdMD5); err != nil {
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
		UserID:    userID,
		Account:   strings.TrimSpace(identifier),
		PwdMD5:    strings.ToUpper(strings.TrimSpace(pwdMD5)),
		Nickname:  nickname,
		IsDefault: true,
		Status:    1,
		LoginMode: loginMode,
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

func (uc *GameAccountUseCase) Verify(ctx context.Context, mode consts.GameLoginMode, identifier, pwdMD5 string) error {
	return uc.mgr.ProbeLogin(ctx, mode, identifier, pwdMD5)
}
