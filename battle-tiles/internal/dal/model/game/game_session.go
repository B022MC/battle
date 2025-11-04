// internal/dal/model/game/game_session.go
package game

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

const TableNameGameSession = "game_session"

type GameSession struct {
	Id                int32                 `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt         time.Time             `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt         time.Time             `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
	DeletedAt         time.Time             `gorm:"column:deleted_at;type:timestamp with time zone" json:"deleted_at"`
	GameCtrlAccountID int32                 `gorm:"column:game_ctrl_account_id;not null" json:"game_ctrl_account_id"`
	UserID            int32                 `gorm:"column:user_id;not null" json:"user_id"`
	HouseGID          int32                 `gorm:"column:house_gid;not null" json:"house_gid"`
	State             string                `gorm:"column:state;type:varchar(20);not null" json:"state"`
	DeviceIP          string                `gorm:"column:device_ip;type:varchar(64);not null;default:''" json:"device_ip"`
	ErrorMsg          string                `gorm:"column:error_msg;type:varchar(255);not null;default:''" json:"error_msg"`
	EndAt             *time.Time            `gorm:"column:end_at;type:timestamp with time zone" json:"end_at"`
	IsDel             soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt" json:"is_del"`
}

func (GameSession) TableName() string { return TableNameGameSession }
