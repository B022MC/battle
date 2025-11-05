package req

// InsertFeeSettleRequest 新增一笔费用结算记录
// @example {"house_gid":20001, "play_group":"A", "amount":300, "feed_at":"2025-10-10T12:00:00Z"}
type InsertFeeSettleRequest struct {
	// 店铺号
	HouseGID int32 `json:"house_gid"  binding:"required"`
	// 圈名/分组
	PlayGroup string `json:"play_group" binding:"required"`
	// 金额（分）
	Amount int32 `json:"amount"     binding:"required"`
	// 结算时间（RFC3339）
	FeedAt string `json:"feed_at"    binding:"required"`
}

// SumFeeSettleRequest 汇总某时间范围的费用
// @example {"house_gid":20001, "play_group":"A", "start_at":"2025-10-01T00:00:00Z", "end_at":"2025-10-08T00:00:00Z"}
type SumFeeSettleRequest struct {
	HouseGID  int32  `json:"house_gid"  binding:"required"`
	PlayGroup string `json:"play_group" binding:"required"`
	StartAt   string `json:"start_at"   binding:"required"`
	EndAt     string `json:"end_at"     binding:"required"`
}

// ListGroupPayoffsRequest 汇总时间范围内各圈费用并计算圈间结转
// @example {"house_gid":20001, "start_at":"2025-10-01T00:00:00Z", "end_at":"2025-10-08T00:00:00Z"}
type ListGroupPayoffsRequest struct {
	HouseGID int32  `json:"house_gid"  binding:"required"`
	StartAt  string `json:"start_at"   binding:"required"`
	EndAt    string `json:"end_at"     binding:"required"`
}
