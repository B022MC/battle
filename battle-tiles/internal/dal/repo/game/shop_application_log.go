package game

import (
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type ShopApplicationLogRepo interface {
	// Create creates a new application log
	Create(ctx context.Context, log *model.GameShopApplicationLog) error

	// ListByHouse lists application logs by house_gid with pagination
	ListByHouse(ctx context.Context, houseGID int32, limit, offset int) ([]*model.GameShopApplicationLog, error)

	// ListByAdmin lists application logs by admin_user_id
	ListByAdmin(ctx context.Context, adminUserID int32, limit, offset int) ([]*model.GameShopApplicationLog, error)

	// CountByAction counts logs by action type for a house
	CountByAction(ctx context.Context, houseGID int32, action string, start, end time.Time) (int64, error)

	// CountByHouse counts total logs for a house
	CountByHouse(ctx context.Context, houseGID int32) (int64, error)
}

type shopApplicationLogRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewShopApplicationLogRepo(data *infra.Data, logger log.Logger) ShopApplicationLogRepo {
	return &shopApplicationLogRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/shop_application_log")),
	}
}

func (r *shopApplicationLogRepo) db(ctx context.Context) *gorm.DB {
	return r.data.GetDBWithContext(ctx)
}

func (r *shopApplicationLogRepo) Create(ctx context.Context, logEntry *model.GameShopApplicationLog) error {
	return r.db(ctx).Create(logEntry).Error
}

func (r *shopApplicationLogRepo) ListByHouse(ctx context.Context, houseGID int32, limit, offset int) ([]*model.GameShopApplicationLog, error) {
	var logs []*model.GameShopApplicationLog
	err := r.db(ctx).
		Where("house_gid = ?", houseGID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

func (r *shopApplicationLogRepo) ListByAdmin(ctx context.Context, adminUserID int32, limit, offset int) ([]*model.GameShopApplicationLog, error) {
	var logs []*model.GameShopApplicationLog
	err := r.db(ctx).
		Where("admin_user_id = ?", adminUserID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

func (r *shopApplicationLogRepo) CountByAction(ctx context.Context, houseGID int32, action string, start, end time.Time) (int64, error) {
	var count int64
	err := r.db(ctx).Model(&model.GameShopApplicationLog{}).
		Where("house_gid = ? AND action = ? AND created_at >= ? AND created_at < ?", houseGID, action, start, end).
		Count(&count).Error
	return count, err
}

func (r *shopApplicationLogRepo) CountByHouse(ctx context.Context, houseGID int32) (int64, error) {
	var count int64
	err := r.db(ctx).Model(&model.GameShopApplicationLog{}).
		Where("house_gid = ?", houseGID).
		Count(&count).Error
	return count, err
}
