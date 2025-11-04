package dbx

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrDBKeyNotFound = errors.New("dbx key not found")
	ErrDBKeyNotEmpty = errors.New("dbx key not empty")
)

type AutoSwitchingDBPlugin struct {
}

func (a *AutoSwitchingDBPlugin) Name() string {
	return "auto-switching-dbx"
}

func (a *AutoSwitchingDBPlugin) Initialize(db *gorm.DB) error {
	db.Callback().Create().Before("gorm:create").Register("dbplugin:before_create", a.beforeCreate)
	db.Callback().Query().Before("gorm:query").Register("dbplugin:before_query", a.beforeQuery)
	db.Callback().Update().Before("gorm:update").Register("dbplugin:before_update", a.beforeUpdate)
	db.Callback().Delete().Before("gorm:delete").Register("dbplugin:before_delete", a.beforeDelete)
	db.Callback().Raw().Before("gorm:raw").Register("dbplugin:before_raw", a.beforeRaw)
	db.Callback().Row().Before("gorm:row").Register("dbplugin:before_row", a.beforeRow)
	return nil
}

func (a *AutoSwitchingDBPlugin) switchingDB(db *gorm.DB) {
	ctx := db.Statement.Context
	if ctx.Value(CtxDBKey) == nil {
		db.AddError(ErrDBKeyNotFound)
		return
	}
	dbKey, ok := ctx.Value(CtxDBKey).(string)
	if !ok {
		db.AddError(ErrDBKeyNotFound)
		return
	}
	if dbKey == "" {
		db.AddError(ErrDBKeyNotEmpty)
		return
	}
	conn, err := ConnPool.GetConn(dbKey)
	if err != nil {
		db.AddError(err)
		return
	}
	// **检查 `ConnPool` 是否已是目标数据库**
	if db.Statement.ConnPool == conn.ConnPool {
		return
	}
	db.Statement.ConnPool = conn.ConnPool
	db.Statement.DB = conn
}

func (a *AutoSwitchingDBPlugin) beforeCreate(db *gorm.DB) {
	a.switchingDB(db)
}
func (a *AutoSwitchingDBPlugin) beforeQuery(db *gorm.DB) {
	a.switchingDB(db)
}
func (a *AutoSwitchingDBPlugin) beforeUpdate(db *gorm.DB) {
	a.switchingDB(db)
}
func (a *AutoSwitchingDBPlugin) beforeDelete(db *gorm.DB) {
	a.switchingDB(db)
}
func (a *AutoSwitchingDBPlugin) beforeRaw(db *gorm.DB) {
	a.switchingDB(db)
}

func (a *AutoSwitchingDBPlugin) beforeRow(db *gorm.DB) {
	a.switchingDB(db)
}

func NewAutoSwitchingDBPlugin() gorm.Plugin {
	return &AutoSwitchingDBPlugin{}
}
