package game

import (
	"time"
)

const TableNameGameCtrlAccount = "game_ctrl_account"

type GameCtrlAccount struct {
	Id           int32      `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt    time.Time  `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
	DeletedAt    time.Time  `gorm:"column:deleted_at;type:timestamp with time zone" json:"deleted_at"`
	LoginMode    int32      `gorm:"column:login_mode;not null" json:"login_mode"` // smallint
	Identifier   string     `gorm:"column:identifier;type:varchar(64);not null" json:"identifier"`
	PwdMD5       string     `gorm:"column:pwd_md5;type:varchar(64);not null" json:"pwd_md5"`
	GameUserID   string     `gorm:"column:game_player_id;type:varchar(32);not null;default:''" json:"game_player_id"` // 改为 game_player_id
	GameID       string     `gorm:"column:game_id;type:varchar(32);not null;default:''" json:"game_id"`
	Status       int32      `gorm:"column:status;not null;default:1" json:"status"`
	LastVerifyAt *time.Time `gorm:"column:last_verify_at;type:timestamp with time zone" json:"last_verify_at"`
	//IsDel        soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt" json:"is_del"`
}

func (GameCtrlAccount) TableName() string { return TableNameGameCtrlAccount }
