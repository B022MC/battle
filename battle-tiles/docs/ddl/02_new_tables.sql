-- ============================================================================
-- New Tables for WeChat Dependency Removal
-- File: 02_new_tables.sql
-- Description: DDL statements for creating new tables
-- ============================================================================

-- ============================================================================
-- 1. game_account_store_binding - Track game account to store bindings
-- ============================================================================

CREATE TABLE IF NOT EXISTS game_account_store_binding (
    id INT PRIMARY KEY AUTO_INCREMENT,
    game_account_id INT NOT NULL COMMENT 'FK to game_account',
    house_gid INT NOT NULL COMMENT 'Store/House game ID from game server',
    bound_by_user_id INT NOT NULL COMMENT 'FK to basic_user - who created this binding',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT 'Binding status: active, inactive',
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
COMMENT='Game account to store binding table - one game account can only bind to one store';

-- ============================================================================
-- 2. game_member - Track game members (players) within stores
-- ============================================================================

CREATE TABLE IF NOT EXISTS game_member (
    id INT PRIMARY KEY AUTO_INCREMENT,
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    game_id INT NOT NULL COMMENT 'Player game ID from game server',
    game_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Player game nickname',
    group_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Group/circle name within store',
    balance INT NOT NULL DEFAULT 0 COMMENT 'Player balance in cents',
    credit INT NOT NULL DEFAULT 0 COMMENT 'Credit limit in cents',
    forbid BOOLEAN NOT NULL DEFAULT false COMMENT 'Whether player is forbidden/banned',
    recommender VARCHAR(64) DEFAULT '' COMMENT 'Recommender name',
    use_multi_gids BOOLEAN DEFAULT false COMMENT 'Allow multiple game IDs',
    active_gid INT DEFAULT NULL COMMENT 'Currently active game ID when multi-gid disabled',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Indexes
    UNIQUE KEY uk_game_member_house_game (house_gid, game_id),
    INDEX idx_game_member_house_group (house_gid, group_name),
    INDEX idx_game_member_game_id (game_id),
    INDEX idx_game_member_forbid (forbid),
    INDEX idx_game_member_balance (balance)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Game members (players) within stores - tracks player info and wallet';

-- ============================================================================
-- 3. game_sync_log - Track synchronization operations
-- ============================================================================

CREATE TABLE IF NOT EXISTS game_sync_log (
    id INT PRIMARY KEY AUTO_INCREMENT,
    session_id INT NOT NULL COMMENT 'FK to game_session',
    sync_type VARCHAR(20) NOT NULL COMMENT 'Type: battle_record, member_list, wallet_update, etc.',
    status VARCHAR(20) NOT NULL COMMENT 'Status: success, failed, partial',
    records_synced INT NOT NULL DEFAULT 0 COMMENT 'Number of records successfully synced',
    error_message TEXT DEFAULT NULL COMMENT 'Error details if sync failed',
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    
    -- Constraints
    CONSTRAINT fk_sync_log_session FOREIGN KEY (session_id) REFERENCES game_session(id) ON DELETE CASCADE,
    
    -- Indexes
    INDEX idx_sync_log_session_started (session_id, started_at),
    INDEX idx_sync_log_status (status),
    INDEX idx_sync_log_type (sync_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Synchronization operation logs - tracks sync status and errors';

-- ============================================================================
-- 4. game_recharge_record - Wallet recharge/withdrawal records
-- ============================================================================

CREATE TABLE IF NOT EXISTS game_recharge_record (
    id INT PRIMARY KEY AUTO_INCREMENT,
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    player_id INT NOT NULL COMMENT 'FK to game_member',
    group_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Group name',
    amount INT NOT NULL COMMENT 'Recharge amount in cents (positive = deposit, negative = withdrawal)',
    balance_before INT NOT NULL COMMENT 'Balance before transaction',
    balance_after INT NOT NULL COMMENT 'Balance after transaction',
    operator_user_id INT DEFAULT NULL COMMENT 'FK to basic_user - who performed the operation',
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
-- 5. game_fee_record - Service fee records
-- ============================================================================

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

-- ============================================================================
-- 6. game_fee_settle - Fee settlement records
-- ============================================================================

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
-- 7. game_deleted_member - Soft-deleted members archive
-- ============================================================================

CREATE TABLE IF NOT EXISTS game_deleted_member (
    id INT PRIMARY KEY AUTO_INCREMENT,
    house_gid INT NOT NULL COMMENT 'Store/House game ID',
    group_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Group name',
    player_id INT NOT NULL COMMENT 'Original player ID',
    game_id INT NOT NULL COMMENT 'Player game ID',
    game_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Player game nickname',
    game_ids VARCHAR(255) DEFAULT '' COMMENT 'Comma-separated list of all game IDs',
    balance INT NOT NULL DEFAULT 0 COMMENT 'Balance at deletion time',
    reason VARCHAR(255) DEFAULT '' COMMENT 'Deletion reason',
    deleted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_deleted_member_house_gid (house_gid),
    INDEX idx_deleted_member_group (house_gid, group_name),
    INDEX idx_deleted_member_player_id (player_id),
    INDEX idx_deleted_member_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Archive of deleted members for audit purposes';

-- ============================================================================
-- 8. game_notice - Notification/notice records
-- ============================================================================

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

-- ============================================================================
-- 9. game_push_stat - Push notification statistics
-- ============================================================================

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
COMMENT='Push notification statistics by date and game type';

-- ============================================================================
-- End of new tables
-- ============================================================================

