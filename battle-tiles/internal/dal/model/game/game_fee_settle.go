package game

import "time"

const TableNameGameFeeSettle = "game_fee_settle"

// GameFeeSettle 本地费用结算记录
type GameFeeSettle struct {
	Id        int32     `gorm:"primaryKey;column:id" json:"id"`
	HouseGID  int32     `gorm:"column:house_gid;not null" json:"house_gid"`
	PlayGroup string    `gorm:"column:play_group;type:varchar(32);not null;default:''" json:"play_group"`
	Amount    int32     `gorm:"column:amount;not null" json:"amount"` // 分
	FeedAt    time.Time `gorm:"column:feed_at;type:timestamp with time zone;not null" json:"feed_at"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
}

func (GameFeeSettle) TableName() string { return TableNameGameFeeSettle }
