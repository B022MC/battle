package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 机器人配置
type Config struct {
	Account AccountConfig `yaml:"account"`
	Game    GameConfig    `yaml:"game"`
	Bot     BotConfig     `yaml:"bot"`
	Plaza   PlazaConfig   `yaml:"plaza"`
}

// AccountConfig 账号配置
type AccountConfig struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	LoginMode string `yaml:"login_mode"` // "account" 或 "mobile"
}

// GameConfig 游戏配置
type GameConfig struct {
	HouseGID   int `yaml:"house_gid"`    // 房间ID
	GameUserID int `yaml:"game_user_id"` // 游戏用户ID
}

// BotConfig 机器人行为配置
type BotConfig struct {
	AutoJoinTable    bool   `yaml:"auto_join_table"`    // 自动加入桌台
	AutoPlay         bool   `yaml:"auto_play"`          // 自动打牌
	PlayDelayMin     int    `yaml:"play_delay_min_ms"`  // 出牌最小延迟(毫秒)
	PlayDelayMax     int    `yaml:"play_delay_max_ms"`  // 出牌最大延迟(毫秒)
	MaxGamesPerDay   int    `yaml:"max_games_per_day"`  // 每天最大游戏局数
	ActiveHoursStart int    `yaml:"active_hours_start"` // 活跃时间开始(小时)
	ActiveHoursEnd   int    `yaml:"active_hours_end"`   // 活跃时间结束(小时)
	Strategy         string `yaml:"strategy"`           // 打牌策略: "random", "conservative", "aggressive"
}

// PlazaConfig 游戏服务器配置
type PlazaConfig struct {
	Server82      string `yaml:"server_82"`         // 登录服务器地址
	Server87Host  string `yaml:"server_87_host"`    // 游戏服务器主机
	KeepAlive     int    `yaml:"keepalive_seconds"` // 保活时间(秒)
	AutoReconnect bool   `yaml:"auto_reconnect"`    // 自动重连
}

// Load 从文件加载配置
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 设置默认值
	if cfg.Bot.PlayDelayMin == 0 {
		cfg.Bot.PlayDelayMin = 1000
	}
	if cfg.Bot.PlayDelayMax == 0 {
		cfg.Bot.PlayDelayMax = 3000
	}
	if cfg.Bot.MaxGamesPerDay == 0 {
		cfg.Bot.MaxGamesPerDay = 100
	}
	if cfg.Bot.ActiveHoursStart == 0 {
		cfg.Bot.ActiveHoursStart = 8
	}
	if cfg.Bot.ActiveHoursEnd == 0 {
		cfg.Bot.ActiveHoursEnd = 22
	}
	if cfg.Bot.Strategy == "" {
		cfg.Bot.Strategy = "random"
	}
	if cfg.Plaza.Server82 == "" {
		cfg.Plaza.Server82 = "androidsc.foxuc.com:8200"
	}
	if cfg.Plaza.Server87Host == "" {
		cfg.Plaza.Server87Host = "newbgp.foxuc.com"
	}
	if cfg.Plaza.KeepAlive == 0 {
		cfg.Plaza.KeepAlive = 30
	}

	return &cfg, nil
}
