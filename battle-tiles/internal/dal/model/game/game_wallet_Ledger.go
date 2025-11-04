package game

import "time"

const TableNameGameWalletLedger = "game_wallet_ledger"

type GameWalletLedger struct {
	Id             int32     `gorm:"primaryKey;column:id" json:"id"`
	HouseGID       int32     `gorm:"column:house_gid;not null" json:"house_gid"`
	MemberID       int32     `gorm:"column:member_id;not null" json:"member_id"`
	ChangeAmount   int32     `gorm:"column:change_amount;not null" json:"change_amount"`
	BalanceBefore  int32     `gorm:"column:balance_before;not null" json:"balance_before"`
	BalanceAfter   int32     `gorm:"column:balance_after;not null" json:"balance_after"`
	Type           int32     `gorm:"column:type;not null" json:"type"` // 1=上分 2=下分 3=强制下分 4=调整
	Reason         string    `gorm:"column:reason;type:text;not null;default:''" json:"reason"`
	OperatorUserID int32     `gorm:"column:operator_user_id;not null" json:"operator_user_id"`
	BizNo          string    `gorm:"column:biz_no;type:varchar(64);not null" json:"biz_no"`
	CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
}

func (GameWalletLedger) TableName() string { return TableNameGameWalletLedger }
