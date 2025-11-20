-- ============================================================================
-- 添加唯一约束：一个游戏账号只能绑定一个店铺
-- ============================================================================
-- 业务规则：一个用户 → 一个游戏账号 → 一个店铺 → 一个圈子
-- ============================================================================

-- 检查当前是否有违反约束的数据
SELECT 
    game_account_id,
    COUNT(*) as house_count
FROM game_account_house
GROUP BY game_account_id
HAVING COUNT(*) > 1
ORDER BY house_count DESC;

-- 如果上面的查询有结果，说明有账号绑定了多个店铺，需要先清理

-- ============================================================================
-- 添加唯一约束
-- ============================================================================

-- 方式1：使用 ALTER TABLE 添加唯一约束
DO $$ 
BEGIN
    -- 尝试删除旧的唯一约束（如果是组合约束）
    BEGIN
        ALTER TABLE game_account_house DROP CONSTRAINT IF EXISTS uk_account_house_unique;
    EXCEPTION
        WHEN OTHERS THEN NULL;
    END;
    
    -- 尝试删除旧索引
    BEGIN
        DROP INDEX IF EXISTS uk_account_house_unique;
    EXCEPTION
        WHEN OTHERS THEN NULL;
    END;
    
    RAISE NOTICE '已删除旧约束';
END $$;

-- 添加新的唯一约束：game_account_id 必须唯一
-- 这样一个游戏账号只能绑定一个店铺
ALTER TABLE game_account_house 
ADD CONSTRAINT uk_game_account_id_unique UNIQUE (game_account_id);

-- 验证约束
SELECT 
    conname as constraint_name,
    contype as constraint_type,
    pg_get_constraintdef(oid) as constraint_definition
FROM pg_constraint
WHERE conrelid = 'game_account_house'::regclass
  AND conname = 'uk_game_account_id_unique';

-- ============================================================================
-- 注释
-- ============================================================================
COMMENT ON CONSTRAINT uk_game_account_id_unique ON game_account_house IS 
'唯一约束：一个游戏账号只能绑定一个店铺（一个用户只能在一个店铺的一个圈子下）';

-- ============================================================================
-- 清理违反约束的数据（如果需要）
-- ============================================================================
-- 如果有账号绑定了多个店铺，可以执行以下SQL保留最新的一条记录

-- 查看需要清理的数据
-- SELECT * FROM game_account_house
-- WHERE game_account_id IN (
--     SELECT game_account_id
--     FROM game_account_house
--     GROUP BY game_account_id
--     HAVING COUNT(*) > 1
-- )
-- ORDER BY game_account_id, id;

-- 删除重复的记录，保留 id 最大的一条（最新的）
-- DELETE FROM game_account_house
-- WHERE id NOT IN (
--     SELECT MAX(id)
--     FROM game_account_house
--     GROUP BY game_account_id
-- );
