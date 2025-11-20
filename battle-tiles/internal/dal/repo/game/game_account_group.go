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
	// Create 创建游戏账号圈子关系
	Create(ctx context.Context, group *model.GameAccountGroup) error
	
	// GetByGameAccountAndHouse 根据游戏账号ID和店铺GID查询
	GetByGameAccountAndHouse(ctx context.Context, gameAccountID, houseGID int32) (*model.GameAccountGroup, error)
	
	// ListByGameAccount 根据游戏账号ID查询所有圈子
	ListByGameAccount(ctx context.Context, gameAccountID int32) ([]*model.GameAccountGroup, error)
	
	// ListByHouse 根据店铺GID查询所有游戏账号圈子关系
	ListByHouse(ctx context.Context, houseGID int32) ([]*model.GameAccountGroup, error)
	
	// ListByGroup 根据圈子ID查询所有游戏账号
	ListByGroup(ctx context.Context, groupID int32) ([]*model.GameAccountGroup, error)
	
	// UpdateStatus 更新状态
	UpdateStatus(ctx context.Context, id int32, status string) error
	
	// Delete 删除关系（软删除或硬删除）
	Delete(ctx context.Context, id int32) error
	
	// DeleteByGameAccountAndHouse 根据游戏账号和店铺删除
	DeleteByGameAccountAndHouse(ctx context.Context, gameAccountID, houseGID int32) error
	
	// ExistsByGameAccountAndHouse 检查游戏账号是否已在某店铺的圈子中
	ExistsByGameAccountAndHouse(ctx context.Context, gameAccountID, houseGID int32) (bool, error)
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

func (r *gameAccountGroupRepo) GetByGameAccountAndHouse(ctx context.Context, gameAccountID, houseGID int32) (*model.GameAccountGroup, error) {
	var group model.GameAccountGroup
	err := r.db(ctx).Where("game_account_id = ? AND house_gid = ?", gameAccountID, houseGID).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *gameAccountGroupRepo) ListByGameAccount(ctx context.Context, gameAccountID int32) ([]*model.GameAccountGroup, error) {
	var groups []*model.GameAccountGroup
	err := r.db(ctx).Where("game_account_id = ? AND status = ?", gameAccountID, model.AccountGroupStatusActive).
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

func (r *gameAccountGroupRepo) Delete(ctx context.Context, id int32) error {
	return r.db(ctx).Delete(&model.GameAccountGroup{}, id).Error
}

func (r *gameAccountGroupRepo) DeleteByGameAccountAndHouse(ctx context.Context, gameAccountID, houseGID int32) error {
	return r.db(ctx).Where("game_account_id = ? AND house_gid = ?", gameAccountID, houseGID).
		Delete(&model.GameAccountGroup{}).Error
}

func (r *gameAccountGroupRepo) ExistsByGameAccountAndHouse(ctx context.Context, gameAccountID, houseGID int32) (bool, error) {
	var count int64
	err := r.db(ctx).Model(&model.GameAccountGroup{}).
		Where("game_account_id = ? AND house_gid = ? AND status = ?", gameAccountID, houseGID, model.AccountGroupStatusActive).
		Count(&count).Error
	return count > 0, err
}

