# 使用游戏账号查询战绩

## 问题描述

之前的实现是使用 `game_user_id`（游戏用户ID）来查询战绩，但实际上应该直接使用 `account.Account`（用户绑定的游戏账号）来查询。

## 修复方案

### 1. 数据库层修改

**文件**：`battle-tiles/internal/dal/repo/game/battle_record.go`

#### 新增接口方法

```go
// 按游戏账号查询战绩列表
ListByPlayerGameName(ctx context.Context, houseGID int32, playerGameName string, groupID *int32, start, end *time.Time, page, size int32) ([]*model.GameBattleRecord, int64, error)

// 按游戏账号获取统计数据
GetPlayerStatsByGameName(ctx context.Context, houseGID int32, playerGameName string, groupID *int32, start, end *time.Time) (totalGames int64, totalScore int, totalFee int, err error)
```

#### 实现细节

这两个方法与原有的 `ListByPlayer` 和 `GetPlayerStats` 类似，但查询条件改为：

```go
// 原来：按 player_game_id 查询
Where("house_gid = ? AND player_game_id = ?", houseGID, playerGameID)

// 现在：按 player_game_name 查询
Where("house_gid = ? AND player_game_name = ?", houseGID, playerGameName)
```

### 2. 业务逻辑层修改

**文件**：`battle-tiles/internal/biz/game/battle_record.go`

#### ListMyBattleRecords

```go
// 修改前：使用 game_user_id
if account.GameUserID == "" {
    return []*model.GameBattleRecord{}, 0, nil
}
return uc.repo.ListByPlayer(ctx, houseGID, account.GameUserID, ...)

// 修改后：使用游戏账号
if account.Account == "" {
    return []*model.GameBattleRecord{}, 0, nil
}
return uc.repo.ListByPlayerGameName(ctx, houseGID, account.Account, ...)
```

#### GetMyBattleStats

```go
// 修改前：使用 game_user_id
if account.GameUserID == "" {
    return 0, 0, 0, nil
}
return uc.repo.GetPlayerStats(ctx, houseGID, account.GameUserID, ...)

// 修改后：使用游戏账号
if account.Account == "" {
    return 0, 0, 0, nil
}
return uc.repo.GetPlayerStatsByGameName(ctx, houseGID, account.Account, ...)
```

#### GetMyBalances

```go
// 修改前：检查 game_user_id
if account.GameUserID == "" {
    return []interface{}{}, nil
}

// 修改后：检查游戏账号
if account.Account == "" {
    return []interface{}{}, nil
}
```

## 数据流程

```
用户登录
    ↓
获取 game_account 记录
    ↓
提取 account.Account（游戏账号，如 "11439919"）
    ↓
查询 game_battle_record 表
WHERE house_gid = ? AND player_game_name = ?
    ↓
返回战绩列表
```

## 关键字段说明

| 字段 | 表 | 说明 |
|------|-----|------|
| `account.Account` | game_account | 用户绑定的游戏账号（如 "11439919"） |
| `player_game_name` | game_battle_record | 战绩记录中的玩家游戏账号 |
| `game_user_id` | game_account | 游戏服务器返回的用户ID（不再用于查询） |
| `player_game_id` | game_battle_record | 战绩记录中的玩家游戏ID（不再用于查询） |

## 修改的文件

1. `battle-tiles/internal/dal/repo/game/battle_record.go`
   - 添加 `ListByPlayerGameName` 方法
   - 添加 `GetPlayerStatsByGameName` 方法

2. `battle-tiles/internal/biz/game/battle_record.go`
   - 修改 `ListMyBattleRecords` 方法
   - 修改 `GetMyBattleStats` 方法
   - 修改 `GetMyBalances` 方法

## 测试步骤

### 1. 编译
```bash
cd battle-tiles
go build -o battle-tiles ./cmd/main.go
```

### 2. 启动服务
```bash
./battle-tiles
```

### 3. 用户登录
```bash
curl -X POST http://localhost:8080/basic/login \
  -H "Content-Type: application/json" \
  -d '{"username": "test1", "password": "password"}'
```

### 4. 查询战绩
```bash
curl -X POST http://localhost:8080/battle-query/my/battles \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"house_gid": 58959, "page": 1, "size": 20}'
```

### 5. 查询统计
```bash
curl -X POST http://localhost:8080/battle-query/my/stats \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"house_gid": 58959}'
```

## 验证

### 数据库验证

```sql
-- 查看用户的游戏账号
SELECT id, user_id, account, game_user_id 
FROM game_account 
WHERE account != '' 
LIMIT 10;

-- 查看战绩数据（按游戏账号）
SELECT * FROM game_battle_record 
WHERE player_game_name = '11439919' 
LIMIT 10;
```

## 优势

1. **直接查询**：不需要依赖 `game_user_id` 字段
2. **兼容性**：与游戏服务器的账号系统直接对应
3. **可靠性**：游戏账号是用户输入的，不会为空
4. **简化**：减少了数据转换的步骤

## 注意事项

- 确保 `game_account.account` 字段不为空
- 确保 `game_battle_record.player_game_name` 字段与 `game_account.account` 一致
- 如果历史数据中 `player_game_name` 为空，需要进行数据迁移

