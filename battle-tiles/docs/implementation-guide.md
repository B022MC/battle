# Implementation Guide: WeChat Dependency Removal

## Overview
This guide provides step-by-step instructions for implementing the migration from WeChat-based authentication to a role-based account binding system.

## Completed Work

### 1. Schema Design ✅
- Created comprehensive DDL files in `docs/ddl/`:
  - `00_complete_schema.sql` - Complete database schema
  - `01_schema_modifications.sql` - Modifications to existing tables
  - `02_new_tables.sql` - New table definitions
- Created migration plan document: `docs/schema-migration-plan.md`

### 2. Model Layer Updates ✅
Updated Go models in `internal/dal/model/`:

#### Basic Module
- **basic_user.go**: Added `Role` and `GameNickname` fields
  - Role constants: `UserRoleSuperAdmin`, `UserRoleStoreAdmin`, `UserRoleRegularUser`
  - Helper methods: `IsSuperAdmin()`, `IsStoreAdmin()`, `IsRegularUser()`

#### Game Module
- **game_account.go**: Added verification fields
  - `GameUserID`, `VerifiedAt`, `VerificationStatus`
  - Verification status constants
  - Helper method: `IsVerified()`

- **game_shop_admin.go**: Added game account binding
  - `GameAccountID`, `IsExclusive` fields
  - Admin role constants
  - Helper methods: `IsAdmin()`, `IsOperator()`

- **game_session.go**: Added auto-sync fields
  - `AutoSyncEnabled`, `LastSyncAt`, `SyncStatus`, `GameAccountID`
  - Session state and sync status constants
  - Helper methods: `IsActive()`, `IsSyncing()`

- **game_battle_record.go**: Added player dimension fields
  - `PlayerID`, `PlayerGameID`, `PlayerGameName`, `GroupName`
  - `Score`, `Fee`, `Factor`, `PlayerBalance`

#### New Models Created
- **game_account_store_binding.go**: Game account to store bindings
- **game_member.go**: Player/member management
- **game_sync_log.go**: Synchronization operation logs
- **game_recharge_record.go**: Wallet transaction records

## Remaining Implementation Tasks

### 3. Business Logic Layer

#### 3.1 User Registration with Game Account Validation
**Location**: `internal/biz/user_biz.go` (or create new file)

**Required Functions**:
```go
// RegisterUserWithGameAccount registers a new user with game account validation
func (b *UserBiz) RegisterUserWithGameAccount(ctx context.Context, req *RegisterRequest) (*BasicUser, error)

// ValidateGameAccount validates game account credentials with game server
func (b *UserBiz) ValidateGameAccount(ctx context.Context, account, password string) (*GameAccountInfo, error)

// GetGameNickname retrieves nickname from game server
func (b *UserBiz) GetGameNickname(ctx context.Context, account, password string) (string, error)
```

**Implementation Steps**:
1. Create game server API client (port from `passing-dragonfly/waiter/plaza/mobilelogon.go`)
2. Implement `GetUserInfoByMobile` and `GetUserInfoByAccount` functions
3. Validate game account during registration
4. Retrieve and store game nickname
5. Create game account binding automatically

#### 3.2 Game Account Management (Super Admin)
**Location**: `internal/biz/game_account_biz.go`

**Required Functions**:
```go
// AddGameAccount adds a new game account for super admin
func (b *GameAccountBiz) AddGameAccount(ctx context.Context, userID int32, account, password string) (*GameAccount, error)

// BindGameAccountToStore binds a game account to a store
// Business Rule: One game account can only bind to ONE store
func (b *GameAccountBiz) BindGameAccountToStore(ctx context.Context, gameAccountID, houseGID, userID int32) error

// UnbindGameAccountFromStore unbinds a game account from a store
func (b *GameAccountBiz) UnbindGameAccountFromStore(ctx context.Context, gameAccountID int32) error

// ValidateStoreBinding validates that game account is not already bound to another store
func (b *GameAccountBiz) ValidateStoreBinding(ctx context.Context, gameAccountID int32) error
```

