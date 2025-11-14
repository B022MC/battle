-- ============================================================================
-- Battle Tiles 数据库迁移脚本
-- 功能: 支持用户在不同圈子有独立余额
-- 日期: 2025-11-15
-- 版本: v1.0
-- ============================================================================

-- 说明:
-- 1. 修改 game_member 表,增加 group_id 字段,支持用户在不同圈子有独立记录
-- 2. 修改 game_member_wallet 表,增加 group_id 字段,支持每个圈子独立的钱包
-- 3. 修改唯一索引,确保 (house_gid, game_id, group_id) 唯一
-- 4. 迁移现有数据,为每个成员创建对应圈子的记录

-- ============================================================================
-- 第一步: 备份现有数据
-- ============================================================================

-- 备份 game_member 表
CREATE TABLE IF NOT EXISTS game_member_backup_20251115 AS 
SELECT * FROM game_member;

-- 备份 game_member_wallet 表
CREATE TABLE IF NOT EXISTS game_member_wallet_backup_20251115 AS 
SELECT * FROM game_member_wallet;

-- ============================================================================
-- 第二步: 修改 game_member 表结构
-- ============================================================================

-- 2.1 增加 group_id 字段
ALTER TABLE game_member 
ADD COLUMN IF NOT EXISTS group_id INT DEFAULT NULL;

-- 2.2 为 group_id 添加注释
COMMENT ON COLUMN game_member.group_id IS '圈子ID,关联 game_shop_group.id';

-- 2.3 为现有数据填充 group_id
-- 根据 group_name 查找对应的 group_id
UPDATE game_member gm
SET group_id = (
    SELECT gsg.id 
    FROM game_shop_group gsg 
    WHERE gsg.house_gid = gm.house_gid 
      AND gsg.group_name = gm.group_name
    LIMIT 1
)
WHERE gm.group_id IS NULL AND gm.group_name IS NOT NULL AND gm.group_name != '';

-- 2.4 删除旧的唯一索引
DROP INDEX IF EXISTS uk_game_member_house_game;

-- 2.5 创建新的唯一索引 (house_gid, game_id, group_id)
CREATE UNIQUE INDEX uk_game_member_house_game_group 
ON game_member(house_gid, game_id, group_id);

-- 2.6 创建 group_id 的普通索引
CREATE INDEX IF NOT EXISTS idx_game_member_group_id 
ON game_member(group_id);

-- 2.7 添加外键约束(可选,根据实际需求决定是否启用)
-- ALTER TABLE game_member 
-- ADD CONSTRAINT fk_game_member_group 
-- FOREIGN KEY (group_id) REFERENCES game_shop_group(id) ON DELETE SET NULL;

-- ============================================================================
-- 第三步: 修改 game_member_wallet 表结构
-- ============================================================================

-- 3.1 增加 group_id 字段
ALTER TABLE game_member_wallet 
ADD COLUMN IF NOT EXISTS group_id INT DEFAULT NULL;

-- 3.2 为 group_id 添加注释
COMMENT ON COLUMN game_member_wallet.group_id IS '圈子ID,关联 game_shop_group.id';

-- 3.3 为现有数据填充 group_id
-- 从 game_member 表获取 group_id
UPDATE game_member_wallet gmw
SET group_id = (
    SELECT gm.group_id 
    FROM game_member gm 
    WHERE gm.id = gmw.member_id
    LIMIT 1
)
WHERE gmw.group_id IS NULL;

-- 3.4 创建唯一索引 (house_gid, member_id, group_id)
CREATE UNIQUE INDEX IF NOT EXISTS uk_game_member_wallet_house_member_group 
ON game_member_wallet(house_gid, member_id, group_id);

-- 3.5 创建 group_id 的普通索引
CREATE INDEX IF NOT EXISTS idx_game_member_wallet_group_id 
ON game_member_wallet(group_id);

-- ============================================================================
-- 第四步: 数据迁移 - 为用户在多个圈子创建独立记录
-- ============================================================================

-- 说明: 如果一个用户(game_id)在同一个店铺(house_gid)加入了多个圈子,
--       需要为每个圈子创建独立的 game_member 和 game_member_wallet 记录

