package resp

import "time"

// WalletVO 钱包信息
// @example {"house_gid":20001, "member_id":1001, "balance":5000, "forbid":false, "limit_min":0}
type WalletVO struct {
	HouseGID int32 `json:"house_gid"`
	MemberID int32 `json:"member_id"`
	// 余额（分）
	Balance int32 `json:"balance"`
	// 是否冻结
	Forbid bool `json:"forbid"`
	// 个性额度下限（分）
	LimitMin  int32     `json:"limit_min"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy int32     `json:"updated_by"`
}

// LedgerVO 资金流水
// @example {"id":1,"house_gid":20001,"member_id":1001,"change_amount":100,"type":1,"reason":"上分","biz_no":"A1","created_at":"2025-10-10T00:00:00Z"}
type LedgerVO struct {
	ID       int32 `json:"id"`
	HouseGID int32 `json:"house_gid"`
	MemberID int32 `json:"member_id"`
	// 变动额（分），下分为负数
	ChangeAmount  int32 `json:"change_amount"`
	BalanceBefore int32 `json:"balance_before"`
	BalanceAfter  int32 `json:"balance_after"`
	// 1上分 2下分 3强制下分 4调整
	Type           int32     `json:"type"`
	Reason         string    `json:"reason"`
	OperatorUserID int32     `json:"operator_user_id"`
	BizNo          string    `json:"biz_no"`
	CreatedAt      time.Time `json:"created_at"`
}

// HouseSettingsVO 店铺设置
// @example {"house_gid":20001, "fees_json":"[]", "share_fee":true, "push_credit":1000}
type HouseSettingsVO struct {
	HouseGID int32 `json:"house_gid"`
	// 运费配置（原始JSON）
	FeesJSON string `json:"fees_json"`
	// 分运开关
	ShareFee bool `json:"share_fee"`
	// 推送额度（分）
	PushCredit int32 `json:"push_credit"`
}
