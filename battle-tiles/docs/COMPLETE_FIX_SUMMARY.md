# 完整修复总结

## 问题概述

用户在注册时绑定了游戏账号，但查询战绩时无法获取数据。根本原因是两个问题的组合：

1. **代码问题**：注册时没有保存 `game_user_id`
2. **查询问题**：查询战绩时使用了错误的参数

## 修复方案

### 第一步：修复注册时的 game_user_id 保存

**文件**：`battle-tiles/internal/biz/game/game_account.go`

**问题**：`BindSingle` 方法只调用 `Verify()` 验证账号，没有获取游戏用户ID

**修复**：
```go
// 修改前：只验证，不获取用户ID
if err := uc.Verify(ctx, mode, identifier, pwdMD5); err != nil {
    return nil, err
}

// 修改后：验证并获取游戏用户信息
info, err := uc.mgr.ProbeLoginWithInfo(ctx, mode, identifier, pwdMD5)
if err != nil {
    return nil, err
}

// 保存游戏用户 ID
a := &model.GameAccount{
    // ...
    GameUserID: fmt.Sprintf("%d", info.UserID), // ✅ 保存游戏用户 ID
}
```

### 第二步：修复历史数据

**文件**：`battle-tiles/internal/biz/game/game_account.go`

**添加方法**：`FixEmptyGameUserID()`

**功能**：
- 查询所有 `game_user_id` 为空的记录
- 逐个调用游戏服务器获取用户ID
- 更新数据库

**使用方式**：
```bash
curl -X POST http://localhost:8080/game/accounts/fix-empty-game-user-id \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### 第三步：修复查询战绩的参数

**文件**：`battle-tiles/internal/biz/game/battle_record.go`

**问题**：查询战绩时使用 `account.Account`（游戏账号）而不是 `account.GameUserID`（游戏用户ID）

**修复的方法**：

#### 1. ListMyBattleRecords
```go
// 修改前
return uc.repo.ListByPlayer(ctx, houseGID, account.Account, GroupID, start, end, page, size)

// 修改后
if account.GameUserID == "" {
    uc.log.Warnf("User %d has no game_user_id", userID)
    return []*model.GameBattleRecord{}, 0, nil
}
return uc.repo.ListByPlayer(ctx, houseGID, account.GameUserID, GroupID, start, end, page, size)
```

#### 2. GetMyBattleStats
```go
// 修改前
return uc.repo.GetPlayerStats(ctx, houseGID, account.Account, groupID, start, end)

// 修改后
if account.GameUserID == "" {
    uc.log.Warnf("User %d has no game_user_id", userID)
    return 0, 0, 0, nil
}
return uc.repo.GetPlayerStats(ctx, houseGID, account.GameUserID, groupID, start, end)
```

#### 3. GetMyBalances
```go
// 修改前：复杂的解析逻辑
var playerGameID int32
if ok, err := parseGameUserID(account.GameUserID, &playerGameID); !ok || err != nil {
    return nil, fmt.Errorf("游戏账号ID格式错误")
}

// 修改后：简化逻辑
if account.GameUserID == "" {
    uc.log.Warnf("User %d has no game_user_id", userID)
    return []interface{}{}, nil
}
```

## 修改的文件列表

### 代码修改
1. `battle-tiles/internal/biz/game/game_account.go`
   - 修改 `BindSingle` 方法
   - 添加 `FixEmptyGameUserID` 方法

2. `battle-tiles/internal/dal/repo/game/game_account.go`
   - 添加 `ListByGameUserIDEmpty` 方法
   - 添加 `Update` 方法

3. `battle-tiles/internal/service/game/game_account.go`
   - 添加 `FixEmptyGameUserID` 处理函数
   - 添加管理员路由

4. `battle-tiles/pkg/plugin/middleware/rbac.go`
   - 添加 `AdminOnly` 中间件

5. `battle-tiles/internal/biz/game/battle_record.go`
   - 修改 `ListMyBattleRecords` 方法
   - 修改 `GetMyBattleStats` 方法
   - 修改 `GetMyBalances` 方法

### 文档
1. `battle-tiles/docs/FIX_EMPTY_GAME_USER_ID.md` - 数据修复指南
2. `battle-tiles/docs/QUICK_FIX_GUIDE.md` - 快速修复指南
3. `battle-tiles/docs/BATTLE_QUERY_FIX.md` - 查询战绩修复文档
4. `battle-tiles/docs/QUERY_BATTLES_QUICK_FIX.md` - 查询战绩快速指南
5. `battle-tiles/docs/COMPLETE_FIX_SUMMARY.md` - 本文档

## 完整修复流程

### 1. 编译并启动服务
```bash
cd battle-tiles
go build -o battle-tiles ./cmd/main.go
./battle-tiles
```

### 2. 修复历史数据（可选）
```bash
# 获取管理员 token
TOKEN=$(curl -X POST http://localhost:8080/basic/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}' | jq -r '.data.token')

# 修复空的 game_user_id
curl -X POST http://localhost:8080/game/accounts/fix-empty-game-user-id \
  -H "Authorization: Bearer $TOKEN"
```

### 3. 测试查询战绩
```bash
# 用户登录
TOKEN=$(curl -X POST http://localhost:8080/basic/login \
  -H "Content-Type: application/json" \
  -d '{"username": "test1", "password": "password"}' | jq -r '.data.token')

# 查询战绩
curl -X POST http://localhost:8080/battle-query/my/battles \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"house_gid": 58959, "page": 1, "size": 20}'
```

## 验证修复

### 检查点

1. **新用户注册**
   - 注册时绑定游戏账号
   - 检查 `game_account` 表中 `game_user_id` 是否被填充

2. **查询战绩**
   - 调用 `/battle-query/my/battles`
   - 应该返回用户的战绩列表

3. **查询统计**
   - 调用 `/battle-query/my/stats`
   - 应该返回统计数据

4. **历史数据**
   - 运行修复接口
   - 检查 `game_account` 表中空的 `game_user_id` 是否被填充

## 关键改动总结

| 问题 | 原因 | 修复 |
|------|------|------|
| 新用户无法查询战绩 | 注册时没有保存 `game_user_id` | 调用 `ProbeLoginWithInfo` 获取并保存 |
| 历史用户无法查询战绩 | `game_user_id` 为空 | 提供修复接口自动填充 |
| 查询战绩返回空结果 | 使用错误的查询参数 | 改为使用 `game_user_id` |

## 后续改进

1. 实现 `GetMyBalances` 的实际逻辑
2. 添加更多的查询过滤条件
3. 优化查询性能
4. 添加更详细的日志记录

