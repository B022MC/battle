// internal/biz/game/game_stats.go
package game

import (
	repo "battle-tiles/internal/dal/repo/game"
	"battle-tiles/internal/dal/resp"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

type GameStatsUseCase struct {
	repo repo.GameStatsRepo
	log  *log.Helper
}

func NewGameStatsUseCase(r repo.GameStatsRepo, logger log.Logger) *GameStatsUseCase {
	return &GameStatsUseCase{
		repo: r,
		log:  log.NewHelper(log.With(logger, "module", "usecase/stats")),
	}
}

func (uc *GameStatsUseCase) Today(ctx context.Context, houseGID int) (*resp.StatsVO, error) {
	if houseGID <= 0 {
		return nil, errors.New("invalid house_gid")
	}
	start, end := dayRange(time.Now())
	return uc.aggregate(ctx, houseGID, start, end)
}

// Yesterday 昨日统计 [yesterday 00:00, today 00:00)
func (uc *GameStatsUseCase) Yesterday(ctx context.Context, houseGID int) (*resp.StatsVO, error) {
	if houseGID <= 0 {
		return nil, errors.New("invalid house_gid")
	}
	todayStart, tomorrowStart := dayRange(time.Now())
	start := todayStart.AddDate(0, 0, -1)
	end := todayStart
	// 保留 end 作为上限（不含），与 Today/Week 一致
	_ = tomorrowStart
	return uc.aggregate(ctx, houseGID, start, end)
}

// LastWeek 上周统计 [lastMonday 00:00, thisMonday 00:00)
func (uc *GameStatsUseCase) LastWeek(ctx context.Context, houseGID int) (*resp.StatsVO, error) {
	if houseGID <= 0 {
		return nil, errors.New("invalid house_gid")
	}
	now := time.Now()
	// 计算本周一与上周一
	y, m, d := now.Date()
	loc := now.Location()
	todayStart := time.Date(y, m, d, 0, 0, 0, 0, loc)
	weekday := int(todayStart.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	thisMonday := todayStart.AddDate(0, 0, 1-weekday)
	lastMonday := thisMonday.AddDate(0, 0, -7)
	return uc.aggregate(ctx, houseGID, lastMonday, thisMonday)
}

func (uc *GameStatsUseCase) Week(ctx context.Context, houseGID int) (*resp.StatsVO, error) {
	if houseGID <= 0 {
		return nil, errors.New("invalid house_gid")
	}
	// 近7天（含今日） => [todayStart-6d, tomorrowStart)
	todayStart, tomorrowStart := dayRange(time.Now())
	start := todayStart.AddDate(0, 0, -6)
	end := tomorrowStart
	return uc.aggregate(ctx, houseGID, start, end)
}

// MemberToday/Yesterday/ThisWeek: 简单复用 ledger 聚合并过滤 member_id
func (uc *GameStatsUseCase) MemberToday(ctx context.Context, houseGID, memberID int) (*resp.StatsVO, error) {
	if houseGID <= 0 || memberID <= 0 {
		return nil, errors.New("invalid args")
	}
	start, end := dayRange(time.Now())
	ledger, err := uc.repo.AggregateLedgerByMember(ctx, houseGID, memberID, start, end)
	if err != nil {
		return nil, err
	}
	walletTotal, err := uc.repo.GetWalletByMember(ctx, houseGID, memberID)
	if err != nil {
		return nil, err
	}
	out := &resp.StatsVO{
		HouseGID:   houseGID,
		RangeStart: start,
		RangeEnd:   end,
		Ledger: resp.StatsLedgerVO{
			Income:          ledger.Income,
			Payout:          ledger.Payout,
			Adjust:          ledger.Adjust,
			Net:             ledger.Net,
			Records:         ledger.Records,
			MembersInvolved: ledger.MembersInvolved,
		},
		Wallet: resp.StatsWalletVO{BalanceTotal: walletTotal, Members: 1, LowBalanceMembers: 0},
		Sess:   resp.StatsSessVO{},
	}
	return out, nil
}

// MemberYesterday 昨日
func (uc *GameStatsUseCase) MemberYesterday(ctx context.Context, houseGID, memberID int) (*resp.StatsVO, error) {
	if houseGID <= 0 || memberID <= 0 {
		return nil, errors.New("invalid args")
	}
	todayStart, _ := dayRange(time.Now())
	start := todayStart.AddDate(0, 0, -1)
	end := todayStart
	ledger, err := uc.repo.AggregateLedgerByMember(ctx, houseGID, memberID, start, end)
	if err != nil {
		return nil, err
	}
	walletTotal, err := uc.repo.GetWalletByMember(ctx, houseGID, memberID)
	if err != nil {
		return nil, err
	}
	out := &resp.StatsVO{
		HouseGID:   houseGID,
		RangeStart: start,
		RangeEnd:   end,
		Ledger: resp.StatsLedgerVO{
			Income:          ledger.Income,
			Payout:          ledger.Payout,
			Adjust:          ledger.Adjust,
			Net:             ledger.Net,
			Records:         ledger.Records,
			MembersInvolved: ledger.MembersInvolved,
		},
		Wallet: resp.StatsWalletVO{BalanceTotal: walletTotal, Members: 1, LowBalanceMembers: 0},
		Sess:   resp.StatsSessVO{},
	}
	return out, nil
}

// MemberThisWeek 本周（近7日含今天）
func (uc *GameStatsUseCase) MemberThisWeek(ctx context.Context, houseGID, memberID int) (*resp.StatsVO, error) {
	if houseGID <= 0 || memberID <= 0 {
		return nil, errors.New("invalid args")
	}
	todayStart, tomorrowStart := dayRange(time.Now())
	start := todayStart.AddDate(0, 0, -6)
	end := tomorrowStart
	ledger, err := uc.repo.AggregateLedgerByMember(ctx, houseGID, memberID, start, end)
	if err != nil {
		return nil, err
	}
	walletTotal, err := uc.repo.GetWalletByMember(ctx, houseGID, memberID)
	if err != nil {
		return nil, err
	}
	out := &resp.StatsVO{
		HouseGID:   houseGID,
		RangeStart: start,
		RangeEnd:   end,
		Ledger: resp.StatsLedgerVO{
			Income:          ledger.Income,
			Payout:          ledger.Payout,
			Adjust:          ledger.Adjust,
			Net:             ledger.Net,
			Records:         ledger.Records,
			MembersInvolved: ledger.MembersInvolved,
		},
		Wallet: resp.StatsWalletVO{BalanceTotal: walletTotal, Members: 1, LowBalanceMembers: 0},
		Sess:   resp.StatsSessVO{},
	}
	return out, nil
}

// MemberLastWeek 上周（上周一至本周一）
func (uc *GameStatsUseCase) MemberLastWeek(ctx context.Context, houseGID, memberID int) (*resp.StatsVO, error) {
	if houseGID <= 0 || memberID <= 0 {
		return nil, errors.New("invalid args")
	}
	now := time.Now()
	y, m, d := now.Date()
	loc := now.Location()
	todayStart := time.Date(y, m, d, 0, 0, 0, 0, loc)
	weekday := int(todayStart.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	thisMonday := todayStart.AddDate(0, 0, 1-weekday)
	lastMonday := thisMonday.AddDate(0, 0, -7)
	ledger, err := uc.repo.AggregateLedgerByMember(ctx, houseGID, memberID, lastMonday, thisMonday)
	if err != nil {
		return nil, err
	}
	walletTotal, err := uc.repo.GetWalletByMember(ctx, houseGID, memberID)
	if err != nil {
		return nil, err
	}
	out := &resp.StatsVO{
		HouseGID:   houseGID,
		RangeStart: lastMonday,
		RangeEnd:   thisMonday,
		Ledger: resp.StatsLedgerVO{
			Income:          ledger.Income,
			Payout:          ledger.Payout,
			Adjust:          ledger.Adjust,
			Net:             ledger.Net,
			Records:         ledger.Records,
			MembersInvolved: ledger.MembersInvolved,
		},
		Wallet: resp.StatsWalletVO{BalanceTotal: walletTotal, Members: 1, LowBalanceMembers: 0},
		Sess:   resp.StatsSessVO{},
	}
	return out, nil
}

func (uc *GameStatsUseCase) aggregate(ctx context.Context, houseGID int, start, end time.Time) (*resp.StatsVO, error) {
	ledger, err := uc.repo.AggregateLedger(ctx, houseGID, start, end)
	if err != nil {
		return nil, err
	}
	wallet, err := uc.repo.AggregateWallet(ctx, houseGID)
	if err != nil {
		return nil, err
	}
	active, err := uc.repo.CountActiveSessions(ctx, houseGID)
	if err != nil {
		return nil, err
	}

	out := &resp.StatsVO{
		HouseGID:   houseGID,
		RangeStart: start,
		RangeEnd:   end,
		Ledger: resp.StatsLedgerVO{
			Income:          ledger.Income,
			Payout:          ledger.Payout,
			Adjust:          ledger.Adjust,
			Net:             ledger.Net,
			Records:         ledger.Records,
			MembersInvolved: ledger.MembersInvolved,
		},
		Wallet: resp.StatsWalletVO{
			BalanceTotal:      wallet.BalanceTotal,
			Members:           wallet.Members,
			LowBalanceMembers: wallet.LowBalanceMembers,
		},
		Sess: resp.StatsSessVO{
			Active: active,
		},
	}
	return out, nil
}

// aggregateByMembers 用成员集合限制统计范围（用于按圈统计，不改动表结构）
func (uc *GameStatsUseCase) aggregateByMembers(ctx context.Context, houseGID int, memberIDs []int, start, end time.Time) (*resp.StatsVO, error) {
	ledger, err := uc.repo.AggregateLedgerByMembers(ctx, houseGID, memberIDs, start, end)
	if err != nil {
		return nil, err
	}
	wallet, err := uc.repo.AggregateWalletByMembers(ctx, houseGID, memberIDs)
	if err != nil {
		return nil, err
	}
	active, err := uc.repo.CountActiveSessions(ctx, houseGID)
	if err != nil {
		return nil, err
	}
	out := &resp.StatsVO{
		HouseGID:   houseGID,
		RangeStart: start,
		RangeEnd:   end,
		Ledger: resp.StatsLedgerVO{
			Income:          ledger.Income,
			Payout:          ledger.Payout,
			Adjust:          ledger.Adjust,
			Net:             ledger.Net,
			Records:         ledger.Records,
			MembersInvolved: ledger.MembersInvolved,
		},
		Wallet: resp.StatsWalletVO{
			BalanceTotal:      wallet.BalanceTotal,
			Members:           wallet.Members,
			LowBalanceMembers: wallet.LowBalanceMembers,
		},
		Sess: resp.StatsSessVO{Active: active},
	}
	return out, nil
}

// AggregateByMembers 对外暴露，供 service 使用
func (uc *GameStatsUseCase) AggregateByMembers(ctx context.Context, houseGID int, memberIDs []int, start, end time.Time) (*resp.StatsVO, error) {
	return uc.aggregateByMembers(ctx, houseGID, memberIDs, start, end)
}

func (uc *GameStatsUseCase) ListActiveByHouse(ctx context.Context) ([]resp.ActiveByHouseVO, error) {
	rows, err := uc.repo.ListActiveSessionsByHouse(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]resp.ActiveByHouseVO, 0, len(rows))
	for _, r := range rows {
		out = append(out, resp.ActiveByHouseVO{HouseGID: r.HouseGID, Active: r.Active})
	}
	return out, nil
}
