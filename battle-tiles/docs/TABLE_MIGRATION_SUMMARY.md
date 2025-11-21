# game_account_group 表迁移总结

## 问题根源

系统中存在两个圈子成员相关的表，导致数据不一致和查询错误：

1. **game_member** - 实际在用的表，有数据
2. **game_account_group** - 重构时引入的表，但一直是空的

多处代码错误地查询了空表 `game_account_group`，导致功能异常。

## 受影响的接口和功能

### 1. 战绩同步 ❌→✅
**文件**: `internal/biz/game/battle_record.go`

```go
// ❌ 错误：查询空表
accountGroup, err := uc.accountGroupRepo.GetActiveByGamePlayerAndHouse(ctx, gamePlayerID, houseGID)

// ✅ 修复：查询 game_member
member, err := uc.memberRepo.GetByGameID(ctx, houseGID, int32(player.UserGameID))
if member.GroupID == nil || *member.GroupID == 0 {
    continue // 没有圈子，跳过
}
```

**影响**: 所有战绩都被跳过，无法计费

---

### 2. GET /api/groups/my ❌→✅
**文件**: `internal/dal/repo/game/shop_group_member.go` - `CountMembers`

```go
// ❌ 错误：统计空表
Table("game_account_group").Where("group_id = ? AND status = ?", groupID, "active")

// ✅ 修复：统计 game_member
Table("game_member").Where("group_id = ?", groupID)
```

**影响**: 圈子成员数量始终显示为 0

---

### 3. POST /api/groups/members/list ❌→✅
**文件**: `internal/dal/repo/game/shop_group_member.go` - `ListMembersByGroup`

```go
// ❌ 错误：从空表查询
Table("game_account_group gag").
    Joins("JOIN game_account ga ON ga.id = gag.game_account_id").
    Joins("JOIN basic_user u ON u.id = ga.user_id")

// ✅ 修复：从 game_member 查询
Table("game_member gm").
    Joins("JOIN game_account ga ON CAST(gm.game_id AS VARCHAR) = ga.game_player_id").
    Joins("JOIN basic_user u ON u.id = ga.user_id")
```

**影响**: 圈子成员列表为空

---

### 4. POST /api/groups/my/list ❌→✅
**文件**: `internal/dal/repo/game/shop_group_member.go` - `ListGroupsByUser`

```go
// ❌ 错误
Table("game_account ga").
    Joins("JOIN game_account_group gag ON ga.id = gag.game_account_id").
    Joins("JOIN game_shop_group g ON g.id = gag.group_id")

// ✅ 修复
Table("game_account ga").
    Joins("JOIN game_member gm ON CAST(gm.game_id AS VARCHAR) = ga.game_player_id").
    Joins("JOIN game_shop_group g ON g.id = gm.group_id")
```

**影响**: 用户的圈子列表为空

---

### 5. ListGroupsByUserAndHouse ❌→✅
**文件**: `internal/dal/repo/game/shop_group_member.go` - `ListGroupsByUserAndHouse`

同上，查询逻辑修复为使用 `game_member` 表。

---

## 表结构对比

### game_member（正确的表）✅

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int4 | 主键 |
| house_gid | int4 | 店铺GID |
| **game_id** | **int4** | **游戏玩家ID（直接匹配战绩）** |
| game_name | varchar(64) | 游戏昵称 |
| **group_id** | **int4** | **圈子ID（用于计费）** |
| group_name | varchar(64) | 圈子名称 |
| balance | int4 | 余额 |
| credit | int4 | 信用额度 |

**数据状态**: ✅ 有数据
```sql
SELECT COUNT(*) FROM game_member WHERE house_gid = 60870;  -- 1 row
```

### game_account_group（废弃的表）❌

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int4 | 主键 |
| game_player_id | varchar(32) | 游戏玩家ID（新添加） |
| house_gid | int4 | 店铺GID |
| group_id | int4 | 圈子ID |
| status | varchar(20) | 状态 |

**数据状态**: ❌ 空表
```sql
SELECT COUNT(*) FROM game_account_group;  -- 0 rows
```

---

## 关联关系

### 新的查询路径

