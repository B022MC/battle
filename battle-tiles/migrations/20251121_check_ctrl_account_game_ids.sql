-- ============================================================================
-- 检查和修复 game_ctrl_account 表中缺失的 game_player_id 和 game_id
-- ============================================================================
-- 问题：旧的中控账号可能没有 game_player_id 和 game_id 字段的值
-- 导致无法启动会话："game_player_id not set, please verify account first"
-- ============================================================================

-- 步骤1：查看哪些中控账号缺少 game_player_id 或 game_id
SELECT 
    id,
    login_mode,
    identifier,
    CASE 
        WHEN game_player_id = '' OR game_player_id IS NULL THEN '❌ 缺少'
        ELSE game_player_id 
    END as game_player_id_status,
    CASE 
        WHEN game_id = '' OR game_id IS NULL THEN '❌ 缺少'
        ELSE game_id
    END as game_id_status,
    status,
    last_verify_at,
    created_at
FROM game_ctrl_account
WHERE deleted_at IS NULL
ORDER BY id;

-- 步骤2：查看有问题的中控账号详情
SELECT 
    id,
    login_mode,
    identifier,
    pwd_md5,
    game_player_id,
    game_id,
    status,
    last_verify_at
FROM game_ctrl_account
WHERE deleted_at IS NULL
  AND (game_player_id = '' OR game_player_id IS NULL OR game_id = '' OR game_id IS NULL)
ORDER BY id;

-- ============================================================================
-- 解决方案
-- ============================================================================

-- 方案1（推荐）：通过 API 重新验证账号
-- 使用以下 API 重新验证中控账号，系统会自动填充 game_player_id 和 game_id：
-- POST /api/shops/ctrlAccounts/:id/verify
-- 或者
-- POST /api/shops/ctrlAccounts
-- Body: { "login_mode": 1, "identifier": "账号", "pwd_md5": "密码MD5", "status": 1 }

-- 方案2：手动更新（如果你知道正确的值）
-- UPDATE game_ctrl_account
-- SET game_player_id = '6851471',  -- UserID (实际是Plaza API返回的玩家ID)
--     game_id = '22948349',         -- GameID  
--     last_verify_at = NOW(),
--     updated_at = NOW()
-- WHERE id = 1 AND deleted_at IS NULL;

-- 方案3：批量清理无效的中控账号（慎用）
-- 如果某些中控账号已经不再使用，可以删除
-- UPDATE game_ctrl_account
-- SET deleted_at = NOW()
-- WHERE (game_player_id = '' OR game_player_id IS NULL)
--   AND status = 0  -- 只删除已禁用的
--   AND deleted_at IS NULL;

-- ============================================================================
-- 验证修复结果
-- ============================================================================

-- 查看所有中控账号的状态
SELECT 
    id,
    identifier,
    game_player_id,
    game_id,
    status,
    CASE
        WHEN game_player_id != '' AND game_id != '' THEN '✅ 正常'
        WHEN status = 0 THEN '⚠️ 已禁用（缺少ID）'
        ELSE '❌ 需要重新验证'
    END as verification_status,
    last_verify_at
FROM game_ctrl_account
WHERE deleted_at IS NULL
ORDER BY id;

-- ============================================================================
-- 注意事项
-- ============================================================================
-- 1. game_player_id 存储 Plaza API 返回的玩家ID (实际来自 UserID 字段)
-- 2. game_id 存储 Plaza API 返回的游戏ID (实际来自 GameID 字段)
-- 3. 这两个字段在调用 Plaza API 验证账号时自动填充
-- 4. 如果没有这两个字段，会话无法启动
-- 5. 推荐使用方案1通过API重新验证，确保数据准确
-- ============================================================================

-- 示例：账号 22590031 的正确数据（来自你的数据）
-- identifier = '22590031'
-- game_player_id = '6851471'  (UserID from Plaza API)
-- game_id = '22948349'         (GameID from Plaza API)
