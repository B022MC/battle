# Battle Tiles 数据库表快速参考

## 📑 表列表总览

| 序号 | 表名 | 中文名称 | 用途 |
|------|------|----------|------|
| 1 | `basic_user` | 基础用户表 | 用户认证和基本信息 |
| 2 | `basic_role` | 基础角色表 | 角色定义和管理 |
| 3 | `basic_menu` | 基础菜单表 | 系统菜单配置 |
| 4 | `basic_role_menu_rel` | 角色菜单关联表 | 角色权限关联 |
| 5 | `basic_user_role_rel` | 用户角色关联表 | 用户角色关联 |
| 6 | `base_platform` | 平台表 | 云平台配置 |
| 7 | `game_account` | 游戏账号表 | 用户游戏账号 |
| 8 | `game_ctrl_account` | 中控账号表 | 超管游戏账号 |
| 9 | `game_account_house` | 中控账号店铺绑定表 | 中控账号与店铺关联 |
| 10 | `game_account_store_binding` | 游戏账号店铺绑定表 | 用户账号与店铺关联 |
| 11 | `game_session` | 游戏会话表 | 登录会话管理 |
| 12 | `game_sync_log` | 游戏同步日志表 | 数据同步记录 |
| 13 | `game_shop_admin` | 店铺管理员表 | 店铺管理员配置 |
| 14 | `game_house_settings` | 店铺设置表 | 店铺配置信息 |
| 15 | `game_member` | 游戏成员表 | 店铺玩家信息 |
| 16 | `game_member_wallet` | 游戏成员钱包表 | 玩家钱包管理 |
| 17 | `game_member_rule` | 游戏成员规则表 | 玩家特殊规则 |
| 18 | `game_battle_record` | 游戏战绩表 | 对战记录 |
| 19 | `game_recharge_record` | 充值记录表 | 充值提现记录 |
| 20 | `game_fee_settle` | 费用结算表 | 费用结算记录 |

## 🔍 常用查询示例

### 1. 用户和角色相关查询

```sql
-- 查询所有角色
SELECT id, code, name, parent_id, enable, created_at
FROM basic_role
WHERE is_deleted = FALSE
ORDER BY id;

-- 查询启用的角色
SELECT id, code, name, remark
FROM basic_role
WHERE is_deleted = FALSE AND enable = TRUE
ORDER BY id;

-- 查询用户及其角色
SELECT
    u.id,
    u.username,
    u.nick_name,
    u.role AS user_role,
    r.code AS role_code,
    r.name AS role_name
FROM basic_user u
LEFT JOIN basic_user_role_rel urr ON u.id = urr.user_id
LEFT JOIN basic_role r ON urr.role_id = r.id
WHERE u.is_del = 0 AND (r.is_deleted = FALSE OR r.id IS NULL);

-- 查询所有超级管理员
SELECT id, username, nick_name, game_nickname, created_at
FROM basic_user
WHERE role = 'super_admin' AND is_del = 0;

-- 查询所有店铺管理员
SELECT id, username, nick_name, role, created_at
FROM basic_user
WHERE role = 'store_admin' AND is_del = 0;

-- 查询用户的游戏账号
SELECT 
    u.username,
    u.nick_name,
    ga.account,
    ga.nickname AS game_nickname,
    ga.verification_status,
    ga.created_at
FROM basic_user u
LEFT JOIN game_account ga ON u.id = ga.user_id
WHERE u.id = ? AND ga.is_del = 0;
```

### 2. 中控账号相关查询

```sql
-- 查询所有中控账号及其绑定的店铺
SELECT 
    gca.id,
    gca.identifier,
    gca.status,
    gca.last_verify_at,
    COUNT(gah.id) AS house_count
FROM game_ctrl_account gca
LEFT JOIN game_account_house gah ON gca.id = gah.game_account_id
WHERE gca.deleted_at IS NULL
GROUP BY gca.id;

-- 查询某个中控账号的所有店铺
SELECT 
    gah.house_gid,
    gah.status,
    gah.is_default,
    gah.created_at
FROM game_account_house gah
WHERE gah.game_account_id = ?
ORDER BY gah.is_default DESC, gah.created_at DESC;
```

### 3. 会话相关查询