```
战绩同步:
game_player_id (战绩) → game_member.game_id → game_member.group_id

成员列表查询:
basic_user.id → game_account.user_id → game_account.game_player_id 
             → game_member.game_id → game_member.group_id → game_shop_group

用户圈子查询:
basic_user.id → game_account.user_id → game_account.game_player_id 
             → game_member.game_id → game_member.group_id → game_shop_group
```

### 关键JOIN条件

```sql
-- game_member 通过 game_id 关联 game_account
JOIN game_account ga ON CAST(gm.game_id AS VARCHAR) = ga.game_player_id

-- game_account 关联 basic_user
JOIN basic_user u ON u.id = ga.user_id

-- game_member 关联 game_shop_group
JOIN game_shop_group g ON g.id = gm.group_id
```

**注意**: 需要将 `game_id` (int4) 转换为 VARCHAR 才能与 `game_player_id` (varchar) 关联。

---

## 修复文件清单

| 文件 | 修改内容 |
|------|----------|
| `internal/biz/game/battle_record.go` | 战绩同步改为查询 `game_member` |
| `internal/dal/repo/game/shop_group_member.go` | 5个方法全部改为查询 `game_member` |
| `migrations/20251122_add_game_player_id_to_account_group.sql` | 临时添加字段避免报错 |

---

## 数据验证

### 验证成员数据
```sql
-- 圈子 6 的成员
SELECT gm.game_id, gm.game_name, gm.group_id, gm.group_name,
       ga.game_player_id, ga.user_id, u.username
FROM game_member gm
JOIN game_account ga ON CAST(gm.game_id AS VARCHAR) = ga.game_player_id
LEFT JOIN basic_user u ON u.id = ga.user_id
WHERE gm.group_id = 6;
```

**结果示例**:
| game_id | game_name | group_id | user_id | username |
|---------|-----------|----------|---------|----------|
| 22953243 | mc | 6 | 7 | test110 |

### 验证成员数量
```sql
SELECT group_id, COUNT(*) as member_count 
FROM game_member 
WHERE house_gid = 60870 AND group_id IS NOT NULL 
GROUP BY group_id;
```

---

## 部署清单

- [x] 1. 添加 `game_player_id` 字段到 `game_account_group` 表
- [x] 2. 修改战绩同步逻辑
- [x] 3. 修改成员统计逻辑
- [x] 4. 修改成员列表查询逻辑
- [x] 5. 修改用户圈子查询逻辑
- [x] 6. 重新编译部署
- [x] 7. 验证服务运行
- [ ] 8. 前端测试各个接口
- [ ] 9. 观察战绩同步日志

---

## 后续建议

### 短期
1. **测试验证**: 测试所有修改的接口，确保功能正常
2. **监控日志**: 观察战绩同步是否正常计费
3. **数据一致性**: 确认 `game_member` 表的数据完整性

### 中期
1. **删除废弃表**: 考虑删除或重命名 `game_account_group` 表
2. **优化JOIN**: 评估是否需要在 `game_member` 表添加 `game_player_id` (varchar) 字段避免 CAST
3. **文档完善**: 更新数据库设计文档，明确表的用途

### 长期
1. **统一数据模型**: 确定 `game_member` 为主表，所有逻辑都基于此表
2. **数据迁移计划**: 如果确实需要 `game_account_group`，制定数据同步方案
3. **代码审查**: 排查是否还有其他地方使用了错误的表

---

## 测试用例

### 1. 测试圈子成员数量
```bash
curl -X POST http://192.168.31.56:8000/groups/my \
  -H "Content-Type: application/json" \
  -d '{"house_gid": 60870}'

# 期望: member_count = 1 (之前是 0)
```

### 2. 测试成员列表
```bash
curl -X POST http://192.168.31.56:8000/groups/members/list \
  -H "Content-Type: application/json" \
  -d '{"group_id": 6, "page": 1, "size": 20}'

# 期望: items 数组有 1 个用户 (之前是空数组)
```

### 3. 测试战绩同步
等待几分钟，查看日志：
```bash
tail -f logs/battle-tiles.log | grep -E "(Mapped|valid players)"

# 期望: "Mapped N valid players" N > 0 (之前是 0)
```

---

**修复日期**: 2025-11-22  
**修复人员**: AI Assistant  
**状态**: ✅ 已部署，待测试验证
