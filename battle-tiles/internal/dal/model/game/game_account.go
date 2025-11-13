// internal/dal/model/game/game_account.go
package game

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

const TableNameGameAccount = "game_account"

// Verification status constants
const (
	VerificationStatusPending  = "pending"  // Not yet verified
	VerificationStatusVerified = "verified" // Successfully verified
	VerificationStatusFailed   = "failed"   // Verification failed
)

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
	LastLoginAt        *time.Time            `gorm:"column:last_login_at;type:timestamp with time zone" json:"last_login_at"`
	LoginMode          string                `gorm:"column:login_mode;type:varchar(10);not null;default:'account'" json:"login_mode"`
	CtrlAccountID      *int32                `gorm:"column:ctrl_account_id" json:"ctrl_account_id"`
	GameUserID         string                `gorm:"column:game_user_id;type:varchar(32);default:'';index:idx_game_account_game_user_id" json:"game_user_id"`
	VerifiedAt         *time.Time            `gorm:"column:verified_at;type:timestamp with time zone" json:"verified_at"`
	VerificationStatus string                `gorm:"column:verification_status;type:varchar(20);default:'pending';index:idx_game_account_verification" json:"verification_status"`
}

// IsVerified checks if the game account has been verified
func (ga *GameAccount) IsVerified() bool {
	return ga.VerificationStatus == VerificationStatusVerified
}

func (GameAccount) TableName() string { return TableNameGameAccount }
