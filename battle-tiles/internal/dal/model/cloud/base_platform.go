package cloud

import "time"

const TableNameBasePlatform = "base_platform"

// BasePlatform mapped from table <base_platform>
type BasePlatform struct {
	Platform  string     `gorm:"column:platform;type:character varying(255);not null;comment:平台" json:"platform"`
	Name      string     `gorm:"column:name;type:character varying(255);not null;comment:名称" json:"name"`                                                                                                                  // 名称
	DBName    string     `gorm:"column:db_name;type:character varying(255);not null;comment:数据库名称" json:"db_name"`                                                                                                         // 机构数据库名称
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp(6) with time zone;not null;default:now();comment:创建时间" json:"created_at" time_format:"2006-01-02 15:04:05" time_utc:"false" format:"2006-01-02 15:04:05"` // 创建时间
}

// TableName BaseBed's table name
func (*BasePlatform) TableName() string {
	return TableNameBasePlatform
}
