# 查询战绩修复文档

## 问题描述

查询战绩接口在查询时使用了错误的参数：
- **错误做法**：使用 `account.Account`（游戏账号，如 "11439919"）
- **正确做法**：使用 `account.GameUserID`（游戏用户ID，如 "12345"）

这导致用户无法查询到自己的战绩，即使 `game_user_id` 字段已经被正确填充。

## 根本原因

在 `battle-tiles/internal/biz/game/battle_record.go` 中：

1. `ListMyBattleRecords` 方法传递的是 `account.Account` 而不是 `account.GameUserID`
2. `GetMyBattleStats` 方法也有同样的问题
3. 数据库中的 `game_battle_record` 表存储的是 `player_game_id`（游戏用户ID），而不是游戏账号

## 修复方案

### 修改的方法

#### 1. ListMyBattleRecords（查询用户战绩列表）

**修改前**：
```go
return uc.repo.ListByPlayer(ctx, houseGID, account.Account, GroupID, start, end, page, size)
```

**修改后**：
```go
// 检查是否有 game_user_id
if account.GameUserID == "" {
    uc.log.Warnf("User %d has no game_user_id", userID)
    return []*model.GameBattleRecord{}, 0, nil
}

// 查询战绩（使用 game_user_id 而不是 account）
return uc.repo.ListByPlayer(ctx, houseGID, account.GameUserID, GroupID, start, end, page, size)
```

#### 2. GetMyBattleStats（查询用户战绩统计）

**修改前**：
```go
return uc.repo.GetPlayerStats(ctx, houseGID, account.Account, groupID, start, end)
```

**修改后**：
```go
// 检查是否有 game_user_id
if account.GameUserID == "" {
    uc.log.Warnf("User %d has no game_user_id", userID)
    return 0, 0, 0, nil
}

// 查询统计（使用 game_user_id 而不是 account）
return uc.repo.GetPlayerStats(ctx, houseGID, account.GameUserID, groupID, start, end)
```

#### 3. GetMyBalances（查询用户余额）

**修改前**：
```go
// 解析 game_user_id 为整数
var playerGameID int32
if ok, err := parseGameUserID(account.GameUserID, &playerGameID); !ok || err != nil {
    uc.log.Errorf("Failed to parse game_user_id %s: %v", account.GameUserID, err)
    return nil, fmt.Errorf("游戏账号ID格式错误")
}

// 返回空列表（暂时实现）
return []interface{}{}, nil
```

**修改后**：
```go
// 检查是否有 game_user_id
if account.GameUserID == "" {
    uc.log.Warnf("User %d has no game_user_id", userID)
    return []interface{}{}, nil
}

// 返回空列表（暂时实现）
// TODO: 实现实际的余额查询逻辑
return []interface{}{}, nil
```

## 关键改动

1. **使用 `account.GameUserID` 而不是 `account.Account`**
   - `account.Account` 是用户输入的游戏账号（如手机号或账号名）
   - `account.GameUserID` 是游戏服务器返回的用户ID（数字形式）
   - 数据库中的 `game_battle_record` 表使用 `player_game_id` 字段存储游戏用户ID

2. **添加 `game_user_id` 为空的检查**
   - 如果用户的 `game_user_id` 为空，直接返回空结果
   - 避免查询到错误的数据

3. **移除不必要的解析逻辑**
   - `game_user_id` 已经是字符串形式的数字，不需要额外解析

## 数据流程

```
用户登录
  ↓
获取用户绑定的游戏账号 (game_account 表)
  ↓
获取 game_user_id (游戏服务器返回的用户ID)
  ↓
使用 game_user_id 查询战绩 (game_battle_record 表)
  ↓
返回战绩列表
```

## 验证修复

修复后，用户应该能够：

1. **查询自己的战绩列表**
   - 调用 `POST /battle-query/my/battles`
   - 返回该用户的所有战绩记录

2. **查询自己的战绩统计**
   - 调用 `POST /battle-query/my/stats`
   - 返回总局数、总成绩、总费用等统计数据

3. **查询自己的余额**
   - 调用 `POST /battle-query/my/balances`
   - 返回各个圈子的余额信息

## 相关文件

- `battle-tiles/internal/biz/game/battle_record.go` - 业务逻辑层
- `battle-tiles/internal/dal/repo/game/battle_record.go` - 数据访问层
- `battle-tiles/internal/dal/model/game/game_battle_record.go` - 数据模型

## 后续改进

1. 实现 `GetMyBalances` 的实际逻辑
2. 添加更多的查询过滤条件
3. 优化查询性能（添加索引等）

