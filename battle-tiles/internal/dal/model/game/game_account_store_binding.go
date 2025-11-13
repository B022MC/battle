package game

import "time"

const TableNameGameAccountStoreBinding = "game_account_store_binding"

// Binding status constants
const (
	BindingStatusActive   = "active"   // Active binding
	BindingStatusInactive = "inactive" // Inactive binding
)

// GameAccountStoreBinding represents the binding between a game account and a store
// Business Rule: One game account can only bind to ONE store
type GameAccountStoreBinding struct {
	Id             int32     `gorm:"primaryKey;column:id" json:"id"`
	GameAccountID  int32     `gorm:"column:game_account_id;not null;uniqueIndex:uk_game_account_house" json:"game_account_id"`
	HouseGID       int32     `gorm:"column:house_gid;not null;uniqueIndex:uk_game_account_house;index:idx_gasb_house_gid" json:"house_gid"`
	BoundByUserID  int32     `gorm:"column:bound_by_user_id;not null;index:idx_gasb_bound_by_user" json:"bound_by_user_id"`
	Status         string    `gorm:"column:status;type:varchar(20);not null;default:'active';index:idx_gasb_status" json:"status"`
	CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
}

func (GameAccountStoreBinding) TableName() string { return TableNameGameAccountStoreBinding }

// IsActive checks if the binding is active
func (b *GameAccountStoreBinding) IsActive() bool {
	return b.Status == BindingStatusActive
}

