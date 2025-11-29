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

	GetHouseStatsDetail(ctx context.Context, houseGID int32, start, end *time.Time) (*HouseStatsDetail, error)
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

	return totalGames, totalScore, totalFee, activeMembers, nil
}

// GetHouseStats 获取店铺统计数据
func (r *battleRecordRepo) GetHouseStats(ctx context.Context, houseGID int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, err error) {
	db := r.db(ctx).Model(&model.GameBattleRecord{}).
		Where("house_gid = ?", houseGID)

	if start != nil {
		db = db.Where("battle_at >= ?", *start)
	}
	if end != nil {
		db = db.Where("battle_at < ?", *end)
	}

	// 缁熻鎬诲眬鏁?
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

// GroupPayoff 圈子间转账记录
type GroupPayoff struct {
	Group string `json:"group"`
	Value int    `json:"value"` // 正值表示收入，负值表示支出
}

// HouseStatsDetail 店铺完整统计数据
type HouseStatsDetail struct {
	HouseGID       int32         `json:"house_gid"`
	TotalGames     int64         `json:"total_games"`     // 总局数
	TotalScore     int           `json:"total_score"`     // 总得分
	TotalFee       int           `json:"total_fee"`       // 总手续费
	RechargeShang  int           `json:"recharge_shang"`  // 上分
	RechargeXia    int           `json:"recharge_xia"`    // 下分
	BalancePay     int           `json:"balance_pay"`     // 待提（正余额）
	BalanceTake    int           `json:"balance_take"`    // 欠费（负余额）
	BalancePayoffs []GroupPayoff `json:"balance_payoffs"` // 圈子间转账
	FeePayoffs     []GroupPayoff `json:"fee_payoffs"`     // 运费分成
}

// GetHouseStatsDetail 获取店铺完整统计数据（包括充值、余额、圈子间转账等）
func (r *battleRecordRepo) GetHouseStatsDetail(ctx context.Context, houseGID int32, start, end *time.Time) (*HouseStatsDetail, error) {
	stats := &HouseStatsDetail{
		HouseGID:       houseGID,
		BalancePayoffs: []GroupPayoff{},
		FeePayoffs:     []GroupPayoff{},
	}

	// 1. 获取对局统计（总局数、总得分、总手续费）
	var err error
	stats.TotalGames, stats.TotalScore, stats.TotalFee, err = r.GetHouseStats(ctx, houseGID, start, end)
	if err != nil {
		return nil, err
	}

	// 2. 获取充值统计（上分/下分）
	rechargeDB := r.db(ctx).Model(&model.GameRechargeRecord{}).
		Where("house_gid = ?", houseGID)
	if start != nil {
		rechargeDB = rechargeDB.Where("recharged_at >= ?", *start)
	}
	if end != nil {
		rechargeDB = rechargeDB.Where("recharged_at < ?", *end)
	}

	// 上分（正值）
	var rechargeShang struct{ Total int }
	if err := rechargeDB.Where("amount > 0").Select("COALESCE(SUM(amount), 0) as total").Scan(&rechargeShang).Error; err != nil {
		return nil, err
	}
	stats.RechargeShang = rechargeShang.Total

	// 下分（负值）
	var rechargeXia struct{ Total int }
	if err := rechargeDB.Where("amount < 0").Select("COALESCE(SUM(amount), 0) as total").Scan(&rechargeXia).Error; err != nil {
		return nil, err
	}
	stats.RechargeXia = rechargeXia.Total

	// 3. 获取余额统计（待提/欠费）
	// 待提：正余额总和
	var balancePay struct{ Total int }
	if err := r.db(ctx).Model(&model.GameAccount{}).
		Where("house_gid = ? AND balance > 0", houseGID).
		Select("COALESCE(SUM(balance), 0) as total").
		Scan(&balancePay).Error; err != nil {
		return nil, err
	}
	stats.BalancePay = balancePay.Total

	// 欠费：负余额总和
	var balanceTake struct{ Total int }
	if err := r.db(ctx).Model(&model.GameAccount{}).
		Where("house_gid = ? AND balance < 0", houseGID).
		Select("COALESCE(SUM(balance), 0) as total").
		Scan(&balanceTake).Error; err != nil {
		return nil, err
	}
	stats.BalanceTake = balanceTake.Total

	// 4. 计算圈子间转账（收圈/补圈）
	balancePayoffs, err := r.calculateGroupBalancePayoffs(ctx, houseGID, start, end)
	if err != nil {
		return nil, err
	}
	stats.BalancePayoffs = balancePayoffs

	// 5. 计算运费分成（收运/补运）
	feePayoffs, err := r.calculateGroupFeePayoffs(ctx, houseGID, start, end)
	if err != nil {
		return nil, err
	}
	stats.FeePayoffs = feePayoffs

	return stats, nil
}

// calculateGroupBalancePayoffs 计算圈子间转账（跨圈对局时的输赢结算）
func (r *battleRecordRepo) calculateGroupBalancePayoffs(ctx context.Context, houseGID int32, start, end *time.Time) ([]GroupPayoff, error) {
	// 获取所有房间
	type RoomInfo struct {
		RoomUID int64 `json:"room_uid"`
	}
	var rooms []RoomInfo

	roomDB := r.db(ctx).Model(&model.GameBattleRecord{}).
		Where("house_gid = ?", houseGID).
		Distinct("room_uid")

	if start != nil {
		roomDB = roomDB.Where("battle_at >= ?", *start)
	}
	if end != nil {
		roomDB = roomDB.Where("battle_at < ?", *end)
	}

	if err := roomDB.Find(&rooms).Error; err != nil {
		return nil, err
	}

	// 按圈子统计每个房间的输赢
	groupBalances := make(map[string]map[string]int) // group -> peer_group -> amount

	for _, room := range rooms {
		type BattleResult struct {
			PlayerID  int32  `json:"player_id"`
			Score     int    `json:"score"`
			PlayGroup string `json:"play_group"`
		}
		var battles []BattleResult

		battleDB := r.db(ctx).Model(&model.GameBattleRecord{}).
			Where("house_gid = ? AND room_uid = ?", houseGID, room.RoomUID).
			Select("player_id, score, play_group")

		if start != nil {
			battleDB = battleDB.Where("battle_at >= ?", *start)
		}
		if end != nil {
			battleDB = battleDB.Where("battle_at < ?", *end)
		}

		if err := battleDB.Find(&battles).Error; err != nil {
			return nil, err
		}

		// 按圈子分组计算输赢
		groupScores := make(map[string]int)
		for _, b := range battles {
			groupScores[b.PlayGroup] += b.Score
		}

		// 计算圈子间转账（简化算法：只记录总输赢）
		for group, score := range groupScores {
			if score == 0 {
				continue
			}

			if groupBalances[group] == nil {
				groupBalances[group] = make(map[string]int)
			}

			// 找出与该圈子有输赢关系的其他圈子
			for otherGroup, otherScore := range groupScores {
				if group == otherGroup {
					continue
				}

				// 如果当前圈赢钱，其他圈输钱
				if score > 0 && otherScore < 0 {
					transferAmount := min(score, -otherScore)
					groupBalances[group][otherGroup] += transferAmount
				}
			}
		}
	}

	// 转换为 GroupPayoff 数组
	var payoffs []GroupPayoff
	for _, peers := range groupBalances {
		for peerGroup, amount := range peers {
			if amount != 0 {
				payoffs = append(payoffs, GroupPayoff{
					Group: peerGroup,
					Value: amount,
				})
			}
		}
	}

	return payoffs, nil
}

// calculateGroupFeePayoffs 计算运费分成
func (r *battleRecordRepo) calculateGroupFeePayoffs(ctx context.Context, houseGID int32, start, end *time.Time) ([]GroupPayoff, error) {
	type FeeResult struct {
		PlayGroup string `json:"play_group"`
		TotalFee  int    `json:"total_fee"`
	}
	var results []FeeResult

	feeDB := r.db(ctx).Model(&model.GameFeeSettle{}).
		Where("house_gid = ?", houseGID).
		Select("play_group, COALESCE(SUM(fee_amount), 0) as total_fee").
		Group("play_group")

	if start != nil {
		feeDB = feeDB.Where("settled_at >= ?", *start)
	}
	if end != nil {
		feeDB = feeDB.Where("settled_at < ?", *end)
	}

	if err := feeDB.Find(&results).Error; err != nil {
		return nil, err
	}

	var payoffs []GroupPayoff
	for _, r := range results {
		if r.TotalFee != 0 {
			payoffs = append(payoffs, GroupPayoff{
				Group: r.PlayGroup,
				Value: r.TotalFee,
			})
		}
	}

	return payoffs, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
