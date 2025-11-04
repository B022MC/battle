// internal/dal/model/game/game_shop_admin.go
package game

import "time"

const TableNameGameShopAdmin = "game_shop_admin"

type GameShopAdmin struct {
	Id        int32      `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp with time zone" json:"deleted_at"` // 物理删为主；该列存在但不启用gorm软删
	HouseGID  int32      `gorm:"column:house_gid;not null" json:"house_gid"`
	UserID    int32      `gorm:"column:user_id;not null" json:"user_id"`
	Role      string     `gorm:"column:role;type:varchar(20);not null" json:"role"` // admin | operator
}

func (GameShopAdmin) TableName() string { return TableNameGameShopAdmin }
