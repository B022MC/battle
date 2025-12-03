# 费用规则检查功能调试指南

## 问题分析

从您的日志来看，程序已经启动，但没有看到玩家坐下的事件。这可能是以下几个原因：

### 1. 会话未启动

**现象**：看到服务器启动日志，但没有看到 `[费用检查] 创建组合handler` 的日志

**原因**：需要通过 API 启动中控会话，才能连接到游戏服务器

**解决方法**：
```bash
# 调用会话启动接口
POST /game/accounts/sessionStart
{
  "ctrl_account_id": <中控账号ID>,
  "house_gid": <店铺GID>
}
```

启动成功后应该看到：
```
INFO [费用检查] 创建组合handler: userID=xxx, houseGID=xxx
INFO module=utils/plaza msg=MDM_MB_SERVER_LIST/SUB_MB_LIST_SERVER ...
```

### 2. 日志级别调整

我已经将关键日志从 `Debugf` 提升到 `Infof`，现在您应该能看到以下日志：

#### 会话启动时
```
INFO [费用检查] 创建组合handler: userID=1, houseGID=12345
```

#### 玩家坐下时
```
INFO [费用检查] 收到玩家坐下事件: gameID=xxx, mappedNum=1, userID=1, houseGID=12345
INFO [费用检查] 桌子信息: gameID=xxx, kindID=101, baseScore=1, mappedNum=1
INFO [费用检查] 玩家信息: gameID=xxx, balance=1000, roomCredit=0, playerCredit=0, effectiveCredit=0
INFO [费用检查] 检查费用资格: gameID=xxx, balance=1000, requiredFee=500
INFO [费用检查] 费用检查通过: gameID=xxx, balance=1000 >= fee=500
INFO [费用检查] 所有检查通过: gameID=xxx, balance=1000 >= credit=0, keep table
```

#### 余额不足时
```
WARN HandlePlayerSitDown: insufficient balance for fee, gameID=xxx, balance=100 < fee=500, dismissing table
```

### 3. 检查配置

确保以下配置正确：

#### 3.1 中控账号配置
- 中控账号需要先验证（调用 `/game/accounts/verify`）
- 验证成功后会设置 `game_player_id`
- 启动会话时会检查这个字段，如果为空会报错

#### 3.2 店铺费用规则配置
```sql
-- 查询店铺费用配置
SELECT house_gid, fees_json 
FROM game_house_settings 
WHERE house_gid = <你的店铺GID>;
```

费用规则格式：
```json
{
  "rules": [
    {
      "threshold": 100,
      "fee": 1000,
      "kind": 0,
      "base": 0
    }
  ]
}
```

#### 3.3 玩家成员信息
```sql
-- 查询玩家是否在店铺成员列表中
SELECT * FROM game_member 
WHERE house_gid = <店铺GID> 
AND game_id = <游戏ID>;
```

## 完整测试流程

### 步骤 1: 创建中控账号
```bash
POST /shops/ctrlAccounts
{
  "login_mode": 1,  # 1=账号登录, 2=手机登录
  "identifier": "test001",
  "password": "123456"
}
```

### 步骤 2: 验证中控账号
```bash
POST /game/accounts/verify
{
  "login_mode": 1,
  "identifier": "test001",
  "password": "123456"
}
```

### 步骤 3: 绑定中控账号到店铺
```bash
POST /shops/ctrlAccounts/bind
{
  "ctrl_account_id": <中控账号ID>,
  "house_gid": <店铺GID>
}
```

### 步骤 4: 配置店铺费用规则
```bash
POST /shops/fees/set
{
  "house_gid": <店铺GID>,
  "fees_json": {
    "rules": [
      {
        "threshold": 0,
        "fee": 1000,  # 最低需要 1000 分（10元）
        "kind": 0,
        "base": 0
      }
    ]
  }
}
```

### 步骤 5: 启动会话
```bash
POST /game/accounts/sessionStart
{
  "ctrl_account_id": <中控账号ID>,
  "house_gid": <店铺GID>
}
```

