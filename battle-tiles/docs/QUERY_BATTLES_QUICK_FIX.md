# 查询战绩快速修复指南

## 一句话总结

修复了查询战绩时使用错误参数的问题：从使用 `account.Account`（游戏账号）改为使用 `account.GameUserID`（游戏用户ID）。

## 问题症状

用户登录后，调用以下接口返回空结果或错误：
- `POST /battle-query/my/battles` - 查询我的战绩
- `POST /battle-query/my/stats` - 查询我的战绩统计
- `POST /battle-query/my/balances` - 查询我的余额

## 修复内容

### 修改文件
`battle-tiles/internal/biz/game/battle_record.go`

### 修改的方法

#### 1. ListMyBattleRecords（第 101-124 行）
- **问题**：使用 `account.Account` 查询战绩
- **修复**：改为使用 `account.GameUserID`
- **添加**：检查 `game_user_id` 是否为空

#### 2. GetMyBattleStats（第 126-149 行）
- **问题**：使用 `account.Account` 查询统计
- **修复**：改为使用 `account.GameUserID`
- **添加**：检查 `game_user_id` 是否为空

#### 3. GetMyBalances（第 179-202 行）
- **问题**：不必要的解析逻辑
- **修复**：简化逻辑，直接检查 `game_user_id` 是否为空
- **移除**：`parseGameUserID` 函数调用

## 关键概念

### 游戏账号 vs 游戏用户ID

| 字段 | 说明 | 示例 | 存储位置 |
|------|------|------|---------|
| `account.Account` | 用户输入的游戏账号 | "11439919" 或 "13800138000" | `game_account.account` |
| `account.GameUserID` | 游戏服务器返回的用户ID | "12345" | `game_account.game_user_id` |
| `player_game_id` | 战绩中存储的玩家ID | "12345" | `game_battle_record.player_game_id` |

### 数据库查询

战绩表 `game_battle_record` 中的查询条件：
```sql
WHERE house_gid = ? AND player_game_id = ?
```

所以必须使用 `game_user_id` 而不是 `account`。

## 测试步骤

### 1. 编译并启动服务
```bash
cd battle-tiles
go build -o battle-tiles ./cmd/main.go
./battle-tiles
```

### 2. 用户登录
```bash
curl -X POST http://localhost:8080/basic/login \
  -H "Content-Type: application/json" \
  -d '{"username": "test1", "password": "password"}'
```

获取 `token`。

### 3. 查询战绩
```bash
curl -X POST http://localhost:8080/battle-query/my/battles \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "house_gid": 58959,
    "page": 1,
    "size": 20
  }'
```

**预期结果**：返回该用户的战绩列表

### 4. 查询统计
```bash
curl -X POST http://localhost:8080/battle-query/my/stats \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"house_gid": 58959}'
```

**预期结果**：返回统计数据（总局数、总成绩、总费用）

## 验证修复

### 检查数据库

```sql
-- 查看用户的 game_user_id
SELECT id, user_id, account, game_user_id 
FROM game_account 
WHERE user_id = 4;

-- 查看该用户的战绩
SELECT * FROM game_battle_record 
WHERE player_game_id = '12345' 
LIMIT 10;
```

### 查看日志

修复后的日志应该显示：
```
[INFO] User 4 successfully queried battles: total=10
```

而不是：
```
[WARN] User 4 has no game_user_id
```

## 相关文档

- 详细文档：`battle-tiles/docs/BATTLE_QUERY_FIX.md`
- 数据修复：`battle-tiles/docs/FIX_EMPTY_GAME_USER_ID.md`
- 代码修复：`battle-tiles/internal/biz/game/battle_record.go`

