# Battle Tiles 数据库 DDL 说明文档

## 📋 文件说明

本目录包含 Battle Tiles 项目的完整数据库结构定义（DDL）。

### 文件列表

- **`ddl_postgresql.sql`** - PostgreSQL 版本的完整 DDL，包含所有表定义和中文注释

## 🗄️ 数据库模块说明

### 1. 基础用户模块 (Basic User Module)

包含用户认证、角色管理、菜单权限等核心功能。

**主要表：**
- `basic_user` - 基础用户表
- `basic_role` - 基础角色表
- `basic_menu` - 菜单表
- `basic_role_menu_rel` - 角色菜单关联表
- `basic_user_role_rel` - 用户角色关联表

**用户角色：**
- `super_admin` - 超级管理员（可管理多个游戏账号）
- `store_admin` - 店铺管理员（独占管理一个店铺）
- `user` - 普通用户（绑定游戏账号）

### 2. 游戏账号模块 (Game Account Module)

管理用户的游戏账号绑定和中控账号。

**主要表：**
- `game_account` - 游戏账号表（用户绑定的游戏账号）
- `game_ctrl_account` - 中控账号表（超级管理员管理的游戏账号）
- `game_account_house` - 中控账号店铺绑定表
- `game_account_store_binding` - 游戏账号店铺绑定表

**业务规则：**
- 用户注册时必须绑定游戏账号
- 游戏账号需要通过游戏服务器验证
- 一个游戏账号只能绑定一个店铺
- 超级管理员可以管理多个中控账号

### 3. 游戏会话模块 (Game Session Module)

管理游戏登录会话和数据同步。

**主要表：**
- `game_session` - 游戏会话表
- `game_sync_log` - 游戏同步日志表

**会话状态：**
- `active` - 活跃
- `inactive` - 未活跃
- `error` - 错误

**同步类型：**
- `battle_record` - 战绩同步（5秒间隔）
- `member_list` - 成员列表同步（30秒间隔）
- `wallet_update` - 钱包更新同步（10秒间隔）
- `room_list` - 房间列表同步
- `group_member` - 圈成员同步

### 4. 店铺管理模块 (Shop Management Module)

管理店铺、店铺管理员和店铺设置。

**主要表：**
- `game_shop_admin` - 店铺管理员表
- `game_house_settings` - 店铺设置表

**管理员角色：**
- `admin` - 管理员（完全权限）
- `operator` - 操作员（有限权限）

### 5. 游戏成员模块 (Game Member Module)

管理店铺内的玩家成员。

**主要表：**
- `game_member` - 游戏成员表
- `game_member_wallet` - 游戏成员钱包表
- `game_member_rule` - 游戏成员规则表（VIP、多号等）

**成员属性：**
- 余额（单位：分）
- 信用额度
- 禁用状态
- VIP 状态
- 多号权限

### 6. 游戏战绩模块 (Game Battle Record Module)

记录游戏对战数据。

**主要表：**
- `game_battle_record` - 游戏战绩表

**记录内容：**
- 房间信息
- 玩家列表（JSON 格式）
- 得分、服务费
- 结算比例
- 玩家余额

### 7. 充值记录模块 (Recharge Record Module)

记录玩家的充值和提现操作。

**主要表：**
- `game_recharge_record` - 充值记录表

**记录内容：**
- 充值金额（正数=充值，负数=提现）
- 操作前后余额
- 操作人
- 操作时间

### 8. 费用结算模块 (Fee Settlement Module)

记录费用结算信息。

**主要表：**
- `game_fee_settle` - 费用结算表

## 🚀 使用方法

### 1. 创建数据库

```bash
# 连接到 PostgreSQL
psql -U postgres

# 创建数据库
CREATE DATABASE battle_tiles;

# 切换到新数据库
\c battle_tiles
```

### 2. 执行 DDL

```bash
# 方法 1: 使用 psql 命令行
psql -U postgres -d battle_tiles -f ddl_postgresql.sql

# 方法 2: 在 psql 中执行
\c battle_tiles
\i ddl_postgresql.sql
```

### 3. 验证表结构

