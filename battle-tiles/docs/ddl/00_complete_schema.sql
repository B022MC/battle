-- ============================================================================
-- Complete Database Schema for battle-tiles
-- File: 00_complete_schema.sql
-- Description: Complete DDL for all tables (existing + new)
-- Version: 2.0 (After WeChat Dependency Removal Migration)
-- ============================================================================

-- ============================================================================
-- BASIC MODULE - User Authentication and Management
-- ============================================================================

-- ----------------------------------------------------------------------------
-- basic_user - System users with role-based access control
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS basic_user (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL COMMENT 'Login username',
    password VARCHAR(255) NOT NULL COMMENT 'Hashed password',
    salt VARCHAR(50) NOT NULL COMMENT 'Password salt',
    wechat_id VARCHAR(64) DEFAULT NULL COMMENT 'WeChat ID (legacy, optional)',
    nick_name VARCHAR(50) DEFAULT NULL COMMENT 'User display name',
    role VARCHAR(20) NOT NULL DEFAULT 'user' COMMENT 'User role: super_admin, store_admin, user',
    game_nickname VARCHAR(64) DEFAULT NULL COMMENT 'Nickname from game account (for regular users)',
    avatar VARCHAR(255) DEFAULT NULL COMMENT 'Avatar URL',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Indexes
    UNIQUE KEY uk_basic_user_username (username),
    INDEX idx_basic_user_wechat_id (wechat_id),
    INDEX idx_basic_user_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='System users with role-based access control';

-- ============================================================================
-- GAME MODULE - Game Accounts and Store Management
-- ============================================================================

-- ----------------------------------------------------------------------------
-- game_account - Game accounts linked to users
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_account (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL COMMENT 'FK to basic_user',
    account VARCHAR(64) NOT NULL COMMENT 'Game account (phone number)',
    pwd_md5 VARCHAR(64) NOT NULL COMMENT 'MD5 hashed game password',
    nickname VARCHAR(64) DEFAULT NULL COMMENT 'Game nickname',
    ctrl_account_id INT DEFAULT NULL COMMENT 'FK to game_ctrl_account (for super admin)',
    game_user_id VARCHAR(32) DEFAULT '' COMMENT 'Game server user ID',
    verified_at TIMESTAMP WITH TIME ZONE DEFAULT NULL COMMENT 'When account was verified',
    verification_status VARCHAR(20) DEFAULT 'pending' COMMENT 'Status: pending, verified, failed',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT fk_game_account_user FOREIGN KEY (user_id) REFERENCES basic_user(id) ON DELETE CASCADE,
    CONSTRAINT fk_game_account_ctrl FOREIGN KEY (ctrl_account_id) REFERENCES game_ctrl_account(id) ON DELETE SET NULL,
    
    -- Indexes
    INDEX idx_game_account_user_id (user_id),
    INDEX idx_game_account_account (account),
    INDEX idx_game_account_ctrl (ctrl_account_id),
    INDEX idx_game_account_verification (verification_status),
    INDEX idx_game_account_game_user_id (game_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Game accounts linked to system users';

-- ----------------------------------------------------------------------------
-- game_ctrl_account - Control accounts for game server access (super admin)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_ctrl_account (
    id INT PRIMARY KEY AUTO_INCREMENT,
    login_mode INT NOT NULL COMMENT 'Login mode: 1=phone, 2=account',
    identifier VARCHAR(64) NOT NULL COMMENT 'Login identifier (phone or account)',
    pwd_md5 VARCHAR(64) NOT NULL COMMENT 'MD5 hashed password',
    game_user_id VARCHAR(32) DEFAULT '' COMMENT 'Game server user ID',
    game_id VARCHAR(32) DEFAULT '' COMMENT 'Game ID',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Indexes
    UNIQUE KEY uk_ctrl_account_identifier (identifier),
    INDEX idx_ctrl_account_game_user_id (game_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Control accounts for super admin game server access';

-- ----------------------------------------------------------------------------
-- game_account_house - Game account to house/store bindings (regular users)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_account_house (
    id INT PRIMARY KEY AUTO_INCREMENT,
    game_account_id INT NOT NULL COMMENT 'FK to game_account',
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    is_default BOOLEAN NOT NULL DEFAULT false COMMENT 'Is default store for this account',
    status INT NOT NULL DEFAULT 1 COMMENT 'Status: 1=active, 0=inactive',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT fk_gah_game_account FOREIGN KEY (game_account_id) REFERENCES game_account(id) ON DELETE CASCADE,
    
    -- Indexes
    INDEX idx_gah_game_account (game_account_id),
    INDEX idx_gah_house_gid (house_gid),
    INDEX idx_gah_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Game account to store bindings for regular users';

-- ----------------------------------------------------------------------------
-- game_ctrl_account_house - Control account to house/store bindings (super admin)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_ctrl_account_house (
    id INT PRIMARY KEY AUTO_INCREMENT,
    game_account_id INT NOT NULL COMMENT 'FK to game_ctrl_account',
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    is_default BOOLEAN NOT NULL DEFAULT false COMMENT 'Is default store',
    status INT NOT NULL DEFAULT 1 COMMENT 'Status: 1=active, 0=inactive',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT fk_gcah_ctrl_account FOREIGN KEY (game_account_id) REFERENCES game_ctrl_account(id) ON DELETE CASCADE,
    
    -- Indexes
    INDEX idx_gcah_game_account (game_account_id),
    INDEX idx_gcah_house_gid (house_gid),
    INDEX idx_gcah_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Control account to store bindings for super admin';

-- ----------------------------------------------------------------------------
-- game_account_store_binding - NEW: Explicit game account to store bindings
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_account_store_binding (
    id INT PRIMARY KEY AUTO_INCREMENT,
    game_account_id INT NOT NULL COMMENT 'FK to game_account',
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    bound_by_user_id INT NOT NULL COMMENT 'FK to basic_user - who created binding',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT 'Status: active, inactive',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT fk_gasb_game_account FOREIGN KEY (game_account_id) REFERENCES game_account(id) ON DELETE CASCADE,
    CONSTRAINT fk_gasb_bound_by_user FOREIGN KEY (bound_by_user_id) REFERENCES basic_user(id) ON DELETE CASCADE,
    
    -- Indexes
    UNIQUE KEY uk_game_account_house (game_account_id, house_gid),
    INDEX idx_gasb_house_gid (house_gid),
    INDEX idx_gasb_bound_by_user (bound_by_user_id),
    INDEX idx_gasb_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Game account to store binding - one game account can only bind to one store';

-- ============================================================================
-- STORE ADMINISTRATION MODULE
-- ============================================================================

-- ----------------------------------------------------------------------------
-- game_shop_admin - Store administrators
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_shop_admin (
    id INT PRIMARY KEY AUTO_INCREMENT,
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    user_id INT NOT NULL COMMENT 'FK to basic_user',
    role VARCHAR(20) NOT NULL COMMENT 'Role: admin, operator',
    game_account_id INT DEFAULT NULL COMMENT 'FK to game_account - which game account admin uses',
    is_exclusive BOOLEAN DEFAULT true COMMENT 'Exclusive binding (admin can only admin one store)',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT fk_shop_admin_user FOREIGN KEY (user_id) REFERENCES basic_user(id) ON DELETE CASCADE,
    CONSTRAINT fk_shop_admin_game_account FOREIGN KEY (game_account_id) REFERENCES game_account(id) ON DELETE SET NULL,
    
    -- Indexes
    INDEX idx_shop_admin_house_gid (house_gid),
    INDEX idx_shop_admin_user_id (user_id),
    INDEX idx_shop_admin_game_account (game_account_id),
    INDEX idx_shop_admin_exclusive (user_id, is_exclusive)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Store administrators with role-based permissions';

-- ----------------------------------------------------------------------------
-- game_shop_group_admin - Group administrators within stores
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_shop_group_admin (
    id INT PRIMARY KEY AUTO_INCREMENT,
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    group_id INT NOT NULL COMMENT 'Group ID within store',
    user_id INT NOT NULL COMMENT 'FK to basic_user',
    role VARCHAR(20) NOT NULL COMMENT 'Role: admin, operator',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    
    -- Constraints
    CONSTRAINT fk_group_admin_user FOREIGN KEY (user_id) REFERENCES basic_user(id) ON DELETE CASCADE,
    
    -- Indexes
    INDEX idx_group_admin_house_gid (house_gid),
    INDEX idx_group_admin_group_id (group_id),
    INDEX idx_group_admin_user_id (user_id),
    INDEX idx_group_admin_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Group administrators within stores';

-- ============================================================================
-- SESSION AND SYNC MODULE
-- ============================================================================

-- ----------------------------------------------------------------------------
-- game_session - Active game sessions with auto-sync
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_session (
    id INT PRIMARY KEY AUTO_INCREMENT,
    game_ctrl_account_id INT NOT NULL COMMENT 'FK to game_ctrl_account',
    user_id INT NOT NULL COMMENT 'FK to basic_user',
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    state VARCHAR(20) NOT NULL COMMENT 'Session state: active, inactive, error',
    device_ip VARCHAR(45) DEFAULT NULL COMMENT 'Device IP address',
    auto_sync_enabled BOOLEAN DEFAULT true COMMENT 'Auto-sync enabled',
    last_sync_at TIMESTAMP WITH TIME ZONE DEFAULT NULL COMMENT 'Last successful sync',
    sync_status VARCHAR(20) DEFAULT 'idle' COMMENT 'Sync status: idle, syncing, error',
    game_account_id INT DEFAULT NULL COMMENT 'FK to game_account',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT fk_session_ctrl_account FOREIGN KEY (game_ctrl_account_id) REFERENCES game_ctrl_account(id) ON DELETE CASCADE,
    CONSTRAINT fk_session_user FOREIGN KEY (user_id) REFERENCES basic_user(id) ON DELETE CASCADE,
    CONSTRAINT fk_session_game_account FOREIGN KEY (game_account_id) REFERENCES game_account(id) ON DELETE CASCADE,
    
    -- Indexes
    INDEX idx_session_ctrl_account (game_ctrl_account_id),
    INDEX idx_session_user_id (user_id),
    INDEX idx_session_house_gid (house_gid),
    INDEX idx_session_state (state),
    INDEX idx_session_sync_status (sync_status),
    INDEX idx_session_game_account (game_account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Active game sessions with auto-sync capabilities';

-- ----------------------------------------------------------------------------
-- game_sync_log - NEW: Synchronization operation logs
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_sync_log (
    id INT PRIMARY KEY AUTO_INCREMENT,
    session_id INT NOT NULL COMMENT 'FK to game_session',
    sync_type VARCHAR(20) NOT NULL COMMENT 'Type: battle_record, member_list, wallet_update',
    status VARCHAR(20) NOT NULL COMMENT 'Status: success, failed, partial',
    records_synced INT NOT NULL DEFAULT 0 COMMENT 'Number of records synced',
    error_message TEXT DEFAULT NULL COMMENT 'Error details if failed',
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,

    -- Constraints
    CONSTRAINT fk_sync_log_session FOREIGN KEY (session_id) REFERENCES game_session(id) ON DELETE CASCADE,

    -- Indexes
    INDEX idx_sync_log_session_started (session_id, started_at),
    INDEX idx_sync_log_status (status),
    INDEX idx_sync_log_type (sync_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Synchronization operation logs';

-- ============================================================================
-- GAME MEMBERS AND WALLET MODULE
-- ============================================================================

-- ----------------------------------------------------------------------------
-- game_member - NEW: Game members (players) within stores
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_member (
    id INT PRIMARY KEY AUTO_INCREMENT,
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    game_id INT NOT NULL COMMENT 'Player game ID from game server',
    game_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Player game nickname',
    group_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Group/circle name',
    balance INT NOT NULL DEFAULT 0 COMMENT 'Player balance in cents',
    credit INT NOT NULL DEFAULT 0 COMMENT 'Credit limit in cents',
    forbid BOOLEAN NOT NULL DEFAULT false COMMENT 'Whether player is forbidden',
    recommender VARCHAR(64) DEFAULT '' COMMENT 'Recommender name',
    use_multi_gids BOOLEAN DEFAULT false COMMENT 'Allow multiple game IDs',
    active_gid INT DEFAULT NULL COMMENT 'Active game ID when multi-gid disabled',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Indexes
    UNIQUE KEY uk_game_member_house_game (house_gid, game_id),
    INDEX idx_game_member_house_group (house_gid, group_name),
    INDEX idx_game_member_game_id (game_id),
    INDEX idx_game_member_forbid (forbid),
    INDEX idx_game_member_balance (balance)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Game members (players) within stores';

-- ----------------------------------------------------------------------------
-- game_member_wallet - Member wallet management
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_member_wallet (
    id INT PRIMARY KEY AUTO_INCREMENT,
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    member_id INT NOT NULL COMMENT 'FK to game_member',
    balance INT NOT NULL DEFAULT 0 COMMENT 'Current balance in cents',
    forbid BOOLEAN NOT NULL DEFAULT false COMMENT 'Whether wallet is frozen',
    limit_min INT NOT NULL DEFAULT 0 COMMENT 'Minimum balance limit',
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by INT NOT NULL DEFAULT 0 COMMENT 'FK to basic_user - who updated',

    -- Constraints
    CONSTRAINT fk_wallet_member FOREIGN KEY (member_id) REFERENCES game_member(id) ON DELETE CASCADE,

    -- Indexes
    INDEX idx_wallet_house_gid (house_gid),
    INDEX idx_wallet_member_id (member_id),
    INDEX idx_wallet_forbid (forbid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Member wallet management';

-- ----------------------------------------------------------------------------
-- game_recharge_record - NEW: Wallet recharge/withdrawal records
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_recharge_record (
    id INT PRIMARY KEY AUTO_INCREMENT,
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    player_id INT NOT NULL COMMENT 'FK to game_member',
    group_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Group name',
    amount INT NOT NULL COMMENT 'Amount in cents (positive=deposit, negative=withdrawal)',
    balance_before INT NOT NULL COMMENT 'Balance before transaction',
    balance_after INT NOT NULL COMMENT 'Balance after transaction',
    operator_user_id INT DEFAULT NULL COMMENT 'FK to basic_user - operator',
    recharged_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_recharge_player FOREIGN KEY (player_id) REFERENCES game_member(id) ON DELETE CASCADE,
    CONSTRAINT fk_recharge_operator FOREIGN KEY (operator_user_id) REFERENCES basic_user(id) ON DELETE SET NULL,

    -- Indexes
    INDEX idx_recharge_house_gid (house_gid),
    INDEX idx_recharge_player (player_id),
    INDEX idx_recharge_recharged_at (recharged_at),
    INDEX idx_recharge_house_group (house_gid, group_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Wallet recharge and withdrawal records';

-- ============================================================================
-- BATTLE RECORDS MODULE
-- ============================================================================

-- ----------------------------------------------------------------------------
-- game_battle_record - Battle/game records (player dimension)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_battle_record (
    id INT PRIMARY KEY AUTO_INCREMENT,
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    group_id INT NOT NULL COMMENT 'Group ID',
    room_uid INT NOT NULL COMMENT 'Room unique ID (MappedNum)',
    kind_id INT NOT NULL COMMENT 'Game kind/type ID',
    base_score INT NOT NULL COMMENT 'Base score',
    battle_at TIMESTAMP WITH TIME ZONE NOT NULL COMMENT 'Battle timestamp',
    players_json TEXT NOT NULL COMMENT 'Players data in JSON format',
    player_id INT DEFAULT NULL COMMENT 'Internal player/member ID',
    player_game_id INT DEFAULT NULL COMMENT 'Player game ID',
    player_game_name VARCHAR(64) DEFAULT '' COMMENT 'Player game nickname',
    group_name VARCHAR(64) DEFAULT '' COMMENT 'Group name',
    score INT DEFAULT 0 COMMENT 'Player score (positive=win, negative=loss)',
    fee INT DEFAULT 0 COMMENT 'Service fee in cents',
    factor DECIMAL(10,4) DEFAULT 1.0000 COMMENT 'Settlement factor',
    player_balance INT DEFAULT 0 COMMENT 'Player balance after settlement',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_battle_player FOREIGN KEY (player_id) REFERENCES game_member(id) ON DELETE SET NULL,

    -- Indexes
    INDEX idx_battle_house_gid (house_gid),
    INDEX idx_battle_room_uid (room_uid),
    INDEX idx_battle_at (battle_at),
    INDEX idx_battle_player_id (player_id),
    INDEX idx_battle_player_game_id (player_game_id),
    INDEX idx_battle_group (house_gid, group_name),
    INDEX idx_battle_kind_id (kind_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Battle/game records organized by player dimension';

-- ----------------------------------------------------------------------------
-- game_fee_record - NEW: Service fee records
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_fee_record (
    id INT PRIMARY KEY AUTO_INCREMENT,
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    room_id INT NOT NULL COMMENT 'Room/table ID',
    group_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Group name',
    amount INT NOT NULL COMMENT 'Fee amount in cents',
    fee_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Indexes
    INDEX idx_fee_house_gid (house_gid),
    INDEX idx_fee_room_id (room_id),
    INDEX idx_fee_at (fee_at),
    INDEX idx_fee_house_group (house_gid, group_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Service fee records per room/table';

-- ----------------------------------------------------------------------------
-- game_fee_settle - NEW: Fee settlement records
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_fee_settle (
    id INT PRIMARY KEY AUTO_INCREMENT,
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    room_id INT NOT NULL COMMENT 'Room/table ID',
    group_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Group name',
    amount INT NOT NULL COMMENT 'Settlement amount in cents',
    fee_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Indexes
    INDEX idx_fee_settle_house_gid (house_gid),
    INDEX idx_fee_settle_room_id (room_id),
    INDEX idx_fee_settle_at (fee_at),
    INDEX idx_fee_settle_house_group (house_gid, group_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Fee settlement records';

-- ============================================================================
-- NOTIFICATION AND STATISTICS MODULE
-- ============================================================================

-- ----------------------------------------------------------------------------
-- game_notice - NEW: Notification/notice records
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_notice (
    id INT PRIMARY KEY AUTO_INCREMENT,
    type INT NOT NULL COMMENT 'Notice type code',
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    player_id INT DEFAULT NULL COMMENT 'FK to game_member',
    group_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Group name',
    manager_uuid VARCHAR(64) DEFAULT '' COMMENT 'Manager UUID (legacy)',
    text TEXT NOT NULL COMMENT 'Notice text content',
    create_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_notice_player FOREIGN KEY (player_id) REFERENCES game_member(id) ON DELETE CASCADE,

    -- Indexes
    INDEX idx_notice_type (type),
    INDEX idx_notice_house_gid (house_gid),
    INDEX idx_notice_player (player_id),
    INDEX idx_notice_group (house_gid, group_name),
    INDEX idx_notice_create_time (create_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Notification and notice records';

-- ----------------------------------------------------------------------------
-- game_push_stat - NEW: Push notification statistics
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_push_stat (
    id INT PRIMARY KEY AUTO_INCREMENT,
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    date DATE NOT NULL COMMENT 'Statistics date',
    game_kind INT NOT NULL COMMENT 'Game type/kind ID',
    base_score INT NOT NULL COMMENT 'Base score',
    count INT NOT NULL DEFAULT 0 COMMENT 'Number of pushes',
    amount BIGINT NOT NULL DEFAULT 0 COMMENT 'Total amount',

    -- Indexes
    UNIQUE KEY uk_push_stat (house_gid, date, game_kind, base_score),
    INDEX idx_push_stat_house_gid (house_gid),
    INDEX idx_push_stat_date (date),
    INDEX idx_push_stat_game_kind (game_kind)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Push notification statistics';

-- ----------------------------------------------------------------------------
-- game_deleted_member - NEW: Soft-deleted members archive
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS game_deleted_member (
    id INT PRIMARY KEY AUTO_INCREMENT,
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    group_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Group name',
    player_id INT NOT NULL COMMENT 'Original player ID',
    game_id INT NOT NULL COMMENT 'Player game ID',
    game_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Player game nickname',
    game_ids VARCHAR(255) DEFAULT '' COMMENT 'Comma-separated game IDs',
    balance INT NOT NULL DEFAULT 0 COMMENT 'Balance at deletion',
    reason VARCHAR(255) DEFAULT '' COMMENT 'Deletion reason',
    deleted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Indexes
    INDEX idx_deleted_member_house_gid (house_gid),
    INDEX idx_deleted_member_group (house_gid, group_name),
    INDEX idx_deleted_member_player_id (player_id),
    INDEX idx_deleted_member_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Archive of deleted members for audit';

-- ============================================================================
-- End of Complete Schema
-- ============================================================================