```sql
-- 查询活跃会话
SELECT 
    gs.id,
    gs.house_gid,
    gca.identifier AS ctrl_account,
    gs.state,
    gs.sync_status,
    gs.last_sync_at,
    gs.created_at
FROM game_session gs
JOIN game_ctrl_account gca ON gs.game_ctrl_account_id = gca.id
WHERE gs.state = 'active' AND gs.is_del = 0
ORDER BY gs.created_at DESC;

-- 查询某个店铺的所有会话
SELECT 
    gs.id,
    gca.identifier AS ctrl_account,
    gs.state,
    gs.sync_status,
    gs.auto_sync_enabled,
    gs.last_sync_at,
    gs.created_at
FROM game_session gs
JOIN game_ctrl_account gca ON gs.game_ctrl_account_id = gca.id
WHERE gs.house_gid = ? AND gs.is_del = 0
ORDER BY gs.created_at DESC;

-- 查询同步日志
SELECT 
    gsl.id,
    gsl.session_id,
    gsl.sync_type,
    gsl.status,
    gsl.records_synced,
    gsl.error_message,
    gsl.started_at,
    gsl.completed_at,
    EXTRACT(EPOCH FROM (gsl.completed_at - gsl.started_at)) AS duration_seconds
FROM game_sync_log gsl
WHERE gsl.session_id = ?
ORDER BY gsl.started_at DESC
LIMIT 100;
```

### 4. 店铺管理相关查询

```sql
-- 查询店铺管理员
SELECT 
    gsa.id,
    gsa.house_gid,
    u.username,
    u.nick_name,
    gsa.role,
    gsa.is_exclusive,
    gsa.created_at
FROM game_shop_admin gsa
JOIN basic_user u ON gsa.user_id = u.id
WHERE gsa.house_gid = ? AND gsa.deleted_at IS NULL;

-- 查询用户管理的店铺
SELECT 
    gsa.house_gid,
    gsa.role,
    gsa.is_exclusive,
    ghs.share_fee,
    ghs.push_credit,
    gsa.created_at
FROM game_shop_admin gsa
LEFT JOIN game_house_settings ghs ON gsa.house_gid = ghs.house_gid
WHERE gsa.user_id = ? AND gsa.deleted_at IS NULL;
```

### 5. 成员相关查询

```sql
-- 查询店铺成员列表
SELECT 
    gm.id,
    gm.game_id,
    gm.game_name,
    gm.group_name,
    gm.balance,
    gm.credit,
    gm.forbid,
    gm.created_at
FROM game_member gm
WHERE gm.house_gid = ?
ORDER BY gm.created_at DESC;

-- 查询成员余额
SELECT 
    gm.game_id,
    gm.game_name,
    gm.balance / 100.0 AS balance_yuan,
    gm.credit / 100.0 AS credit_yuan,
    gmw.balance / 100.0 AS wallet_balance_yuan
FROM game_member gm
LEFT JOIN game_member_wallet gmw ON gm.id = gmw.member_id
WHERE gm.house_gid = ? AND gm.game_id = ?;

-- 查询禁用成员
SELECT 
    gm.game_id,
    gm.game_name,
    gm.group_name,
    gm.balance / 100.0 AS balance_yuan,
    gm.updated_at
FROM game_member gm
WHERE gm.house_gid = ? AND gm.forbid = TRUE;
```

### 6. 战绩相关查询

```sql
-- 查询最近战绩
SELECT 
    gbr.id,
    gbr.house_gid,
    gbr.room_uid,
    gbr.kind_id,
    gbr.base_score,
    gbr.battle_at,
    gbr.player_game_id,
    gbr.player_game_name,
    gbr.score,
    gbr.fee / 100.0 AS fee_yuan,
    gbr.player_balance / 100.0 AS balance_yuan
FROM game_battle_record gbr
WHERE gbr.house_gid = ?
ORDER BY gbr.battle_at DESC
LIMIT 100;

-- 查询某个玩家的战绩
SELECT 
    gbr.battle_at,
    gbr.room_uid,
    gbr.kind_id,
    gbr.score,
    gbr.fee / 100.0 AS fee_yuan,
    gbr.player_balance / 100.0 AS balance_yuan
FROM game_battle_record gbr
WHERE gbr.house_gid = ? AND gbr.player_game_id = ?
ORDER BY gbr.battle_at DESC
LIMIT 50;

-- 统计某个店铺的战绩汇总
SELECT 
    DATE(gbr.battle_at) AS battle_date,
    COUNT(*) AS total_battles,
    SUM(gbr.score) AS total_score,
    SUM(gbr.fee) / 100.0 AS total_fee_yuan
FROM game_battle_record gbr
WHERE gbr.house_gid = ?
    AND gbr.battle_at >= NOW() - INTERVAL '30 days'
GROUP BY DATE(gbr.battle_at)
ORDER BY battle_date DESC;
```

