-- ============================================================================
-- Schema Modifications for WeChat Dependency Removal
-- File: 01_schema_modifications.sql
-- Description: DDL statements for modifying existing tables
-- ============================================================================

-- ============================================================================
-- 1. Modify basic_user table - Add role and game nickname
-- ============================================================================

ALTER TABLE basic_user
ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user' COMMENT 'User role: super_admin, store_admin, user',
ADD COLUMN game_nickname VARCHAR(64) DEFAULT NULL COMMENT 'Nickname from game account (for regular users)',
ADD INDEX idx_basic_user_role (role);

COMMENT ON COLUMN basic_user.role IS 'User role: super_admin (can manage multiple game accounts), store_admin (exclusive to one store), user (regular user)';
COMMENT ON COLUMN basic_user.game_nickname IS 'Game nickname retrieved during registration for regular users';

-- ============================================================================
-- 2. Modify game_account table - Add verification fields
-- ============================================================================

ALTER TABLE game_account
ADD COLUMN game_user_id VARCHAR(32) DEFAULT '' COMMENT 'Game server user ID',
ADD COLUMN verified_at TIMESTAMP WITH TIME ZONE DEFAULT NULL COMMENT 'When account was verified with game server',
ADD COLUMN verification_status VARCHAR(20) DEFAULT 'pending' COMMENT 'Verification status: pending, verified, failed',
ADD INDEX idx_game_account_verification (verification_status),
ADD INDEX idx_game_account_game_user_id (game_user_id);

COMMENT ON COLUMN game_account.game_user_id IS 'User ID from game server (retrieved during verification)';
COMMENT ON COLUMN game_account.verified_at IS 'Timestamp when account was successfully verified with game server';
COMMENT ON COLUMN game_account.verification_status IS 'Status: pending (not verified), verified (successfully verified), failed (verification failed)';

-- ============================================================================
-- 3. Modify game_shop_admin table - Add game account reference
-- ============================================================================

ALTER TABLE game_shop_admin
ADD COLUMN game_account_id INT DEFAULT NULL COMMENT 'FK to game_account - which game account this admin uses',
ADD COLUMN is_exclusive BOOLEAN DEFAULT true COMMENT 'Whether this is an exclusive binding (store admin can only admin one store)',
ADD CONSTRAINT fk_game_shop_admin_game_account FOREIGN KEY (game_account_id) REFERENCES game_account(id) ON DELETE SET NULL,
ADD INDEX idx_game_shop_admin_game_account (game_account_id),
ADD INDEX idx_game_shop_admin_exclusive (user_id, is_exclusive);

COMMENT ON COLUMN game_shop_admin.game_account_id IS 'Game account used by this admin for store management';
COMMENT ON COLUMN game_shop_admin.is_exclusive IS 'If true, user can only be admin of ONE store under ONE game account';

-- ============================================================================
-- 4. Modify game_session table - Add auto-sync fields
-- ============================================================================

ALTER TABLE game_session
ADD COLUMN auto_sync_enabled BOOLEAN DEFAULT true COMMENT 'Whether auto-sync is enabled for this session',
ADD COLUMN last_sync_at TIMESTAMP WITH TIME ZONE DEFAULT NULL COMMENT 'Last successful sync time',
ADD COLUMN sync_status VARCHAR(20) DEFAULT 'idle' COMMENT 'Sync status: idle, syncing, error',
ADD COLUMN game_account_id INT DEFAULT NULL COMMENT 'FK to game_account',
ADD CONSTRAINT fk_game_session_game_account FOREIGN KEY (game_account_id) REFERENCES game_account(id) ON DELETE CASCADE,
ADD INDEX idx_game_session_sync_status (sync_status),
ADD INDEX idx_game_session_game_account (game_account_id);

COMMENT ON COLUMN game_session.auto_sync_enabled IS 'If true, session will automatically sync game records';
COMMENT ON COLUMN game_session.last_sync_at IS 'Timestamp of last successful synchronization';
COMMENT ON COLUMN game_session.sync_status IS 'Current sync status: idle (not syncing), syncing (in progress), error (sync failed)';
COMMENT ON COLUMN game_session.game_account_id IS 'Game account associated with this session';

-- ============================================================================
-- 5. Modify game_battle_record table - Add player dimension fields
-- ============================================================================

ALTER TABLE game_battle_record
ADD COLUMN player_id INT DEFAULT NULL COMMENT 'Internal player/member ID',
ADD COLUMN player_game_id INT DEFAULT NULL COMMENT 'Player game ID from game server',
ADD COLUMN player_game_name VARCHAR(64) DEFAULT '' COMMENT 'Player game nickname',
ADD COLUMN group_name VARCHAR(64) DEFAULT '' COMMENT 'Group name within store',
ADD COLUMN score INT DEFAULT 0 COMMENT 'Player score/points in this battle',
ADD COLUMN fee INT DEFAULT 0 COMMENT 'Service fee charged',
ADD COLUMN factor DECIMAL(10,4) DEFAULT 1.0000 COMMENT 'Settlement factor',
ADD COLUMN player_balance INT DEFAULT 0 COMMENT 'Player balance after settlement',
ADD INDEX idx_game_battle_record_player_id (player_id),
ADD INDEX idx_game_battle_record_player_game_id (player_game_id),
ADD INDEX idx_game_battle_record_group (house_gid, group_name),
ADD INDEX idx_game_battle_record_battle_at (battle_at);

COMMENT ON COLUMN game_battle_record.player_id IS 'Internal member ID (FK to game_member)';
COMMENT ON COLUMN game_battle_record.player_game_id IS 'Player game ID from game server';
COMMENT ON COLUMN game_battle_record.player_game_name IS 'Player nickname from game server';
COMMENT ON COLUMN game_battle_record.group_name IS 'Group/circle name within the store';
COMMENT ON COLUMN game_battle_record.score IS 'Player score in this battle (positive = win, negative = loss)';
COMMENT ON COLUMN game_battle_record.fee IS 'Service fee charged for this battle (in cents)';
COMMENT ON COLUMN game_battle_record.factor IS 'Settlement factor applied to score';
COMMENT ON COLUMN game_battle_record.player_balance IS 'Player wallet balance after this battle settlement (in cents)';

-- ============================================================================
-- 6. Add constraints for business rules
-- ============================================================================

-- Constraint: Store admin can only be admin of ONE store when is_exclusive = true
-- This will be enforced at application level with validation:
-- SELECT COUNT(*) FROM game_shop_admin WHERE user_id = ? AND is_exclusive = true
-- Result must be <= 1

-- Constraint: Regular users must have at least one verified game account
-- This will be enforced at application level during registration and account deletion

-- ============================================================================
-- End of schema modifications
-- ============================================================================

