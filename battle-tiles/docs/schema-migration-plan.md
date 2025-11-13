# Schema Migration Plan: WeChat Dependency Removal

## Overview
Migrate functionality from passing-dragonfly/waiter to battle-tiles, eliminating WeChat dependency and implementing a multi-level account binding system.

## Current Schema Analysis

### Existing Tables in battle-tiles:
1. **basic_user** - System users with username/password auth
2. **game_account** - Game accounts linked to users
3. **game_ctrl_account** - Control accounts for game server access
4. **game_ctrl_account_house** - Binding between ctrl accounts and stores
5. **game_shop_admin** - Store administrators
6. **game_shop_group_admin** - Group administrators within stores
7. **game_session** - Active game sessions
8. **game_battle_record** - Battle/game records
9. **game_member_wallet** - Member wallet balances

### Legacy Tables in passing-dragonfly/waiter:
1. **player** - Players with WeChat keys
2. **gamer** - Player-game ID mappings
3. **house** - Stores/茶馆 with game accounts
4. **manager** - Store managers with WeChat keys
5. **battle_record** - Battle records (per player)
6. **recharge_record** - Wallet transactions
7. **fee_record** - Fee records

## Requirements

### 1. User Role System
- **Super Administrator**: Can bind multiple game accounts, each to one store ID
- **Store Administrator**: Exclusive to ONE store under ONE game account
- **Regular User**: MUST bind a game account during registration

### 2. Account Binding Rules
- Super Admin → Multiple Game Accounts → Each Game Account → One Store ID
- Store Admin → ONE Store under ONE Game Account (exclusive)
- Regular User → Must have at least one game account bound

### 3. Registration Flow
- Regular users MUST provide game account credentials during registration
- System validates game account via game server API
- System retrieves game nickname and stores as user profile nickname
- Game account binding is created automatically

### 4. Auto-sync Store Session
- When super admin adds game account + store binding
- System automatically creates a game session
- Session connects to game server (plaza.Session equivalent)
- Begins automatic synchronization of game records
- Data organized by player dimension

## Schema Changes Required

### 1. Modify `basic_user` Table
**Add Fields:**
- `role` VARCHAR(20) NOT NULL DEFAULT 'user' - User role (super_admin, store_admin, user)
- `game_nickname` VARCHAR(64) - Nickname from game account (for regular users)

**Indexes:**
- Add index on `role`

### 2. Modify `game_account` Table
**Add Fields:**
- `game_user_id` VARCHAR(32) - Game server user ID
- `verified_at` TIMESTAMP - When account was verified with game server
- `verification_status` VARCHAR(20) DEFAULT 'pending' - pending, verified, failed

**Business Rules:**
- Regular users must have at least one verified game account
- Super admins can have multiple game accounts

### 3. Create `game_account_store_binding` Table
**Purpose:** Track which game accounts are bound to which stores

**Fields:**
- `id` INT PRIMARY KEY AUTO_INCREMENT
- `game_account_id` INT NOT NULL - FK to game_account
- `house_gid` INT NOT NULL - Store/House game ID
- `bound_by_user_id` INT NOT NULL - FK to basic_user (who created binding)
- `status` VARCHAR(20) DEFAULT 'active' - active, inactive
- `created_at` TIMESTAMP
- `updated_at` TIMESTAMP

**Indexes:**
- UNIQUE(game_account_id, house_gid)
- INDEX(house_gid)
- INDEX(bound_by_user_id)

**Constraints:**
- One game_account can only bind to one house_gid
- Only super_admin can create bindings

### 4. Modify `game_shop_admin` Table
**Add Fields:**
- `game_account_id` INT - FK to game_account (which game account this admin uses)
- `is_exclusive` BOOLEAN DEFAULT true - Whether this is an exclusive binding

**Business Rules:**
- Store admin can only be admin for ONE store under ONE game account
- Check: SELECT COUNT(*) FROM game_shop_admin WHERE user_id = ? AND is_exclusive = true MUST be <= 1

### 5. Modify `game_session` Table
**Add Fields:**
- `auto_sync_enabled` BOOLEAN DEFAULT true - Whether auto-sync is enabled
- `last_sync_at` TIMESTAMP - Last successful sync time
- `sync_status` VARCHAR(20) - idle, syncing, error
- `game_account_id` INT - FK to game_account

**Purpose:**
- Track active sessions for auto-sync
- Monitor sync status and health

### 6. Create `game_member` Table
**Purpose:** Track game members (players) within stores

**Fields:**
- `id` INT PRIMARY KEY AUTO_INCREMENT
- `house_gid` INT NOT NULL - Store ID
- `game_id` INT NOT NULL - Player's game ID
- `game_name` VARCHAR(64) - Player's game nickname
- `group_name` VARCHAR(64) - Group within store
- `balance` INT DEFAULT 0 - Player balance (in cents)
- `credit` INT DEFAULT 0 - Credit limit
- `forbid` BOOLEAN DEFAULT false - Whether player is forbidden
- `created_at` TIMESTAMP
- `updated_at` TIMESTAMP

