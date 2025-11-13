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
	SaveBatchWithDedup(ctx context.Context, list []*model.GameBattleRecord) (int, error)
	List(ctx context.Context, houseGID int32, groupID *int32, gameID *int32, start, end *time.Time, page, size int32) ([]*model.GameBattleRecord, int64, error)
	ListByPlayer(ctx context.Context, houseGID int32, playerGameID int32, start, end *time.Time, page, size int32) ([]*model.GameBattleRecord, int64, error)
	GetPlayerStats(ctx context.Context, houseGID int32, playerGameID int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, err error)
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

// SaveBatchWithDedup 批量保存战绩，自动去重（根据 battle_at + player_game_id 判断）
func (r *battleRecordRepo) SaveBatchWithDedup(ctx context.Context, list []*model.GameBattleRecord) (int, error) {
	if len(list) == 0 {
		return 0, nil
	}

	saved := 0
	for _, record := range list {
		// 检查是否已存在
		var existing model.GameBattleRecord
		err := r.db(ctx).Where("battle_at = ? AND player_game_id = ? AND house_gid = ?",
			record.BattleAt, record.PlayerGameID, record.HouseGID).
			First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			// 不存在，插入
			if err := r.db(ctx).Create(record).Error; err != nil {
				r.log.Errorf("save battle record failed: %v", err)
				continue
			}
			saved++
		} else if err != nil {
			r.log.Errorf("check battle record existence failed: %v", err)
			continue
		}
		// 已存在，跳过
	}

	return saved, nil
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

// ListByPlayer 查询指定玩家的战绩
func (r *battleRecordRepo) ListByPlayer(ctx context.Context, houseGID int32, playerGameID int32, start, end *time.Time, page, size int32) ([]*model.GameBattleRecord, int64, error) {
	db := r.db(ctx).Model(&model.GameBattleRecord{}).
		Where("house_gid = ? AND player_game_id = ?", houseGID, playerGameID)

	if start != nil {
		db = db.Where("battle_at >= ?", *start)
	}
	if end != nil {
		db = db.Where("battle_at < ?", *end)
	}

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

// GetPlayerStats 获取玩家统计数据
func (r *battleRecordRepo) GetPlayerStats(ctx context.Context, houseGID int32, playerGameID int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, err error) {
	db := r.db(ctx).Model(&model.GameBattleRecord{}).
		Where("house_gid = ? AND player_game_id = ?", houseGID, playerGameID)

	if start != nil {
		db = db.Where("battle_at >= ?", *start)
	}
	if end != nil {
		db = db.Where("battle_at < ?", *end)
	}

	// 统计总局数
	if err = db.Count(&totalGames).Error; err != nil {
		return 0, 0, 0, err
	}

	// 统计总分数和总费用
	type Result struct {
		TotalScore int
		TotalFee   int
	}
	var result Result
	err = db.Select("COALESCE(SUM(score), 0) as total_score, COALESCE(SUM(fee), 0) as total_fee").
		Scan(&result).Error
	if err != nil {
		return 0, 0, 0, err
	}

	return totalGames, result.TotalScore, result.TotalFee, nil
}
