# 快速开始指南

## 第一步：初始化项目

运行初始化脚本（Windows PowerShell）：

```powershell
.\setup.ps1
```

或手动执行：

```powershell
# 1. 复制配置文件
cp config.yaml.example config.yaml

# 2. 安装依赖
go mod download
go mod tidy
```

## 第二步：复制Plaza协议代码

### 方案A：自动复制（推荐）

如果 `battle-tiles` 和 `battle-bot` 在同一目录下，运行：

```powershell
.\setup.ps1
```

### 方案B：手动复制

```powershell
# 创建目录
mkdir internal\plaza\game
mkdir internal\plaza\consts

# 复制文件
xcopy /Y ..\battle-tiles\internal\utils\plaza\*.go internal\plaza\
xcopy /Y ..\battle-tiles\internal\dal\vo\game\*.go internal\plaza\game\
xcopy /Y ..\battle-tiles\internal\consts\*.go internal\plaza\consts\
```

### 修改Import路径

复制后需要修改所有文件的import路径：

**查找并替换：**
- `battle-tiles/internal/utils/plaza` → `battle-bot/internal/plaza`
- `battle-tiles/internal/dal/vo/game` → `battle-bot/internal/plaza/game`
- `battle-tiles/internal/consts` → `battle-bot/internal/plaza/consts`

可以使用VS Code的全局查找替换功能（Ctrl+Shift+H）。

## 第三步：配置机器人

编辑 `config.yaml`：

```yaml
account:
  username: "你的账号或手机号"
  password: "你的密码"       # 填写明文即可，程序会自动转MD5
  login_mode: "account"     # account(账号登录) 或 mobile(手机号登录)

game:
  house_gid: 123456         # 从battle-tiles获取房间ID
  game_user_id: 0           # 首次运行后自动获取

bot:
  auto_join_table: true
  auto_play: true
  strategy: "random"
```

### 账号配置说明

- `username`: 游戏账号或手机号
- `password`: 登录密码（**填写明文即可**，程序启动时会自动转换为大写MD5）
- `login_mode`: 登录方式，`account`(账号登录) 或 `mobile`(手机号登录)

### 如何获取 house_gid？

1. 登录 battle-tiles 后台管理系统
2. 进入 "店铺管理" 或 "房间管理"
3. 查看房间ID（通常是一个6位数字）

### 如何获取 game_user_id？

首次运行时设置为 0，程序会在登录成功后自动打印：

```
✅ 登录成功！
游戏用户ID: 789012
```

将这个ID填入 `config.yaml` 的 `game_user_id` 字段。

## 第四步：运行机器人

### 方式1：直接运行

```bash
go run cmd/bot/main.go
```

### 方式2：编译后运行

```bash
# 编译
go build -o battle-bot.exe cmd/bot/main.go

# 运行
./battle-bot.exe
```

### 方式3：使用Makefile

```bash
# 运行
make run

# 或编译后运行
make build
./battle-bot.exe
```

## 预期输出

成功启动后应该看到：

```
🤖 四川游戏家园机器人已启动...
账号: 你的账号
房间: 123456
✅ 登录成功！
成员列表更新: 15个成员
房间列表更新: 3个房间
发现可用桌台: 1001 (玩家数: 2)
```

## 故障排查

### 1. 登录失败

```
❌ 登录失败！
```

**可能原因：**
- 账号或密码错误
- login_mode 设置错误（账号登录应为 "account"，手机号登录应为 "mobile"）
- 服务器地址配置错误

**解决方法：**
- 检查 config.yaml 中的账号密码
- 确认 login_mode 设置正确
- 检查网络连接

### 2. 编译错误

```
could not import battle-bot/internal/plaza
```

**可能原因：**
- 未复制 plaza 协议代码
- import 路径未修改

**解决方法：**
- 运行 `.\setup.ps1` 自动复制
- 或手动复制并修改 import 路径

### 3. 房间列表为空

```
房间列表更新: 0个房间
```

**可能原因：**
- house_gid 配置错误
- 当前房间没有活跃桌台

**解决方法：**
- 确认 house_gid 正确
- 登录游戏客户端确认房间是否有桌台

## 开发模式

### 启用调试日志

修改代码添加更详细的日志输出：

```go
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

### 测试单个功能

```bash
# 只测试登录
go run cmd/bot/main.go -config=config.yaml

# 按 Ctrl+C 停止后查看日志
```

## 下一步

1. **实现打牌逻辑**：编辑 `internal/bot/strategy.go`
2. **添加牌型识别**：实现 `CardAnalyzer`
3. **优化策略**：改进 AI 决策算法
4. **添加统计功能**：记录游戏数据

## 参考资料

- [README.md](README.md) - 完整项目文档
- [battle-tiles](../battle-tiles) - 协议参考实现
- [config.yaml.example](config.yaml.example) - 配置示例

## 技术支持

遇到问题？
1. 检查日志输出
2. 查看 [README.md](README.md) 的故障排查部分
3. 提交 Issue

祝你使用愉快！🎉
