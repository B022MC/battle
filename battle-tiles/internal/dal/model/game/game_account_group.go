// internal/dal/model/game/game_account_group.go
package game

import "time"

const TableNameGameAccountGroup = "game_account_group"

// GameAccountGroup 游戏玩家圈子关系表（用于战绩同步和计费）
type GameAccountGroup struct {
	Id            int32     `gorm:"primaryKey;column:id" json:"id"`
	GameAccountID int32     `gorm:"column:game_account_id;not null" json:"game_account_id"`                // 游戏账号ID（外键）
	GamePlayerID  string    `gorm:"column:game_player_id;type:varchar(32);not null" json:"game_player_id"` // 游戏内玩家ID
	HouseGID      int32     `gorm:"column:house_gid;not null" json:"house_gid"`
	GroupID       int32     `gorm:"column:group_id;not null" json:"group_id"`
	GroupName     string    `gorm:"column:group_name;type:varchar(64);not null;default:''" json:"group_name"`
	AdminUserID   *int32    `gorm:"column:admin_user_id" json:"admin_user_id,omitempty"` // 可选，圈主用户ID
	Status        string    `gorm:"column:status;type:varchar(20);not null;default:'active'" json:"status"`
	JoinedAt      time.Time `gorm:"column:joined_at;type:timestamp with time zone;not null;default:now()" json:"joined_at"`
	CreatedAt     time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
}

func (GameAccountGroup) TableName() string { return TableNameGameAccountGroup }

// Status constants
const (
	AccountGroupStatusActive   = "active"
	AccountGroupStatusInactive = "inactive"
)
