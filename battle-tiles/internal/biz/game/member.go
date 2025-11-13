package game

import (
	"context"

	basicModel "battle-tiles/internal/dal/model/basic"
	"battle-tiles/internal/dal/repo/basic"
	"battle-tiles/internal/dal/repo/game"
	utils2 "battle-tiles/pkg/utils"

	"github.com/go-kratos/kratos/v2/log"
)

type MemberUseCase struct {
	userRepo      basic.BasicUserRepo
	shopAdminRepo game.GameShopAdminRepo
	log           *log.Helper
}

func NewMemberUseCase(
	userRepo basic.BasicUserRepo,
	shopAdminRepo game.GameShopAdminRepo,
	logger log.Logger,
) *MemberUseCase {
	return &MemberUseCase{
		userRepo:      userRepo,
		shopAdminRepo: shopAdminRepo,
		log:           log.NewHelper(log.With(logger, "module", "usecase/member")),
	}
}

// ListAllUsers 查看所有用户（超级管理员和店铺管理员都可以查看）
func (uc *MemberUseCase) ListAllUsers(ctx context.Context, page, size int32, keyword string) ([]*basicModel.BasicUser, int64, error) {
	// 使用 CORM 的分页查询
	db := uc.userRepo.DB(ctx).Model(&basicModel.BasicUser{})

	// 如果有关键词，添加搜索条件
	if keyword != "" {
		db = db.Where("username LIKE ? OR nick_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 使用 PageParam
	pageParam := &utils2.PageParam{
		PageNo:   page,
		PageSize: size,
	}

	users, total, err := uc.userRepo.ListPage(ctx, db, pageParam)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByID 根据ID获取用户信息
func (uc *MemberUseCase) GetUserByID(ctx context.Context, userID int32) (*basicModel.BasicUser, error) {
	return uc.userRepo.SelectOneByPK(ctx, userID)
}

// CheckIsShopAdmin 检查用户是否是店铺管理员
func (uc *MemberUseCase) CheckIsShopAdmin(ctx context.Context, houseGID int32, userID int32) (bool, error) {
	return uc.shopAdminRepo.Exists(ctx, houseGID, userID)
}

// ListShopAdmins 获取店铺的所有管理员
func (uc *MemberUseCase) ListShopAdmins(ctx context.Context, houseGID int32) ([]*basicModel.BasicUser, error) {
	// 先获取管理员记录
	admins, err := uc.shopAdminRepo.ListByHouse(ctx, houseGID)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	users := make([]*basicModel.BasicUser, 0, len(admins))
	for _, admin := range admins {
		user, err := uc.userRepo.SelectOneByPK(ctx, admin.UserID)
		if err != nil {
			uc.log.Warnf("failed to get user %d: %v", admin.UserID, err)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}
