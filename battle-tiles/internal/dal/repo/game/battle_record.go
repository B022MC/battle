package game

import (
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type BattleRecordRepo interface {
	SaveBatch(ctx context.Context, list []*model.GameBattleRecord) error
	List(ctx context.Context, houseGID int32, groupID *int32, gameID *int32, start, end *time.Time, page, size int32) ([]*model.GameBattleRecord, int64, error)
}

type battleRecordRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewBattleRecordRepo(data *infra.Data, logger log.Logger) BattleRecordRepo {
	return &battleRecordRepo{data: data, log: log.NewHelper(log.With(logger, "module", "repo/battle_record"))}
}

func (r *battleRecordRepo) db(ctx context.Context) *gorm.DB { return r.data.GetDBWithContext(ctx) }

func (r *battleRecordRepo) SaveBatch(ctx context.Context, list []*model.GameBattleRecord) error {
	if len(list) == 0 {
		return nil
	}
	return r.db(ctx).Create(&list).Error
}

func (r *battleRecordRepo) List(ctx context.Context, houseGID int32, groupID *int32, gameID *int32, start, end *time.Time, page, size int32) ([]*model.GameBattleRecord, int64, error) {
	db := r.db(ctx).Model(&model.GameBattleRecord{}).Where("house_gid = ?", houseGID)
	if groupID != nil {
		db = db.Where("group_id = ?", *groupID)
	}
	if start != nil {
		db = db.Where("battle_at >= ?", *start)
	}
	if end != nil {
		db = db.Where("battle_at < ?", *end)
	}
	if gameID != nil {
		db = db.Where("players_json LIKE ?", "%\"game_id\":%")
	} // 粗糙过滤，前端再细分
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
	var out []*model.GameBattleRecord
	if err := db.Order("battle_at DESC").Limit(int(size)).Offset(int(offset)).Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, total, nil
}
