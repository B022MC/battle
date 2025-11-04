package game

import "time"

const TableNameGameMemberWallet = "game_member_wallet"

type GameMemberWallet struct {
	Id        int32     `gorm:"primaryKey;column:id" json:"id"`
	HouseGID  int32     `gorm:"column:house_gid;not null" json:"house_gid"`
	MemberID  int32     `gorm:"column:member_id;not null" json:"member_id"`
	Balance   int32     `gorm:"column:balance;not null;default:0" json:"balance"`
	Forbid    bool      `gorm:"column:forbid;not null;default:false" json:"forbid"`
	LimitMin  int32     `gorm:"column:limit_min;not null;default:0" json:"limit_min"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
	UpdatedBy int32     `gorm:"column:updated_by;not null;default:0" json:"updated_by"`
}

func (GameMemberWallet) TableName() string { return TableNameGameMemberWallet }
