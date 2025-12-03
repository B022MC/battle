package game

import (
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type GameMemberRepo interface {
	// 根据 game_id 和 group_id 查询成员
	GetByGameIDAndGroup(ctx context.Context, houseGID int32, gameID int32, groupID *int32) (*model.GameMember, error)

	// 根据 game_id 查询成员，取第一个结果
	GetByGameID(ctx context.Context, houseGID int32, gameID int32) (*model.GameMember, error)

	// 静默模式查询（不输出SQL日志），用于定时任务
	GetByGameIDQuiet(ctx context.Context, houseGID int32, gameID int32) (*model.GameMember, error)

	// 鏌ヨ鎴愬憳鍦ㄦ墍鏈夊湀瀛愮殑璁板綍
	ListByGameID(ctx context.Context, houseGID int32, gameID int32) ([]*model.GameMember, error)

	// 鏌ヨ鍦堝瓙鐨勬墍鏈夋垚鍛?
	ListByGroup(ctx context.Context, houseGID int32, groupID int32, page, size int32) ([]*model.GameMember, int64, error)

	// 鏍规嵁 member_id 鏌ヨ鎴愬憳
	GetByID(ctx context.Context, memberID int32) (*model.GameMember, error)

	// 创建或更新成员记录
	Create(ctx context.Context, member *model.GameMember) error

	// 更新成员的圈子（用于拉圈功能）
	UpdateGroup(ctx context.Context, houseGID int32, gameID int32, newGroupID int32, newGroupName string) error

	// 查询店铺下所有成员（包括没有圈子的）
	ListByHouse(ctx context.Context, houseGID int32, page, size int32) ([]*model.GameMember, int64, error)

	// 查询店铺下没有圈子的成员
	ListByHouseWithoutGroup(ctx context.Context, houseGID int32) ([]*model.GameMember, error)

	// 清除指定圈子所有成员的圈子归属（将 group_id 设为 NULL）
	ClearGroupForAllMembers(ctx context.Context, groupID int32) error

	// 删除成员（通过 game_id 删除）
	DeleteByGameID(ctx context.Context, houseGID int32, gameID int32) error

	// 设置成员置顶状态
	SetPinStatus(ctx context.Context, houseGID int32, gameID int32, isPinned bool, pinOrder int32) error

	// 更新成员备注
	UpdateRemark(ctx context.Context, houseGID int32, gameID int32, remark string) error

	// 更新成员禁用状态
	UpdateForbid(ctx context.Context, houseGID int32, gameID int32, forbid bool) error

	// 查询店铺下所有被禁用的成员
	ListForbiddenByHouse(ctx context.Context, houseGID int32) ([]*model.GameMember, error)

	// UpsertForbid 更新禁用状态，如果不存在则创建
	UpsertForbid(ctx context.Context, houseGID int32, gameID int32, gameName string, forbid bool) error
}

type gameMemberRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewGameMemberRepo(data *infra.Data, logger log.Logger) GameMemberRepo {
	return &gameMemberRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/game_member")),
	}
}

func (r *gameMemberRepo) db(ctx context.Context) *gorm.DB      { return r.data.GetDBWithContext(ctx) }
func (r *gameMemberRepo) dbQuiet(ctx context.Context) *gorm.DB { return r.data.GetDBQuiet(ctx) }