**预期日志**：
```
INFO [费用检查] 创建组合handler: userID=xxx, houseGID=xxx
INFO module=utils/plaza msg=MDM_MB_SERVER_LIST ...
```

### 步骤 6: 模拟玩家加入

让一个玩家加入店铺并坐到桌子上。

**如果玩家余额 >= 1000**：
```
INFO [费用检查] 收到玩家坐下事件: gameID=xxx, mappedNum=1, ...
INFO [费用检查] 桌子信息: gameID=xxx, kindID=101, baseScore=1, ...
INFO [费用检查] 玩家信息: gameID=xxx, balance=5000, ...
INFO [费用检查] 检查费用资格: gameID=xxx, balance=5000, requiredFee=1000
INFO [费用检查] 费用检查通过: gameID=xxx, balance=5000 >= fee=1000
INFO [费用检查] 所有检查通过: gameID=xxx, balance=5000 >= credit=0, keep table
```

**如果玩家余额 < 1000**：
```
INFO [费用检查] 收到玩家坐下事件: gameID=xxx, mappedNum=1, ...
INFO [费用检查] 桌子信息: gameID=xxx, kindID=101, baseScore=1, ...
INFO [费用检查] 玩家信息: gameID=xxx, balance=500, ...
INFO [费用检查] 检查费用资格: gameID=xxx, balance=500, requiredFee=1000
WARN HandlePlayerSitDown: insufficient balance for fee, gameID=xxx, balance=500 < fee=1000, dismissing table
```

玩家会被踢出桌子。

## 常见问题排查

### Q1: 看不到 `[费用检查]` 相关日志

**可能原因**：
1. 会话未启动 → 检查是否调用了 `/game/accounts/sessionStart`
2. `creditHandler` 为空 → 检查 Wire 依赖注入是否正确

**检查方法**：
```bash
# 查看是否有 "创建组合handler" 的日志
grep "创建组合handler" logs/battle-tiles-*.log

# 查看是否有 "creditHandler为空" 的警告
grep "creditHandler为空" logs/battle-tiles-*.log
```

### Q2: 看到坐下事件，但没有执行检查

**可能原因**：
1. 店铺未配置费用规则
2. 费用规则格式错误
3. 玩家不在成员列表中

**检查方法**：
```bash
# 查看完整的错误日志
grep "HandlePlayerSitDown" logs/battle-tiles-*.log | grep -E "WARN|ERROR"
```

### Q3: 玩家符合条件但还是被踢出

**可能原因**：
1. 房间额度检查失败（在费用检查之前）
2. 玩家余额计算错误

**检查方法**：
```bash
# 查看玩家余额和要求
grep "玩家信息" logs/battle-tiles-*.log
```

### Q4: 如果 `creditHandler` 为空

**现象**：日志显示 `[警告] creditHandler为空，未启用费用检查`

**解决方法**：
1. 检查 Wire 配置是否正确生成
2. 重新生成 Wire 代码：
```bash
cd /Users/b022mc/project/battle/battle-tiles
wire gen ./cmd/go-kgin-platform
```
3. 重新编译：
```bash
go build ./cmd/go-kgin-platform
```

## 直接测试方法

如果您想快速测试，可以：

1. **启动程序**
2. **启动会话**（上面步骤 5）
3. **观察日志**：
   - 如果看到 `[费用检查] 创建组合handler`，说明事件处理器已经注入
   - 如果看到 `[警告] creditHandler为空`，说明依赖注入有问题

4. **让玩家加入桌子**
   - 如果看到 `[费用检查] 收到玩家坐下事件`，说明事件正常触发
   - 如果没有看到，检查游戏服务器连接是否正常

## 监控命令

实时监控日志：
```bash
tail -f /Users/b022mc/project/battle/battle-tiles/logs/battle-tiles-*.log | grep -E "\[费用检查\]|HandlePlayerSitDown|OnUserSitDown"
```

这样可以实时看到所有费用检查相关的日志。
