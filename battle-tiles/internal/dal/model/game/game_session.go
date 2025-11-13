// internal/dal/model/game/game_session.go
package game

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

const TableNameGameSession = "game_session"

// Session state constants
const (
	SessionStateActive   = "active"   // Session is active
	SessionStateInactive = "inactive" // Session is inactive
	SessionStateError    = "error"    // Session encountered an error
)

// Sync status constants
const (
	SyncStatusIdle    = "idle"    // Not syncing
	SyncStatusSyncing = "syncing" // Currently syncing
	SyncStatusError   = "error"   // Sync error
)

type GameSession struct {
	Id                int32                 `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt         time.Time             `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt         time.Time             `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
	DeletedAt         time.Time             `gorm:"column:deleted_at;type:timestamp with time zone" json:"deleted_at"`
	GameCtrlAccountID int32                 `gorm:"column:game_ctrl_account_id;not null" json:"game_ctrl_account_id"`
	UserID            int32                 `gorm:"column:user_id;not null" json:"user_id"`
	HouseGID          int32                 `gorm:"column:house_gid;not null" json:"house_gid"`
	State             string                `gorm:"column:state;type:varchar(20);not null;index:idx_session_state" json:"state"`
	DeviceIP          string                `gorm:"column:device_ip;type:varchar(64);not null;default:''" json:"device_ip"`
	ErrorMsg          string                `gorm:"column:error_msg;type:varchar(255);not null;default:''" json:"error_msg"`
	EndAt             *time.Time            `gorm:"column:end_at;type:timestamp with time zone" json:"end_at"`
	IsDel             soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt" json:"is_del"`
	AutoSyncEnabled   bool                  `gorm:"column:auto_sync_enabled;default:true" json:"auto_sync_enabled"`
	LastSyncAt        *time.Time            `gorm:"column:last_sync_at;type:timestamp with time zone" json:"last_sync_at"`
	SyncStatus        string                `gorm:"column:sync_status;type:varchar(20);default:'idle';index:idx_session_sync_status" json:"sync_status"`
	GameAccountID     *int32                `gorm:"column:game_account_id;index:idx_session_game_account" json:"game_account_id"`
}

// IsActive checks if session is active
func (s *GameSession) IsActive() bool {
	return s.State == SessionStateActive
}

// IsSyncing checks if session is currently syncing
func (s *GameSession) IsSyncing() bool {
	return s.SyncStatus == SyncStatusSyncing
}

func (GameSession) TableName() string { return TableNameGameSession }
