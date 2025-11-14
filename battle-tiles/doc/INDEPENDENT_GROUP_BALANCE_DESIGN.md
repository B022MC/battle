# 独立圈子余额设计文档

## 📋 概述

本文档描述了 Battle Tiles 系统中"用户在不同圈子有独立余额"功能的设计和实现。

### 设计目标

- **独立余额**: 同一个用户在同一个店铺的不同圈子中,拥有独立的余额
- **数据隔离**: 不同圈子的战绩和资金互不影响
- **向后兼容**: 支持现有数据的平滑迁移

### 业务场景示例

```
用户: 张三 (game_id: 12345)
店铺: 欢乐茶馆 (house_gid: 100)

圈子A (VIP圈):
  - 余额: 1000元
  - 战绩: 50局,总得分 +500

圈子B (普通圈):
  - 余额: 500元
  - 战绩: 30局,总得分 -200

说明: 张三在VIP圈和普通圈的余额、战绩完全独立
```

## 🗄️ 数据模型变更

### 1. game_member 表变更

#### 变更前

```sql
CREATE TABLE game_member (
    id SERIAL PRIMARY KEY,
    house_gid INT NOT NULL,
    game_id INT NOT NULL,
    game_name VARCHAR(64) NOT NULL,
    group_name VARCHAR(64) NOT NULL DEFAULT '',
    balance INT NOT NULL DEFAULT 0,
    ...
    CONSTRAINT uk_game_member_house_game UNIQUE (house_gid, game_id)
);
```

**问题**: 一个用户(game_id)在一个店铺(house_gid)只能有一条记录,无法支持多圈子独立余额。

#### 变更后

```sql
CREATE TABLE game_member (
    id SERIAL PRIMARY KEY,
    house_gid INT NOT NULL,
    game_id INT NOT NULL,
    game_name VARCHAR(64) NOT NULL,
    group_id INT,  -- 新增字段
    group_name VARCHAR(64) NOT NULL DEFAULT '',
    balance INT NOT NULL DEFAULT 0,
    ...
    CONSTRAINT uk_game_member_house_game_group UNIQUE (house_gid, game_id, group_id)
);

CREATE INDEX idx_game_member_group_id ON game_member(group_id);
```

**改进**: 
- 增加 `group_id` 字段,关联 `game_shop_group.id`
- 唯一索引改为 `(house_gid, game_id, group_id)`,支持同一用户在不同圈子有多条记录
- 保留 `group_name` 字段用于冗余显示

### 2. game_member_wallet 表变更

#### 变更前

```sql
CREATE TABLE game_member_wallet (
    id SERIAL PRIMARY KEY,
    house_gid INT NOT NULL,
    member_id INT NOT NULL,
    balance INT NOT NULL DEFAULT 0,
    ...
);
```

**问题**: 钱包与 member_id 绑定,无法支持同一成员在不同圈子的独立钱包。

#### 变更后

```sql
CREATE TABLE game_member_wallet (
    id SERIAL PRIMARY KEY,
    house_gid INT NOT NULL,
    member_id INT NOT NULL,
    group_id INT,  -- 新增字段
    balance INT NOT NULL DEFAULT 0,
    ...
    CONSTRAINT uk_game_member_wallet_house_member_group UNIQUE (house_gid, member_id, group_id)
);

CREATE INDEX idx_game_member_wallet_group_id ON game_member_wallet(group_id);
```

**改进**:
- 增加 `group_id` 字段
- 唯一索引改为 `(house_gid, member_id, group_id)`,支持同一成员在不同圈子有独立钱包

### 3. game_battle_record 表

**无需变更**: 该表已经有 `group_id` 和 `group_name` 字段,天然支持按圈子记录战绩。

```sql
CREATE TABLE game_battle_record (
    id SERIAL PRIMARY KEY,
    house_gid INT NOT NULL,
    group_id INT NOT NULL,  -- 已有字段
    group_name VARCHAR(64),  -- 已有字段
    player_id INT,
    player_game_id INT,
    score INT,
    ...
);
```

## 📊 数据关系图

```
game_shop_group (圈子表)
    ├── id (圈子ID)
    ├── house_gid (店铺ID)
    └── group_name (圈子名称)
         │
         ├─── game_member (成员表)
         │     ├── id (成员ID)
         │     ├── house_gid (店铺ID)
         │     ├── game_id (游戏ID)
         │     ├── group_id (圈子ID) ← 新增
         │     └── balance (余额)
         │          │
         │          └─── game_member_wallet (钱包表)
         │                ├── member_id (成员ID)
         │                ├── group_id (圈子ID) ← 新增
         │                └── balance (余额)
         │
         └─── game_battle_record (战绩表)
               ├── group_id (圈子ID)
               ├── player_game_id (玩家游戏ID)
               └── score (得分)
```

## 🔄 数据迁移策略

### 迁移步骤

1. **备份数据**: 创建 `game_member_backup_20251115` 和 `game_member_wallet_backup_20251115`
2. **增加字段**: 为两个表增加 `group_id` 字段
3. **填充数据**: 根据 `group_name` 查找对应的 `group_id` 并填充
4. **修改索引**: 删除旧的唯一索引,创建新的唯一索引
5. **创建记录**: 为在多个圈子的用户创建独立的成员和钱包记录
6. **验证数据**: 检查数据完整性和一致性

