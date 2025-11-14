// internal/dal/model/game/game_player_record.go
package game

import (
	"time"
)

const TableNameGamePlayerRecord = "game_player_record"

// GamePlayerRecord stores per-player battle data for efficient querying by player dimension.
// One row per player per battle.
type GamePlayerRecord struct {
	Id        int32     `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`

	// References
	BattleRecordID int32 `gorm:"column:battle_record_id;not null" json:"battle_record_id"` // FK -> game_battle_record.id
	HouseGID       int32 `gorm:"column:house_gid;not null" json:"house_gid"`               // Redundant for faster filtering

	// Player identifiers
	PlayerGID     int64  `gorm:"column:player_gid;not null" json:"player_gid"`  // Game-level unique player ID (from platform)
	GameAccountID *int32 `gorm:"column:game_account_id" json:"game_account_id"` // Optional mapping to local game_account

	// Snapshot fields for querying without joining game_battle_record
	GroupID    int32     `gorm:"column:group_id;not null" json:"group_id"`
	RoomUID    int32     `gorm:"column:room_uid;not null" json:"room_uid"`
	KindID     int32     `gorm:"column:kind_id;not null" json:"kind_id"`
	BaseScore  int32     `gorm:"column:base_score;not null" json:"base_score"`
	ScoreDelta int32     `gorm:"column:score_delta;not null;default:0" json:"score_delta"`
	IsWinner   bool      `gorm:"column:is_winner;not null;default:false" json:"is_winner"`
	BattleAt   time.Time `gorm:"column:battle_at;type:timestamp with time zone;not null" json:"battle_at"`

	// Raw extra data from platform for future-proofing
	MetaJSON string `gorm:"column:meta_json;type:jsonb;not null;default:'{}'" json:"meta_json"`
}

func (GamePlayerRecord) TableName() string { return TableNameGamePlayerRecord }
