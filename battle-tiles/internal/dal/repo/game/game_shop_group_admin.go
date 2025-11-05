package game

import (
	"context"

	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"

	"github.com/go-kratos/kratos/v2/log"
)

type GameShopGroupAdminRepo interface {
	Exists(ctx context.Context, houseGID, groupID, userID int32) (bool, error)
	ListByUser(ctx context.Context, userID int32) ([]*model.GameShopGroupAdmin, error)
	ListByHouse(ctx context.Context, houseGID int32) ([]*model.GameShopGroupAdmin, error)
	Upsert(ctx context.Context, in *model.GameShopGroupAdmin) error
}

type gameShopGroupAdminRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewShopGroupAdminRepo(data *infra.Data, logger log.Logger) GameShopGroupAdminRepo {
	return &gameShopGroupAdminRepo{data: data, log: log.NewHelper(log.With(logger, "module", "repo/game_shop_group_admin"))}
}

func (r *gameShopGroupAdminRepo) Exists(ctx context.Context, houseGID, groupID, userID int32) (bool, error) {
	var cnt int64
	if err := r.data.GetDBWithContext(ctx).Model(&model.GameShopGroupAdmin{}).
		Where("house_gid = ? AND group_id = ? AND user_id = ?", houseGID, groupID, userID).
		Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (r *gameShopGroupAdminRepo) ListByUser(ctx context.Context, userID int32) ([]*model.GameShopGroupAdmin, error) {
	var out []*model.GameShopGroupAdmin
	err := r.data.GetDBWithContext(ctx).
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&out).Error
	return out, err
}

func (r *gameShopGroupAdminRepo) ListByHouse(ctx context.Context, houseGID int32) ([]*model.GameShopGroupAdmin, error) {
	var out []*model.GameShopGroupAdmin
	err := r.data.GetDBWithContext(ctx).
		Where("house_gid = ?", houseGID).
		Order("id DESC").
		Find(&out).Error
	return out, err
}

// Upsert 以 (house_gid,user_id) 唯一，更新 group_id/role
func (r *gameShopGroupAdminRepo) Upsert(ctx context.Context, in *model.GameShopGroupAdmin) error {
	db := r.data.GetDBWithContext(ctx)
	// 若无唯一约束，这里采用先查后更/插
	var existing model.GameShopGroupAdmin
	err := db.Where("house_gid = ? AND user_id = ?", in.HouseGID, in.UserID).First(&existing).Error
	if err == nil {
		existing.GroupID = in.GroupID
		existing.Role = in.Role
		return db.Save(&existing).Error
	}
	return db.Create(in).Error
}