### 迁移脚本

详见: `doc/migration_independent_group_balance.sql`

### 数据迁移示例

**迁移前**:
```
game_member:
| id | house_gid | game_id | game_name | group_name | balance |
|----|-----------|---------|-----------|------------|---------|
| 1  | 100       | 12345   | 张三      | VIP圈      | 1000    |

game_shop_group_member:
| id | group_id | user_id |
|----|----------|---------|
| 1  | 10       | 1       |  -- 用户1在VIP圈(group_id=10)
| 2  | 11       | 1       |  -- 用户1在普通圈(group_id=11)
```

**迁移后**:
```
game_member:
| id | house_gid | game_id | game_name | group_id | group_name | balance |
|----|-----------|---------|-----------|----------|------------|---------|
| 1  | 100       | 12345   | 张三      | 10       | VIP圈      | 1000    |
| 2  | 100       | 12345   | 张三      | 11       | 普通圈     | 0       |

game_member_wallet:
| id | house_gid | member_id | group_id | balance |
|----|-----------|-----------|----------|---------|
| 1  | 100       | 1         | 10       | 1000    |
| 2  | 100       | 2         | 11       | 0       |
```

## 🔍 查询示例

### 1. 查询用户在所有圈子的余额

```sql
SELECT 
    gm.game_id,
    gm.game_name,
    gsg.group_name,
    gm.balance,
    gmw.balance AS wallet_balance
FROM game_member gm
INNER JOIN game_shop_group gsg ON gsg.id = gm.group_id
LEFT JOIN game_member_wallet gmw ON gmw.member_id = gm.id AND gmw.group_id = gm.group_id
WHERE gm.house_gid = ? AND gm.game_id = ?
ORDER BY gsg.group_name;
```

### 2. 查询某个圈子的所有成员余额

```sql
SELECT 
    gm.game_id,
    gm.game_name,
    gm.balance,
    gmw.balance AS wallet_balance
FROM game_member gm
LEFT JOIN game_member_wallet gmw ON gmw.member_id = gm.id AND gmw.group_id = gm.group_id
WHERE gm.house_gid = ? AND gm.group_id = ?
ORDER BY gm.balance DESC;
```

### 3. 查询用户在某个圈子的战绩

```sql
SELECT 
    gbr.battle_at,
    gbr.room_uid,
    gbr.score,
    gbr.fee,
    gbr.player_balance
FROM game_battle_record gbr
WHERE gbr.house_gid = ? 
  AND gbr.player_game_id = ?
  AND gbr.group_id = ?
ORDER BY gbr.battle_at DESC
LIMIT 20;
```

### 4. 统计用户在某个圈子的战绩汇总

```sql
SELECT 
    COUNT(*) AS total_games,
    SUM(score) AS total_score,
    SUM(fee) AS total_fee,
    AVG(score) AS avg_score
FROM game_battle_record
WHERE house_gid = ? 
  AND player_game_id = ?
  AND group_id = ?
  AND battle_at >= ?
  AND battle_at < ?;
```

## ⚠️ 注意事项

### 1. 数据一致性

- **余额同步**: `game_member.balance` 和 `game_member_wallet.balance` 需要保持同步
- **圈子关联**: 确保 `group_id` 正确关联到 `game_shop_group.id`
- **战绩记录**: 新增战绩时必须指定正确的 `group_id`

### 2. 业务逻辑变更

- **余额查询**: 需要指定 `group_id` 参数
- **充值提现**: 需要指定操作的圈子
- **战绩统计**: 需要按圈子分组统计

### 3. 性能优化

- **索引优化**: 已为 `group_id` 创建索引,提高查询性能
- **分页查询**: 使用 LIMIT/OFFSET 进行分页
- **避免全表扫描**: 查询时始终带上 `house_gid` 和 `group_id` 条件

## 🚀 后续开发任务

### 1. Repository 层修改

- [ ] 修改 `GameMemberRepo` 接口,增加 `group_id` 参数
- [ ] 修改 `GameMemberWalletRepo` 接口,增加 `group_id` 参数
- [ ] 修改 `BattleRecordRepo` 查询方法,支持按圈子筛选

### 2. UseCase 层修改

- [ ] 修改余额查询逻辑,支持多圈子
- [ ] 修改充值提现逻辑,指定圈子
- [ ] 修改战绩统计逻辑,按圈子分组

### 3. Service 层修改

- [ ] 实现用户查询自己在所有圈子的余额 API
- [ ] 实现用户查询自己在某个圈子的战绩 API
- [ ] 实现管理员查询圈子成员余额 API
- [ ] 实现管理员查询圈子战绩统计 API

### 4. 前端适配

- [ ] 余额显示支持多圈子切换
- [ ] 战绩查询支持圈子筛选
- [ ] 充值提现支持选择圈子

## 📝 版本历史

| 版本 | 日期 | 作者 | 说明 |
|------|------|------|------|
| v1.0 | 2025-11-15 | AI Assistant | 初始版本,支持独立圈子余额 |

## 🔗 相关文档

- [数据库迁移脚本](./migration_independent_group_balance.sql)
- [表结构参考](./TABLE_REFERENCE.md)
- [DDL 说明文档](./README_DDL.md)

