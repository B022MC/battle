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

	// 按游戏用户ID查询（保留用于向后兼容）
	ListByPlayer(ctx context.Context, houseGID int32, playerGameID interface{}, groupID *int32, start, end *time.Time, page, size int32) ([]*model.GameBattleRecord, int64, error)

	// 按游戏账号查询
	ListByPlayerGameName(ctx context.Context, houseGID int32, playerGameName string, groupID *int32, start, end *time.Time, page, size int32) ([]*model.GameBattleRecord, int64, error)

	// 按游戏用户ID获取统计（保留用于向后兼容）
	GetPlayerStats(ctx context.Context, houseGID int32, playerGameID interface{}, groupID *int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, err error)

	// 按游戏账号获取统计
	GetPlayerStatsByGameName(ctx context.Context, houseGID int32, playerGameName string, groupID *int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, err error)

	GetGroupStats(ctx context.Context, houseGID int32, groupID int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, activeMembers int64, err error)

	GetHouseStats(ctx context.Context, houseGID int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, err error)
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

// SaveBatchWithDedup 鎵归噺淇濆瓨鎴樼哗锛岃嚜鍔ㄥ幓閲嶏紙鏍规嵁 battle_at + player_game_id 鍒ゆ柇锛?
func (r *battleRecordRepo) SaveBatchWithDedup(ctx context.Context, list []*model.GameBattleRecord) (int, error) {
	if len(list) == 0 {
		return 0, nil
	}

	saved := 0
	for _, record := range list {
		// 妫€鏌ユ槸鍚﹀凡瀛樺湪
		var existing model.GameBattleRecord
		err := r.db(ctx).Where("battle_at = ? AND player_game_id = ? AND house_gid = ?",
			record.BattleAt, record.PlayerGameID, record.HouseGID).
			First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			// 涓嶅瓨鍦紝鎻掑叆
			if err := r.db(ctx).Create(record).Error; err != nil {
				r.log.Errorf("save battle record failed: %v", err)
				continue
			}
			saved++
		} else if err != nil {
			r.log.Errorf("check battle record existence failed: %v", err)
			continue
		}
		// 宸插瓨鍦紝璺宠繃
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
	} // 绮楃硻杩囨护锛屽墠绔啀缁嗗垎
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

func (r *battleRecordRepo) ListByPlayer(ctx context.Context, houseGID int32, playerGameID interface{}, groupID *int32, start, end *time.Time, page, size int32) ([]*model.GameBattleRecord, int64, error) {
	// 支持 int32 和 string 两种类型
	var db *gorm.DB
	if str, ok := playerGameID.(string); ok {
		db = r.db(ctx).Model(&model.GameBattleRecord{}).
			Where("house_gid = ? AND player_game_name = ?", houseGID, str)
	} else {
		db = r.db(ctx).Model(&model.GameBattleRecord{}).
			Where("house_gid = ? AND player_game_id = ?", houseGID, playerGameID)
	}

	if groupID != nil {
		db = db.Where("group_id = ?", *groupID)
	}

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

// ListByPlayerGameName 按游戏账号查询战绩
func (r *battleRecordRepo) ListByPlayerGameName(ctx context.Context, houseGID int32, playerGameName string, groupID *int32, start, end *time.Time, page, size int32) ([]*model.GameBattleRecord, int64, error) {
	db := r.db(ctx).Model(&model.GameBattleRecord{}).
		Where("house_gid = ? AND player_game_ID = ?", houseGID, playerGameName)

	if groupID != nil {
		db = db.Where("group_id = ?", *groupID)
	}

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

// GetPlayerStats 鑾峰彇鐜╁缁熻鏁版嵁
func (r *battleRecordRepo) GetPlayerStats(ctx context.Context, houseGID int32, playerGameID interface{}, groupID *int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, err error) {
	// 支持 int32 和 string 两种类型
	var db *gorm.DB
	if str, ok := playerGameID.(string); ok {
		db = r.db(ctx).Model(&model.GameBattleRecord{}).
			Where("house_gid = ? AND player_game_name = ?", houseGID, str)
	} else {
		db = r.db(ctx).Model(&model.GameBattleRecord{}).
			Where("house_gid = ? AND player_game_id = ?", houseGID, playerGameID)
	}

	// 鏂板: 鏀寔鎸夊湀瀛愮瓫閫?
	if groupID != nil {
		db = db.Where("group_id = ?", *groupID)
	}

	if start != nil {
		db = db.Where("battle_at >= ?", *start)
	}
	if end != nil {
		db = db.Where("battle_at < ?", *end)
	}

	// 缁熻鎬诲眬鏁?
	if err = db.Count(&totalGames).Error; err != nil {
		return 0, 0, 0, err
	}

	// 缁熻鎬诲垎鏁板拰鎬昏垂鐢?
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

// GetPlayerStatsByGameName 按游戏账号获取玩家统计数据
func (r *battleRecordRepo) GetPlayerStatsByGameName(ctx context.Context, houseGID int32, playerGameName string, groupID *int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, err error) {
	db := r.db(ctx).Model(&model.GameBattleRecord{}).
		Where("house_gid = ? AND player_game_name = ?", houseGID, playerGameName)

	if groupID != nil {
		db = db.Where("group_id = ?", *groupID)
	}

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

// GetGroupStats 鑾峰彇鍦堝瓙缁熻鏁版嵁
func (r *battleRecordRepo) GetGroupStats(ctx context.Context, houseGID int32, groupID int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, activeMembers int64, err error) {
	db := r.db(ctx).Model(&model.GameBattleRecord{}).
		Where("house_gid = ? AND group_id = ?", houseGID, groupID)

	if start != nil {
		db = db.Where("battle_at >= ?", *start)
	}
	if end != nil {
		db = db.Where("battle_at < ?", *end)
	}

	// 缁熻鎬诲眬鏁?
	if err = db.Count(&totalGames).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	type Result struct {
		TotalScore int
		TotalFee   int
	}
	var result Result
	err = db.Select("COALESCE(SUM(score), 0) as total_score, COALESCE(SUM(fee), 0) as total_fee").
		Scan(&result).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}

	err = r.db(ctx).Model(&model.GameBattleRecord{}).
		Where("house_gid = ? AND group_id = ?", houseGID, groupID).
		Distinct("player_game_id").
		Count(&activeMembers).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return totalGames, result.TotalScore, result.TotalFee, activeMembers, nil
}

// GetHouseStats 鑾峰彇搴楅摵缁熻鏁版嵁
func (r *battleRecordRepo) GetHouseStats(ctx context.Context, houseGID int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, err error) {
	db := r.db(ctx).Model(&model.GameBattleRecord{}).
		Where("house_gid = ?", houseGID)

	if start != nil {
		db = db.Where("battle_at >= ?", *start)
	}
	if end != nil {
		db = db.Where("battle_at < ?", *end)
	}

	// 缁熻鎬诲眬鏁?
	if err = db.Count(&totalGames).Error; err != nil {
		return 0, 0, 0, err
	}

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
