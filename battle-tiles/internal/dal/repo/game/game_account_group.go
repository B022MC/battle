// internal/dal/repo/game/game_account_group.go
package game

import (
	"context"

	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type GameAccountGroupRepo interface {
	// Create 创建游戏玩家圈子关系
	Create(ctx context.Context, group *model.GameAccountGroup) error

	// GetByGamePlayerAndHouse 根据游戏玩家ID和店铺GID查询
	GetByGamePlayerAndHouse(ctx context.Context, gamePlayerID string, houseGID int32) (*model.GameAccountGroup, error)

	// GetActiveByGamePlayerAndHouse 根据游戏玩家ID和店铺GID查询活跃圈子（用于战绩同步）
	GetActiveByGamePlayerAndHouse(ctx context.Context, gamePlayerID string, houseGID int32) (*model.GameAccountGroup, error)

	// ListByGamePlayer 根据游戏玩家ID查询所有圈子
	ListByGamePlayer(ctx context.Context, gamePlayerID string) ([]*model.GameAccountGroup, error)

	// ListByHouse 根据店铺GID查询所有游戏玩家圈子关系
	ListByHouse(ctx context.Context, houseGID int32) ([]*model.GameAccountGroup, error)

	// ListByGroup 根据圈子ID查询所有游戏玩家
	ListByGroup(ctx context.Context, groupID int32) ([]*model.GameAccountGroup, error)

	// UpdateStatus 更新状态
	UpdateStatus(ctx context.Context, id int32, status string) error

	// UpdateGroupByGamePlayerAndHouse 更新游戏玩家在某店铺的圈子（用于拉圈）
	UpdateGroupByGamePlayerAndHouse(ctx context.Context, gamePlayerID string, houseGID, newGroupID int32, newGroupName string) error

	// DeactivateAllByGroup 将指定圈子的所有游戏玩家设为无圈子状态（status=inactive）
	DeactivateAllByGroup(ctx context.Context, groupID int32) error

	// Delete 删除关系（软删除或硬删除）
	Delete(ctx context.Context, id int32) error

	// DeleteByGamePlayerAndHouse 根据游戏玩家和店铺删除
	DeleteByGamePlayerAndHouse(ctx context.Context, gamePlayerID string, houseGID int32) error

	// ExistsByGamePlayerAndHouse 检查游戏玩家是否已在某店铺的圈子中
	ExistsByGamePlayerAndHouse(ctx context.Context, gamePlayerID string, houseGID int32) (bool, error)
}

type gameAccountGroupRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewGameAccountGroupRepo(data *infra.Data, logger log.Logger) GameAccountGroupRepo {
	return &gameAccountGroupRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/game_account_group")),
	}
}

func (r *gameAccountGroupRepo) db(ctx context.Context) *gorm.DB {
	return r.data.GetDBWithContext(ctx)
}

func (r *gameAccountGroupRepo) Create(ctx context.Context, group *model.GameAccountGroup) error {
	return r.db(ctx).Create(group).Error
}

