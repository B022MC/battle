package game

import (
	"context"

	basicModel "battle-tiles/internal/dal/model/basic"
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type ShopGroupMemberRepo interface {
	// AddMember 添加成员到圈子
	AddMember(ctx context.Context, m *model.GameShopGroupMember) error
	// RemoveMember 从圈子移除成员
	RemoveMember(ctx context.Context, groupID int32, userID int32) error
	// IsMember 检查用户是否是圈子成员
	IsMember(ctx context.Context, groupID int32, userID int32) (bool, error)
	// ListMembersByGroup 获取圈子的所有成员（返回用户信息）
	ListMembersByGroup(ctx context.Context, groupID int32, page, size int32) ([]*basicModel.BasicUser, int64, error)
	// ListGroupsByUser 获取用户加入的所有圈子
	ListGroupsByUser(ctx context.Context, userID int32) ([]*model.GameShopGroup, error)
	// ListGroupsByUserAndHouse 获取用户在指定店铺下加入的所有圈子
	ListGroupsByUserAndHouse(ctx context.Context, userID int32, houseGID int32) ([]*model.GameShopGroup, error)
	// CountMembers 统计圈子成员数量
	CountMembers(ctx context.Context, groupID int32) (int64, error)
	// BatchAddMembers 批量添加成员
	BatchAddMembers(ctx context.Context, members []*model.GameShopGroupMember) error
}

type shopGroupMemberRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewShopGroupMemberRepo(data *infra.Data, logger log.Logger) ShopGroupMemberRepo {
	return &shopGroupMemberRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/shop_group_member")),
	}
}

func (r *shopGroupMemberRepo) db(ctx context.Context) *gorm.DB {
	return r.data.GetDBWithContext(ctx)
}

func (r *shopGroupMemberRepo) AddMember(ctx context.Context, m *model.GameShopGroupMember) error {
	// 使用 FirstOrCreate 实现幂等性
	return r.db(ctx).
		Where("group_id = ? AND user_id = ?", m.GroupID, m.UserID).
		FirstOrCreate(m).Error
}

func (r *shopGroupMemberRepo) RemoveMember(ctx context.Context, groupID int32, userID int32) error {
	return r.db(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Delete(&model.GameShopGroupMember{}).Error
}

func (r *shopGroupMemberRepo) IsMember(ctx context.Context, groupID int32, userID int32) (bool, error) {
	var cnt int64
	// 从 game_account_group 查询（通过 game_account 关联）
	err := r.db(ctx).
		Table("game_account ga").
		Joins("JOIN game_account_group gag ON ga.id = gag.game_account_id").
		Where("ga.user_id = ? AND gag.group_id = ? AND gag.status = ?", userID, groupID, "active").
		Unscoped(). // 禁用自动软删除过滤
		Count(&cnt).Error
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (r *shopGroupMemberRepo) ListMembersByGroup(ctx context.Context, groupID int32, page, size int32) ([]*basicModel.BasicUser, int64, error) {
	var total int64
	var users []*basicModel.BasicUser

	// 从 game_account_group 统计总数
	if err := r.db(ctx).
		Table("game_account_group gag").
		Joins("JOIN game_account ga ON ga.id = gag.game_account_id").
		Where("gag.group_id = ? AND gag.status = ? AND ga.user_id IS NOT NULL", groupID, "active").
		Unscoped(). // 禁用自动软删除过滤
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 从 game_account_group 查询成员列表
	offset := (page - 1) * size
	err := r.db(ctx).
		Table("game_account_group gag").
		Select("u.*").
		Joins("JOIN game_account ga ON ga.id = gag.game_account_id").
		Joins("JOIN basic_user u ON u.id = ga.user_id").
		Where("gag.group_id = ? AND gag.status = ? AND u.is_del = 0", groupID, "active").
		Order("gag.joined_at DESC").
		Limit(int(size)).
		Offset(int(offset)).
		Unscoped(). // 禁用自动软删除过滤
		Find(&users).Error

	return users, total, err
}

func (r *shopGroupMemberRepo) ListGroupsByUser(ctx context.Context, userID int32) ([]*model.GameShopGroup, error) {
	var groups []*model.GameShopGroup
	// 从 game_account_group 查询（通过 game_account 关联）
	err := r.db(ctx).
		Table("game_account ga").
		Select("g.*").
		Joins("JOIN game_account_group gag ON ga.id = gag.game_account_id").
		Joins("JOIN game_shop_group g ON g.id = gag.group_id").
		Where("ga.user_id = ? AND gag.status = ? AND g.is_active = ?", userID, "active", true).
		Order("gag.joined_at DESC").
		Unscoped(). // 禁用自动软删除过滤
		Find(&groups).Error
	return groups, err
}

func (r *shopGroupMemberRepo) ListGroupsByUserAndHouse(ctx context.Context, userID int32, houseGID int32) ([]*model.GameShopGroup, error) {
	var groups []*model.GameShopGroup
	// 从 game_account_group 查询（通过 game_account 关联）
	err := r.db(ctx).
		Table("game_account ga").
		Select("g.*").
		Joins("JOIN game_account_group gag ON ga.id = gag.game_account_id").
		Joins("JOIN game_shop_group g ON g.id = gag.group_id").
		Where("ga.user_id = ? AND gag.house_gid = ? AND gag.status = ? AND g.is_active = ?", userID, houseGID, "active", true).
		Order("gag.joined_at DESC").
		Unscoped(). // 禁用自动软删除过滤
		Find(&groups).Error
	return groups, err
}

func (r *shopGroupMemberRepo) CountMembers(ctx context.Context, groupID int32) (int64, error) {
	var cnt int64
	// 从 game_account_group 统计成员数量
	err := r.db(ctx).
		Table("game_account_group").
		Where("group_id = ? AND status = ?", groupID, "active").
		Unscoped(). // 禁用自动软删除过滤
		Count(&cnt).Error
	return cnt, err
}

func (r *shopGroupMemberRepo) BatchAddMembers(ctx context.Context, members []*model.GameShopGroupMember) error {
	if len(members) == 0 {
		return nil
	}

	// 使用事务批量插入，忽略重复
	return r.db(ctx).Transaction(func(tx *gorm.DB) error {
		for _, m := range members {
			if err := tx.Where("group_id = ? AND user_id = ?", m.GroupID, m.UserID).
				FirstOrCreate(m).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
