package game

import (
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
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
	// OnConflict DO UPDATE SET ...
	// 仅更新这些字段
	// 注意：gorm 的 OnConflict 需要导入 clause，但为了尽量少引入，这里用 Save 语义：若存在主键则更新，否则插入。
	).Save(in).Error
}

func (r *houseSettingsRepo) Get(ctx context.Context, houseGID int32) (*model.GameHouseSettings, error) {
	var m model.GameHouseSettings
	if err := r.db(ctx).Where("house_gid = ?", houseGID).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
