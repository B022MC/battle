package game

import "time"

const TableNameGameMember = "game_member"

// GameMember represents a player/member within a store and group

type GameMember struct {
	Id           int32     `gorm:"primaryKey;column:id" json:"id"`
	HouseGID     int32     `gorm:"column:house_gid;not null;uniqueIndex:uk_game_member_house_game;index:idx_game_member_house_group" json:"house_gid"`
	GameID       int32     `gorm:"column:game_id;not null;uniqueIndex:uk_game_member_house_game;index:idx_game_member_game_id" json:"game_id"`
	GameName     string    `gorm:"column:game_name;type:varchar(64);not null;default:''" json:"game_name"`
	GroupID      *int32    `gorm:"column:group_id;index:idx_game_member_group_id" json:"group_id"` // 圈子ID,关联 game_shop_group.id（一个玩家在一个店铺只能属于一个圈子）
	GroupName    string    `gorm:"column:group_name;type:varchar(64);not null;default:'';index:idx_game_member_house_group" json:"group_name"`
	Balance      int32     `gorm:"column:balance;not null;default:0;index:idx_game_member_balance" json:"balance"`
	Credit       int32     `gorm:"column:credit;not null;default:0" json:"credit"`
	Forbid       bool      `gorm:"column:forbid;not null;default:false;index:idx_game_member_forbid" json:"forbid"`
	Recommender  string    `gorm:"column:recommender;type:varchar(64);default:''" json:"recommender"`
	UseMultiGids bool      `gorm:"column:use_multi_gids;default:false" json:"use_multi_gids"`
	ActiveGid    *int32    `gorm:"column:active_gid" json:"active_gid"`
	IsPinned     bool      `gorm:"column:is_pinned;default:false;index:idx_game_member_pinned" json:"is_pinned"` // 是否置顶
	PinOrder     int32     `gorm:"column:pin_order;default:0;index:idx_game_member_pin_order" json:"pin_order"`  // 置顶排序（数字越小越靠前）
	Remark       string    `gorm:"column:remark;type:varchar(255);default:''" json:"remark"`                     // 管理员备注
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
}

func (GameMember) TableName() string { return TableNameGameMember }

// IsForbidden checks if the member is forbidden/banned
func (m *GameMember) IsForbidden() bool {
	return m.Forbid
}

// GetBalanceInYuan returns balance in yuan (balance is stored in cents)
func (m *GameMember) GetBalanceInYuan() float64 {
	return float64(m.Balance) / 100.0
}
