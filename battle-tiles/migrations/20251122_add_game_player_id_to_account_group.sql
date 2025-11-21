-- 快速修复：添加 game_player_id 字段到 game_account_group 表
-- 保留原有 game_account_id 字段，实现平滑迁移

BEGIN;

-- 1. 添加新字段 game_player_id
ALTER TABLE game_account_group 
ADD COLUMN IF NOT EXISTS game_player_id VARCHAR(32);

-- 2. 从 game_account 表迁移数据填充 game_player_id
UPDATE game_account_group gag
SET game_player_id = ga.game_player_id
FROM game_account ga
WHERE gag.game_account_id = ga.id
  AND gag.game_player_id IS NULL
  AND ga.game_player_id IS NOT NULL
  AND ga.game_player_id != '';

-- 3. 创建索引
CREATE INDEX IF NOT EXISTS idx_game_account_group_game_player_id 
ON game_account_group(game_player_id);

CREATE INDEX IF NOT EXISTS idx_game_account_group_player_house 
ON game_account_group(game_player_id, house_gid);

-- 4. 添加注释
COMMENT ON COLUMN game_account_group.game_player_id IS '游戏玩家ID（来自游戏内部）';

-- 5. 统计信息
DO $$
DECLARE
    total_count INTEGER;
    filled_count INTEGER;
    empty_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO total_count FROM game_account_group;
    SELECT COUNT(*) INTO filled_count FROM game_account_group WHERE game_player_id IS NOT NULL;
    SELECT COUNT(*) INTO empty_count FROM game_account_group WHERE game_player_id IS NULL;
    
    RAISE NOTICE '==============================================';
    RAISE NOTICE '迁移完成统计：';
    RAISE NOTICE '总记录数: %', total_count;
    RAISE NOTICE '已填充 game_player_id: %', filled_count;
    RAISE NOTICE '未填充 game_player_id: %', empty_count;
    RAISE NOTICE '==============================================';
    
    IF empty_count > 0 THEN
        RAISE WARNING '有 % 条记录的 game_player_id 为空，这些记录的 game_account 可能缺少 game_player_id', empty_count;
    END IF;
END $$;

COMMIT;

-- 使用说明：
-- 1. 此脚本添加 game_player_id 字段但不删除 game_account_id
-- 2. 两个字段可以同时存在，实现平滑过渡
-- 3. 旧代码和新代码可以共存一段时间
-- 4. 确认稳定后可以考虑删除 game_account_id 字段
