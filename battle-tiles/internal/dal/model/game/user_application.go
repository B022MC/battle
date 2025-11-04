package game

import "time"

const TableNameUserApplication = "game_user_application"

// UserApplication 平台侧用户发起的申请（与游戏端消息分离维护）
type UserApplication struct {
	Id        int32     `gorm:"primaryKey;column:id" json:"id"`
	HouseGID  int32     `gorm:"column:house_gid;not null" json:"house_gid"`
	Applicant int32     `gorm:"column:applicant;not null" json:"applicant"` // 平台用户ID
	Type      int32     `gorm:"column:type;not null" json:"type"`           // 1=admin申请, 2=入圈申请
	AdminUID  int32     `gorm:"column:admin_user_id;not null;default:0" json:"admin_user_id"`
	Note      string    `gorm:"column:note;type:text" json:"note"`
	Status    int32     `gorm:"column:status;not null;default:0" json:"status"` // 0待审,1通过,2拒绝
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (UserApplication) TableName() string { return TableNameUserApplication }
