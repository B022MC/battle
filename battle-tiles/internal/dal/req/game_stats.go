package req

// StatsBaseRequest 统计公共入参
// @example {"house_gid":20001}
type StatsBaseRequest struct {
	// 店铺号（圈ID/HouseGID）
	HouseGID int `json:"house_gid" binding:"required"`
	// 圈ID（GroupID，可选；不传则表示按店铺汇总）
	GroupID *int `json:"group_id,omitempty"`
}

// MemberStatsRequest 成员维度统计入参
// @example {"house_gid":20001, "member_id":1001}
type MemberStatsRequest struct {
	// 店铺号（圈ID/HouseGID）
	HouseGID int `json:"house_gid" binding:"required"`
	// 成员ID（平台 member_id）
	MemberID int `json:"member_id" binding:"required"`
}

// BattleDetailRequest 用户战绩明细（按 group+period 拉取，若提供 game_id 则在服务端过滤）
// @example {"house_gid":20001, "group_id":20001, "period":"today", "game_id":123456}
type BattleDetailRequest struct {
	// 店铺号（圈ID/HouseGID）
	HouseGID int `json:"house_gid" binding:"required"`
	// 圈ID（GroupID）
	GroupID int `json:"group_id" binding:"required"`
	// 过滤的用户 gameId（可选）
	GameID *int `json:"game_id"`
	// 查询周期：today|yesterday|thisweek
	Period string `json:"period" binding:"required,oneof=today yesterday thisweek"`
}
