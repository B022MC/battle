-- ============================================================================
-- 游戏账号入圈系统 - 数据库迁移脚本
-- 日期: 2025-11-20
-- 描述: 创建游戏账号圈子关系表，实现游戏账号入圈机制
-- ============================================================================

BEGIN;

-- ============================================================================
-- 1. 修改 game_account 表，user_id 改为可选
-- ============================================================================

-- 检查并修改 user_id 字段
DO $$
BEGIN
    -- 移除 NOT NULL 约束（如果存在）
    ALTER TABLE game_account ALTER COLUMN user_id DROP NOT NULL;
    RAISE NOTICE 'game_account.user_id 已改为可选';
EXCEPTION
    WHEN OTHERS THEN
        RAISE NOTICE 'game_account.user_id 可能已经是可选的';
END $$;

-- 更新字段注释
COMMENT ON COLUMN game_account.user_id IS '关联的用户ID（可选，用于用户反向查询圈子和战绩）';

-- ============================================================================
-- 2. 创建 game_account_group 表（游戏账号圈子关系）
-- ============================================================================

-- 删除旧表（如果存在）
DROP TABLE IF EXISTS game_account_group CASCADE;

-- 创建新表
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
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 创建唯一约束（一个游戏账号在一个店铺只能属于一个圈子）
-- 先删除可能存在的旧约束
DO $$
BEGIN
    -- 尝试删除旧约束
    BEGIN
        ALTER TABLE game_account_group DROP CONSTRAINT IF EXISTS uk_game_account_house;
    EXCEPTION
        WHEN OTHERS THEN NULL;
    END;
    
    -- 尝试删除同名索引
    BEGIN
        DROP INDEX IF EXISTS uk_game_account_house;
    EXCEPTION
        WHEN OTHERS THEN NULL;
    END;
    
    -- 添加新约束
    ALTER TABLE game_account_group 
    ADD CONSTRAINT uk_game_account_house UNIQUE (game_account_id, house_gid);
    
    RAISE NOTICE '唯一约束 uk_game_account_house 已创建';
EXCEPTION
    WHEN duplicate_table THEN
        RAISE NOTICE '唯一约束 uk_game_account_house 已存在';
END $$;

-- 表注释
COMMENT ON TABLE game_account_group IS '游戏账号圈子关系表（游戏账号入圈记录）';

-- 列注释
COMMENT ON COLUMN game_account_group.id IS '关系ID';
COMMENT ON COLUMN game_account_group.game_account_id IS '游戏账号ID（game_account.id）';
COMMENT ON COLUMN game_account_group.house_gid IS '店铺GID';
COMMENT ON COLUMN game_account_group.group_id IS '圈子ID（game_shop_group.id）';
COMMENT ON COLUMN game_account_group.group_name IS '圈子名称（冗余字段，方便查询）';
COMMENT ON COLUMN game_account_group.admin_user_id IS '圈主用户ID（店铺管理员）';
COMMENT ON COLUMN game_account_group.approved_by_user_id IS '审批通过的管理员ID';
COMMENT ON COLUMN game_account_group.status IS '状态：active=激活, inactive=未激活';
COMMENT ON COLUMN game_account_group.joined_at IS '加入时间';
COMMENT ON COLUMN game_account_group.created_at IS '创建时间';
COMMENT ON COLUMN game_account_group.updated_at IS '更新时间';

-- 创建索引（使用 IF NOT EXISTS 避免重复创建错误）
CREATE INDEX IF NOT EXISTS idx_account_group_game_account ON game_account_group(game_account_id);
CREATE INDEX IF NOT EXISTS idx_account_group_house ON game_account_group(house_gid);
CREATE INDEX IF NOT EXISTS idx_account_group_group ON game_account_group(group_id);
CREATE INDEX IF NOT EXISTS idx_account_group_admin ON game_account_group(admin_user_id);
CREATE INDEX IF NOT EXISTS idx_account_group_status ON game_account_group(status);
CREATE INDEX IF NOT EXISTS idx_account_group_house_status ON game_account_group(house_gid, status);

-- 创建更新触发器
DROP TRIGGER IF EXISTS update_game_account_group_updated_at ON game_account_group;
CREATE TRIGGER update_game_account_group_updated_at
    BEFORE UPDATE ON game_account_group
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- ============================================================================
-- 3. 修改 game_battle_record 表，player_id 改为可选
-- ============================================================================

-- 检查并修改 player_id 字段
DO $$
BEGIN
    ALTER TABLE game_battle_record ALTER COLUMN player_id DROP NOT NULL;
    RAISE NOTICE 'game_battle_record.player_id 已改为可选';