```sql
-- 查看所有表
\dt

-- 查看表结构
\d basic_user
\d game_account
\d game_session

-- 查看表注释
SELECT 
    c.relname AS table_name,
    obj_description(c.oid) AS table_comment
FROM pg_class c
JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = 'public' AND c.relkind = 'r'
ORDER BY c.relname;

-- 查看字段注释
SELECT 
    a.attname AS column_name,
    col_description(a.attrelid, a.attnum) AS column_comment
FROM pg_attribute a
JOIN pg_class c ON a.attrelid = c.oid
WHERE c.relname = 'basic_user' AND a.attnum > 0;
```

## 🔧 维护建议

### 1. 定期备份

```bash
# 全量备份
pg_dump -h localhost -U postgres -d battle_tiles -F c -f backup_$(date +%Y%m%d).dump

# 恢复备份
pg_restore -h localhost -U postgres -d battle_tiles backup_20251112.dump
```

### 2. 性能优化

```sql
-- 分析表统计信息
ANALYZE basic_user;
ANALYZE game_account;
ANALYZE game_session;
ANALYZE game_battle_record;

-- 查看表大小
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- 查看索引使用情况
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;
```

### 3. 清理软删除数据

```sql
-- 清理 90 天前的软删除数据
DELETE FROM basic_user WHERE is_del = 1 AND deleted_at < NOW() - INTERVAL '90 days';
DELETE FROM game_account WHERE is_del = 1 AND deleted_at < NOW() - INTERVAL '90 days';
DELETE FROM game_session WHERE is_del = 1 AND deleted_at < NOW() - INTERVAL '90 days';

-- 清理后执行 VACUUM
VACUUM ANALYZE basic_user;
VACUUM ANALYZE game_account;
VACUUM ANALYZE game_session;
```

## 📊 数据字典

### 金额单位说明

所有金额字段统一使用 **分** 作为单位：
- `balance` - 余额（分）
- `credit` - 信用额度（分）
- `amount` - 金额（分）
- `fee` - 服务费（分）
- `push_credit` - 推送额度（分）

**转换公式：**
- 1 元 = 100 分
- 显示时需要除以 100

### 时间字段说明

所有时间字段使用 `TIMESTAMP WITH TIME ZONE` 类型：
- `created_at` - 创建时间（自动设置）
- `updated_at` - 更新时间（触发器自动更新）
- `deleted_at` - 删除时间（软删除）
- `last_login_at` - 最后登录时间
- `last_sync_at` - 最后同步时间

### 软删除说明

使用双重软删除机制：
- `is_del` - 软删除标记（0=未删除，1=已删除）
- `deleted_at` - 删除时间戳

**查询时需要过滤：**
```sql
SELECT * FROM basic_user WHERE is_del = 0;
SELECT * FROM game_account WHERE is_del = 0;
```

## 🔐 安全建议

### 1. 创建专用数据库用户

```sql
-- 创建只读用户
CREATE USER battle_readonly WITH PASSWORD 'secure_password_here';
GRANT CONNECT ON DATABASE battle_tiles TO battle_readonly;
GRANT USAGE ON SCHEMA public TO battle_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO battle_readonly;

-- 创建读写用户
CREATE USER battle_readwrite WITH PASSWORD 'secure_password_here';
GRANT CONNECT ON DATABASE battle_tiles TO battle_readwrite;
GRANT USAGE ON SCHEMA public TO battle_readwrite;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO battle_readwrite;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO battle_readwrite;
```

### 2. 启用 SSL 连接

在 `postgresql.conf` 中配置：
```
ssl = on
ssl_cert_file = 'server.crt'
ssl_key_file = 'server.key'
```

### 3. 配置连接限制

```sql
-- 限制用户连接数
ALTER USER battle_readwrite CONNECTION LIMIT 50;
ALTER USER battle_readonly CONNECTION LIMIT 20;
```

## 📝 更新日志

### 2025-11-12
- ✅ 创建完整的 PostgreSQL DDL
- ✅ 添加所有表的中文注释
- ✅ 添加索引优化
- ✅ 添加触发器自动更新 updated_at
- ✅ 添加维护脚本和备份建议

## 🤝 贡献指南

如需修改数据库结构：

1. 在 `battle-tiles/internal/dal/model` 中修改 Go 模型
2. 更新 `ddl_postgresql.sql` 文件
3. 添加迁移脚本（如果需要）
4. 更新本文档

## 📞 联系方式

如有问题，请联系开发团队。

