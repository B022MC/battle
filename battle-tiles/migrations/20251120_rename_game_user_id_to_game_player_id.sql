-- 重命名 game_user_id 为 game_player_id 以避免字段命名歧义
-- 
-- 歧义说明：
--   - user_id: 系统用户ID（basic_user.id）
--   - game_user_id: 游戏玩家ID（容易误解为"游戏账号的用户ID"）
--   - game_player_id: 游戏玩家ID（更清晰的命名）

BEGIN;

-- ============================================================================
-- 1. 重命名 game_account 表的 game_user_id 字段
-- ============================================================================
ALTER TABLE game_account 
RENAME COLUMN game_user_id TO game_player_id;

-- 更新字段注释
COMMENT ON COLUMN game_account.game_player_id IS '游戏服务器返回的玩家ID（用于标识游戏中的玩家）';

-- 重命名索引
ALTER INDEX idx_game_account_game_user_id RENAME TO idx_game_account_game_player_id;

-- ============================================================================
-- 2. 重命名 game_ctrl_account 表的 game_user_id 字段
-- ============================================================================
ALTER TABLE game_ctrl_account 
RENAME COLUMN game_user_id TO game_player_id;

-- 更新字段注释
COMMENT ON COLUMN game_ctrl_account.game_player_id IS '游戏服务器返回的玩家ID';

-- 重命名索引
ALTER INDEX idx_ctrl_account_game_user_id RENAME TO idx_ctrl_account_game_player_id;

-- ============================================================================
-- 验证
-- ============================================================================
DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '=== 字段重命名完成 ===';
    RAISE NOTICE 'game_account.game_user_id => game_account.game_player_id';
    RAISE NOTICE 'game_ctrl_account.game_user_id => game_ctrl_account.game_player_id';
    RAISE NOTICE '';
END $$;

COMMIT;
