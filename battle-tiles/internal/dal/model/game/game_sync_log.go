package game

import "time"

const TableNameGameSyncLog = "game_sync_log"

// Sync type constants
const (
	SyncTypeBattleRecord = "battle_record"  // Battle record synchronization
	SyncTypeMemberList   = "member_list"    // Member list synchronization
	SyncTypeWalletUpdate = "wallet_update"  // Wallet update synchronization
	SyncTypeRoomList     = "room_list"      // Room list synchronization
	SyncTypeGroupMember  = "group_member"   // Group member synchronization
)

// Sync status constants
const (
	SyncStatusSuccess = "success" // Sync completed successfully
	SyncStatusFailed  = "failed"  // Sync failed
	SyncStatusPartial = "partial" // Sync partially completed
)

// GameSyncLog tracks synchronization operations and their status
type GameSyncLog struct {
	Id            int32      `gorm:"primaryKey;column:id" json:"id"`
	SessionID     int32      `gorm:"column:session_id;not null;index:idx_sync_log_session_started" json:"session_id"`
	SyncType      string     `gorm:"column:sync_type;type:varchar(20);not null;index:idx_sync_log_type" json:"sync_type"`
	Status        string     `gorm:"column:status;type:varchar(20);not null;index:idx_sync_log_status" json:"status"`
	RecordsSynced int32      `gorm:"column:records_synced;not null;default:0" json:"records_synced"`
	ErrorMessage  string     `gorm:"column:error_message;type:text" json:"error_message"`
	StartedAt     time.Time  `gorm:"column:started_at;type:timestamp with time zone;not null;index:idx_sync_log_session_started" json:"started_at"`
	CompletedAt   *time.Time `gorm:"column:completed_at;type:timestamp with time zone" json:"completed_at"`
}

func (GameSyncLog) TableName() string { return TableNameGameSyncLog }

// IsSuccess checks if sync was successful
func (l *GameSyncLog) IsSuccess() bool {
	return l.Status == SyncStatusSuccess
}

// IsFailed checks if sync failed
func (l *GameSyncLog) IsFailed() bool {
	return l.Status == SyncStatusFailed
}

// Duration returns the duration of the sync operation
func (l *GameSyncLog) Duration() *time.Duration {
	if l.CompletedAt == nil {
		return nil
	}
	d := l.CompletedAt.Sub(l.StartedAt)
	return &d
}

