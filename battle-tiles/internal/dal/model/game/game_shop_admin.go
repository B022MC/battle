// internal/dal/model/game/game_shop_admin.go
package game

import "time"

const TableNameGameShopAdmin = "game_shop_admin"

// Admin role constants
const (
	AdminRoleAdmin    = "admin"    // Full administrator
	AdminRoleOperator = "operator" // Operator with limited permissions
)

type GameShopAdmin struct {
	Id        int32      `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp with time zone" json:"deleted_at"` // 物理删为主；该列存在但不启用gorm软删
	HouseGID       int32      `gorm:"column:house_gid;not null;index:idx_shop_admin_house_gid" json:"house_gid"`
	UserID         int32      `gorm:"column:user_id;not null;index:idx_shop_admin_user_id,idx_shop_admin_exclusive" json:"user_id"`
	Role           string     `gorm:"column:role;type:varchar(20);not null" json:"role"` // admin | operator
	GameAccountID  *int32     `gorm:"column:game_account_id;index:idx_shop_admin_game_account" json:"game_account_id"`
	IsExclusive    bool       `gorm:"column:is_exclusive;default:true;index:idx_shop_admin_exclusive" json:"is_exclusive"`
}

// IsAdmin checks if this is a full administrator
func (sa *GameShopAdmin) IsAdmin() bool {
	return sa.Role == AdminRoleAdmin
}

// IsOperator checks if this is an operator
func (sa *GameShopAdmin) IsOperator() bool {
	return sa.Role == AdminRoleOperator
}

func (GameShopAdmin) TableName() string { return TableNameGameShopAdmin }
