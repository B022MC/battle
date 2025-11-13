package game

import "time"

const TableNameGameRechargeRecord = "game_recharge_record"

// GameRechargeRecord represents wallet recharge/withdrawal records
type GameRechargeRecord struct {
	Id             int32     `gorm:"primaryKey;column:id" json:"id"`
	HouseGID       int32     `gorm:"column:house_gid;not null;index:idx_recharge_house_gid,idx_recharge_house_group" json:"house_gid"`
	PlayerID       int32     `gorm:"column:player_id;not null;index:idx_recharge_player" json:"player_id"`
	GroupName      string    `gorm:"column:group_name;type:varchar(64);not null;default:'';index:idx_recharge_house_group" json:"group_name"`
	Amount         int32     `gorm:"column:amount;not null" json:"amount"` // Positive = deposit, Negative = withdrawal
	BalanceBefore  int32     `gorm:"column:balance_before;not null" json:"balance_before"`
	BalanceAfter   int32     `gorm:"column:balance_after;not null" json:"balance_after"`
	OperatorUserID *int32    `gorm:"column:operator_user_id" json:"operator_user_id"`
	RechargedAt    time.Time `gorm:"column:recharged_at;type:timestamp with time zone;not null;index:idx_recharge_recharged_at" json:"recharged_at"`
	CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
}

func (GameRechargeRecord) TableName() string { return TableNameGameRechargeRecord }

// IsDeposit checks if this is a deposit transaction
func (r *GameRechargeRecord) IsDeposit() bool {
	return r.Amount > 0
}

// IsWithdrawal checks if this is a withdrawal transaction
func (r *GameRechargeRecord) IsWithdrawal() bool {
	return r.Amount < 0
}

// GetAmountInYuan returns amount in yuan
func (r *GameRechargeRecord) GetAmountInYuan() float64 {
	return float64(r.Amount) / 100.0
}

