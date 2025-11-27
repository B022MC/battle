package infra

import (
	"battle-tiles/internal/conf"
	"battle-tiles/internal/consts"
	cloudModel "battle-tiles/internal/dal/model/cloud"
	"battle-tiles/internal/infra/plaza"
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	pdb "battle-tiles/pkg/plugin/dbx"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewPSQL,
	NewDBMap,
	//NewData,
	NewRedis,
	NewRdbMap,
	NewData,
	plaza.NewManager,
)

// Data .
type Data struct {
	ConfData     *conf.Data
	log          *log.Helper
	DB           *gorm.DB
	DBMap        map[string]*gorm.DB
	RDB          *redis.Client
	RDBMap       map[string]*redis.Client
	PlazaManager plaza.Manager
}

// NewData .
func NewData(
	c *conf.Data,
	logger log.Logger,
	rdb *redis.Client,
	db *gorm.DB,
	dBMap map[string]*gorm.DB,
	RDBMap map[string]*redis.Client,
	pm plaza.Manager,
) (*Data, func(), error) {
	logging := log.NewHelper(log.With(logger, "module", "infra"))
	cleanup := func() {
		rdb.Close()
		pm.StopAll()
		logging.Info("closing the infra resources")
	}
	// 初始化连接池
	cdb, cok := dBMap[consts.CloudPlatformDB]
	if !cok {
		panic("cloud cloud db not found")
	}
	pdb.InitConnPool(c.Database.Source, cdb)
	refreshDBSource()
	rdb.AddHook(pdb.NewWithPlatformKeyHook())
	return &Data{
		ConfData:     c,
		log:          logging,
		RDB:          rdb,
		RDBMap:       RDBMap,
		DB:           db,
		DBMap:        dBMap,
		PlazaManager: pm,
	}, cleanup, nil
}

func (data *Data) GetDBWithContext(ctx context.Context) *gorm.DB {
	if ctx == nil { // 如果ctx为nil，则创建一个默认的上下文 解决迁移遗漏的
		ctx = context.Background()
	}
	db := data.DB.WithContext(ctx)
	if ctx.Value(pdb.CtxDBKey) == nil {
		db.AddError(errors.New("db key not found"))
		return db
	}
	dbKey, ok := ctx.Value(pdb.CtxDBKey).(string)
	if !ok {
		db.AddError(errors.New("db key type error"))
		return db
	}
	if dbKey == "" {
		db.AddError(errors.New("db key is empty"))
		return db
	}
	conn, err := pdb.ConnPool.GetConn(dbKey)
	if err != nil {
		db.AddError(err)
		return db
	}

	// 检查是否需要静默模式
	if quiet, ok := ctx.Value(quietDBKey).(bool); ok && quiet {
		return conn.WithContext(ctx).Session(&gorm.Session{
			Logger: gormLogger.New(
				zap.NewStdLog(zap.L()),
				gormLogger.Config{
					SlowThreshold: 100 * time.Millisecond,
					LogLevel:      gormLogger.Silent, // 静默模式，不输出SQL日志
					Colorful:      false,
				},
			),
		})
	}

	return conn.WithContext(ctx)
}

// NewRedis redis 连接实例
func NewRedis(c *conf.Data) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       int(c.Redis.Db),
	})
	if statusCmd := client.Ping(context.Background()); statusCmd != nil {
		_, err := statusCmd.Result()
		if err != nil {
			panic(err)
		}
		return client
	}
	return nil
}

// NewPSQL psql 连接实例
func NewPSQL(c *conf.Data, logConf *conf.Log) *gorm.DB {
	// 全局设置为中国时区，且统一 GORM 的时间来源
	loc, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = loc

	db, err := gorm.Open(postgres.Open(c.Database.Source), &gorm.Config{
		Logger: gormLogger.New(
			zap.NewStdLog(zap.L()),
			gormLogger.Config{
				SlowThreshold: 100 * time.Millisecond,
				LogLevel:      gormLogger.Info, // 默认显示所有SQL日志
				Colorful:      true,
			},
		),

		CreateBatchSize: 500,
		NowFunc:         func() time.Time { return time.Now().In(loc) },
	})
	if err != nil {
		panic(err)
	}

	// 设置当前连接会话时区，确保数据库侧 now()/timestamp 为北京时间
	_ = db.Exec("SET TIME ZONE 'Asia/Shanghai'").Error
	return db
}

