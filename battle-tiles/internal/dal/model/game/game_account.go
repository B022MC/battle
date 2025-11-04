// internal/dal/model/game/game_account.go
package game

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

const TableNameGameAccount = "game_account"

type GameAccount struct {
	Id            int32                 `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt     time.Time             `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt     time.Time             `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
	DeletedAt     time.Time             `gorm:"column:deleted_at;type:timestamp with time zone" json:"deleted_at"`
	IsDel         soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt" json:"is_del"`
	UserID        int32                 `gorm:"column:user_id;not null" json:"user_id"`
	Account       string                `gorm:"column:account;type:varchar(64);not null" json:"account"`
	PwdMD5        string                `gorm:"column:pwd_md5;type:varchar(64);not null" json:"pwd_md5"`
	Nickname      string                `gorm:"column:nickname;type:varchar(64);not null;default:''" json:"nickname"`
	IsDefault     bool                  `gorm:"column:is_default;not null;default:false" json:"is_default"`
	Status        int32                 `gorm:"column:status;not null;default:1" json:"status"`
	LastLoginAt   *time.Time            `gorm:"column:last_login_at;type:timestamp with time zone" json:"last_login_at"`
	LoginMode     string                `gorm:"column:login_mode;type:varchar(10);not null;default:'account'" json:"login_mode"`
	CtrlAccountID *int32                `gorm:"column:ctrl_account_id" json:"ctrl_account_id"`
}

func (GameAccount) TableName() string { return TableNameGameAccount }
