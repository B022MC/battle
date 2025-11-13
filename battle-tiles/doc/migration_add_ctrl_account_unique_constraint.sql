-- Migration: Add unique constraint to game_ctrl_account table
-- Date: 2025-11-12
-- Description: Add unique constraint on (login_mode, identifier) to support ON CONFLICT clause

-- 添加唯一约束到 game_ctrl_account 表
-- 这个约束确保同一登录方式下的账号标识是唯一的
CREATE UNIQUE INDEX IF NOT EXISTS uk_ctrl_account_login_identifier 
ON game_ctrl_account(login_mode, identifier);

-- 验证约束是否创建成功
SELECT 
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = 'game_ctrl_account' 
  AND indexname = 'uk_ctrl_account_login_identifier';

