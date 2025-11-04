// internal/dal/repo/game/shop_admin_repo.go
package game

import (
	"context"

	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GameShopAdminRepo interface {
	Assign(ctx context.Context, m *model.GameShopAdmin) error
	Revoke(ctx context.Context, houseGID int32, userID int32) error
	Exists(ctx context.Context, houseGID int32, userID int32) (bool, error)
	ListByHouse(ctx context.Context, houseGID int32) ([]*model.GameShopAdmin, error)
}

type gameShopAdminRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewShopAdminRepo(data *infra.Data, logger log.Logger) GameShopAdminRepo {
	return &gameShopAdminRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/game_shop_admin")),
	}
}

func (r *gameShopAdminRepo) db(ctx context.Context) *gorm.DB { return r.data.GetDBWithContext(ctx) }

// Assign：幂等UPSERT（唯一键 house_gid + user_id），更新 role / updated_at
func (r *gameShopAdminRepo) Assign(ctx context.Context, m *model.GameShopAdmin) error {
	return r.db(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "house_gid"}, {Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"role", "updated_at"}),
		}).
		Create(m).Error
}

func (r *gameShopAdminRepo) Revoke(ctx context.Context, houseGID int32, userID int32) error {
	return r.db(ctx).
		Where("house_gid = ? AND user_id = ?", houseGID, userID).
		Delete(&model.GameShopAdmin{}).Error
}

func (r *gameShopAdminRepo) Exists(ctx context.Context, houseGID int32, userID int32) (bool, error) {
	var cnt int64
	if err := r.db(ctx).
		Model(&model.GameShopAdmin{}).
		Where("house_gid = ? AND user_id = ?", houseGID, userID).
		Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (r *gameShopAdminRepo) ListByHouse(ctx context.Context, houseGID int32) ([]*model.GameShopAdmin, error) {
	var out []*model.GameShopAdmin
	err := r.db(ctx).
		Where("house_gid = ?", houseGID).
		Order("id DESC").
		Find(&out).Error
	return out, err
}
