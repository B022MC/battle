// internal/biz/game/game_account_group.go
package game

import (
	"context"
	"fmt"

	model "battle-tiles/internal/dal/model/game"
	repo "battle-tiles/internal/dal/repo/game"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type GameAccountGroupUseCase struct {
	accountRepo      repo.GameAccountRepo
	accountGroupRepo repo.GameAccountGroupRepo
	shopGroupRepo    repo.ShopGroupRepo
	log              *log.Helper
}

func NewGameAccountGroupUseCase(
	accountRepo repo.GameAccountRepo,
	accountGroupRepo repo.GameAccountGroupRepo,
	shopGroupRepo repo.ShopGroupRepo,
	logger log.Logger,
) *GameAccountGroupUseCase {
	return &GameAccountGroupUseCase{
		accountRepo:      accountRepo,
		accountGroupRepo: accountGroupRepo,
		shopGroupRepo:    shopGroupRepo,
		log:              log.NewHelper(log.With(logger, "module", "usecase/game_account_group")),
	}
}

// FindOrCreateGameAccount 根据游戏玩家ID查找或创建游戏账号
// 如果游戏账号不存在，创建一个未绑定用户的游戏账号
func (uc *GameAccountGroupUseCase) FindOrCreateGameAccount(
	ctx context.Context,
	gamePlayerID string,
	account string,
	nickname string,
) (*model.GameAccount, error) {
	// 先尝试查找
	acc, err := uc.accountRepo.GetByGamePlayerID(ctx, gamePlayerID)
	if err == nil {
		return acc, nil
	}

	// 如果不是记录不存在的错误，返回错误
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 创建新的游戏账号（user_id 为 NULL）
	newAcc := &model.GameAccount{
		UserID:       nil, // 未绑定用户
		GamePlayerID: gamePlayerID,
		Account:      account,
		Nickname:     nickname,
		Status:       1,
		LoginMode:    "account",
	}

	if err := uc.accountRepo.Create(ctx, newAcc); err != nil {
		return nil, err
	}

	return newAcc, nil
}

// EnsureGroupForAdmin 确保管理员有圈子，如果没有则创建
func (uc *GameAccountGroupUseCase) EnsureGroupForAdmin(
	ctx context.Context,
	houseGID int32,
	adminUserID int32,
	adminNickname string,
) (*model.GameShopGroup, error) {
	// 先查询是否已有圈子
	group, err := uc.shopGroupRepo.GetByAdmin(ctx, houseGID, adminUserID)
	if err == nil {
		return group, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 创建默认圈子
	defaultGroupName := fmt.Sprintf("%s的圈子", adminNickname)
	newGroup := &model.GameShopGroup{
		HouseGID:    houseGID,
		GroupName:   defaultGroupName,
		AdminUserID: adminUserID,
		Description: "自动创建",
		IsActive:    true,
	}

	if err := uc.shopGroupRepo.Create(ctx, newGroup); err != nil {
		return nil, err
	}

	return newGroup, nil
}

// AddGameAccountToGroup 将游戏账号加入圈子
func (uc *GameAccountGroupUseCase) AddGameAccountToGroup(
	ctx context.Context,
	gameAccountID int32,
	houseGID int32,
	groupID int32,
	groupName string,
	adminUserID int32,
	approvedByUserID int32,
) error {
	// 先获取 game_player_id
	gameAccount, err := uc.accountRepo.GetByID(ctx, gameAccountID)
	if err != nil || gameAccount == nil || gameAccount.GamePlayerID == "" {
		return fmt.Errorf("游戏账号 %d 不存在或缺少 game_player_id", gameAccountID)
	}

	// 检查游戏玩家是否已在该店铺的圈子中
	exists, err := uc.accountGroupRepo.ExistsByGamePlayerAndHouse(ctx, gameAccount.GamePlayerID, houseGID)
	if err != nil {
		return err
	}

	if exists {
		// 已存在，更新状态为激活
		existing, err := uc.accountGroupRepo.GetByGamePlayerAndHouse(ctx, gameAccount.GamePlayerID, houseGID)
		if err != nil {
			return err
		}
		return uc.accountGroupRepo.UpdateStatus(ctx, existing.Id, model.AccountGroupStatusActive)
	}

	// 创建新的游戏玩家圈子关系
	adminUserIDPtr := &adminUserID
	accountGroup := &model.GameAccountGroup{
		GameAccountID: gameAccountID,
		GamePlayerID:  gameAccount.GamePlayerID,
		HouseGID:      houseGID,
		GroupID:       groupID,
		GroupName:     groupName,
		AdminUserID:   adminUserIDPtr,
		Status:        model.AccountGroupStatusActive,
	}

	return uc.accountGroupRepo.Create(ctx, accountGroup)
}

// RemoveGameAccountFromGroup 将游戏账号从圈子移除
func (uc *GameAccountGroupUseCase) RemoveGameAccountFromGroup(
	ctx context.Context,
	gameAccountID int32,
	houseGID int32,
) error {
	// 获取 game_player_id
	gameAccount, err := uc.accountRepo.GetByID(ctx, gameAccountID)
	if err != nil || gameAccount == nil || gameAccount.GamePlayerID == "" {
		return fmt.Errorf("游戏账号 %d 不存在或缺少 game_player_id", gameAccountID)
	}
	return uc.accountGroupRepo.DeleteByGamePlayerAndHouse(ctx, gameAccount.GamePlayerID, houseGID)
}

// ListGroupsByGameAccount 查询游戏账号加入的所有圈子
func (uc *GameAccountGroupUseCase) ListGroupsByGameAccount(
	ctx context.Context,
	gameAccountID int32,
) ([]*model.GameAccountGroup, error) {
	// 获取 game_player_id
	gameAccount, err := uc.accountRepo.GetByID(ctx, gameAccountID)
	if err != nil || gameAccount == nil || gameAccount.GamePlayerID == "" {
		return nil, fmt.Errorf("游戏账号 %d 不存在或缺少 game_player_id", gameAccountID)
	}
	return uc.accountGroupRepo.ListByGamePlayer(ctx, gameAccount.GamePlayerID)
}

// ListGroupsByUser 根据用户ID查询用户绑定的游戏账号的所有圈子（反向查询）
func (uc *GameAccountGroupUseCase) ListGroupsByUser(
	ctx context.Context,
	userID int32,
) ([]*model.GameAccountGroup, error) {
	// 1. 查询用户绑定的游戏账号
	gameAccount, err := uc.accountRepo.GetOneByUser(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*model.GameAccountGroup{}, nil
		}
		return nil, err
	}

	// 2. 查询该游戏账号加入的所有圈子（使用 game_player_id）
	if gameAccount.GamePlayerID == "" {
		return nil, fmt.Errorf("游戏账号缺少 game_player_id")
	}
	return uc.accountGroupRepo.ListByGamePlayer(ctx, gameAccount.GamePlayerID)
}

// ListGameAccountsByGroup 查询圈子中的所有游戏账号
func (uc *GameAccountGroupUseCase) ListGameAccountsByGroup(
	ctx context.Context,
	groupID int32,
) ([]*model.GameAccountGroup, error) {
	return uc.accountGroupRepo.ListByGroup(ctx, groupID)
}
