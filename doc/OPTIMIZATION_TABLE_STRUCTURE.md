# 表结构优化：统一使用 game_account_group 作为主表

## 优化日期
2025-11-21

## 优化目标

简化圈子成员管理的表结构，解决数据冗余和同步问题，统一使用 `game_account_group` 作为主表。

---

## 业务规则

1. **游戏账号可以在未绑定用户时入圈**
   - 场景：游戏玩家申请加入店铺 → 管理员审批 → 游戏账号自动加入管理员的圈子
   - 此时游戏账号尚未绑定用户

2. **用户绑定游戏账号后自动在圈内**
   - 用户注册并绑定游戏账号
   - 自动继承游戏账号的圈子关系

3. **一个用户在一个店铺只能加入一个圈子**
   - 通过 `game_account_group` 的唯一约束保证

---

## 优化前的表结构

### 存在的问题

```
写入路径（添加成员）：
- game_shop_group_member (用户维度)
- game_account_group (游戏账号维度)
- game_member (业务数据)
- game_account_house (店铺绑定)

查询路径（查询我的圈子）：
- 从 game_account_group 查询

问题：数据冗余，写入和查询不一致
```

### 冗余表分析

| 表名 | 作用 | 是否必需 | 说明 |
|------|------|----------|------|
| `game_account_group` | 游戏账号-圈子关系 | ✅ 必需 | 支持未绑定用户的游戏账号入圈 |
| `game_member` | 游戏成员业务数据 | ✅ 必需 | 存储余额、积分、禁用状态等 |
| `game_shop_group_member` | 用户-圈子关系 | ❌ 冗余 | 可通过 game_account 反查 |
| `game_account_house` | 游戏账号-店铺绑定 | ❌ 冗余 | house_gid 已在 game_account_group 中 |

---

## 优化后的表结构

### 核心表

#### 1. game_account_group（主表）

```sql
CREATE TABLE game_account_group (
    id SERIAL PRIMARY KEY,
    game_account_id INTEGER NOT NULL,
    house_gid INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    group_name VARCHAR(64) NOT NULL DEFAULT '',
    admin_user_id INTEGER NOT NULL,
    approved_by_user_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_game_account_house UNIQUE (game_account_id, house_gid)
);
```

**作用：**
- 游戏账号-圈子关系（主表）
- 支持未绑定用户的游戏账号入圈
- 包含 `house_gid`，无需单独的 `game_account_house` 表

#### 2. game_member（业务数据表）

```sql
CREATE TABLE game_member (
    id SERIAL PRIMARY KEY,
    house_gid INTEGER NOT NULL,
    game_id INTEGER NOT NULL,  -- 实际是 game_account_id
    game_name VARCHAR(64) NOT NULL DEFAULT '',
    group_id INTEGER,
    group_name VARCHAR(64) NOT NULL DEFAULT '',
    balance INTEGER NOT NULL DEFAULT 0,
    credit INTEGER NOT NULL DEFAULT 0,
    forbid BOOLEAN NOT NULL DEFAULT FALSE,
    recommender VARCHAR(64) DEFAULT '',
    use_multi_gids BOOLEAN DEFAULT FALSE,
    active_gid INTEGER,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_game_member_house_game_group UNIQUE (house_gid, game_id, group_id)
);
```

**作用：**
- 存储游戏成员的业务数据
- 余额、积分、禁用状态等

### 废弃表

#### game_shop_group_member（已废弃写入）

- **原作用**：用户-圈子成员关系
- **为什么废弃**：可以通过 `game_account.user_id` + `game_account_group` 反查
- **迁移策略**：
  - 不再写入新数据
  - 历史数据保留（不删除表）
  - 所有查询改为从 `game_account_group` 读取

#### game_account_house（已废弃写入）

- **原作用**：游戏账号-店铺绑定
- **为什么废弃**：`house_gid` 已在 `game_account_group` 中
- **迁移策略**：
  - 不再写入新数据
  - 历史数据保留
  - 查询店铺时从 `game_account_group` 获取

---

## 代码修改

### 1. 添加成员（AddMembersToGroup）

**修改前：**
```go
// 写入 4 个表
1. game_shop_group_member
2. game_account_group
3. game_member
4. game_account_house
```

**修改后：**
```go
// 只写入 2 个表
1. game_account_group (主表)
2. game_member (业务数据)
```

### 2. 移除成员（RemoveMemberFromGroup）

**修改前：**
```go
// 删除 3 个表
1. game_shop_group_member
2. game_member
3. game_account_house
```

**修改后：**
```go
// 只删除 2 个表
1. game_account_group (主表)
2. game_member (业务数据)
```

### 3. 查询用户的圈子

**修改前：**
```go
// 从 game_shop_group_member 查询
SELECT g.* 
FROM game_shop_group_member gm
JOIN game_shop_group g ON g.id = gm.group_id
WHERE gm.user_id = ? AND g.is_active = ?
```

**修改后：**
```go
// 从 game_account_group 查询（通过 game_account 关联）
SELECT g.* 
FROM game_account ga
JOIN game_account_group gag ON ga.id = gag.game_account_id
JOIN game_shop_group g ON g.id = gag.group_id
WHERE ga.user_id = ? AND gag.status = 'active' AND g.is_active = ?
```

---

## 数据流转

### 场景 1：游戏账号申请入圈（未绑定用户）

```sql
-- 1. 创建游戏账号（如果不存在）
INSERT INTO game_account (game_player_id, account, nickname, ...)
VALUES ('游戏玩家ID', '账号', '昵称', ...);

-- 2. 加入圈子
INSERT INTO game_account_group (game_account_id, group_id, house_gid, ...)
VALUES (游戏账号ID, 管理员圈子ID, 店铺ID, ...);

-- 3. 创建业务数据
INSERT INTO game_member (house_gid, game_id, group_id, ...)
VALUES (店铺ID, 游戏账号ID, 圈子ID, ...);
```

