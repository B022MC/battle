# 修复：添加成员到圈子时数据不同步问题

## 问题描述

用户被添加到圈子后，无法查询到自己的圈子信息。

### 现象
- 管理员添加用户到圈子后，`game_shop_group_member` 表有记录
- 但用户调用"查询我的圈子"接口时返回空
- `game_member` 和 `game_account_group` 表中没有数据

## 根本原因

### 数据表使用情况分析

**正在使用的表：**
1. **game_shop_group** - 店铺圈子表（基础信息）
2. **game_shop_group_member** - 用户-圈子成员表（写入用）
3. **game_account_group** - 游戏账号-圈子关系表（查询用）⚠️
4. **game_member** - 游戏成员信息表（包含余额、积分等）
5. **game_account_house** - 游戏账号-店铺绑定表
6. **game_shop_admin** - 店铺管理员表

**废弃/未使用的表：**
- game_account_store_binding
- game_house_settings
- game_member_rule
- game_recharge_record（可能在其他模块使用）

### 问题核心

系统在架构演进过程中，添加成员和查询成员使用了不同的表：

- **添加成员时**：只写入 `game_shop_group_member` 表
- **查询圈子时**：从 `game_account_group` 表读取

导致写入和读取不一致，用户看不到自己的圈子。

### 代码路径对比

**添加成员：**
```
POST /api/groups/members/add
  → ShopGroupService.AddMembers
  → ShopGroupUseCase.AddMembersToGroup
  → 仅写入 game_shop_group_member
```

**查询我的圈子：**
```
POST /api/groups/my/list
  → ShopGroupService.ListMyGroups
  → GameAccountGroupUseCase.ListGroupsByUser
  → 从 game_account_group 表读取
```

## 解决方案

修改 `AddMembersToGroup` 方法，添加成员时同步创建所有相关记录：

### 1. 数据表关系

```
用户 (basic_user)
  ↓ user_id
游戏账号 (game_account)
  ↓ game_account_id
├─ game_account_group (游戏账号-圈子关系，用于查询)
├─ game_member (游戏成员信息，包含余额积分)
└─ game_account_house (游戏账号-店铺绑定)

同时维护：
└─ game_shop_group_member (用户-圈子成员关系)
```

### 2. 修改内容

#### 2.1 添加依赖注入

在 `ShopGroupUseCase` 中注入 `GameAccountGroupRepo`：

```go
type ShopGroupUseCase struct {
    // ... 其他依赖
    accountGroupRepo  game.GameAccountGroupRepo   // 用于操作 game_account_group 记录
}
```

#### 2.2 完善 AddMembersToGroup 方法

修改后的流程：

1. **验证权限**：检查圈子是否属于该管理员
2. **检查冲突**：确保用户不在该店铺的其他圈子中
3. **创建 game_shop_group_member 记录**
4. **查询用户的游戏账号**
5. **创建 game_account_group 记录**（游戏账号-圈子关系）
6. **创建 game_member 记录**（游戏成员信息）
7. **创建 game_account_house 记录**（账号-店铺绑定）

#### 2.3 新增 Repo 方法

为了支持创建操作，添加了以下方法：

**GameMemberRepo:**
```go
// 创建或更新成员记录
Create(ctx context.Context, member *model.GameMember) error
```

**GameAccountHouseRepo:**
```go
// 创建账号店铺绑定
Create(ctx context.Context, accountHouse *model.GameAccountHouse) error
```

### 3. 实现特点

- **幂等性**：所有创建操作使用 `FirstOrCreate`，避免重复数据
- **容错性**：如果用户没有游戏账号，只创建 `game_shop_group_member`，不阻塞流程
- **日志记录**：创建失败时记录警告日志，便于问题排查

## 影响范围

### 已修改文件
1. `/internal/biz/game/shop_group.go` - 修改 `AddMembersToGroup` 方法
2. `/internal/dal/repo/game/game_member.go` - 添加 `Create` 方法
3. `/internal/dal/repo/game/game_account_house.go` - 添加 `Create` 方法
4. `/cmd/go-kgin-platform/wire_gen.go` - 自动生成（wire gen）

### 不需要修改
- API 层（Service）：接口保持不变
- 前端代码：无需修改
- 数据库结构：无需迁移

## 验证步骤

1. **重启服务**
   ```bash
   # 重新编译并启动
   go build ./cmd/go-kgin-platform
   ```

2. **测试添加成员**
   ```bash
   POST /api/groups/members/add
   {
     "group_id": 123,
     "user_ids": [456]
   }
   ```

3. **验证数据**
   ```sql
   -- 检查所有相关表是否都有记录
   SELECT * FROM game_shop_group_member WHERE user_id = 456;
   SELECT * FROM game_account_group WHERE game_account_id = (SELECT id FROM game_account WHERE user_id = 456);
   SELECT * FROM game_member WHERE game_id = (SELECT id FROM game_account WHERE user_id = 456);
   SELECT * FROM game_account_house WHERE game_account_id = (SELECT id FROM game_account WHERE user_id = 456);
   ```

4. **测试查询圈子**
   ```bash
   POST /api/groups/my/list
   # 应该能看到用户加入的圈子
   ```

## 历史数据处理

如果有历史数据不完整，可以运行以下 SQL 补全：

```sql
-- 补全 game_account_group 记录
INSERT INTO game_account_group (
    game_account_id, house_gid, group_id, group_name, 
    admin_user_id, approved_by_user_id, status
)
SELECT 
    ga.id AS game_account_id,
    gsg.house_gid,
    gsgm.group_id,
    gsg.group_name,
    gsg.admin_user_id,
    gsg.admin_user_id AS approved_by_user_id,
    'active' AS status
FROM game_shop_group_member gsgm
JOIN game_shop_group gsg ON gsgm.group_id = gsg.id
JOIN game_account ga ON gsgm.user_id = ga.user_id
WHERE ga.is_del = 0
  AND gsg.is_active = true
  AND NOT EXISTS (
      SELECT 1 FROM game_account_group gag
      WHERE gag.game_account_id = ga.id
        AND gag.house_gid = gsg.house_gid
  )
ON CONFLICT (game_account_id, house_gid) DO NOTHING;
```

## 相关文档

- [会员移除清理机制](./CLEANUP_ON_MEMBER_REMOVAL.md)
- [游戏账号入圈系统迁移](../battle-tiles/migrations/20251120_game_account_group.sql)

## 总结

本次修复解决了数据写入和读取不一致的问题，通过在添加成员时同步创建所有必要的关联记录，确保了数据的完整性和一致性。修改向后兼容，不影响现有功能。
