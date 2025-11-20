// internal/dal/model/game/game_account_group.go
package game

import "time"

const TableNameGameAccountGroup = "game_account_group"

// GameAccountGroup 游戏账号圈子关系表
type GameAccountGroup struct {
	Id                int32     `gorm:"primaryKey;column:id" json:"id"`
	GameAccountID     int32     `gorm:"column:game_account_id;not null" json:"game_account_id"`
	HouseGID          int32     `gorm:"column:house_gid;not null" json:"house_gid"`
	GroupID           int32     `gorm:"column:group_id;not null" json:"group_id"`
	GroupName         string    `gorm:"column:group_name;type:varchar(64);not null;default:''" json:"group_name"`
	AdminUserID       int32     `gorm:"column:admin_user_id;not null" json:"admin_user_id"`
	ApprovedByUserID  int32     `gorm:"column:approved_by_user_id;not null" json:"approved_by_user_id"`
	Status            string    `gorm:"column:status;type:varchar(20);not null;default:'active'" json:"status"`
	JoinedAt          time.Time `gorm:"column:joined_at;type:timestamp with time zone;not null;default:now()" json:"joined_at"`
	CreatedAt         time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
}

func (GameAccountGroup) TableName() string { return TableNameGameAccountGroup }

// Status constants
const (
	AccountGroupStatusActive   = "active"
	AccountGroupStatusInactive = "inactive"
)

