package game

import "time"

const TableNameGameMemberWallet = "game_member_wallet"

// GameMemberWallet 游戏成员钱包表
type GameMemberWallet struct {
	Id        int32     `gorm:"primaryKey;column:id" json:"id"`
	HouseGID  int32     `gorm:"column:house_gid;not null;uniqueIndex:uk_game_member_wallet_house_member_group" json:"house_gid"`
	GameID    int32     `gorm:"column:game_id;not null;default:0;index:idx_game_member_wallet_game;uniqueIndex:uk_game_member_wallet_house_game_group" json:"game_id"` // 游戏ID，便于按 game_id 查询
	MemberID  int32     `gorm:"column:member_id;not null;uniqueIndex:uk_game_member_wallet_house_member_group" json:"member_id"`
	GroupID   *int32    `gorm:"column:group_id;uniqueIndex:uk_game_member_wallet_house_member_group;uniqueIndex:uk_game_member_wallet_house_game_group;index:idx_game_member_wallet_group_id" json:"group_id"` // 圈子ID,关联 game_shop_group.id
	Balance   int32     `gorm:"column:balance;not null;default:0" json:"balance"`
	Forbid    bool      `gorm:"column:forbid;not null;default:false" json:"forbid"`
	LimitMin  int32     `gorm:"column:limit_min;not null;default:0" json:"limit_min"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
	UpdatedBy int32     `gorm:"column:updated_by;not null;default:0" json:"updated_by"`
}

func (GameMemberWallet) TableName() string { return TableNameGameMemberWallet }