### 7. 充值记录查询

```sql
-- 查询充值记录
SELECT 
    grr.id,
    grr.player_id,
    grr.group_name,
    grr.amount / 100.0 AS amount_yuan,
    grr.balance_before / 100.0 AS balance_before_yuan,
    grr.balance_after / 100.0 AS balance_after_yuan,
    grr.recharged_at,
    CASE 
        WHEN grr.amount > 0 THEN '充值'
        ELSE '提现'
    END AS transaction_type
FROM game_recharge_record grr
WHERE grr.house_gid = ?
ORDER BY grr.recharged_at DESC
LIMIT 100;

-- 统计充值汇总
SELECT 
    DATE(grr.recharged_at) AS recharge_date,
    COUNT(*) AS total_transactions,
    SUM(CASE WHEN grr.amount > 0 THEN grr.amount ELSE 0 END) / 100.0 AS total_deposit_yuan,
    SUM(CASE WHEN grr.amount < 0 THEN ABS(grr.amount) ELSE 0 END) / 100.0 AS total_withdrawal_yuan
FROM game_recharge_record grr
WHERE grr.house_gid = ?
    AND grr.recharged_at >= NOW() - INTERVAL '30 days'
GROUP BY DATE(grr.recharged_at)
ORDER BY recharge_date DESC;
```

## 📊 统计查询

### 系统统计

```sql
-- 用户统计
SELECT 
    role,
    COUNT(*) AS user_count
FROM basic_user
WHERE is_del = 0
GROUP BY role;

-- 游戏账号统计
SELECT 
    verification_status,
    COUNT(*) AS account_count
FROM game_account
WHERE is_del = 0
GROUP BY verification_status;

-- 会话统计
SELECT 
    state,
    COUNT(*) AS session_count
FROM game_session
WHERE is_del = 0
GROUP BY state;
```

### 店铺统计

```sql
-- 店铺成员统计
SELECT 
    house_gid,
    COUNT(*) AS member_count,
    SUM(balance) / 100.0 AS total_balance_yuan,
    COUNT(CASE WHEN forbid = TRUE THEN 1 END) AS forbidden_count
FROM game_member
GROUP BY house_gid
ORDER BY member_count DESC;

-- 店铺战绩统计（最近30天）
SELECT 
    house_gid,
    COUNT(*) AS battle_count,
    SUM(fee) / 100.0 AS total_fee_yuan
FROM game_battle_record
WHERE battle_at >= NOW() - INTERVAL '30 days'
GROUP BY house_gid
ORDER BY battle_count DESC;
```

## 🔧 维护查询

### 数据清理

```sql
-- 查看软删除数据量
SELECT 
    'basic_user' AS table_name,
    COUNT(*) AS deleted_count
FROM basic_user WHERE is_del = 1
UNION ALL
SELECT 
    'game_account',
    COUNT(*)
FROM game_account WHERE is_del = 1
UNION ALL
SELECT 
    'game_session',
    COUNT(*)
FROM game_session WHERE is_del = 1;

-- 查看可清理的旧数据
SELECT 
    'basic_user' AS table_name,
    COUNT(*) AS cleanable_count
FROM basic_user 
WHERE is_del = 1 AND deleted_at < NOW() - INTERVAL '90 days'
UNION ALL
SELECT 
    'game_account',
    COUNT(*)
FROM game_account 
WHERE is_del = 1 AND deleted_at < NOW() - INTERVAL '90 days';
```

### 性能监控

```sql
-- 查看表大小
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS total_size,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) AS table_size,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename) - pg_relation_size(schemaname||'.'||tablename)) AS index_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- 查看慢查询（需要启用 pg_stat_statements 扩展）
-- CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    max_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 20;
```

## 📝 注意事项

1. **金额单位**：所有金额字段使用分为单位，显示时需要除以 100
2. **软删除**：查询时需要过滤 `is_del = 0` 或 `deleted_at IS NULL`
3. **时区**：所有时间字段使用 `TIMESTAMP WITH TIME ZONE`
4. **索引**：复杂查询前检查是否有合适的索引
5. **性能**：大数据量查询时使用 `LIMIT` 限制结果集

## 🔗 相关文档

- [完整 DDL 文件](./ddl_postgresql.sql)
- [DDL 使用说明](./README_DDL.md)

