# 最终修复：使用游戏账号查询战绩

## 问题陈述

用户要求：**直接使用用户绑定的游戏账号查询游戏战绩，而不是用游戏账号ID查询**

## 解决方案

### 核心改动

从使用 `game_user_id` 改为使用 `account.Account`（游戏账号）

```
修改前：account.GameUserID → player_game_id
修改后：account.Account → player_game_name
```

### 修改的文件

#### 1. battle-tiles/internal/dal/repo/game/battle_record.go

**新增两个方法**：

```go
// 按游戏账号查询战绩列表
ListByPlayerGameName(ctx context.Context, houseGID int32, playerGameName string, 
    groupID *int32, start, end *time.Time, page, size int32) 
    ([]*model.GameBattleRecord, int64, error)

// 按游戏账号获取统计数据
GetPlayerStatsByGameName(ctx context.Context, houseGID int32, playerGameName string, 
    groupID *int32, start, end *time.Time) 
    (totalGames int64, totalScore int, totalFee int, err error)
```

**查询条件**：
```sql
WHERE house_gid = ? AND player_game_name = ?
```

#### 2. battle-tiles/internal/biz/game/battle_record.go

**修改三个方法**：

| 方法 | 修改内容 |
|------|--------|
| `ListMyBattleRecords` | 使用 `account.Account` 和 `ListByPlayerGameName` |
| `GetMyBattleStats` | 使用 `account.Account` 和 `GetPlayerStatsByGameName` |
| `GetMyBalances` | 检查 `account.Account` 而不是 `game_user_id` |

## 代码示例

### ListMyBattleRecords

```go
// 修改前
if account.GameUserID == "" {
    return []*model.GameBattleRecord{}, 0, nil
}
return uc.repo.ListByPlayer(ctx, houseGID, account.GameUserID, ...)

// 修改后
if account.Account == "" {
    return []*model.GameBattleRecord{}, 0, nil
}
return uc.repo.ListByPlayerGameName(ctx, houseGID, account.Account, ...)
```

### GetMyBattleStats

```go
// 修改前
if account.GameUserID == "" {
    return 0, 0, 0, nil
}
return uc.repo.GetPlayerStats(ctx, houseGID, account.GameUserID, ...)

// 修改后
if account.Account == "" {
    return 0, 0, 0, nil
}
return uc.repo.GetPlayerStatsByGameName(ctx, houseGID, account.Account, ...)
```

## 数据流程

```
用户登录
    ↓
获取 game_account 记录
    ↓
提取 account.Account（游戏账号，如 "11439919"）
    ↓
调用 ListByPlayerGameName 或 GetPlayerStatsByGameName
    ↓
查询 game_battle_record 表
WHERE house_gid = ? AND player_game_name = ?
    ↓
返回战绩列表或统计数据
```

## 优势

| 优势 | 说明 |
|------|------|
| **直接查询** | 不需要依赖 game_user_id 字段 |
| **可靠性** | 游戏账号是用户输入的，不会为空 |
| **兼容性** | 与游戏服务器的账号系统直接对应 |
| **简化** | 减少了数据转换的步骤 |

## 关键字段

| 字段 | 表 | 说明 |
|------|-----|------|
| `account.Account` | game_account | 用户绑定的游戏账号（如 "11439919"） |
| `player_game_name` | game_battle_record | 战绩记录中的玩家游戏账号 |
| `game_user_id` | game_account | 游戏服务器返回的用户ID（不再用于查询） |
| `player_game_id` | game_battle_record | 战绩记录中的玩家游戏ID（不再用于查询） |

## 快速测试

### 1. 编译
```bash
cd battle-tiles
go build -o battle-tiles ./cmd/main.go
```

### 2. 启动
```bash
./battle-tiles
```

### 3. 登录
```bash
curl -X POST http://localhost:8080/basic/login \
  -H "Content-Type: application/json" \
  -d '{"username": "test1", "password": "password"}'
```

### 4. 查询战绩
```bash
curl -X POST http://localhost:8080/battle-query/my/battles \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"house_gid": 58959, "page": 1, "size": 20}'
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

## 相关文档

- `QUERY_BY_GAME_ACCOUNT.md` - 详细说明
- `GAME_ACCOUNT_QUERY_QUICK_REF.md` - 快速参考

## 总结

✅ 已完成使用游戏账号查询战绩的修复
✅ 新增两个数据库查询方法
✅ 修改三个业务逻辑方法
✅ 提供完整的文档和测试指南

