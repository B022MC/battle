package game

import (
	"context"
	"fmt"
	"time"

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

// AddMembersToGroup 添加成员到圈子
// 当用户被添加到圈子时，需要创建以下记录：
// 1. game_account_group - 游戏账号圈子关系（主表，用于查询）
// 2. game_member - 游戏成员业务数据（余额、积分等）
//
// 注意：用户必须先绑定游戏账号才能加入圈子
func (uc *ShopGroupUseCase) AddMembersToGroup(ctx context.Context, groupID int32, adminUserID int32, userIDs []int32) error {
	// 检查圈子是否属于该管理员
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.AdminUserID != adminUserID {
		return fmt.Errorf("无权操作该圈子")
	}

	// 预先验证所有用户
	for _, userID := range userIDs {
		// 1. 验证用户必须在平台注册
		user, err := uc.userRepo.SelectOneByPK(ctx, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("用户 %d 不存在，必须是已注册的平台用户", userID)
			}
			return fmt.Errorf("查询用户 %d 失败: %w", userID, err)
		}

		// 2. 验证用户角色（不能添加超级管理员或店铺管理员）
		if user.IsSuperAdmin() {
			return fmt.Errorf("不能将超级管理员添加为圈子成员")
		}
		if user.IsStoreAdmin() {
			return fmt.Errorf("不能将店铺管理员添加为圈子成员")
		}

		// 3. 查询用户在该店铺下的所有圈子
		userGroups, err := uc.memberRepo.ListGroupsByUserAndHouse(ctx, userID, group.HouseGID)
		if err != nil {
			return fmt.Errorf("查询用户圈子失败: %w", err)
		}

		// 4. 如果用户已经在该店铺的某个圈子中，不能再加入其他圈子
		if len(userGroups) > 0 {
			// 检查是否是当前圈子
			for _, existingGroup := range userGroups {
				if existingGroup.Id != groupID {
					return fmt.Errorf("用户已在该店铺的其他圈子中（圈子：%s），不能重复加入", existingGroup.GroupName)
				}
			}
		}
	}

	// 批量添加成员（分别处理每个用户，确保创建所有相关记录）
	for _, userID := range userIDs {
		// 1. 查询或创建游戏账号
		gameAccount, err := uc.accountRepo.GetOneByUser(ctx, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// 用户没有游戏账号，自动创建虚拟游戏账号
				uc.log.Infof("用户 %d 没有绑定游戏账号，自动创建虚拟账号", userID)

				// 查询用户信息（已在前面验证过，这里只需获取用户名）
				user, err := uc.userRepo.SelectOneByPK(ctx, userID)
				if err != nil {
					// 理论上不会到这里，因为前面已经验证过
					uc.log.Errorf("查询用户 %d 信息失败: %v", userID, err)
					return fmt.Errorf("系统错误：无法获取用户信息")
				}

				// 创建虚拟游戏账号
				virtualAccount := &model.GameAccount{
					UserID:       &userID,
					Account:      user.Username,                                           // 使用平台用户名
					GamePlayerID: fmt.Sprintf("virtual_%d_%d", userID, time.Now().Unix()), // 虚拟游戏ID
					PwdMD5:       "",                                                      // 虚拟账号无密码
				}

				if err := uc.accountRepo.Create(ctx, virtualAccount); err != nil {
					uc.log.Errorf("创建虚拟游戏账号失败: %v", err)
					return fmt.Errorf("创建虚拟游戏账号失败: %w", err)
				}

				gameAccount = virtualAccount
				uc.log.Infof("已为用户 %d 创建虚拟游戏账号 (ID: %d)", userID, gameAccount.Id)
			} else {
				uc.log.Errorf("查询用户 %d 的游戏账号失败: %v", userID, err)
				return fmt.Errorf("查询游戏账号失败: %w", err)
			}
		}

		// 2. 创建 game_account_group 记录（游戏账号圈子关系，主表）
		accountGroup := &model.GameAccountGroup{
			GameAccountID:    gameAccount.Id,
			HouseGID:         group.HouseGID,
			GroupID:          groupID,
			GroupName:        group.GroupName,
			AdminUserID:      adminUserID,
			ApprovedByUserID: adminUserID,
			Status:           model.AccountGroupStatusActive,
		}

		// 检查是否已存在
		exists, err := uc.accountGroupRepo.ExistsByGameAccountAndHouse(ctx, gameAccount.Id, group.HouseGID)
		if err != nil {
			uc.log.Warnf("检查 game_account_group 是否存在失败: %v", err)
			continue
		} else if exists {
			// 更新状态为激活
			existing, err := uc.accountGroupRepo.GetByGameAccountAndHouse(ctx, gameAccount.Id, group.HouseGID)
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

		// 3. 创建或更新 game_member 记录（业务数据）
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

		uc.log.Infof("成功添加用户 %d 到圈子 %d 并创建所有关联记录", userID, groupID)
	}

	return nil
}

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
			userGroups, err := uc.memberRepo.ListGroupsByUserAndHouse(ctx, *gameAccount.UserID, group.HouseGID)
			if err == nil && len(userGroups) > 0 {
				for _, existingGroup := range userGroups {
					if existingGroup.Id != groupID {
						uc.log.Warnf("游戏账号 %d 已在其他圈子中（圈子：%s），跳过", gameAccountID, existingGroup.GroupName)
						continue
					}
				}
			}
		}

		// 3. 创建 game_account_group 记录
		accountGroup := &model.GameAccountGroup{
			GameAccountID:    gameAccount.Id,
			HouseGID:         group.HouseGID,
			GroupID:          groupID,
			GroupName:        group.GroupName,
			AdminUserID:      adminUserID,
			ApprovedByUserID: adminUserID,
			Status:           model.AccountGroupStatusActive,
		}

		// 检查是否已存在
		exists, err := uc.accountGroupRepo.ExistsByGameAccountAndHouse(ctx, gameAccount.Id, group.HouseGID)
		if err != nil {
			uc.log.Warnf("检查 game_account_group 是否存在失败: %v", err)
			continue
		} else if exists {
			// 更新状态为激活
			existing, err := uc.accountGroupRepo.GetByGameAccountAndHouse(ctx, gameAccount.Id, group.HouseGID)
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

	// 2. 删除 game_account_group 记录（游戏账号圈子关系，主表）
	if err := uc.accountGroupRepo.DeleteByGameAccountAndHouse(ctx, gameAccount.Id, group.HouseGID); err != nil {
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