func NewDBMap(c *conf.Data, logConf *conf.Log, logger log.Logger) map[string]*gorm.DB {
	helper := log.NewHelper(logger)

	DBMap := make(map[string]*gorm.DB, len(c.DatabaseList))
	for _, database := range c.DatabaseList {
		switch database.Driver {
		case "cloudsqlpostgres":
			// 重复设置时区无副作用
			loc, _ := time.LoadLocation("Asia/Shanghai")
			time.Local = loc
			db, err := gorm.Open(postgres.Open(database.Source), &gorm.Config{
				Logger: gormLogger.New(
					zap.NewStdLog(zap.L()),
					gormLogger.Config{
						SlowThreshold: 100 * time.Millisecond,
						LogLevel:      gormLogger.Info, // 默认显示所有SQL日志
						Colorful:      false,
					},
				),

				CreateBatchSize: 500,
				NowFunc:         func() time.Time { return time.Now().In(loc) },
			})
			if err != nil {
				helper.Errorw("connect",
					"driver", database.Driver,
					"alias", database.Alias,
					"dsn", database.Source,
					"err", err,
				)
				continue
			}
			// 设置连接会话时区
			_ = db.Exec("SET TIME ZONE 'Asia/Shanghai'").Error
			DBMap[database.Alias] = db
		}
	}
	return DBMap
}
func NewRdbMap(c *conf.Data, logger log.Logger) (map[string]*redis.Client, error) {
	helper := log.NewHelper(logger)
	RDBMap := make(map[string]*redis.Client, len(c.RedisList))

	for _, rdbConf := range c.RedisList {
		var client *redis.Client

		switch rdbConf.Alias {
		case "redis_device_collect":
			client = redis.NewClient(&redis.Options{
				Addr:     rdbConf.Addr,
				Password: rdbConf.Password,
				DB:       int(rdbConf.Db),
			})
		}

		pong, err := client.Ping(context.Background()).Result()
		if err != nil {
			helper.Errorw("redis connect error",
				zap.String("alias", rdbConf.Alias),
				zap.String("addr", rdbConf.Addr),
				zap.Error(err),
			)
			return nil, err
		}
		helper.Info("redis ping successful",
			zap.String("alias", rdbConf.Alias),
			zap.String("addr", rdbConf.Addr),
			zap.String("reply", pong),
		)
		RDBMap[rdbConf.Alias] = client
	}

	return RDBMap, nil
}

// quietDBKey 用于在 context 中标记是否使用静默模式
type contextKey string

const quietDBKey contextKey = "quiet_db"

// WithQuietDB 返回一个带静默标记的 context
func WithQuietDB(ctx context.Context) context.Context {
	return context.WithValue(ctx, quietDBKey, true)
}

// GetDBQuiet 获取一个静默模式的DB（不输出SQL日志），用于定时任务等场景
func (d *Data) GetDBQuiet(ctx context.Context) *gorm.DB {
	return d.GetDBWithContext(ctx).Session(&gorm.Session{
		Logger: gormLogger.New(
			zap.NewStdLog(zap.L()),
			gormLogger.Config{
				SlowThreshold: 100 * time.Millisecond,
				LogLevel:      gormLogger.Silent, // 静默模式，不输出SQL日志
				Colorful:      false,
			},
		),
	})
}

func refreshDBSource() {
	sb := sqlbuilder.NewSelectBuilder()
	// 使用 platform 字段作为连接池的别名（与请求头 Platform 一致）
	sb.Select("platform as alias", "db_name").
		From(cloudModel.TableNameBasePlatform)
	sql, args := sb.Build()

	pdb.ConnPool.Refresh(sql)
	go pdb.ConnPool.RefreshDB(sql, args...)
}
