package resp

import (
	"fmt"
	"time"
)

// RoomCreditLimitItem 房间额度限制项
type RoomCreditLimitItem struct {
	Id           int32     `json:"id"`
	HouseGID     int32     `json:"house_gid"`
	GroupName    string    `json:"group_name"`
	GameKind     int32     `json:"game_kind"`
	GameKindName string    `json:"game_kind_name,omitempty"` // 游戏类型名称（如"红中"）
	BaseScore    int32     `json:"base_score"`
	CreditLimit  int32     `json:"credit_limit"` // 单位：分
	CreditYuan   float64   `json:"credit_yuan"`  // 单位：元
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    int32     `json:"updated_by"`
}

// RoomCreditLimitListResponse 房间额度限制列表响应
type RoomCreditLimitListResponse struct {
	Total int32                  `json:"total"`
	Items []*RoomCreditLimitItem `json:"items"`
}

// RoomCreditLimitResponse 房间额度限制单条响应
type RoomCreditLimitResponse struct {
	*RoomCreditLimitItem
}

// CheckPlayerCreditResponse 检查玩家额度响应
type CheckPlayerCreditResponse struct {
	CanEnter        bool    `json:"can_enter"`        // 是否可以进入
	PlayerBalance   int32   `json:"player_balance"`   // 玩家余额（分）
	RequiredCredit  int32   `json:"required_credit"`  // 需要的额度（分）
	PlayerCredit    int32   `json:"player_credit"`    // 玩家个人额度调整（分）
	EffectiveCredit int32   `json:"effective_credit"` // 有效额度要求（分）
	BalanceYuan     float64 `json:"balance_yuan"`     // 余额（元）
	RequiredYuan    float64 `json:"required_yuan"`    // 需要的额度（元）
}

// FormatCreditDisplay 格式化额度显示（例如："🈲 100/红中/5" 或 "🈲 100"）
func FormatCreditDisplay(creditLimit int32, gameKindName string, baseScore int32) string {
	creditYuan := float64(creditLimit) / 100.0
	if gameKindName != "" && baseScore > 0 {
		return fmt.Sprintf("🈲 %.0f/%s/%d", creditYuan, gameKindName, baseScore)
	} else if gameKindName == "" && baseScore == 0 {
		return fmt.Sprintf("🈲 %.0f", creditYuan)
	}
	return fmt.Sprintf("🈲 %.0f", creditYuan)
}
