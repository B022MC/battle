package game

import (
	"context"
	"fmt"

	model "battle-tiles/internal/dal/model/game"

	"gorm.io/gorm"
)

type RoomCreditLimitRepo interface {
	// Upsert 插入或更新房间额度限制
	Upsert(ctx context.Context, limit *model.GameRoomCreditLimit) error

	// Get 获取特定的房间额度限制
	Get(ctx context.Context, houseGID int32, groupName string, gameKind int32, baseScore int32) (*model.GameRoomCreditLimit, error)

	// List 获取指定店铺的所有额度限制
	List(ctx context.Context, houseGID int32) ([]*model.GameRoomCreditLimit, error)

	// ListByGroup 获取指定店铺和圈子的额度限制
	ListByGroup(ctx context.Context, houseGID int32, groupName string) ([]*model.GameRoomCreditLimit, error)

	// Delete 删除房间额度限制
	Delete(ctx context.Context, houseGID int32, groupName string, gameKind int32, baseScore int32) error

	// GetCreditLimit 按优先级查找房间额度限制（不包含玩家个人额度）
	// 优先级：圈子+游戏类型+底分 > 圈子默认 > 全局+游戏类型+底分 > 全局默认
	GetCreditLimit(ctx context.Context, houseGID int32, groupName string, gameKind int32, baseScore int32) (int32, bool)
}

type roomCreditLimitRepo struct {
	db *gorm.DB
}

func NewRoomCreditLimitRepo(db *gorm.DB) RoomCreditLimitRepo {
	return &roomCreditLimitRepo{db: db}
}

func (r *roomCreditLimitRepo) Upsert(ctx context.Context, limit *model.GameRoomCreditLimit) error {
	return r.db.WithContext(ctx).
		Where("house_gid = ? AND group_name = ? AND game_kind = ? AND base_score = ?",
			limit.HouseGID, limit.GroupName, limit.GameKind, limit.BaseScore).
		Assign(map[string]interface{}{
			"credit_limit": limit.CreditLimit,
			"updated_by":   limit.UpdatedBy,
		}).
		FirstOrCreate(limit).Error
}

func (r *roomCreditLimitRepo) Get(ctx context.Context, houseGID int32, groupName string, gameKind int32, baseScore int32) (*model.GameRoomCreditLimit, error) {
	var limit model.GameRoomCreditLimit
	err := r.db.WithContext(ctx).
		Where("house_gid = ? AND group_name = ? AND game_kind = ? AND base_score = ?",
			houseGID, groupName, gameKind, baseScore).
		First(&limit).Error
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

func (r *roomCreditLimitRepo) List(ctx context.Context, houseGID int32) ([]*model.GameRoomCreditLimit, error) {
	var limits []*model.GameRoomCreditLimit
	err := r.db.WithContext(ctx).
		Where("house_gid = ?", houseGID).
		Order("group_name, game_kind, base_score").
		Find(&limits).Error
	return limits, err
}

func (r *roomCreditLimitRepo) ListByGroup(ctx context.Context, houseGID int32, groupName string) ([]*model.GameRoomCreditLimit, error) {
	var limits []*model.GameRoomCreditLimit
	err := r.db.WithContext(ctx).
		Where("house_gid = ? AND group_name = ?", houseGID, groupName).
		Order("game_kind, base_score").
		Find(&limits).Error
	return limits, err
}

func (r *roomCreditLimitRepo) Delete(ctx context.Context, houseGID int32, groupName string, gameKind int32, baseScore int32) error {
	return r.db.WithContext(ctx).
		Where("house_gid = ? AND group_name = ? AND game_kind = ? AND base_score = ?",
			houseGID, groupName, gameKind, baseScore).
		Delete(&model.GameRoomCreditLimit{}).Error
}

// GetCreditLimit 按优先级查找房间额度限制
// 查找顺序（与 passing-dragonfly 一致）：
// 1. 圈子+游戏类型+底分 (group-kind-base)
// 2. 圈子默认 (group-0-0)
// 3. 全局+游戏类型+底分 (-kind-base)
// 4. 全局默认 (0-0)
// 5. 兜底值：99999元
func (r *roomCreditLimitRepo) GetCreditLimit(ctx context.Context, houseGID int32, groupName string, gameKind int32, baseScore int32) (int32, bool) {
	var limit model.GameRoomCreditLimit

	// 1. 圈子+游戏类型+底分
	if groupName != "" {
		err := r.db.WithContext(ctx).
			Where("house_gid = ? AND group_name = ? AND game_kind = ? AND base_score = ?",
				houseGID, groupName, gameKind, baseScore).
			First(&limit).Error
		if err == nil {
			return limit.CreditLimit, true
		}
	}

	// 2. 圈子默认
	if groupName != "" {
		err := r.db.WithContext(ctx).
			Where("house_gid = ? AND group_name = ? AND game_kind = 0 AND base_score = 0",
				houseGID, groupName).
			First(&limit).Error
		if err == nil {
			return limit.CreditLimit, true
		}
	}

	// 3. 全局+游戏类型+底分
	err := r.db.WithContext(ctx).
		Where("house_gid = ? AND group_name = '' AND game_kind = ? AND base_score = ?",
			houseGID, gameKind, baseScore).
		First(&limit).Error
	if err == nil {
		return limit.CreditLimit, true
	}

	// 4. 全局默认
	err = r.db.WithContext(ctx).
		Where("house_gid = ? AND group_name = '' AND game_kind = 0 AND base_score = 0",
			houseGID).
		First(&limit).Error
	if err == nil {
		return limit.CreditLimit, true
	}

	// 5. 兜底值（建议配置为一个很大的数）
	return 99999 * 100, true
}

// FormatCreditKey 格式化额度限制的键（用于日志和调试）
func FormatCreditKey(groupName string, gameKind int32, baseScore int32) string {
	if groupName == "" {
		return fmt.Sprintf("-%d-%d", gameKind, baseScore)
	}
	return fmt.Sprintf("%s-%d-%d", groupName, gameKind, baseScore)
}
