// internal/dal/model/game/game_shop_application_log.go
package game

import "time"

const TableNameGameShopApplicationLog = "game_shop_application_log"

// Action constants
const (
	ApplicationActionApproved = "approved"
	ApplicationActionRejected = "rejected"
)

// GameShopApplicationLog 店铺申请操作日志
// 记录管理员对玩家申请的处理操作，申请数据本身存储在 Plaza Session 内存中
type GameShopApplicationLog struct {
	ID           int64     `gorm:"primaryKey;column:id" json:"id"`
	HouseGID     int32     `gorm:"column:house_gid;not null;index:idx_app_log_house" json:"house_gid"`
	ApplierGID   int32     `gorm:"column:applier_gid;not null" json:"applier_gid"`
	ApplierGName string    `gorm:"column:applier_gname;type:varchar(100);not null" json:"applier_gname"`
	Action       string    `gorm:"column:action;type:varchar(20);not null;index:idx_app_log_action" json:"action"` // approved | rejected
	AdminUserID  int32     `gorm:"column:admin_user_id;not null;index:idx_app_log_admin" json:"admin_user_id"`
	AdminGameID  *int32    `gorm:"column:admin_game_id" json:"admin_game_id"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:now();index:idx_app_log_created" json:"created_at"`
}

func (GameShopApplicationLog) TableName() string {
	return TableNameGameShopApplicationLog
}

// IsApproved returns true if the action is approved
func (log *GameShopApplicationLog) IsApproved() bool {
	return log.Action == ApplicationActionApproved
}

// IsRejected returns true if the action is rejected
func (log *GameShopApplicationLog) IsRejected() bool {
	return log.Action == ApplicationActionRejected
}