EXCEPTION
    WHEN OTHERS THEN
        RAISE NOTICE 'game_battle_record.player_id 可能已经是可选的';
END $$;

-- 更新字段注释
COMMENT ON COLUMN game_battle_record.player_id IS '玩家用户ID（可选，用于用户反向查询）';
COMMENT ON COLUMN game_battle_record.player_game_id IS '玩家游戏账号ID（game_account.id，必填）';

-- ============================================================================
-- 4. 数据迁移：将现有的用户-圈子关系迁移到游戏账号-圈子关系
-- ============================================================================

-- 检查是否存在 game_shop_group_member 表
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables 
               WHERE table_name = 'game_shop_group_member') THEN
        
        -- 迁移数据
        INSERT INTO game_account_group (
            game_account_id,
            house_gid,
            group_id,
            group_name,
            admin_user_id,
            approved_by_user_id,
            status,
            joined_at,
            created_at
        )
        SELECT 
            ga.id AS game_account_id,
            gsg.house_gid,
            gsgm.group_id,
            gsg.group_name,
            gsg.admin_user_id,
            gsg.admin_user_id AS approved_by_user_id,
            'active' AS status,
            gsgm.joined_at,
            gsgm.created_at
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
        
        RAISE NOTICE '数据迁移完成';
    ELSE
        RAISE NOTICE 'game_shop_group_member 表不存在，跳过数据迁移';
    END IF;
END $$;

-- ============================================================================
-- 5. 数据验证
-- ============================================================================

DO $$
DECLARE
    v_game_account_count INTEGER;
    v_account_group_count INTEGER;
    v_user_with_account_count INTEGER;
    v_account_without_user_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO v_game_account_count 
    FROM game_account WHERE is_del = 0;
    
    SELECT COUNT(*) INTO v_account_group_count 
    FROM game_account_group WHERE status = 'active';
    
    SELECT COUNT(*) INTO v_user_with_account_count 
    FROM game_account WHERE user_id IS NOT NULL AND is_del = 0;
    
    SELECT COUNT(*) INTO v_account_without_user_count 
    FROM game_account WHERE user_id IS NULL AND is_del = 0;
    
    RAISE NOTICE '';
    RAISE NOTICE '=== 数据迁移统计 ===';
    RAISE NOTICE '游戏账号总数: %', v_game_account_count;
    RAISE NOTICE '游戏账号-圈子关系数: %', v_account_group_count;
    RAISE NOTICE '已绑定用户的游戏账号数: %', v_user_with_account_count;
    RAISE NOTICE '未绑定用户的游戏账号数: %', v_account_without_user_count;
    RAISE NOTICE '==================';
    RAISE NOTICE '';
END $$;

COMMIT;

-- ============================================================================
-- 迁移完成！
-- ============================================================================

-- 常用查询示例：

-- 1. 战绩同步时查询游戏账号和圈子
/*
SELECT 
    ga.id AS game_account_id,
    ga.user_id,
    gag.group_id,
    gag.group_name
FROM game_account ga
LEFT JOIN game_account_group gag 
    ON ga.id = gag.game_account_id 
    AND gag.house_gid = 58959
    AND gag.status = 'active'
WHERE ga.game_user_id = '22805688'
  AND ga.is_del = 0;
*/

-- 2. 用户查询自己的圈子
/*
-- 第一步：获取用户的游戏账号
SELECT id FROM game_account 
WHERE user_id = 123 AND is_del = 0;

-- 第二步：获取游戏账号的圈子
SELECT 
    gag.*,
    gsg.description,
    bu.nick_name AS admin_name
FROM game_account_group gag
JOIN game_shop_group gsg ON gag.group_id = gsg.id
JOIN basic_user bu ON gag.admin_user_id = bu.id
WHERE gag.game_account_id IN (456, 789)
  AND gag.status = 'active';
*/

-- 3. 用户查询自己的战绩
/*
-- 第一步：获取用户的游戏账号
SELECT id FROM game_account 
WHERE user_id = 123 AND is_del = 0;

-- 第二步：查询战绩
SELECT * FROM game_battle_record
WHERE player_game_id IN (456, 789)
ORDER BY battle_at DESC
LIMIT 50;
*/

-- 4. 管理员查询圈内成员的游戏账号
/*
SELECT 
    ga.*,
    gag.joined_at,
    gag.status
FROM game_account_group gag
JOIN game_account ga ON gag.game_account_id = ga.id
WHERE gag.group_id = 1
  AND gag.status = 'active'
ORDER BY gag.joined_at DESC;
*/
