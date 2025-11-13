package game

import "time"

const TableNameGameShopGroupMember = "game_shop_group_member"

// GameShopGroupMember 圈子成员关系表（用户可以加入多个圈子）
type GameShopGroupMember struct {
	Id        int32     `gorm:"primaryKey;column:id" json:"id"`
	GroupID   int32     `gorm:"column:group_id;not null;uniqueIndex:uk_group_member_group_user;index:idx_group_member_group" json:"group_id"`
	UserID    int32     `gorm:"column:user_id;not null;uniqueIndex:uk_group_member_group_user;index:idx_group_member_user" json:"user_id"`
	JoinedAt  time.Time `gorm:"column:joined_at;type:timestamp with time zone;not null;default:now()" json:"joined_at"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
}

func (GameShopGroupMember) TableName() string { return TableNameGameShopGroupMember }

