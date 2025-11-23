package game

import "time"

const TableNameGameRoomCreditLimit = "game_room_credit_limit"

// GameRoomCreditLimit 房间额度限制表
// 用于设置不同游戏类型、底分、圈子对应的进入房间所需的最低余额
// 例如："额度100/1/红中" 表示进入1分底的红中房间需要至少100分余额
type GameRoomCreditLimit struct {
	Id          int32     `gorm:"primaryKey;column:id" json:"id"`
	HouseGID    int32     `gorm:"column:house_gid;not null;uniqueIndex:uk_room_credit_limit;index:idx_room_credit_house" json:"house_gid"`
	GroupName   string    `gorm:"column:group_name;type:varchar(64);not null;default:'';uniqueIndex:uk_room_credit_limit" json:"group_name"`
	GameKind    int32     `gorm:"column:game_kind;not null;default:0;uniqueIndex:uk_room_credit_limit" json:"game_kind"`
	BaseScore   int32     `gorm:"column:base_score;not null;default:0;uniqueIndex:uk_room_credit_limit" json:"base_score"`
	CreditLimit int32     `gorm:"column:credit_limit;not null;default:0" json:"credit_limit"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
	UpdatedBy   int32     `gorm:"column:updated_by;not null;default:0" json:"updated_by"`
}

func (GameRoomCreditLimit) TableName() string {
	return TableNameGameRoomCreditLimit
}

// IsGlobal 判断是否为全局默认设置（0-0）
func (l *GameRoomCreditLimit) IsGlobal() bool {
	return l.GroupName == "" && l.GameKind == 0 && l.BaseScore == 0
}

// IsGroupDefault 判断是否为圈子默认设置（group-0-0）
func (l *GameRoomCreditLimit) IsGroupDefault() bool {
	return l.GroupName != "" && l.GameKind == 0 && l.BaseScore == 0
}

// IsSpecific 判断是否为特定游戏类型和底分的设置
func (l *GameRoomCreditLimit) IsSpecific() bool {
	return l.GameKind > 0 && l.BaseScore > 0
}

// GetCreditInYuan 获取额度（元）
func (l *GameRoomCreditLimit) GetCreditInYuan() float64 {
	return float64(l.CreditLimit) / 100.0
}
