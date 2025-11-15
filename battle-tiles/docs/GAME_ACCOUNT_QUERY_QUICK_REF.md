# 使用游戏账号查询战绩 - 快速参考

## 修改总结

| 项目 | 修改前 | 修改后 |
|------|--------|--------|
| **查询字段** | `account.GameUserID` | `account.Account` |
| **数据库字段** | `player_game_id` | `player_game_name` |
| **方法名** | `ListByPlayer` | `ListByPlayerGameName` |
| **方法名** | `GetPlayerStats` | `GetPlayerStatsByGameName` |
| **优势** | 依赖 game_user_id | 直接使用游戏账号 |

## 代码对比

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

### GetMyBalances

```go
// 修改前
if account.GameUserID == "" {
    return []interface{}{}, nil
}

// 修改后
if account.Account == "" {
    return []interface{}{}, nil
}
```

## 新增方法

### ListByPlayerGameName

```go
func (r *battleRecordRepo) ListByPlayerGameName(
    ctx context.Context, 
    houseGID int32, 
    playerGameName string, 
    groupID *int32, 
    start, end *time.Time, 
    page, size int32
) ([]*model.GameBattleRecord, int64, error)
```

**查询条件**：
```sql
WHERE house_gid = ? AND player_game_name = ?
```

### GetPlayerStatsByGameName

```go
func (r *battleRecordRepo) GetPlayerStatsByGameName(
    ctx context.Context, 
    houseGID int32, 
    playerGameName string, 
    groupID *int32, 
    start, end *time.Time
) (totalGames int64, totalScore int, totalFee int, err error)
```

**查询条件**：
```sql
WHERE house_gid = ? AND player_game_name = ?
```

## 修改的文件

### 1. battle-tiles/internal/dal/repo/game/battle_record.go

**新增方法**：
- `ListByPlayerGameName` - 按游戏账号查询战绩列表
- `GetPlayerStatsByGameName` - 按游戏账号获取统计数据

**接口更新**：
```go
type BattleRecordRepo interface {
    // ... 其他方法 ...
    ListByPlayerGameName(...) ([]*model.GameBattleRecord, int64, error)
    GetPlayerStatsByGameName(...) (int64, int, int, error)
}
```

### 2. battle-tiles/internal/biz/game/battle_record.go

**修改方法**：
- `ListMyBattleRecords` - 使用 `account.Account` 和 `ListByPlayerGameName`
- `GetMyBattleStats` - 使用 `account.Account` 和 `GetPlayerStatsByGameName`
- `GetMyBalances` - 检查 `account.Account` 而不是 `game_user_id`

## 测试命令

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

### 5. 查询统计
```bash
curl -X POST http://localhost:8080/battle-query/my/stats \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"house_gid": 58959}'
```

## 数据库验证

### 查看游戏账号
```sql
SELECT id, user_id, account, game_user_id 
FROM game_account 
WHERE account != '' 
LIMIT 10;
```

### 查看战绩（按游戏账号）
```sql
SELECT * FROM game_battle_record 
WHERE player_game_name = '11439919' 
LIMIT 10;
```

## 关键点

✅ **使用 `account.Account`** - 用户绑定的游戏账号
✅ **查询 `player_game_name`** - 战绩表中的游戏账号字段
✅ **不依赖 `game_user_id`** - 避免空值问题
✅ **直接对应** - 游戏账号与战绩直接匹配

## 常见问题

**Q: 为什么改用游戏账号而不是 game_user_id？**
A: 因为游戏账号是用户输入的，不会为空，而 game_user_id 可能为空或不匹配。

**Q: 旧的方法还能用吗？**
A: 可以，`ListByPlayer` 和 `GetPlayerStats` 方法仍然存在，但不再被使用。

**Q: 需要修改数据库吗？**
A: 不需要，只需确保 `player_game_name` 字段有数据。

**Q: 如何验证修改是否成功？**
A: 调用查询接口，应该能返回战绩列表。

