package game

import (
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HouseSettingsRepo interface {
	Upsert(ctx context.Context, in *model.GameHouseSettings) error
	Get(ctx context.Context, houseGID int32) (*model.GameHouseSettings, error)
}

type houseSettingsRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewHouseSettingsRepo(data *infra.Data, logger log.Logger) HouseSettingsRepo {
	return &houseSettingsRepo{data: data, log: log.NewHelper(log.With(logger, "module", "repo/house_settings"))}
}

func (r *houseSettingsRepo) db(ctx context.Context) *gorm.DB { return r.data.GetDBWithContext(ctx) }

func (r *houseSettingsRepo) Upsert(ctx context.Context, in *model.GameHouseSettings) error {
	// upsert by house_gid
	return r.db(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "house_gid"}},
			DoUpdates: clause.AssignmentColumns([]string{"fees_json", "share_fee", "push_credit", "updated_at", "updated_by"}),
		},
	).Create(in).Error
}

func (r *houseSettingsRepo) Get(ctx context.Context, houseGID int32) (*model.GameHouseSettings, error) {
	var m model.GameHouseSettings
	if err := r.db(ctx).Where("house_gid = ?", houseGID).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