// GetByGameIDAndGroup 鏍规嵁 game_id 鍜?group_id 鏌ヨ鎴愬憳
func (r *gameMemberRepo) GetByGameIDAndGroup(ctx context.Context, houseGID int32, gameID int32, groupID *int32) (*model.GameMember, error) {
	var member model.GameMember
	db := r.db(ctx).Where("house_gid = ? AND game_id = ?", houseGID, gameID)

	if groupID != nil {
		db = db.Where("group_id = ?", *groupID)
	} else {
		db = db.Where("group_id IS NULL")
	}

	if err := db.First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// GetByGameID 根据 game_id 查询成员，取第一个结果
func (r *gameMemberRepo) GetByGameID(ctx context.Context, houseGID int32, gameID int32) (*model.GameMember, error) {
	var member model.GameMember
	if err := r.db(ctx).
		Where("house_gid = ? AND game_id = ?", houseGID, gameID).
		First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// GetByGameIDQuiet 静默模式查询（不输出SSQL日志），用于定时任务
func (r *gameMemberRepo) GetByGameIDQuiet(ctx context.Context, houseGID int32, gameID int32) (*model.GameMember, error) {
	var member model.GameMember
	if err := r.dbQuiet(ctx).
		Where("house_gid = ? AND game_id = ?", houseGID, gameID).
		First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// ListByGameID 鏌ヨ鎴愬憳鍦ㄦ墍鏈夊湀瀛愮殑璁板綍
func (r *gameMemberRepo) ListByGameID(ctx context.Context, houseGID int32, gameID int32) ([]*model.GameMember, error) {
	var members []*model.GameMember
	if err := r.db(ctx).
		Where("house_gid = ? AND game_id = ?", houseGID, gameID).
		Order("group_id").
		Find(&members).Error; err != nil {
		return nil, err
	}

	return members, nil
}

// ListByGroup 鏌ヨ鍦堝瓙鐨勬墍鏈夋垚鍛?
func (r *gameMemberRepo) ListByGroup(ctx context.Context, houseGID int32, groupID int32, page, size int32) ([]*model.GameMember, int64, error) {
	db := r.db(ctx).Model(&model.GameMember{}).
		Where("house_gid = ? AND group_id = ?", houseGID, groupID)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 200 {
		size = 20
	}
	offset := (page - 1) * size

	var members []*model.GameMember
	if err := db.
		Order("balance DESC").
		Limit(int(size)).
		Offset(int(offset)).
		Find(&members).Error; err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

// GetByID 根据 member_id 查询成员
func (r *gameMemberRepo) GetByID(ctx context.Context, memberID int32) (*model.GameMember, error) {
	var member model.GameMember
	if err := r.db(ctx).Where("id = ?", memberID).First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// Create 创建或更新成员记录（使用 FirstOrCreate 实现幂等性）
func (r *gameMemberRepo) Create(ctx context.Context, member *model.GameMember) error {
	return r.db(ctx).
		Where("house_gid = ? AND game_id = ? AND group_id = ?", member.HouseGID, member.GameID, member.GroupID).
		FirstOrCreate(member).Error
}

// UpdateGroup 更新成员的圈子（用于拉圈功能）
// 注意：由于唯一索引包含 group_id，需要先删除旧记录再创建新记录
func (r *gameMemberRepo) UpdateGroup(ctx context.Context, houseGID int32, gameID int32, newGroupID int32, newGroupName string) error {
	return r.db(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 查询现有记录，保留余额等业务数据
		var oldMember model.GameMember
		err := tx.Where("house_gid = ? AND game_id = ?", houseGID, gameID).First(&oldMember).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		// 2. 删除所有该 game_id 在该店铺的记录（可能有多个圈子的历史记录）
		if err := tx.Where("house_gid = ? AND game_id = ?", houseGID, gameID).Delete(&model.GameMember{}).Error; err != nil {
			return err
		}

		// 3. 创建新记录，继承原有的余额等业务数据
		newMember := &model.GameMember{
			HouseGID:    houseGID,
			GameID:      gameID,
			GameName:    oldMember.GameName,
			GroupID:     &newGroupID,
			GroupName:   newGroupName,
			Balance:     oldMember.Balance, // 继承余额
			Credit:      oldMember.Credit,  // 继承信用
			Forbid:      oldMember.Forbid,  // 继承禁用状态
			Recommender: oldMember.Recommender,
		}

		return tx.Create(newMember).Error
	})
}

// ListByHouse 查询店铺下所有成员（分页）
func (r *gameMemberRepo) ListByHouse(ctx context.Context, houseGID int32, page, size int32) ([]*model.GameMember, int64, error) {
	db := r.db(ctx).Model(&model.GameMember{}).Where("house_gid = ?", houseGID)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 200 {
		size = 20
	}
	offset := (page - 1) * size

	var members []*model.GameMember
	if err := db.
		Order("created_at DESC").
		Limit(int(size)).
		Offset(int(offset)).
		Find(&members).Error; err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

// ListByHouseWithoutGroup 查询店铺下没有圈子的成员
func (r *gameMemberRepo) ListByHouseWithoutGroup(ctx context.Context, houseGID int32) ([]*model.GameMember, error) {
	var members []*model.GameMember
	err := r.db(ctx).
		Where("house_gid = ? AND (group_id IS NULL OR group_id = 0)", houseGID).
		Order("created_at DESC").
		Find(&members).Error
	return members, err
}

// ClearGroupForAllMembers 清除指定圈子所有成员的圈子归属
// 用于管理员降级时，将其圈子内的所有游戏账号设为无圈子状态
func (r *gameMemberRepo) ClearGroupForAllMembers(ctx context.Context, groupID int32) error {
	return r.db(ctx).Model(&model.GameMember{}).
		Where("group_id = ?", groupID).
		Updates(map[string]interface{}{
			"group_id":   nil,
			"group_name": "",
		}).Error
}

// DeleteByGameID 删除成员（通过 game_id 删除）
func (r *gameMemberRepo) DeleteByGameID(ctx context.Context, houseGID int32, gameID int32) error {
	return r.db(ctx).
		Where("house_gid = ? AND game_id = ?", houseGID, gameID).
		Delete(&model.GameMember{}).Error
}

// SetPinStatus 设置成员置顶状态
func (r *gameMemberRepo) SetPinStatus(ctx context.Context, houseGID int32, gameID int32, isPinned bool, pinOrder int32) error {
	return r.db(ctx).Model(&model.GameMember{}).
		Where("house_gid = ? AND game_id = ?", houseGID, gameID).
		Updates(map[string]interface{}{
			"is_pinned": isPinned,
			"pin_order": pinOrder,
		}).Error
}

// UpdateRemark 更新成员备注
func (r *gameMemberRepo) UpdateRemark(ctx context.Context, houseGID int32, gameID int32, remark string) error {
	return r.db(ctx).Model(&model.GameMember{}).
		Where("house_gid = ? AND game_id = ?", houseGID, gameID).
		Update("remark", remark).Error
}

// UpdateForbid 更新成员禁用状态（类似 passing-dragonfly 的 _freezeMembers）
func (r *gameMemberRepo) UpdateForbid(ctx context.Context, houseGID int32, gameID int32, forbid bool) error {
	return r.db(ctx).Model(&model.GameMember{}).
		Where("house_gid = ? AND game_id = ?", houseGID, gameID).
		Update("forbid", forbid).Error
}

// ListForbiddenByHouse 查询店铺下所有被禁用的成员
func (r *gameMemberRepo) ListForbiddenByHouse(ctx context.Context, houseGID int32) ([]*model.GameMember, error) {
	var members []*model.GameMember
	err := r.db(ctx).
		Where("house_gid = ? AND forbid = ?", houseGID, true).
		Find(&members).Error
	return members, err
}

// UpsertForbid 更新禁用状态，如果不存在则创建
func (r *gameMemberRepo) UpsertForbid(ctx context.Context, houseGID int32, gameID int32, gameName string, forbid bool) error {
	res := r.db(ctx).Model(&model.GameMember{}).
		Where("house_gid = ? AND game_id = ?", houseGID, gameID).
		Update("forbid", forbid)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 && gameName != "" {
		// Create new
		member := model.GameMember{
			HouseGID: houseGID,
			GameID:   gameID,
			GameName: gameName,
			Forbid:   forbid,
		}
		// 注意：这里不设置 group_id，只作为基本信息记录
		return r.db(ctx).Create(&member).Error
	}
	return nil
}