### 场景 2：用户绑定游戏账号

```sql
-- 更新游戏账号的 user_id
UPDATE game_account 
SET user_id = 用户ID 
WHERE id = 游戏账号ID;

-- 用户自动继承游戏账号的圈子关系（无需额外操作）
```

### 场景 3：管理员添加用户到圈子

```go
// 1. 查询用户的游戏账号
gameAccount := GetOneByUser(userID)

// 2. 创建 game_account_group 记录
INSERT INTO game_account_group (game_account_id, group_id, house_gid, ...)

// 3. 创建 game_member 记录
INSERT INTO game_member (house_gid, game_id, group_id, ...)
```

### 场景 4：查询用户的圈子

```sql
SELECT gag.*, g.group_name, g.description
FROM game_account ga
JOIN game_account_group gag ON ga.id = gag.game_account_id
JOIN game_shop_group g ON g.id = gag.group_id
WHERE ga.user_id = ? 
  AND gag.status = 'active' 
  AND g.is_active = true
ORDER BY gag.joined_at DESC;
```

---

## 修改的文件

### 核心业务逻辑
1. `/internal/biz/game/shop_group.go`
   - 修改 `AddMembersToGroup`：移除 `game_shop_group_member` 和 `game_account_house` 的写入
   - 修改 `RemoveMemberFromGroup`：移除 `game_shop_group_member` 和 `game_account_house` 的删除
   - 移除 `accountHouseRepo` 依赖

2. `/internal/biz/game/game_account.go`
   - 修改 `GetMyHouses`：从 `game_account_group` 查询店铺信息，不再使用 `game_account_house`
   - 添加 `accountGroupRepo` 依赖

### 数据访问层
3. `/internal/dal/repo/game/shop_group_member.go`
   - 修改 `ListGroupsByUser`：从 `game_account_group` 查询
   - 修改 `ListGroupsByUserAndHouse`：从 `game_account_group` 查询
   - 修改 `IsMember`：从 `game_account_group` 查询

### 依赖注入
4. `/cmd/go-kgin-platform/wire_gen.go`
   - 自动生成（`wire gen`）

---

## 优化效果

### 数据一致性
✅ **写入和查询使用同一张表**（`game_account_group`）
✅ **消除数据冗余和同步问题**

### 性能优化
✅ **减少写入操作**：从 4 张表减少到 2 张表
✅ **减少删除操作**：从 3 张表减少到 2 张表
✅ **查询性能不变**：通过 JOIN 查询，有索引支持

### 代码简洁性
✅ **减少依赖**：移除 `accountHouseRepo`
✅ **逻辑清晰**：单一数据源，易于维护

---

## 验证步骤

### 1. 编译测试
```bash
cd /Users/b022mc/project/battle/battle-tiles
go build ./cmd/go-kgin-platform
```

### 2. 功能测试

#### 测试添加成员
```bash
POST /api/groups/members/add
{
  "group_id": 123,
  "user_ids": [456]
}
```

验证数据：
```sql
-- 应该有记录
SELECT * FROM game_account_group WHERE game_account_id = (SELECT id FROM game_account WHERE user_id = 456);
SELECT * FROM game_member WHERE game_id = (SELECT id FROM game_account WHERE user_id = 456);

-- 不应该有新记录
SELECT * FROM game_shop_group_member WHERE user_id = 456 AND created_at > NOW() - INTERVAL '1 minute';
SELECT * FROM game_account_house WHERE game_account_id = (SELECT id FROM game_account WHERE user_id = 456) AND created_at > NOW() - INTERVAL '1 minute';
```

#### 测试查询圈子
```bash
POST /api/groups/my/list
```

应该能正常返回用户加入的圈子。

#### 测试移除成员
```bash
POST /api/groups/members/remove
{
  "group_id": 123,
  "user_id": 456
}
```

验证数据：
```sql
-- 应该被删除
SELECT * FROM game_account_group WHERE game_account_id = (SELECT id FROM game_account WHERE user_id = 456);
SELECT * FROM game_member WHERE game_id = (SELECT id FROM game_account WHERE user_id = 456);
```

---

## 历史数据处理

### 不需要迁移

由于：
1. 新代码已经从 `game_account_group` 读取
2. 历史数据在 `game_account_group` 中已经存在（之前的修复已经同步写入）
3. `game_shop_group_member` 和 `game_account_house` 的历史数据保留，不影响功能

### 可选清理（未来）

如果确认系统稳定运行一段时间后，可以考虑：

```sql
-- 备份表
CREATE TABLE game_shop_group_member_backup AS SELECT * FROM game_shop_group_member;
CREATE TABLE game_account_house_backup AS SELECT * FROM game_account_house;

-- 删除表（谨慎操作）
-- DROP TABLE game_shop_group_member;
-- DROP TABLE game_account_house;
```

---

## 相关文档

- [添加成员数据同步修复](./FIX_ADD_MEMBER_TO_GROUP.md)
- [会员移除清理机制](./CLEANUP_ON_MEMBER_REMOVAL.md)
- [游戏账号入圈系统迁移](../battle-tiles/migrations/20251120_game_account_group.sql)

---

## 总结

本次优化通过统一使用 `game_account_group` 作为主表，成功：

1. ✅ **消除数据冗余**：减少了 2 张冗余表的写入
2. ✅ **提高数据一致性**：写入和查询使用同一数据源
3. ✅ **简化代码逻辑**：减少依赖，易于维护
4. ✅ **保持向后兼容**：API 接口无需修改
5. ✅ **支持业务需求**：完全支持"游戏账号未绑定用户时入圈"的场景

**优化完成！** 🎉
