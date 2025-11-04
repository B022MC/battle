package game

import "time"

const TableNameGameMemberRule = "game_member_rule"

// GameMemberRule 成员规则（VIP/多号/临时解禁）
type GameMemberRule struct {
	Id          int32      `gorm:"primaryKey;column:id" json:"id"`
	HouseGID    int32      `gorm:"column:house_gid;not null" json:"house_gid"`
	MemberID    int32      `gorm:"column:member_id;not null" json:"member_id"`
	VIP         bool       `gorm:"column:vip;not null;default:false" json:"vip"`
	MultiGIDs   bool       `gorm:"column:multi_gids;not null;default:false" json:"multi_gids"`
	TempRelease int32      `gorm:"column:temp_release;not null;default:0" json:"temp_release"` // 临时解禁上限（分），0 为无
	ExpireAt    *time.Time `gorm:"column:expire_at;type:timestamp with time zone" json:"expire_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
	UpdatedBy   int32      `gorm:"column:updated_by;not null;default:0" json:"updated_by"`
}

func (GameMemberRule) TableName() string { return TableNameGameMemberRule }