#### 3.3 Store Admin Management
**Location**: `internal/biz/shop_admin_biz.go`

**Required Functions**:
```go
// AssignStoreAdmin assigns a user as store administrator
// Business Rule: Store admin can only be admin of ONE store under ONE game account
func (b *ShopAdminBiz) AssignStoreAdmin(ctx context.Context, userID, houseGID, gameAccountID int32) error

// ValidateExclusiveBinding validates that user is not already admin of another store
func (b *ShopAdminBiz) ValidateExclusiveBinding(ctx context.Context, userID int32) error

// RemoveStoreAdmin removes a user from store administrator role
func (b *ShopAdminBiz) RemoveStoreAdmin(ctx context.Context, userID, houseGID int32) error

// GetUserAdminStores gets all stores where user is administrator
func (b *ShopAdminBiz) GetUserAdminStores(ctx context.Context, userID int32) ([]*GameShopAdmin, error)
```

#### 3.4 Auto-Sync Session Management
**Location**: `internal/biz/session_biz.go`

**Required Functions**:
```go
// CreateSessionForBinding creates a game session when game account is bound to store
// This triggers automatic synchronization
func (b *SessionBiz) CreateSessionForBinding(ctx context.Context, gameAccountID, houseGID, userID int32) (*GameSession, error)

// StartAutoSync starts automatic synchronization for a session
func (b *SessionBiz) StartAutoSync(ctx context.Context, sessionID int32) error

// StopAutoSync stops automatic synchronization for a session
func (b *SessionBiz) StopAutoSync(ctx context.Context, sessionID int32) error

// GetSyncStatus gets the current sync status for a session
func (b *SessionBiz) GetSyncStatus(ctx context.Context, sessionID int32) (*SyncStatus, error)
```

### 4. Plaza Session Integration

#### 4.1 Port Plaza Session Logic
**Location**: `internal/plaza/` (new package)

**Files to Create**:
- `session.go` - Main session management (port from `waiter/plaza/session.go`)
- `encoder.go` - Protocol encoder/decoder
- `commands.go` - Game server commands
- `parser.go` - Response parsers
- `const.go` - Protocol constants

**Key Components**:
```go
// Session represents a connection to game server
type Session struct {
    userName      string
    userPwd       string
    userID        int
    houseGID      int
    autoReconnect bool
    handler       IPlazaHandler
    // ... connection fields
}

// IPlazaHandler interface for handling game server events
type IPlazaHandler interface {
    OnRoomListUpdated(rooms []*RoomInfo)
    OnMemberListUpdated(members []*MemberInfo)
    OnBattleRecordReceived(record *BattleRecord)
    OnMemberInserted(member *MemberInfo)
    OnMemberDeleted(memberID int)
    OnSessionRestarted(newSession *Session)
}

// NewSession creates a new game server session
func NewSession(userName, userPwd string, userID, houseGID int, handler IPlazaHandler) (*Session, error)
```

#### 4.2 Implement Sync Workers
**Location**: `internal/worker/sync_worker.go`

**Required Functions**:
```go
// SyncWorker manages background synchronization tasks
type SyncWorker struct {
    sessions map[int32]*plaza.Session
    db       *gorm.DB
}

// StartSyncForSession starts sync worker for a session
func (w *SyncWorker) StartSyncForSession(sessionID int32, session *plaza.Session) error

// StopSyncForSession stops sync worker for a session
func (w *SyncWorker) StopSyncForSession(sessionID int32) error

// SyncBattleRecords syncs battle records from game server
func (w *SyncWorker) SyncBattleRecords(ctx context.Context, sessionID int32) error

// SyncMemberList syncs member list from game server
func (w *SyncWorker) SyncMemberList(ctx context.Context, sessionID int32) error
```

### 5. API Layer

#### 5.1 User Management APIs
**Location**: `internal/api/user_api.go`

**Endpoints**:
- `POST /api/v1/auth/register` - Register with game account
- `POST /api/v1/users/{id}/role` - Update user role (super admin only)
- `GET /api/v1/users/{id}/profile` - Get user profile