-- 4.1 查找需要迁移的数据
-- 找出所有在多个圈子的用户
WITH multi_group_users AS (
    SELECT DISTINCT
        gm.house_gid,
        gm.game_id,
        gm.game_name,
        gsgm.group_id
    FROM game_member gm
    INNER JOIN game_shop_group_member gsgm ON gsgm.user_id = (
        SELECT ga.user_id 
        FROM game_account ga 
        WHERE ga.game_id = gm.game_id 
        LIMIT 1
    )
    INNER JOIN game_shop_group gsg ON gsg.id = gsgm.group_id AND gsg.house_gid = gm.house_gid
    WHERE gm.group_id IS NULL OR gm.group_id != gsgm.group_id
)
-- 4.2 为每个圈子创建 game_member 记录
INSERT INTO game_member (
    house_gid, 
    game_id, 
    game_name, 
    group_id,
    group_name, 
    balance, 
    credit, 
    forbid, 
    recommender, 
    use_multi_gids, 
    active_gid,
    created_at,
    updated_at
)
SELECT 
    mgu.house_gid,
    mgu.game_id,
    mgu.game_name,
    mgu.group_id,
    gsg.group_name,
    0 AS balance,  -- 新圈子初始余额为0
    COALESCE(gm.credit, 0) AS credit,
    COALESCE(gm.forbid, false) AS forbid,
    gm.recommender,
    COALESCE(gm.use_multi_gids, false) AS use_multi_gids,
    gm.active_gid,
    NOW() AS created_at,
    NOW() AS updated_at
FROM multi_group_users mgu
INNER JOIN game_shop_group gsg ON gsg.id = mgu.group_id
LEFT JOIN game_member gm ON gm.house_gid = mgu.house_gid AND gm.game_id = mgu.game_id
WHERE NOT EXISTS (
    SELECT 1 FROM game_member gm2 
    WHERE gm2.house_gid = mgu.house_gid 
      AND gm2.game_id = mgu.game_id 
      AND gm2.group_id = mgu.group_id
)
ON CONFLICT (house_gid, game_id, group_id) DO NOTHING;

-- 4.3 为新创建的 game_member 记录创建对应的 game_member_wallet 记录
INSERT INTO game_member_wallet (
    house_gid,
    member_id,
    group_id,
    balance,
    forbid,
    limit_min,
    updated_at,
    updated_by
)
SELECT 
    gm.house_gid,
    gm.id AS member_id,
    gm.group_id,
    0 AS balance,  -- 新圈子初始余额为0
    gm.forbid,
    0 AS limit_min,
    NOW() AS updated_at,
    0 AS updated_by
FROM game_member gm
WHERE gm.group_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM game_member_wallet gmw 
      WHERE gmw.member_id = gm.id 
        AND gmw.group_id = gm.group_id
  )
ON CONFLICT (house_gid, member_id, group_id) DO NOTHING;

-- ============================================================================
-- 第五步: 验证数据完整性
-- ============================================================================

-- 5.1 检查是否有 game_member 记录没有对应的 group_id
SELECT COUNT(*) AS members_without_group_id
FROM game_member
WHERE group_id IS NULL;

-- 5.2 检查是否有 game_member_wallet 记录没有对应的 group_id
SELECT COUNT(*) AS wallets_without_group_id
FROM game_member_wallet
WHERE group_id IS NULL;

-- 5.3 检查是否有 game_member 记录没有对应的 wallet 记录
SELECT COUNT(*) AS members_without_wallet
FROM game_member gm
WHERE gm.group_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM game_member_wallet gmw 
      WHERE gmw.member_id = gm.id 
        AND gmw.group_id = gm.group_id
  );

-- 5.4 统计每个店铺的圈子数量和成员数量
SELECT 
    gsg.house_gid,
    gsg.group_name,
    COUNT(DISTINCT gm.game_id) AS member_count,
    SUM(gmw.balance) AS total_balance
FROM game_shop_group gsg
LEFT JOIN game_member gm ON gm.group_id = gsg.id
LEFT JOIN game_member_wallet gmw ON gmw.member_id = gm.id AND gmw.group_id = gsg.id
GROUP BY gsg.house_gid, gsg.id, gsg.group_name
ORDER BY gsg.house_gid, gsg.group_name;

-- ============================================================================
-- 第六步: 清理和优化
-- ============================================================================

-- 6.1 更新表统计信息
ANALYZE game_member;
ANALYZE game_member_wallet;

-- 6.2 重建索引(可选)
-- REINDEX TABLE game_member;
-- REINDEX TABLE game_member_wallet;

-- ============================================================================
-- 回滚脚本 (如果需要回滚,请执行以下语句)
-- ============================================================================

/*
-- 恢复 game_member 表
DROP TABLE IF EXISTS game_member;
ALTER TABLE game_member_backup_20251115 RENAME TO game_member;

-- 恢复 game_member_wallet 表
DROP TABLE IF EXISTS game_member_wallet;
ALTER TABLE game_member_wallet_backup_20251115 RENAME TO game_member_wallet;

-- 重建原有索引
CREATE UNIQUE INDEX uk_game_member_house_game ON game_member(house_gid, game_id);
*/

-- ============================================================================
-- 迁移完成
-- ============================================================================

