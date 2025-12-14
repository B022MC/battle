package game

import "time"

const TableNameGameBattleRecord = "game_battle_record"

// GameBattleRecord 本地战绩快照（一行代表一局一桌，玩家列表以 JSON 存储）
// 支持按玩家维度查询和统计
type GameBattleRecord struct {
	Id             int32     `gorm:"primaryKey;column:id" json:"id"`
	HouseGID       int32     `gorm:"column:house_gid;not null;index:idx_battle_house_gid,idx_battle_group" json:"house_gid"`
	GroupID        int32     `gorm:"column:group_id;not null" json:"group_id"`
	RoomUID        int32     `gorm:"column:room_uid;not null;index:idx_battle_room_uid" json:"room_uid"` // MappedNum
	KindID         int32     `gorm:"column:kind_id;not null;index:idx_battle_kind_id" json:"kind_id"`
	BaseScore      int32     `gorm:"column:base_score;not null" json:"base_score"`
	BattleAt       time.Time `gorm:"column:battle_at;type:timestamp with time zone;not null;index:idx_battle_at" json:"battle_at"`
	PlayersJSON    string    `gorm:"column:players_json;type:text;not null" json:"players_json"`
	PlayerGameID   *int32    `gorm:"column:player_game_id;index:idx_battle_player_game_id" json:"player_game_id"`
	PlayerGameName string    `gorm:"column:player_game_name;type:varchar(64);default:''" json:"player_game_name"`
	GroupName      string    `gorm:"column:group_name;type:varchar(64);default:'';index:idx_battle_group" json:"group_name"`
	Score          int32     `gorm:"column:score;default:0" json:"score"`
	Fee            int32     `gorm:"column:fee;default:0" json:"fee"`
	Factor         float64   `gorm:"column:factor;type:decimal(10,4);default:1.0000" json:"factor"`
	PlayerBalance  int32     `gorm:"column:player_balance;default:0" json:"player_balance"`
	PlayerCredit   int32     `gorm:"column:player_credit;default:0" json:"player_credit"` // 玩家额度
	CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at"`
}

func (GameBattleRecord) TableName() string { return TableNameGameBattleRecord }