#### 5.2 Game Account Management APIs
**Location**: `internal/api/game_account_api.go`

**Endpoints**:
- `POST /api/v1/game-accounts` - Add game account
- `GET /api/v1/game-accounts` - List user's game accounts
- `POST /api/v1/game-accounts/{id}/verify` - Verify game account
- `DELETE /api/v1/game-accounts/{id}` - Remove game account
- `POST /api/v1/game-accounts/{id}/bind-store` - Bind to store
- `DELETE /api/v1/game-accounts/{id}/unbind-store` - Unbind from store

#### 5.3 Store Admin Management APIs
**Location**: `internal/api/shop_admin_api.go`

**Endpoints**:
- `POST /api/v1/stores/{houseGid}/admins` - Assign store admin
- `DELETE /api/v1/stores/{houseGid}/admins/{userId}` - Remove store admin
- `GET /api/v1/users/{id}/admin-stores` - List stores user is admin of
- `GET /api/v1/stores/{houseGid}/admins` - List store administrators

#### 5.4 Session Management APIs
**Location**: `internal/api/session_api.go`

**Endpoints**:
- `GET /api/v1/sessions` - List active sessions
- `POST /api/v1/sessions/{id}/start` - Start session
- `POST /api/v1/sessions/{id}/stop` - Stop session
- `GET /api/v1/sessions/{id}/sync-status` - Get sync status
- `POST /api/v1/sync/trigger` - Manually trigger sync

### 6. Database Migration

#### 6.1 Create Migration Scripts
**Location**: `migrations/`

**Files to Create**:
1. `001_add_user_roles.up.sql` - Add role fields to basic_user
2. `002_add_game_account_verification.up.sql` - Add verification fields
3. `003_create_game_account_store_binding.up.sql` - Create binding table
4. `004_add_shop_admin_fields.up.sql` - Add game_account_id to shop_admin
5. `005_add_session_sync_fields.up.sql` - Add sync fields to session
6. `006_create_game_member.up.sql` - Create game_member table
7. `007_create_sync_log.up.sql` - Create sync_log table
8. `008_add_battle_record_player_fields.up.sql` - Add player fields to battle_record

#### 6.2 Run Migrations
```bash
# Apply all migrations
go run cmd/migrate/main.go up

# Rollback if needed
go run cmd/migrate/main.go down
```

### 7. Testing

#### 7.1 Unit Tests
- Test user role validation
- Test game account binding rules
- Test store admin exclusive binding logic
- Test session creation and management

#### 7.2 Integration Tests
- Test registration flow with game account
- Test game account verification with game server
- Test auto-sync trigger on binding creation
- Test data synchronization accuracy

#### 7.3 End-to-End Tests
- Complete user registration flow
- Super admin managing multiple game accounts
- Store admin assignment and restrictions
- Auto-sync functionality

## Implementation Order

1. **Phase 1: Database & Models** ✅ (Completed)
   - Schema design
   - Model updates
   - DDL generation

2. **Phase 2: Game Server Integration** (Next)
   - Port plaza session logic
   - Implement game account validation
   - Create sync workers

3. **Phase 3: Business Logic**
   - User registration with game account
   - Game account management
   - Store admin management
   - Session management

4. **Phase 4: API Layer**
   - User management APIs
   - Game account APIs
   - Store admin APIs
   - Session APIs

5. **Phase 5: Testing & Deployment**
   - Unit tests
   - Integration tests
   - Data migration
   - Production deployment

## Next Steps

1. Create plaza package and port session logic from legacy project
2. Implement game account validation service
3. Create sync worker for automatic synchronization
4. Implement business logic layer
5. Create API endpoints
6. Write comprehensive tests
7. Perform data migration

## Notes

- All WeChat-related fields are kept for backward compatibility but are optional
- The system supports gradual migration - existing users can continue using WeChat login
- New users MUST use game account registration
- Super admins have full control over game account and store bindings
- Store admins are restricted to ONE store under ONE game account
- Auto-sync is triggered automatically when game account is bound to store

