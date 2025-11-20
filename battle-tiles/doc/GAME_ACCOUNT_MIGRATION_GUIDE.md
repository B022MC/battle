# 游戏账号入圈系统 - 迁移指南

## 📋 概述

实现游戏账号入圈机制，用户通过游戏账号反向查询圈子和战绩。

---

## 🗄️ 数据库变更

### 1. 修改现有表

#### game_account
```sql
-- user_id 改为可选
ALTER TABLE game_account ALTER COLUMN user_id DROP NOT NULL;
```

#### game_battle_record
```sql
-- player_id 改为可选
ALTER TABLE game_battle_record ALTER COLUMN player_id DROP NOT NULL;
```

### 2. 新增表

#### game_account_group（游戏账号圈子关系）
```sql
CREATE TABLE game_account_group (
    id SERIAL PRIMARY KEY,
    game_account_id INTEGER NOT NULL,
    house_gid INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    group_name VARCHAR(64) NOT NULL,
    admin_user_id INTEGER NOT NULL,
    approved_by_user_id INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT uk_game_account_house UNIQUE (game_account_id, house_gid)
);
```

---

## 🚀 执行迁移

### 方法 1：直接执行（推荐）

```bash
# 进入项目目录
cd battle-tiles

# 执行迁移脚本
psql -U postgres -d battle_db -f migrations/20251120_game_account_group.sql
```

### 方法 2：分步执行

```bash
# 1. 连接数据库
psql -U postgres -d battle_db

# 2. 执行迁移
\i migrations/20251120_game_account_group.sql

# 3. 查看结果
SELECT * FROM game_account_group LIMIT 5;
```

---

## ✅ 验证迁移

### 1. 检查表是否创建成功
```sql
\d game_account_group
```

### 2. 检查字段是否可为空
```sql
SELECT 
    column_name, 
    is_nullable 
FROM information_schema.columns 
WHERE table_name = 'game_account' 
  AND column_name = 'user_id';
```

### 3. 查看数据统计
```sql
SELECT 
    (SELECT COUNT(*) FROM game_account WHERE is_del = 0) AS 游戏账号总数,
    (SELECT COUNT(*) FROM game_account_group WHERE status = 'active') AS 圈子关系数,
    (SELECT COUNT(*) FROM game_account WHERE user_id IS NOT NULL AND is_del = 0) AS 已绑定用户数,
    (SELECT COUNT(*) FROM game_account WHERE user_id IS NULL AND is_del = 0) AS 未绑定用户数;
```

---

## 🔄 核心流程

### 1. 申请入圈（管理员审批）
```go
// 1. 查找或创建游戏账号
gameAccount := GetOrCreateGameAccount(applierGid, applierGname)

// 2. 获取管理员的圈子
group := GetOrCreateAdminGroup(houseGid, adminUserId)

// 3. 游戏账号入圈
INSERT INTO game_account_group (
    game_account_id, house_gid, group_id, group_name,
    admin_user_id, approved_by_user_id
) VALUES (
    gameAccount.ID, houseGid, group.ID, group.GroupName,
    adminUserId, currentAdminId
)
```

### 2. 战绩同步
```go
// 查询游戏账号和圈子信息（一次查询）
SELECT 
    ga.id, ga.user_id, gag.group_id, gag.group_name
FROM game_account ga
LEFT JOIN game_account_group gag 
    ON ga.id = gag.game_account_id 
    AND gag.house_gid = ?
    AND gag.status = 'active'
WHERE ga.game_user_id = ?

// 如果 group_id 不为空，保存战绩
if groupId != nil {
    INSERT INTO game_battle_record (
        player_id,      // 可为NULL
        player_game_id, // 必填
        group_name,
        ...
    )
}
```

### 3. 用户查询圈子
```go
// 第一步：获取用户的游戏账号
gameAccountIds := SELECT id FROM game_account 
                  WHERE user_id = ? AND is_del = 0

// 第二步：获取圈子
groups := SELECT * FROM game_account_group 
          WHERE game_account_id IN (gameAccountIds)
```

### 4. 用户查询战绩
```go
// 第一步：获取用户的游戏账号
gameAccountIds := SELECT id FROM game_account 
                  WHERE user_id = ? AND is_del = 0

// 第二步：查询战绩
battles := SELECT * FROM game_battle_record 
           WHERE player_game_id IN (gameAccountIds)
           ORDER BY battle_at DESC
```

---

## 📝 API 变更

### 1. 申请审批 API
```
POST /shops/game-applications/approve

Body:
{
  "house_gid": 58959,
  "applier_gid": 22805688,
  "applier_gname": "玩家昵称"
}

变更：
- 创建/查找游戏账号（user_id=NULL）
- 创建 game_account_group 记录
- 记录 approved_by_user_id
```

### 2. 战绩同步 API
```
内部接口

变更：
- 使用 LEFT JOIN 查询游戏账号和圈子
- player_id 可为NULL
- 只保存已入圈的游戏账号战绩
```

### 3. 用户绑定游戏账号 API
```
POST /game/account/bind

Body:
{
  "account": "游戏账号",
  "pwd_md5": "密码MD5"
}

变更：
- 查找游戏账号（可能已存在）
- 更新 user_id 字段
```

### 4. 查询圈子 API
```
GET /my/groups

变更：
- 先查询用户的游戏账号
- 再查询游戏账号的圈子
```

### 5. 查询战绩 API
```
GET /shops/my-battles

变更：
- 先查询用户的游戏账号
- 再查询游戏账号的战绩
```

---

## ⚠️ 注意事项

### 1. 事务处理
- 迁移脚本使用 BEGIN/COMMIT 包裹
- 如果出错会自动回滚

### 2. 数据一致性
- 唯一约束：`(game_account_id, house_gid)`
- 一个游戏账号在一个店铺只能属于一个圈子

### 3. 向后兼容
- 已绑定用户的游戏账号保持绑定
- 已有的圈子关系自动迁移

### 4. 错误处理
- 如果表已存在，会先删除再创建
- 如果字段已经可为空，会跳过修改

---

## 🔧 回滚方案

如果需要回滚：

```sql
BEGIN;

-- 删除新表
DROP TABLE IF EXISTS game_account_group CASCADE;

-- 恢复 NOT NULL 约束（如果需要）
-- ALTER TABLE game_account ALTER COLUMN user_id SET NOT NULL;
-- ALTER TABLE game_battle_record ALTER COLUMN player_id SET NOT NULL;

COMMIT;
```

---

## 📊 性能优化

已创建的索引：
- `idx_account_group_game_account` - 游戏账号查询
- `idx_account_group_house` - 店铺查询
- `idx_account_group_group` - 圈子查询
- `idx_account_group_status` - 状态查询
- `idx_account_group_house_status` - 组合查询

---

## 📚 相关文档

- 设计方案：`GAME_ACCOUNT_REDESIGN_V2.md`
- 简明指南：`GAME_ACCOUNT_SIMPLE_GUIDE.md`
- 迁移脚本：`migrations/20251120_game_account_group.sql`

---

## ✅ 迁移检查清单

- [ ] 备份数据库
- [ ] 执行迁移脚本
- [ ] 验证表结构
- [ ] 检查数据统计
- [ ] 测试申请审批流程
- [ ] 测试战绩同步
- [ ] 测试用户查询
- [ ] 更新后端代码
- [ ] 更新前端代码
- [ ] 生产环境部署

---

**更新时间：** 2025-11-20  
**版本：** 1.0

