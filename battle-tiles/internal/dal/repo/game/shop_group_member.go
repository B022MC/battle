package game

import (
	"context"
	"fmt"

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
	var results []map[string]interface{}

	// 从 game_member 统计总数
	if err := r.db(ctx).
		Table("game_member gm").
		Where("gm.group_id = ?", groupID).
		Unscoped().
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 从 game_member 查询成员列表（LEFT JOIN 用户信息）
	offset := (page - 1) * size
	err := r.db(ctx).
		Table("game_member gm").
		Select(`gm.game_id, gm.game_name, gm.balance, gm.credit,
			ga.id as account_id, ga.user_id,
			u.username, u.nick_name, u.avatar`).
		Joins("LEFT JOIN game_account ga ON CAST(gm.game_id AS VARCHAR) = ga.game_player_id").
		Joins("LEFT JOIN basic_user u ON u.id = ga.user_id AND u.is_del = 0").
		Where("gm.group_id = ?", groupID).
		Order("gm.created_at DESC").
		Limit(int(size)).
		Offset(int(offset)).
		Unscoped().
		Find(&results).Error

	if err != nil {
		return nil, 0, err
	}

	// 转换为 BasicUser 结构（未绑定用户的使用游戏昵称）
	users := make([]*basicModel.BasicUser, 0, len(results))
	for _, row := range results {
		user := &basicModel.BasicUser{}

		// 如果有绑定的用户ID，使用用户信息
		if userID, ok := row["user_id"].(int32); ok && userID > 0 {
			user.Id = userID
			if username, ok := row["username"].(string); ok {
				user.Username = username
			}
			if nickName, ok := row["nick_name"].(string); ok {
				user.NickName = nickName
			}
			if avatar, ok := row["avatar"].(string); ok {
				user.Avatar = avatar
			}
		} else {
			// 未绑定用户：使用游戏信息占位
			user.Id = 0 // 特殊标记：未绑定
			if gameName, ok := row["game_name"].(string); ok {
				user.Username = gameName + "(未绑定)"
				user.NickName = gameName
			}
			if gameID, ok := row["game_id"].(int32); ok {
				user.GameNickname = row["game_name"].(string)
				// 将 game_id 保存到 introduction 字段供前端识别
				user.Introduction = fmt.Sprintf("game_id:%d", gameID)
			}
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *shopGroupMemberRepo) ListGroupsByUser(ctx context.Context, userID int32) ([]*model.GameShopGroup, error) {
	var groups []*model.GameShopGroup
	// 从 game_member 查询（通过 game_account 关联）
	err := r.db(ctx).
		Table("game_account ga").
		Select("DISTINCT g.*").
		Joins("JOIN game_member gm ON CAST(gm.game_id AS VARCHAR) = ga.game_player_id").
		Joins("JOIN game_shop_group g ON g.id = gm.group_id").
		Where("ga.user_id = ? AND gm.group_id IS NOT NULL AND g.is_active = ?", userID, true).
		Order("gm.created_at DESC").
		Unscoped(). // 禁用 GORM 自动软删除，避免给 game_member 添加 is_del 条件
		Find(&groups).Error
	return groups, err
}

func (r *shopGroupMemberRepo) ListGroupsByUserAndHouse(ctx context.Context, userID int32, houseGID int32) ([]*model.GameShopGroup, error) {
	var groups []*model.GameShopGroup
	// 从 game_member 查询（通过 game_account 关联）
	err := r.db(ctx).
		Table("game_account ga").
		Select("DISTINCT g.*").
		Joins("JOIN game_member gm ON CAST(gm.game_id AS VARCHAR) = ga.game_player_id").
		Joins("JOIN game_shop_group g ON g.id = gm.group_id").
		Where("ga.user_id = ? AND gm.house_gid = ? AND gm.group_id IS NOT NULL AND g.is_active = ?", userID, houseGID, true).
		Order("gm.created_at DESC").
		Unscoped(). // 禁用 GORM 自动软删除，避免给 game_member 添加 is_del 条件
		Find(&groups).Error
	return groups, err
}

func (r *shopGroupMemberRepo) CountMembers(ctx context.Context, groupID int32) (int64, error) {
	var cnt int64
	// 从 game_member 统计成员数量（实际在用的表）
	err := r.db(ctx).
		Table("game_member").
		Where("group_id = ?", groupID).
		Unscoped(). // 禁用 GORM 自动软删除，game_member 表没有软删除字段
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
