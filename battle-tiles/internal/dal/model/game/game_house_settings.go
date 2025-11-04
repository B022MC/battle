package game

import "time"

const TableNameGameHouseSettings = "game_house_settings"

// GameHouseSettings 保存店铺级设置（运费、分运开关、推送额度等）
type GameHouseSettings struct {
	Id         int32     `gorm:"primaryKey;column:id" json:"id"`
	HouseGID   int32     `gorm:"uniqueIndex:uk_house;column:house_gid;not null" json:"house_gid"`
	FeesJSON   string    `gorm:"column:fees_json;type:text;not null;default:''" json:"fees_json"` // 运费规则 JSON
	ShareFee   bool      `gorm:"column:share_fee;not null;default:false" json:"share_fee"`        // 分运开关
	PushCredit int32     `gorm:"column:push_credit;not null;default:0" json:"push_credit"`        // 推送额度（单位：分）
	UpdatedAt  time.Time `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at"`
	UpdatedBy  int32     `gorm:"column:updated_by;not null;default:0" json:"updated_by"` // 操作人（平台用户ID）
}

func (GameHouseSettings) TableName() string { return TableNameGameHouseSettings }
