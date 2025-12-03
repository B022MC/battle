package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"battle-bot/internal/bot"
	"battle-bot/internal/config"
)

var (
	configPath = flag.String("config", "config.yaml", "配置文件路径")
)

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建机器人
	gameBot, err := bot.NewBot(cfg)
	if err != nil {
		log.Fatalf("创建机器人失败: %v", err)
	}

	// 启动机器人
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := gameBot.Start(ctx); err != nil {
		log.Fatalf("启动机器人失败: %v", err)
	}

	fmt.Println("🤖 四川游戏家园机器人已启动...")
	fmt.Printf("账号: %s\n", cfg.Account.Username)
	fmt.Printf("房间: %d\n", cfg.Game.HouseGID)
	fmt.Println("等待登录...")
	fmt.Println("提示: 如果看到 '✅ 登录成功！' 说明连接正常")

	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\n正在停止机器人...")

	// 优雅关闭
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := gameBot.Stop(shutdownCtx); err != nil {
		log.Printf("停止机器人时出错: %v", err)
	}

	fmt.Println("机器人已停止")
}
