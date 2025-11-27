package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/tx7do/kratos-transport/transport/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"battle-tiles/internal/conf"

	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	_ "go.uber.org/automaxprocs"
)

func init() {
	// 设置时区为中国标准时间 (UTC+8)
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Printf("Failed to load timezone: %v\n", err)
		return
	}
	time.Local = loc
}

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "battle-tiles"
	// Version is the version of the compiled software.
	Version string = "v1.0.0"
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs/config.yaml", "config path, eg: -conf config.yaml")
}

// parseLogLevel 将日志级别字符串转换为zapcore.Level
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel // 默认使用 info 级别
	}
}

func newApp(logger log.Logger, hs *gin.Server, ms transport.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			ms,
		),
	)
}

func main() {
	flag.Parse()

	currentDir, _ := os.Getwd()
	fileLog := path.Join(currentDir, fmt.Sprintf("/logs/%s-%s.log", Name, Version))
	fmt.Println(fileLog)
	dir := filepath.Dir(fileLog)
	if err := os.MkdirAll(dir, 0766); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(fileLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
		return
	}

	writeSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(f), zapcore.AddSync(os.Stdout))
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	encoder := zapcore.NewJSONEncoder(cfg)

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// 从配置中获取日志级别，如果未配置则默认使用 info
	logLevel := zapcore.InfoLevel
	if bc.Log != nil && bc.Log.Level != "" {
		logLevel = parseLogLevel(bc.Log.Level)
	}

	logger := kratoszap.NewLogger(zap.New(zapcore.NewCore(encoder, writeSyncer, logLevel)))

	app, cleanup, err := wireApp(bc.Global, bc.Server, bc.Data, bc.Log, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	if err := app.Run(); err != nil {
		panic(err)
	}

}