**Indexes:**
- UNIQUE(house_gid, game_id)
- INDEX(house_gid, group_name)
- INDEX(game_id)

### 7. Enhance `game_battle_record` Table
**Add Fields:**
- `player_id` INT - Internal player/member ID
- `player_game_id` INT - Player's game ID
- `player_game_name` VARCHAR(64) - Player's game nickname
- `group_name` VARCHAR(64) - Group name
- `score` INT - Score/points
- `fee` INT - Service fee
- `factor` DECIMAL(10,4) - Settlement factor
- `player_balance` INT - Player balance after settlement

**Purpose:**
- Store detailed battle records per player
- Support settlement calculations
- Enable player-dimension queries

### 8. Create `game_sync_log` Table
**Purpose:** Track sync operations and errors

**Fields:**
- `id` INT PRIMARY KEY AUTO_INCREMENT
- `session_id` INT - FK to game_session
- `sync_type` VARCHAR(20) - battle_record, member_list, wallet_update
- `status` VARCHAR(20) - success, failed, partial
- `records_synced` INT - Number of records synced
- `error_message` TEXT - Error details if failed
- `started_at` TIMESTAMP
- `completed_at` TIMESTAMP

**Indexes:**
- INDEX(session_id, started_at)
- INDEX(status)

## Implementation Phases

### Phase 1: Schema Updates
1. Create migration scripts for all table modifications
2. Add new tables
3. Create indexes and constraints
4. Generate complete DDL documentation

### Phase 2: User Role System
1. Add role field to BasicUser model
2. Create role constants and enums
3. Implement role-based middleware/guards
4. Update user creation logic

### Phase 3: Registration & Game Account Binding
1. Create game account validation service
2. Implement registration flow with game account requirement
3. Add game server API client (port from plaza package)
4. Store game nickname in user profile

### Phase 4: Super Admin Game Account Management
1. Create API endpoints for game account management
2. Implement game account → store binding logic
3. Add validation: one game account → one store
4. Create admin UI/API for managing bindings

### Phase 5: Store Admin Exclusive Binding
1. Add validation logic for exclusive store admin binding
2. Implement checks during admin assignment
3. Add API to check current admin bindings
4. Handle edge cases (admin removal, transfer)

### Phase 6: Auto-sync Implementation
1. Port plaza.Session logic from waiter
2. Create session manager service
3. Implement auto-sync trigger on game account + store binding
4. Add background workers for continuous sync
5. Implement sync status monitoring

### Phase 7: Data Migration
1. Create migration scripts for existing data
2. Map WeChat-based data to new schema
3. Validate data integrity
4. Create rollback procedures

## API Endpoints Required

### User Management
- POST /api/v1/auth/register - Register with game account
- POST /api/v1/users/{id}/role - Update user role (super admin only)

### Game Account Management
- POST /api/v1/game-accounts - Add game account
- GET /api/v1/game-accounts - List user's game accounts
- POST /api/v1/game-accounts/{id}/verify - Verify game account
- DELETE /api/v1/game-accounts/{id} - Remove game account

### Store Binding Management
- POST /api/v1/game-accounts/{id}/bind-store - Bind game account to store
- DELETE /api/v1/game-accounts/{id}/unbind-store - Unbind from store
- GET /api/v1/stores/{houseGid}/bindings - List store bindings

### Store Admin Management
- POST /api/v1/stores/{houseGid}/admins - Assign store admin
- DELETE /api/v1/stores/{houseGid}/admins/{userId} - Remove store admin
- GET /api/v1/users/{id}/admin-stores - List stores user is admin of

### Session Management
- GET /api/v1/sessions - List active sessions
- POST /api/v1/sessions/{id}/start - Start session
- POST /api/v1/sessions/{id}/stop - Stop session
- GET /api/v1/sessions/{id}/sync-status - Get sync status

### Sync Management
- POST /api/v1/sync/trigger - Manually trigger sync
- GET /api/v1/sync/logs - Get sync logs
- GET /api/v1/sync/status/{sessionId} - Get sync status

## Testing Strategy

### Unit Tests
- User role validation
- Game account binding rules
- Store admin exclusive binding logic
- Session creation and management

### Integration Tests
- Registration flow with game account
- Game account verification with game server
- Auto-sync trigger on binding creation
- Data synchronization accuracy

### End-to-End Tests
- Complete user registration flow
- Super admin managing multiple game accounts
- Store admin assignment and restrictions
- Auto-sync functionality

## Rollback Plan

1. Keep legacy tables during migration period
2. Implement dual-write during transition
3. Create data validation scripts
4. Maintain rollback scripts for each phase
5. Monitor system health metrics

## Success Criteria

1. ✅ All WeChat dependencies removed
2. ✅ User role system fully functional
3. ✅ Game account binding working correctly
4. ✅ Store admin exclusive binding enforced
5. ✅ Auto-sync operational and reliable
6. ✅ All existing functionality preserved
7. ✅ Complete DDL documentation generated
8. ✅ Data migration successful with validation

