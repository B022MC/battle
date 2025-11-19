// internal/dal/repo/game/game_stats.go
package game

import (
	"battle-tiles/internal/infra"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type GameStatsRepo interface {
	// 资金流水聚合（按时间范围）
	AggregateLedger(ctx context.Context, houseGID int, start, end time.Time) (*LedgerAgg, error)
	// 资金流水聚合（按时间范围 + 成员）
	AggregateLedgerByMember(ctx context.Context, houseGID, memberID int, start, end time.Time) (*LedgerAgg, error)
	// 资金流水聚合（按时间范围 + 成员集合）
	AggregateLedgerByMembers(ctx context.Context, houseGID int, memberIDs []int, start, end time.Time) (*LedgerAgg, error)
	// 钱包现势聚合（不随时间）
	AggregateWallet(ctx context.Context, houseGID int) (*WalletAgg, error)
	// 钱包现势聚合（按成员集合）
	AggregateWalletByMembers(ctx context.Context, houseGID int, memberIDs []int) (*WalletAgg, error)
	// 当前在线会话数
	CountActiveSessions(ctx context.Context, houseGID int) (int64, error)
	// 成员钱包（当前余额）
	GetWalletByMember(ctx context.Context, houseGID, memberID int) (int64, error)
	// 按店铺统计当前在线会话数
	ListActiveSessionsByHouse(ctx context.Context) ([]ActiveByHouse, error)
}

type statsRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewStatsRepo(data *infra.Data, logger log.Logger) GameStatsRepo {
	return &statsRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/stats")),
	}
}

func (r *statsRepo) db(ctx context.Context) *gorm.DB { return r.data.GetDBWithContext(ctx) }

type LedgerAgg struct {
	Income          int64
	Payout          int64
	Adjust          int64
	Net             int64
	Records         int64
	MembersInvolved int64
}

func (r *statsRepo) AggregateLedger(ctx context.Context, houseGID int, start, end time.Time) (*LedgerAgg, error) {
	// 从 game_battle_record 表查询统计数据
	// income: 总分数（score）
	// payout: 总费用（fee）
	// adjust: 0（暂无调整数据）
	// net: 总分数 - 总费用
	// records: 战绩总数
	// members_involved: 参与的玩家数
	raw := `
SELECT
	COALESCE(SUM(score), 0)                                                                                           AS income,
	COALESCE(SUM(fee), 0)                                                                                             AS payout,
	0                                                                                                                 AS adjust,
	COALESCE(SUM(score) - SUM(fee), 0)                                                                                AS net,
	COUNT(*)                                                                                                          AS records,
	COUNT(DISTINCT player_id)                                                                                         AS members_involved
FROM game_battle_record
WHERE house_gid = ? AND battle_at >= ? AND battle_at < ?;
`
	var row LedgerAgg
	if err := r.db(ctx).Raw(raw, houseGID, start, end).Scan(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *statsRepo) AggregateLedgerByMember(ctx context.Context, houseGID, memberID int, start, end time.Time) (*LedgerAgg, error) {
	raw := `
SELECT
    COALESCE(SUM(score), 0)                                                                                           AS income,
    COALESCE(SUM(fee), 0)                                                                                             AS payout,
    0                                                                                                                 AS adjust,
    COALESCE(SUM(score) - SUM(fee), 0)                                                                                AS net,
    COUNT(*)                                                                                                          AS records,
    COUNT(DISTINCT player_id)                                                                                         AS members_involved
FROM game_battle_record
WHERE house_gid = ? AND player_id = ? AND battle_at >= ? AND battle_at < ?;
`
	var row LedgerAgg
	if err := r.db(ctx).Raw(raw, houseGID, memberID, start, end).Scan(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *statsRepo) AggregateLedgerByMembers(ctx context.Context, houseGID int, memberIDs []int, start, end time.Time) (*LedgerAgg, error) {
	if len(memberIDs) == 0 {
		return &LedgerAgg{}, nil
	}
	raw := `
SELECT
    COALESCE(SUM(score), 0)                                                                                           AS income,
    COALESCE(SUM(fee), 0)                                                                                             AS payout,
    0                                                                                                                 AS adjust,
    COALESCE(SUM(score) - SUM(fee), 0)                                                                                AS net,
    COUNT(*)                                                                                                          AS records,
    COUNT(DISTINCT player_id)                                                                                         AS members_involved
FROM game_battle_record
WHERE house_gid = ? AND player_id IN (?) AND battle_at >= ? AND battle_at < ?;
`
	var row LedgerAgg
	if err := r.db(ctx).Raw(raw, houseGID, memberIDs, start, end).Scan(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

type WalletAgg struct {
	BalanceTotal      int64
	Members           int64
	LowBalanceMembers int64
}

func (r *statsRepo) AggregateWallet(ctx context.Context, houseGID int) (*WalletAgg, error) {
	raw := `
SELECT
	COALESCE(SUM(balance), 0)                                                                                           AS balance_total,
	COUNT(*)                                                                                                            AS members,
	COALESCE(SUM(CASE WHEN balance < limit_min THEN 1 ELSE 0 END), 0)                                                  AS low_balance_members
FROM game_member_wallet
WHERE house_gid = ?;
`
	var row WalletAgg
	if err := r.db(ctx).Raw(raw, houseGID).Scan(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *statsRepo) AggregateWalletByMembers(ctx context.Context, houseGID int, memberIDs []int) (*WalletAgg, error) {
	if len(memberIDs) == 0 {
		return &WalletAgg{}, nil
	}
	raw := `
SELECT
    COALESCE(SUM(balance), 0)                                                                                           AS balance_total,
    COUNT(*)                                                                                                            AS members,
    COALESCE(SUM(CASE WHEN balance < limit_min THEN 1 ELSE 0 END), 0)                                                  AS low_balance_members
FROM game_member_wallet
WHERE house_gid = ? AND member_id IN (?);
`
	var row WalletAgg
	if err := r.db(ctx).Raw(raw, houseGID, memberIDs).Scan(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *statsRepo) CountActiveSessions(ctx context.Context, houseGID int) (int64, error) {
	var cnt int64
	err := r.db(ctx).
		Table("game_session").
		Where("house_gid = ?  AND end_at IS NULL", houseGID).
		Count(&cnt).Error
	return cnt, err
}

func (r *statsRepo) GetWalletByMember(ctx context.Context, houseGID, memberID int) (int64, error) {
	var bal int64
	err := r.db(ctx).Table("game_member_wallet").
		Select("COALESCE(balance,0)").
		Where("house_gid = ? AND member_id = ?", houseGID, memberID).
		Scan(&bal).Error
	return bal, err
}

type ActiveByHouse struct {
	HouseGID int   `gorm:"column:house_gid" json:"house_gid"`
	Active   int64 `gorm:"column:active" json:"active"`
}

func (r *statsRepo) ListActiveSessionsByHouse(ctx context.Context) ([]ActiveByHouse, error) {
	var rows []ActiveByHouse
	err := r.db(ctx).
		Table("game_session").
		Select("house_gid, COUNT(*) AS active").
		Where("end_at IS NULL").
		Group("house_gid").
		Order("house_gid").
		Scan(&rows).Error
	return rows, err
}
