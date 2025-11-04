package dbx

type BaseDB struct {
	Alias  string `gorm:"column:alias;type:character varying(128);primaryKey;comment:主键ID" json:"alias"`      // 别名
	DBName string `gorm:"column:db_name;type:character varying(255);not null;comment:机构数据库名称" json:"db_name"` // 机构数据库名称
}
