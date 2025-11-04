package req

type CreditDepositRequest struct {
	HouseGID int32  `json:"house_gid"  binding:"required,gt=0"`
	MemberID int32  `json:"member_id"  binding:"required,gt=0"`
	Amount   int32  `json:"amount"     binding:"required,gt=0"` // 上分正数
	BizNo    string `json:"biz_no"     binding:"required"`      // 幂等键
	Reason   string `json:"reason"`                             // 备注
}

type CreditWithdrawRequest struct {
	HouseGID int32  `json:"house_gid"  binding:"required,gt=0"`
	MemberID int32  `json:"member_id"  binding:"required,gt=0"`
	Amount   int32  `json:"amount"     binding:"required,gt=0"` // 下分正数
	BizNo    string `json:"biz_no"     binding:"required"`
	Reason   string `json:"reason"`
}

type CreditForceWithdrawRequest struct {
	HouseGID int32  `json:"house_gid"  binding:"required,gt=0"`
	MemberID int32  `json:"member_id"  binding:"required,gt=0"`
	Amount   int32  `json:"amount"     binding:"required,gt=0"`
	BizNo    string `json:"biz_no"     binding:"required"`
	Reason   string `json:"reason"`
}

type UpdateMemberLimitRequest struct {
	HouseGID int32 `json:"house_gid"  binding:"required,gt=0"`
	MemberID int32 `json:"member_id"  binding:"required,gt=0"`
	// 可选字段：部分更新
	LimitMin *int32 `json:"limit_min,omitempty"` // 例如 -1000
	Forbid   *bool  `json:"forbid,omitempty"`    // true=禁，false=解
	Reason   string `json:"reason"`
}

// 单人钱包
type GetWalletRequest struct {
	HouseGID int32 `json:"house_gid" binding:"required"`
	MemberID int32 `json:"member_id" binding:"required"`
}

// 钱包列表（余额区间 / 个性额度过滤 / 分页）
type ListWalletsRequest struct {
	HouseGID       int32  `json:"house_gid" binding:"required"`
	MinBalance     *int32 `json:"min_balance" binding:"omitempty"`
	MaxBalance     *int32 `json:"max_balance" binding:"omitempty"`
	HasCustomLimit *bool  `json:"has_custom_limit" binding:"omitempty"` // true: limit_min>0; false: limit_min=0; nil: 不筛
	Page           int32  `json:"page" binding:"omitempty"`
	PageSize       int32  `json:"page_size" binding:"omitempty"`
}

// 流水列表（时间范围 / 类型 / 成员 / 分页）
type ListLedgerRequest struct {
	HouseGID int32   `json:"house_gid" binding:"required"`
	MemberID *int32  `json:"member_id" binding:"omitempty"`
	Type     *int32  `json:"type" binding:"omitempty,oneof=1 2 3 4"` // 1上分 2下分 3强制下分 4调整
	StartAt  *string `json:"start_at" binding:"omitempty"`           // RFC3339 或 "2006-01-02"
	EndAt    *string `json:"end_at" binding:"omitempty"`             // 同上
	Page     int32   `json:"page" binding:"omitempty"`
	PageSize int32   `json:"page_size" binding:"omitempty"`
}

// ---- house settings ----
type SetFeesRequest struct {
	HouseGID int32  `json:"house_gid" binding:"required"`
	FeesJSON string `json:"fees_json" binding:"required"` // JSON 字符串
}

type SetShareFeeRequest struct {
	HouseGID int32 `json:"house_gid" binding:"required"`
	Share    bool  `json:"share"`
}

type SetPushCreditRequest struct {
	HouseGID int32 `json:"house_gid" binding:"required"`
	Credit   int32 `json:"credit"`
}
