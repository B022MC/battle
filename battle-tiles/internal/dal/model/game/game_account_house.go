package game

import "time"

const TableNameGameAccountHouse = "game_account_house"

type GameAccountHouse struct {
	Id            int32     `gorm:"primaryKey;column:id" json:"id"`
	GameAccountID int32     `gorm:"column:game_account_id;not null" json:"game_account_id"`
	HouseGID      int32     `gorm:"column:house_gid;not null" json:"house_gid"`
	IsDefault     bool      `gorm:"column:is_default;not null;default:false" json:"is_default"`
	Status        int32     `gorm:"column:status;not null;default:1" json:"status"`
	CreatedAt     time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
}

func (GameAccountHouse) TableName() string { return TableNameGameAccountHouse }
