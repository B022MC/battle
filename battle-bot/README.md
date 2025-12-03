# Battle Bot - 四川游戏家园机器人

基于Go语言开发的四川游戏家园自动化机器人，用于游戏活跃和测试。

## 功能特性

- ✅ **自动登录**: 支持账号/手机号登录
- ✅ **自动加入游戏**: 智能查找并加入可用桌台
- ✅ **自动打牌**: 根据策略自动出牌
- ✅ **智能延迟**: 模拟真实玩家操作间隔
- ✅ **活跃时段控制**: 可配置机器人活跃时间
- ✅ **游戏次数限制**: 防止过度使用
- ✅ **断线重连**: 自动处理网络异常

## 项目结构

```
battle-bot/
├── cmd/
│   └── bot/
│       └── main.go          # 程序入口
├── internal/
│   ├── bot/
│   │   ├── bot.go           # 机器人核心逻辑
│   │   └── strategy.go      # 打牌策略（待实现）
│   ├── config/
│   │   └── config.go        # 配置管理
│   └── plaza/
│       ├── types.go         # 游戏协议类型定义
│       ├── session.go       # 会话管理（待从battle-tiles复制）
│       ├── encoder.go       # 加密解密（待从battle-tiles复制）
│       └── protocol.go      # 协议解析（待从battle-tiles复制）
├── config.yaml.example      # 配置示例
├── go.mod                   # Go模块定义
└── README.md                # 本文件
```

## 快速开始

### 1. 复制Plaza协议代码

需要从 `battle-tiles` 项目复制以下文件到 `internal/plaza/`:

```bash
# 从battle-tiles复制核心协议文件
cp ../battle-tiles/internal/utils/plaza/*.go internal/plaza/
cp ../battle-tiles/internal/dal/vo/game/*.go internal/plaza/game/
cp ../battle-tiles/internal/consts/const.go internal/plaza/consts/
```

### 2. 配置机器人

```bash
# 复制配置文件
cp config.yaml.example config.yaml

# 编辑配置文件，填写你的账号信息
notepad config.yaml
```

### 3. 安装依赖

```bash
go mod download
```

### 4. 运行机器人

```bash
# 使用默认配置文件运行
go run cmd/bot/main.go

# 指定配置文件
go run cmd/bot/main.go -config=my-config.yaml

# 编译后运行
go build -o battle-bot.exe cmd/bot/main.go
./battle-bot.exe
```

## 配置说明

### 账号配置

- `username`: 游戏账号或手机号
- `password`: 登录密码（明文，程序会自动MD5加密）
- `login_mode`: 登录方式，`account`(账号) 或 `mobile`(手机号)

### 游戏配置

- `house_gid`: 房间ID，可以从battle-tiles后台管理系统获取
- `game_user_id`: 游戏内用户ID，首次登录后会自动获取

### 机器人行为配置

- `auto_join_table`: 是否自动查找并加入桌台
- `auto_play`: 是否自动打牌
- `play_delay_min_ms` / `play_delay_max_ms`: 出牌延迟范围，模拟真实玩家
- `max_games_per_day`: 每天最大游戏局数限制
- `active_hours_start` / `active_hours_end`: 机器人活跃时间段
- `strategy`: 打牌策略
  - `random`: 随机出牌
  - `conservative`: 保守策略（优先保留大牌）
  - `aggressive`: 激进策略（优先出大牌）

### 服务器配置

- `server_82`: 登录服务器地址
- `server_87_host`: 游戏服务器主机
- `keepalive_seconds`: TCP保活间隔
- `auto_reconnect`: 断线是否自动重连

## 开发计划

### ✅ 已完成
- [x] 项目结构搭建
- [x] 配置管理
- [x] 基础机器人框架
- [x] 协议类型定义

### 🚧 进行中
- [ ] 复制并适配plaza协议代码
- [ ] 实现游戏状态解析
- [ ] 实现自动加入桌台逻辑

### 📋 待开发
- [ ] 实现基础出牌策略
- [ ] 添加牌型识别算法
- [ ] 实现AI决策逻辑
- [ ] 添加统计和日志功能
- [ ] Web管理界面

## 技术架构

### 协议层
- 直接使用TCP协议与游戏服务器通信
- 自定义加密算法（参考battle-tiles实现）
- 二进制协议解析

### 机器人层
- 事件驱动架构
- 实现 `IPlazaHandler` 接口处理游戏事件
- 状态机管理游戏流程

### 策略层
- 插件化设计，支持多种打牌策略
- 牌型分析算法
- AI决策引擎

## 注意事项

⚠️ **重要提醒**:

1. **仅供学习研究**: 本项目仅用于技术学习和游戏活跃，请勿用于商业用途
2. **遵守游戏规则**: 使用机器人可能违反游戏服务条款，请自行承担风险
3. **合理使用**: 建议设置合理的延迟和游戏次数限制，避免对游戏服务器造成压力
4. **账号安全**: 配置文件包含密码信息，请妥善保管，不要上传到公共仓库

## 依赖项

- Go 1.21+
- gopkg.in/yaml.v3 - YAML配置解析

## License

MIT License

## 联系方式

如有问题或建议，请提Issue。
