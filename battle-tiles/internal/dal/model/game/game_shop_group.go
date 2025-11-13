package game

import "time"

const TableNameGameShopGroup = "game_shop_group"

// GameShopGroup 店铺圈子表（每个店铺管理员对应一个圈子）
type GameShopGroup struct {
	Id          int32     `gorm:"primaryKey;column:id" json:"id"`
	HouseGID    int32     `gorm:"column:house_gid;not null;index:idx_shop_group_house" json:"house_gid"`
	GroupName   string    `gorm:"column:group_name;type:varchar(64);not null" json:"group_name"`
	AdminUserID int32     `gorm:"column:admin_user_id;not null;index:idx_shop_group_admin" json:"admin_user_id"`
	Description string    `gorm:"column:description;type:text;default:''" json:"description"`
	IsActive    bool      `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
}

func (GameShopGroup) TableName() string { return TableNameGameShopGroup }

