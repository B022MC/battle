// internal/biz/game/game_shop_admin.go
package game

import (
	basicModel "battle-tiles/internal/dal/model/basic"
	model "battle-tiles/internal/dal/model/game"
	basicRepo "battle-tiles/internal/dal/repo/basic"
	repo "battle-tiles/internal/dal/repo/game"
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ShopAdminUseCase struct {
	repo          repo.GameShopAdminRepo
	shopGroupRepo repo.ShopGroupRepo
	basicUserRepo basicRepo.BasicUserRepo
	authRepo      basicRepo.AuthRepo
	log           *log.Helper
}

func NewShopAdminUseCase(
	r repo.GameShopAdminRepo,
	shopGroupRepo repo.ShopGroupRepo,
	basicUserRepo basicRepo.BasicUserRepo,
	authRepo basicRepo.AuthRepo,
	logger log.Logger,
) *ShopAdminUseCase {
	return &ShopAdminUseCase{
		repo:          r,
		shopGroupRepo: shopGroupRepo,
		basicUserRepo: basicUserRepo,
		authRepo:      authRepo,
		log:           log.NewHelper(log.With(logger, "module", "usecase/shop_admin")),
	}
}

func validateRole(role string) (string, error) {
	role = strings.ToLower(strings.TrimSpace(role))
	if role == "" {
		return "admin", nil
	}
	switch role {
	case "admin", "operator":
		return role, nil
	default:
		return "", errors.New("invalid role (must be admin|operator)")
	}
}

func (uc *ShopAdminUseCase) Assign(ctx context.Context, houseGID int32, targetUserID int32, role string) error {
	if houseGID <= 0 || targetUserID <= 0 {
		return errors.New("invalid house_gid or user_id")
	}
	r, err := validateRole(role)
	if err != nil {
		return err
	}

	// 1. 检查用户是否存在
	user, err := uc.basicUserRepo.SelectOneByPK(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return errors.Wrap(err, "查询用户失败")
	}

	// 2. 创建店铺管理员记录
	m := &model.GameShopAdmin{
		HouseGID: int32(houseGID),
		UserID:   targetUserID,
		Role:     r,
	}
	if err := uc.repo.Assign(ctx, m); err != nil {
		return err
	}

	// 3. 更新用户角色为店铺管理员
	// 3.1 更新 basic_user 表的 role 字段（冗余字段）
	user.Role = basicModel.UserRoleStoreAdmin
	if _, err := uc.basicUserRepo.UpdateByPK(ctx, user); err != nil {
		uc.log.Errorf("更新用户角色失败: %v", err)
		// 不回滚，因为主要操作已完成
	}

	// 3.2 更新 RBAC 角色关联（sys_user_role）
	// 注意：RBAC 中的角色 Code 可能是 "shop_admin" 或 "store_admin"，为了兼容性，优先尝试 "shop_admin"
	// 参考 ShopApplicationService 中的用法
	if err := uc.authRepo.EnsureUserHasOnlyRoleByCode(ctx, targetUserID, "shop_admin"); err != nil {
		// 如果 shop_admin 不存在，尝试 store_admin
		uc.log.Warnf("RBAC角色 shop_admin 设置失败，尝试 store_admin: %v", err)
		if err := uc.authRepo.EnsureUserHasOnlyRoleByCode(ctx, targetUserID, basicModel.UserRoleStoreAdmin); err != nil {
			uc.log.Errorf("RBAC角色设置失败: %v", err)
		}
	}

	// 4. 自动创建该管理员的圈子
	groupName := user.NickName + "的圈子"
	if user.NickName == "" {
		groupName = user.Username + "的圈子"
	}

	group := &model.GameShopGroup{
		HouseGID:    houseGID,
		GroupName:   groupName,
		AdminUserID: targetUserID,
		IsActive:    true,
	}

	if err := uc.shopGroupRepo.Create(ctx, group); err != nil {
		uc.log.Errorf("创建圈子失败: %v", err)
		// 不回滚，管理员已创建成功
	}

	uc.log.Infof("成功设置用户 %d 为店铺 %d 的管理员", targetUserID, houseGID)
	return nil
}

func (uc *ShopAdminUseCase) Revoke(ctx context.Context, houseGID int32, targetUserID int32) error {
	if targetUserID <= 0 {
		return errors.New("invalid user_id")
	}

	// 如果没有提供 houseGID，通过 user_id 查找
	actualHouseGID := houseGID
	if actualHouseGID <= 0 {
		// 查询用户的管理员记录
		admins, err := uc.repo.ListByUser(ctx, targetUserID)
		if err != nil {
			return errors.Wrap(err, "查询用户管理员记录失败")
		}
		if len(admins) == 0 {
			return errors.New("用户不是任何店铺的管理员")
		}
		// 取第一个（假设一个用户只能是一个店铺的管理员）
		actualHouseGID = admins[0].HouseGID
	}

	// 1. 删除该管理员的圈子（会级联删除圈子成员）
	groups, err := uc.shopGroupRepo.ListByAdmin(ctx, targetUserID)
	if err != nil {
		uc.log.Errorf("查询管理员圈子失败: %v", err)
	} else {
		for _, group := range groups {
			if group.HouseGID == actualHouseGID {
				if err := uc.shopGroupRepo.Delete(ctx, group.Id); err != nil {
					uc.log.Errorf("删除圈子 %d 失败: %v", group.Id, err)
				}
			}
		}
	}

	// 2. 删除店铺管理员记录
	if err := uc.repo.Revoke(ctx, actualHouseGID, targetUserID); err != nil {
		return err
	}

	// 3. 更新用户角色为普通用户
	user, err := uc.basicUserRepo.SelectOneByPK(ctx, targetUserID)
	if err == nil {
		// 3.1 更新 basic_user 表的 role 字段
		user.Role = basicModel.UserRoleRegularUser
		if _, err := uc.basicUserRepo.UpdateByPK(ctx, user); err != nil {
			uc.log.Errorf("更新用户角色失败: %v", err)
		}

		// 3.2 更新 RBAC 角色关联（恢复为 user）
		if err := uc.authRepo.EnsureUserHasOnlyRoleByCode(ctx, targetUserID, "user"); err != nil {
			uc.log.Errorf("RBAC角色恢复失败: %v", err)
		}
	}

	uc.log.Infof("成功移除用户 %d 的店铺 %d 管理员身份", targetUserID, actualHouseGID)
	return nil
}

func (uc *ShopAdminUseCase) IsAdmin(ctx context.Context, houseGID int32, userID int32) (bool, error) {
	return uc.repo.Exists(ctx, houseGID, userID)
}

func (uc *ShopAdminUseCase) List(ctx context.Context, houseGID int32) ([]*model.GameShopAdmin, error) {
	if houseGID <= 0 {
		return nil, errors.New("invalid house_gid")
	}
	return uc.repo.ListByHouse(ctx, houseGID)
}

func (uc *ShopAdminUseCase) ListByUser(ctx context.Context, userID int32) ([]*model.GameShopAdmin, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user_id")
	}
	return uc.repo.ListByUser(ctx, userID)
}
