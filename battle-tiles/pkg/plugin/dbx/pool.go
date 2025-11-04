package dbx

import (
	"context"
	"fmt"
	"github.com/sethvargo/go-retry"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	llog "log"
	"os"
	"sync"
	"time"
)

var (
	ConnPool *connPool
	once     sync.Once
)

const (
	CtxDBKey = "Platform"
)

type connPool struct {
	source    string
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.RWMutex
	dbMap     map[string]*gorm.DB
	refreshDB *gorm.DB
}

func (p *connPool) connDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.New(
			llog.New(os.Stdout, "\r\n"+dsn+"\r\n", llog.LstdFlags), // 日志输出到标准输出
			gormLogger.Config{
				SlowThreshold: 100 * time.Millisecond, // Slow SQL 阈值
				LogLevel:      gormLogger.Info,        // 日志级别
				Colorful:      true,                   // 彩色日志输出
			},
		),
		CreateBatchSize: 500,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (p *connPool) NewConn(dbName string, alias string) error {
	return retry.Do(p.ctx, retry.WithMaxDuration(time.Second*15, retry.NewFibonacci(time.Second)), func(ctx context.Context) error {
		db, err := p.connDB(fmt.Sprintf("%s dbname=%s", p.source, dbName))
		if err != nil {
			return err
		}

		p.mu.Lock()
		p.dbMap[alias] = db
		p.mu.Unlock()
		fmt.Println("new dbx ", alias, "["+dbName+"]")
		return nil
	})
}

func (p *connPool) GetConn(alias string) (*gorm.DB, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if db, ok := p.dbMap[alias]; ok {
		return db, nil
	} else {
		return nil, fmt.Errorf("dbx %s not found", alias)
	}
}

func (p *connPool) FetchIsExist(sql string) ([]string, error) {
	p.mu.Lock()
	result := make([]string, 0)
	for alias, db := range p.dbMap {
		var isExist bool
		if err := db.Raw(sql).Scan(&isExist).Error; err != nil {
			continue
		}
		if isExist {
			result = append(result, alias)
		}
	}
	p.mu.Unlock()
	return result, nil
}

func (p *connPool) Refresh(sql string, values ...interface{}) {
	list := make([]BaseDB, 0)
	if values == nil {
		if err := p.refreshDB.Raw(sql).Scan(&list).Error; err != nil {
			return
		}
	} else {
		if err := p.refreshDB.Raw(sql, values...).Scan(&list).Error; err != nil {
			return
		}
	}

	for _, db := range list {
		p.NewConn(db.DBName, db.Alias)
	}
}

func (p *connPool) RefreshDB(sql string, values ...interface{}) {
	ticker := time.NewTicker(time.Hour * 12)
	for {
		select {
		case <-p.ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			p.Refresh(sql, values...)
		}

	}
}

func InitConnPool(source string, refreshDB *gorm.DB) {
	once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		ConnPool = &connPool{
			source:    source,
			ctx:       ctx,
			cancel:    cancel,
			mu:        sync.RWMutex{},
			refreshDB: refreshDB,
			dbMap:     make(map[string]*gorm.DB),
		}
	})

}
