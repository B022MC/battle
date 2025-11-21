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
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type GameAccountUseCase struct {
	accRepo            repo.GameAccountRepo
	accCtrlAccountRepo repo.GameCtrlAccountRepo
	accHouseRepo       repo.GameAccountHouseRepo
	accountGroupRepo   repo.GameAccountGroupRepo // 用于查询 game_account_group
	sessRepo           repo.SessionRepo
	mgr                plaza.Manager
	log                *log.Helper
}

func NewGameAccountUseCase(
	acc repo.GameAccountRepo,
	ctrlAcc repo.GameCtrlAccountRepo,
	accHouse repo.GameAccountHouseRepo,
	accountGroup repo.GameAccountGroupRepo,
	sess repo.SessionRepo,
	mgr plaza.Manager,
	logger log.Logger,
) *GameAccountUseCase {
	return &GameAccountUseCase{
		accRepo:            acc,
		accCtrlAccountRepo: ctrlAcc,
		accHouseRepo:       accHouse,
		accountGroupRepo:   accountGroup,
		sessRepo:           sess,
		mgr:                mgr,
		log:                log.NewHelper(log.With(logger, "module", "usecase/game_account")),
	}
}

// 只绑定"我的"账号 普通用户才使用这个方法
func (uc *GameAccountUseCase) BindSingle(ctx context.Context, userID int32, mode consts.GameLoginMode, identifier, pwdMD5, nickname string) (*model.GameAccount, error) {
	// 探活并获取游戏用户信息
	info, err := uc.mgr.ProbeLoginWithInfo(ctx, mode, identifier, pwdMD5)
	if err != nil {
		return nil, err
	}

	if _, err := uc.accRepo.GetOneByUser(ctx, userID); err == nil {
		return nil, errors.New("you have already bound a game account")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 检查游戏账号是否已存在
	existingAccount, err := uc.accRepo.GetByAccount(ctx, identifier)
	if err == nil && existingAccount != nil {
		// 游戏账号已存在
		if existingAccount.UserID != nil && *existingAccount.UserID != 0 {
			// 已被其他用户绑定
			return nil, errors.New("this game account is already bound to another user")
		}
		// 游戏账号存在但未绑定用户，更新绑定并同步 GamePlayerID
		existingAccount.UserID = &userID
		existingAccount.IsDefault = true
		existingAccount.GamePlayerID = cast.ToString(info.GameID) // 同步更新 GamePlayerID
		existingAccount.Nickname = nickname                       // 同步更新昵称
		if err := uc.accRepo.Update(ctx, existingAccount); err != nil {
			return nil, err
		}
		return existingAccount, nil
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 游戏账号不存在，创建新账号
	loginMode := "account"
	if mode == consts.GameLoginModeMobile {
		loginMode = "mobile"
	}
	a := &model.GameAccount{
		UserID:       &userID, // 指针类型
		Account:      strings.TrimSpace(identifier),
		PwdMD5:       strings.ToUpper(strings.TrimSpace(pwdMD5)),
		Nickname:     nickname,
		IsDefault:    true,
		Status:       1,
		LoginMode:    loginMode,
		GamePlayerID: cast.ToString(info.GameID), // 保存游戏账号 ID
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

func (uc *GameAccountUseCase) GetMyHouses(ctx context.Context, userID int32) (*model.GameAccountHouse, error) {
	// 先获取用户的游戏账号
	acc, err := uc.accRepo.GetOneByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 从 game_account_group 查询店铺信息（一个用户在一个店铺只能加入一个圈子）
	// 使用 game_player_id 查询
	if acc.GamePlayerID == "" {
		return nil, fmt.Errorf("游戏账号缺少 game_player_id")
	}
	accountGroups, err := uc.accountGroupRepo.ListByGamePlayer(ctx, acc.GamePlayerID)
	if err != nil {
		return nil, err
	}

	if len(accountGroups) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// 取第一个激活的圈子关系，转换为 GameAccountHouse 格式
	for _, ag := range accountGroups {
		if ag.Status == model.AccountGroupStatusActive {
			// 转换为 GameAccountHouse 结构（保持 API 响应格式不变）
			// 使用 game_account.id (acc.Id)，而非 GamePlayerID
			return &model.GameAccountHouse{
				Id:            ag.Id,
				GameAccountID: acc.Id, // 使用当前 game_account 的 ID
				HouseGID:      ag.HouseGID,
				IsDefault:     true,
				Status:        1,
			}, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

func (uc *GameAccountUseCase) Verify(ctx context.Context, mode consts.GameLoginMode, identifier, pwdMD5 string) error {
	return uc.mgr.ProbeLogin(ctx, mode, identifier, pwdMD5)
}
