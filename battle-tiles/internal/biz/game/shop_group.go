package game

import (
	"context"
	"fmt"

	basicModel "battle-tiles/internal/dal/model/basic"
	model "battle-tiles/internal/dal/model/game"
	basicRepo "battle-tiles/internal/dal/repo/basic"
	"battle-tiles/internal/dal/repo/game"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type ShopGroupUseCase struct {
	groupRepo        game.ShopGroupRepo
	memberRepo       game.ShopGroupMemberRepo // 用于查询圈子成员列表
	shopAdminRepo    game.GameShopAdminRepo
	gameMemberRepo   game.GameMemberRepo       // 用于操作 game_member 记录
	accountRepo      game.GameAccountRepo      // 用于查询 game_account
	accountGroupRepo game.GameAccountGroupRepo // 用于操作 game_account_group 记录（主表）
	userRepo         basicRepo.BasicUserRepo   // 用于查询平台用户信息
	log              *log.Helper
}

func NewShopGroupUseCase(
	groupRepo game.ShopGroupRepo,
	memberRepo game.ShopGroupMemberRepo,
	shopAdminRepo game.GameShopAdminRepo,
	gameMemberRepo game.GameMemberRepo,
	accountRepo game.GameAccountRepo,
	accountGroupRepo game.GameAccountGroupRepo,
	userRepo basicRepo.BasicUserRepo,
	logger log.Logger,
) *ShopGroupUseCase {
	return &ShopGroupUseCase{
		groupRepo:        groupRepo,
		memberRepo:       memberRepo,
		shopAdminRepo:    shopAdminRepo,
		gameMemberRepo:   gameMemberRepo,
		accountRepo:      accountRepo,
		accountGroupRepo: accountGroupRepo,
		userRepo:         userRepo,
		log:              log.NewHelper(log.With(logger, "module", "usecase/shop_group")),
	}
}

// CreateGroup 创建圈子（店铺管理员创建）
func (uc *ShopGroupUseCase) CreateGroup(ctx context.Context, houseGID int32, adminUserID int32, groupName, description string) (*model.GameShopGroup, error) {
	// 检查用户是否是该店铺的管理员
	isAdmin, err := uc.shopAdminRepo.Exists(ctx, houseGID, adminUserID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, fmt.Errorf("用户不是该店铺的管理员")
	}

	// 检查是否已经有圈子
	existing, err := uc.groupRepo.GetByAdmin(ctx, houseGID, adminUserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existing != nil {
		return existing, nil // 已存在，返回现有圈子
	}

	// 创建圈子
	group := &model.GameShopGroup{
		HouseGID:    houseGID,
		GroupName:   groupName,
		AdminUserID: adminUserID,
		Description: description,
		IsActive:    true,
	}

	if err := uc.groupRepo.Create(ctx, group); err != nil {
		return nil, err
	}

	return group, nil
}

// GetMyGroup 获取我的圈子（店铺管理员）
func (uc *ShopGroupUseCase) GetMyGroup(ctx context.Context, houseGID int32, adminUserID int32) (*model.GameShopGroup, error) {
	return uc.groupRepo.GetByAdmin(ctx, houseGID, adminUserID)
}

// ListGroupsByHouse 获取店铺下的所有圈子
func (uc *ShopGroupUseCase) ListGroupsByHouse(ctx context.Context, houseGID int32) ([]*model.GameShopGroup, error) {
	return uc.groupRepo.ListByHouse(ctx, houseGID)
}

// UpdateGroup 更新圈子信息
func (uc *ShopGroupUseCase) UpdateGroup(ctx context.Context, groupID int32, adminUserID int32, groupName, description string) error {
	// 检查圈子是否属于该管理员
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.AdminUserID != adminUserID {
		return fmt.Errorf("无权修改该圈子")
	}

	group.GroupName = groupName
	group.Description = description

	return uc.groupRepo.Update(ctx, group)
}

// 已删除 AddMembersToGroup 方法
// 现在统一使用 /shops/members/pull-to-group 接口进行拉圈操作

// AddGameAccountsToGroup 通过游戏账号ID添加成员到圈子（用于添加游戏内已存在的玩家）
// 这个方法不需要平台用户ID，直接通过游戏账号ID添加
func (uc *ShopGroupUseCase) AddGameAccountsToGroup(ctx context.Context, groupID int32, adminUserID int32, gameAccountIDs []int32) error {
	// 检查圈子是否属于该管理员
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.AdminUserID != adminUserID {
		return fmt.Errorf("无权操作该圈子")
	}

	// 批量添加游戏账号到圈子
	for _, gameAccountID := range gameAccountIDs {
		// 1. 查询游戏账号
		gameAccount, err := uc.accountRepo.GetByID(ctx, gameAccountID)
		if err != nil {
			uc.log.Warnf("游戏账号 %d 不存在，跳过", gameAccountID)
			continue
		}

		// 2. 检查该游戏账号是否已经在该店铺的某个圈子中
		if gameAccount.UserID != nil {
			// 检查 game_player_id
			if gameAccount.GamePlayerID == "" {
				uc.log.Warnf("游戏账号 %d 缺少 game_player_id，跳过", gameAccount.Id)
				continue
			}

			// 检查是否已经在圈子中
			exists, err := uc.accountGroupRepo.ExistsByGamePlayerAndHouse(ctx, gameAccount.GamePlayerID, group.HouseGID)
			if err != nil {
				uc.log.Warnf("检查游戏账号 %d 是否在圈子中失败: %v", gameAccount.Id, err)
				continue
			}
			if exists {
				uc.log.Infof("游戏账号 %d 已在圈子 %d 中，跳过", gameAccount.Id, groupID)
				continue
			}

			// 创建 game_account_group 记录（使用 game_player_id）
			adminUserIDPtr := &adminUserID
			accountGroup := &model.GameAccountGroup{
				GamePlayerID: gameAccount.GamePlayerID,
				HouseGID:     group.HouseGID,
				GroupID:      groupID,
				GroupName:    group.GroupName,
				AdminUserID:  adminUserIDPtr,
				Status:       model.AccountGroupStatusActive,
			}

			// 检查是否已存在
			exists, err = uc.accountGroupRepo.ExistsByGamePlayerAndHouse(ctx, gameAccount.GamePlayerID, group.HouseGID)
			if err != nil {
				uc.log.Warnf("检查 game_account_group 是否存在失败: %v", err)
				continue
			} else if exists {
				// 更新状态为激活
				existing, err := uc.accountGroupRepo.GetByGamePlayerAndHouse(ctx, gameAccount.GamePlayerID, group.HouseGID)
				if err == nil {
					uc.accountGroupRepo.UpdateStatus(ctx, existing.Id, model.AccountGroupStatusActive)
				}
			} else {
				// 创建新记录
				if err := uc.accountGroupRepo.Create(ctx, accountGroup); err != nil {
					uc.log.Errorf("创建 game_account_group 失败: %v", err)
					continue
				}
			}

		} else {
			// user_id 为空的情况，也需要使用 game_player_id
			if gameAccount.GamePlayerID == "" {
				uc.log.Warnf("游戏账号 %d 缺少 game_player_id，跳过", gameAccount.Id)
				continue
			}

			// 检查是否已存在
			exists, err := uc.accountGroupRepo.ExistsByGamePlayerAndHouse(ctx, gameAccount.GamePlayerID, group.HouseGID)
			if err != nil {
				uc.log.Warnf("检查 game_account_group 是否存在失败: %v", err)
				continue
			} else if exists {
				// 更新状态为激活
				existing, err := uc.accountGroupRepo.GetByGamePlayerAndHouse(ctx, gameAccount.GamePlayerID, group.HouseGID)
				if err == nil {
					uc.accountGroupRepo.UpdateStatus(ctx, existing.Id, model.AccountGroupStatusActive)
				}
			} else {
				// 创建新记录
				adminUserIDPtr := &adminUserID
				accountGroup := &model.GameAccountGroup{
					GamePlayerID: gameAccount.GamePlayerID,
					HouseGID:     group.HouseGID,
					GroupID:      groupID,
					GroupName:    group.GroupName,
					AdminUserID:  adminUserIDPtr,
					Status:       model.AccountGroupStatusActive,
				}
				if err := uc.accountGroupRepo.Create(ctx, accountGroup); err != nil {
					uc.log.Errorf("创建 game_account_group 失败: %v", err)
					continue
				}
			}
		}

		// 4. 创建或更新 game_member 记录
		gameMember := &model.GameMember{
			HouseGID:  group.HouseGID,
			GameID:    gameAccount.Id,
			GameName:  gameAccount.Nickname,
			GroupID:   &groupID,
			GroupName: group.GroupName,
			Balance:   0,
			Credit:    0,
			Forbid:    false,
		}
		if err := uc.gameMemberRepo.Create(ctx, gameMember); err != nil {
			uc.log.Warnf("创建 game_member 失败: %v", err)
		}

		uc.log.Infof("成功添加游戏账号 %d 到圈子 %d", gameAccountID, groupID)
	}

	return nil
}

// RemoveMemberFromGroup 从圈子移除成员
// 当用户被踢出圈子时，需要清理相关联的记录：
// 1. game_account_group - 游戏账号圈子关系（主表）
// 2. game_member - 游戏成员业务数据
func (uc *ShopGroupUseCase) RemoveMemberFromGroup(ctx context.Context, groupID int32, adminUserID int32, userID int32) error {
	// 检查圈子是否属于该管理员
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.AdminUserID != adminUserID {
		return fmt.Errorf("无权操作该圈子")
	}

	// 1. 查询用户的游戏账号
	gameAccount, err := uc.accountRepo.GetOneByUser(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.log.Warnf("用户 %d 没有绑定游戏账号，无法移除", userID)
			return fmt.Errorf("用户没有绑定游戏账号")
		}
		return fmt.Errorf("查询游戏账号失败: %w", err)
	}

	// 删除 game_account_group 记录（使用 game_player_id）
	if gameAccount.GamePlayerID == "" {
		uc.log.Warnf("游戏账号 %d 缺少 game_player_id，无法移除", gameAccount.Id)
		return fmt.Errorf("游戏账号缺少 game_player_id")
	}
	if err := uc.accountGroupRepo.DeleteByGamePlayerAndHouse(ctx, gameAccount.GamePlayerID, group.HouseGID); err != nil {
		uc.log.Errorf("删除 game_account_group 失败: %v", err)
		return fmt.Errorf("移除圈子成员失败: %w", err)
	}

	// 3. 删除 game_member 记录（游戏成员业务数据）
	if err := uc.gameMemberRepo.DeleteByGameID(ctx, group.HouseGID, gameAccount.Id); err != nil {
		uc.log.Warnf("删除 game_member 失败: %v", err)
		// 不阻塞流程
	}

	uc.log.Infof("成功移除用户 %d 并清理关联记录", userID)
	return nil
}

// ListGroupMembers 获取圈子成员列表
func (uc *ShopGroupUseCase) ListGroupMembers(ctx context.Context, groupID int32, page, size int32) ([]*basicModel.BasicUser, int64, error) {
	return uc.memberRepo.ListMembersByGroup(ctx, groupID, page, size)
}

// ListMyGroups 获取我加入的所有圈子（通过游戏账号反向查询）
// 注意：此方法已废弃，请使用 GameAccountGroupUseCase.ListGroupsByUser
func (uc *ShopGroupUseCase) ListMyGroups(ctx context.Context, userID int32) ([]*model.GameShopGroup, error) {
	// 保留旧逻辑以兼容
	return uc.memberRepo.ListGroupsByUser(ctx, userID)
}

// GetGroupWithMemberCount 获取圈子信息（包含成员数量）
func (uc *ShopGroupUseCase) GetGroupWithMemberCount(ctx context.Context, groupID int32) (*model.GameShopGroup, int64, error) {
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return nil, 0, err
	}

	count, err := uc.memberRepo.CountMembers(ctx, groupID)
	if err != nil {
		return nil, 0, err
	}

	return group, count, nil
}

// EnsureGroupForAdmin 确保店铺管理员有圈子（自动创建）
func (uc *ShopGroupUseCase) EnsureGroupForAdmin(ctx context.Context, houseGID int32, adminUserID int32, defaultGroupName string) (*model.GameShopGroup, error) {
	// 先查询是否已有圈子
	group, err := uc.groupRepo.GetByAdmin(ctx, houseGID, adminUserID)
	if err == nil {
		return group, nil // 已存在
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 不存在，创建默认圈子
	return uc.CreateGroup(ctx, houseGID, adminUserID, defaultGroupName, "")
}
