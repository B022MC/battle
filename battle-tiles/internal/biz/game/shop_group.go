package game

import (
	"context"
	"fmt"

	basicModel "battle-tiles/internal/dal/model/basic"
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/dal/repo/game"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type ShopGroupUseCase struct {
	groupRepo     game.ShopGroupRepo
	memberRepo    game.ShopGroupMemberRepo
	shopAdminRepo game.GameShopAdminRepo
	log           *log.Helper
}

func NewShopGroupUseCase(
	groupRepo game.ShopGroupRepo,
	memberRepo game.ShopGroupMemberRepo,
	shopAdminRepo game.GameShopAdminRepo,
	logger log.Logger,
) *ShopGroupUseCase {
	return &ShopGroupUseCase{
		groupRepo:     groupRepo,
		memberRepo:    memberRepo,
		shopAdminRepo: shopAdminRepo,
		log:           log.NewHelper(log.With(logger, "module", "usecase/shop_group")),
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
func (uc *ShopGroupUseCase) AddMembersToGroup(ctx context.Context, groupID int32, adminUserID int32, userIDs []int32) error {
	// 检查圈子是否属于该管理员
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.AdminUserID != adminUserID {
		return fmt.Errorf("无权操作该圈子")
	}

	// 检查每个用户是否已经在该店铺的其他圈子中
	for _, userID := range userIDs {
		// 查询用户在该店铺下的所有圈子
		userGroups, err := uc.memberRepo.ListGroupsByUserAndHouse(ctx, userID, group.HouseGID)
		if err != nil {
			return fmt.Errorf("查询用户圈子失败: %w", err)
		}

		// 如果用户已经在该店铺的某个圈子中，不能再加入其他圈子
		if len(userGroups) > 0 {
			// 检查是否是当前圈子
			for _, existingGroup := range userGroups {
				if existingGroup.Id != groupID {
					return fmt.Errorf("用户已在该店铺的其他圈子中（圈子：%s），不能重复加入", existingGroup.GroupName)
				}
			}
		}
	}

	// 批量添加成员
	members := make([]*model.GameShopGroupMember, 0, len(userIDs))
	for _, userID := range userIDs {
		members = append(members, &model.GameShopGroupMember{
			GroupID: groupID,
			UserID:  userID,
		})
	}

	return uc.memberRepo.BatchAddMembers(ctx, members)
}

// RemoveMemberFromGroup 从圈子移除成员
func (uc *ShopGroupUseCase) RemoveMemberFromGroup(ctx context.Context, groupID int32, adminUserID int32, userID int32) error {
	// 检查圈子是否属于该管理员
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.AdminUserID != adminUserID {
		return fmt.Errorf("无权操作该圈子")
	}

	return uc.memberRepo.RemoveMember(ctx, groupID, userID)
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
