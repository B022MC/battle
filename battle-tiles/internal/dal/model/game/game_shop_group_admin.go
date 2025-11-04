package game

import "time"

const TableNameGameShopGroupAdmin = "game_shop_group_admin"

type GameShopGroupAdmin struct {
	Id        int32      `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp with time zone" json:"deleted_at"`
	HouseGID  int32      `gorm:"column:house_gid;not null" json:"house_gid"`
	GroupID   int32      `gorm:"column:group_id;not null" json:"group_id"`
	UserID    int32      `gorm:"column:user_id;not null" json:"user_id"`
	Role      string     `gorm:"column:role;type:varchar(20);not null" json:"role"` // admin | operator
}

func (GameShopGroupAdmin) TableName() string { return TableNameGameShopGroupAdmin }