func (r *gameAccountGroupRepo) GetByGamePlayerAndHouse(ctx context.Context, gamePlayerID string, houseGID int32) (*model.GameAccountGroup, error) {
	var group model.GameAccountGroup
	err := r.db(ctx).Where("game_player_id = ? AND house_gid = ?", gamePlayerID, houseGID).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// GetActiveByGamePlayerAndHouse 查询游戏玩家在指定店铺的活跃圈子
// 用于战绩同步时确定当前圈子归属
func (r *gameAccountGroupRepo) GetActiveByGamePlayerAndHouse(ctx context.Context, gamePlayerID string, houseGID int32) (*model.GameAccountGroup, error) {
	var group model.GameAccountGroup
	err := r.db(ctx).Where("game_player_id = ? AND house_gid = ? AND status = ?",
		gamePlayerID, houseGID, model.AccountGroupStatusActive).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *gameAccountGroupRepo) ListByGamePlayer(ctx context.Context, gamePlayerID string) ([]*model.GameAccountGroup, error) {
	var groups []*model.GameAccountGroup
	err := r.db(ctx).Where("game_player_id = ? AND status = ?", gamePlayerID, model.AccountGroupStatusActive).
		Order("created_at DESC").
		Find(&groups).Error
	return groups, err
}

func (r *gameAccountGroupRepo) ListByHouse(ctx context.Context, houseGID int32) ([]*model.GameAccountGroup, error) {
	var groups []*model.GameAccountGroup
	err := r.db(ctx).Where("house_gid = ? AND status = ?", houseGID, model.AccountGroupStatusActive).
		Order("created_at DESC").
		Find(&groups).Error
	return groups, err
}

func (r *gameAccountGroupRepo) ListByGroup(ctx context.Context, groupID int32) ([]*model.GameAccountGroup, error) {
	var groups []*model.GameAccountGroup
	err := r.db(ctx).Where("group_id = ? AND status = ?", groupID, model.AccountGroupStatusActive).
		Order("created_at DESC").
		Find(&groups).Error
	return groups, err
}

func (r *gameAccountGroupRepo) UpdateStatus(ctx context.Context, id int32, status string) error {
	return r.db(ctx).Model(&model.GameAccountGroup{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdateGroupByGamePlayerAndHouse 更新游戏玩家在某店铺的圈子（用于拉圈）
// 逻辑：将旧圈子设为 inactive，创建或激活新圈子
func (r *gameAccountGroupRepo) UpdateGroupByGamePlayerAndHouse(ctx context.Context, gamePlayerID string, houseGID, newGroupID int32, newGroupName string) error {
	return r.db(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 将该玩家在该店铺的所有圈子关系设为 inactive
		if err := tx.Model(&model.GameAccountGroup{}).
			Where("game_player_id = ? AND house_gid = ?", gamePlayerID, houseGID).
			Update("status", model.AccountGroupStatusInactive).Error; err != nil {
			return err
		}

		// 2. 查找是否已存在该圈子的记录
		var existing model.GameAccountGroup
		err := tx.Where("game_player_id = ? AND house_gid = ? AND group_id = ?",
			gamePlayerID, houseGID, newGroupID).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			// 3a. 不存在，创建新记录
			newRecord := &model.GameAccountGroup{
				GamePlayerID: gamePlayerID,
				HouseGID:     houseGID,
				GroupID:      newGroupID,
				GroupName:    newGroupName,
				Status:       model.AccountGroupStatusActive,
			}
			return tx.Create(newRecord).Error
		} else if err != nil {
			return err
		} else {
			// 3b. 已存在，激活它
			return tx.Model(&model.GameAccountGroup{}).
				Where("id = ?", existing.Id).
				Updates(map[string]interface{}{
					"status":     model.AccountGroupStatusActive,
					"group_name": newGroupName,
				}).Error
		}
	})
}

// DeactivateAllByGroup 将指定圈子的所有游戏账号设为无圈子状态
// 用于管理员降级时，清空其管理的圈子
func (r *gameAccountGroupRepo) DeactivateAllByGroup(ctx context.Context, groupID int32) error {
	return r.db(ctx).Model(&model.GameAccountGroup{}).
		Where("group_id = ?", groupID).
		Update("status", model.AccountGroupStatusInactive).Error
}

func (r *gameAccountGroupRepo) Delete(ctx context.Context, id int32) error {
	return r.db(ctx).Delete(&model.GameAccountGroup{}, id).Error
}

func (r *gameAccountGroupRepo) DeleteByGamePlayerAndHouse(ctx context.Context, gamePlayerID string, houseGID int32) error {
	return r.db(ctx).Where("game_player_id = ? AND house_gid = ?", gamePlayerID, houseGID).
		Delete(&model.GameAccountGroup{}).Error
}

func (r *gameAccountGroupRepo) ExistsByGamePlayerAndHouse(ctx context.Context, gamePlayerID string, houseGID int32) (bool, error) {
	var count int64
	err := r.db(ctx).Model(&model.GameAccountGroup{}).
		Where("game_player_id = ? AND house_gid = ? AND status = ?", gamePlayerID, houseGID, model.AccountGroupStatusActive).
		Count(&count).Error
	return count > 0, err
}
