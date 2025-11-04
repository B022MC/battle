// internal/dal/resp/stats.go
package resp

import "time"

// StatsVO 统计结果
// @description 统计时间范围为 [range_start, range_end)
type StatsVO struct {
	// 店铺号
	HouseGID int `json:"house_gid"`
	// 起始时间（含）
	RangeStart time.Time `json:"range_start"`
	// 结束时间（不含）
	RangeEnd time.Time `json:"range_end"`

	Ledger StatsLedgerVO `json:"ledger"`
	Wallet StatsWalletVO `json:"wallet"`
	Sess   StatsSessVO   `json:"session"`
}

// StatsLedgerVO 账本聚合
// @example {"income":1000,"payout":800,"adjust":0,"net":200,"records":10,"members_involved":3}
type StatsLedgerVO struct {
	// 上分总额（正数）
	Income int64 `json:"income"`
	// 下分总额（正数）
	Payout int64 `json:"payout"`
	// 调整（可正可负）
	Adjust int64 `json:"adjust"`
	// 净变动
	Net int64 `json:"net"`
	// 流水条数
	Records int64 `json:"records"`
	// 参与流水成员数
	MembersInvolved int64 `json:"members_involved"`
}

// StatsWalletVO 钱包聚合
// @example {"balance_total":10000,"members":50,"low_balance_members":2}
type StatsWalletVO struct {
	BalanceTotal      int64 `json:"balance_total"`
	Members           int64 `json:"members"`
	LowBalanceMembers int64 `json:"low_balance_members"`
}

// StatsSessVO 会话统计
// @example {"active":3}
type StatsSessVO struct {
	// 当前在线会话数
	Active int64 `json:"active"`
}

// MemberStatsVO 成员维度统计
// @example {"house_gid":20001,"member_id":1001,"range_start":"2025-10-10T00:00:00Z","range_end":"2025-10-11T00:00:00Z","ledger":{"income":100},"balance":500}
type MemberStatsVO struct {
	HouseGID   int       `json:"house_gid"`
	MemberID   int       `json:"member_id"`
	RangeStart time.Time `json:"range_start"`
	RangeEnd   time.Time `json:"range_end"`

	Ledger StatsLedgerVO `json:"ledger"`
	// 当前余额
	Balance int64 `json:"balance"`
}

// ActiveByHouseVO 按店铺的会话活跃数
type ActiveByHouseVO struct {
	HouseGID int   `json:"house_gid"`
	Active   int64 `json:"active"`
}
