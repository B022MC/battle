package game

import "time"

const TableNameGameBattleRecord = "game_battle_record"

// GameBattleRecord 本地战绩快照（一行代表一局一桌，玩家列表以 JSON 存储）
type GameBattleRecord struct {
	Id          int32     `gorm:"primaryKey;column:id" json:"id"`
	HouseGID    int32     `gorm:"column:house_gid;not null" json:"house_gid"`
	GroupID     int32     `gorm:"column:group_id;not null" json:"group_id"`
	RoomUID     int32     `gorm:"column:room_uid;not null" json:"room_uid"` // MappedNum
	KindID      int32     `gorm:"column:kind_id;not null" json:"kind_id"`
	BaseScore   int32     `gorm:"column:base_score;not null" json:"base_score"`
	BattleAt    time.Time `gorm:"column:battle_at;type:timestamp with time zone;not null" json:"battle_at"`
	PlayersJSON string    `gorm:"column:players_json;type:text;not null" json:"players_json"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
}

func (GameBattleRecord) TableName() string { return TableNameGameBattleRecord }
