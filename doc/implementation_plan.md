# Project Migration and User Management System Implementation Plan

## Executive Summary

This document outlines the comprehensive implementation plan for migrating the legacy WeChat-dependent system to a modern game-account-based user management system. The migration involves database schema updates, backend API enhancements, and frontend interface modifications across three projects: `battle-tiles` (backend), `battle-reusables` (frontend), and `passing-dragonfly` (legacy).

### Key Objectives
1. Remove WeChat dependencies from user management
2. Enable super administrators to manage multiple game accounts with store bindings
3. Enforce game account binding during user registration
4. Implement automatic store session activation and game record synchronization
5. Maintain strict access control for store administrators

---

## Phase 1: Database Schema Analysis and Updates

### 1.1 Current Database Structure Analysis

**Existing Tables:**
- `basic_user`: Core user authentication and profile data
- `game_account`: Links users to game accounts (one-to-many)
- `game_ctrl_account`: System-level game account for automation
- `game_account_house`: Links game accounts to stores/houses
- `game_shop_admin`: Links users to stores with role assignments
- `game_session`: Manages active game sessions
- `game_battle_record`: Stores battle records (one row per game)
- `game_player_record`: Player-dimension records (to be created/verified)

**Current Limitations:**
1. `basic_user.wechat_id` field represents legacy WeChat dependency
2. No constraint preventing multiple store bindings per game account
3. No automatic session activation mechanism
4. Registration doesn't enforce game account binding
5. No automatic game record synchronization

### 1.2 Database Schema Modifications

#### 1.2.1 Remove WeChat Dependencies
**Table:** `basic_user`
- **Action:** Mark `wechat_id` field as deprecated (keep for data migration, remove later)
- **Rationale:** Maintain backward compatibility during transition period
- **Migration Strategy:** 
  - Phase 1: Make field nullable and optional
  - Phase 2: Migrate existing WeChat-based users to game accounts
  - Phase 3: Drop column after full migration

#### 1.2.2 Enforce Game Account to Store Binding Constraint
**Table:** `game_account_house`
- **Action:** Add unique constraint on `game_account_id`
- **Constraint:** `uk_game_account_house_account`
- **Effect:** Each game account can only be bound to ONE store
- **Impact:** Super admins can have multiple game accounts, each bound to different stores

#### 1.2.3 Enforce Store Admin Single-Store Constraint
**Table:** `game_shop_admin`
- **Action:** Add partial unique index on `user_id` where `deleted_at IS NULL`
- **Constraint:** `uk_game_shop_admin_user_active`
- **Effect:** Each user can only be an active admin for ONE store at a time
- **Rationale:** Prevents cross-store management conflicts

#### 1.2.4 Create/Verify Player-Dimension Record Table
**Table:** `game_player_record`
- **Purpose:** Store game records organized by player for efficient querying
- **Key Fields:**
  - `battle_record_id`: Links to `game_battle_record`
  - `house_gid`: Store identifier
  - `player_gid`: Game player ID
  - `game_account_id`: Links to `game_account` (nullable for unregistered players)
  - `score_delta`, `is_winner`, `battle_at`: Game statistics
  - `meta_json`: Flexible JSONB field for additional data
- **Indexes:**
  - `idx_gpr_house_battleat`: Query by store and time
  - `idx_gpr_player_battleat`: Query by player and time
  - `idx_gpr_account_battleat`: Query by account and time
  - `idx_gpr_group_room`: Query by group and room

#### 1.2.5 Add Foreign Key Constraints
- `game_account.user_id` → `basic_user.id`
- `game_player_record.battle_record_id` → `game_battle_record.id`
- `game_player_record.game_account_id` → `game_account.id`

### 1.3 DDL Documentation
All DDL scripts are documented in `/doc/user_management.ddl`

---

## Phase 2: Backend Implementation (battle-tiles)

### 2.1 User Registration Enhancement

#### 2.1.1 Update Registration Request Model
**File:** `internal/dal/req/basic_user.go`

**New Structure:**
```go
type RegisterRequest struct {
    Username        string `json:"username" binding:"required"`
    Password        string `json:"password" binding:"required"`
    GameLoginMode   string `json:"game_login_mode" binding:"required,oneof=account mobile"`
    GameAccount     string `json:"game_account" binding:"required"`
    GamePasswordMD5 string `json:"game_password_md5" binding:"required"`
    // NickName and Avatar will be fetched from game profile
}
```

#### 2.1.2 Update Registration Business Logic
**File:** `internal/biz/basic/basic_login.go`

**New Registration Flow:**
1. Validate username uniqueness
2. **NEW:** Validate game account credentials via game API
3. **NEW:** Fetch game profile (nickname, avatar, game_user_id)
4. Create `basic_user` record with game nickname
5. **NEW:** Create `game_ctrl_account` record
6. **NEW:** Create `game_account` record linked to user
7. Assign default role
8. Return JWT token with user info

**Key Changes:**
- Add game account validation step using existing `GameAccountUseCase.Verify()` method
- Fetch and use game nickname as user nickname
- Create game account binding during registration
- Handle validation errors gracefully

### 2.2 Super Administrator Game Account Management

#### 2.2.1 New API Endpoints
**File:** `internal/service/game/game_ctrl_account.go` (extend existing)

**Endpoints to Add:**
1. `POST /game/admin/ctrl-accounts/bind` - Bind new game account to store
2. `GET /game/admin/ctrl-accounts` - List all bound game accounts
3. `DELETE /game/admin/ctrl-accounts/:id` - Unbind game account
4. `PUT /game/admin/ctrl-accounts/:id/house` - Update store binding

#### 2.2.2 Bind Game Account to Store Logic
**Process:**
1. Verify user has super admin role (check JWT claims)
2. Validate game account credentials via game API
3. Check if game account already bound (enforce one-store-per-account)
4. Create `game_ctrl_account` record
5. Create `game_account_house` binding
6. **NEW:** Automatically start store session
7. **NEW:** Trigger game record synchronization
8. Return success response

#### 2.2.3 Automatic Session Activation
**File:** `internal/biz/game/ctrl_session.go` (extend existing)

**New Method:** `AutoStartSessionOnBind(ctx, userID, ctrlAccountID, houseGID)`
- Called automatically when game account is bound to store
- Creates `game_session` record with state="online"
- Initializes session monitoring
- Triggers background sync job

### 2.3 Game Record Synchronization Service

#### 2.3.1 Sync Service Architecture
**File:** `internal/service/game/game_record_sync.go` (NEW)

**Components:**
1. **Sync Trigger:** Activated on session start
2. **Data Fetcher:** Queries game API for battle records
3. **Data Transformer:** Converts battle records to player-dimension records
4. **Data Persister:** Saves to `game_player_record` table
5. **Progress Tracker:** Monitors sync status

#### 2.3.2 Sync Process Flow
1. Session starts → Trigger sync job
2. Fetch battle records from game API (paginated)
3. For each battle record:
   - Parse players JSON from `game_battle_record.players_json`
   - Create one `game_player_record` per player
   - Link to `game_account` if player is registered
4. Update sync status
5. Handle errors and retries

#### 2.3.3 Background Job Implementation
**Use existing `asynq` task queue**
**File:** `internal/biz/asynq.go` (extend existing)

**New Task:** `TaskSyncGameRecords`
- Payload: `{ctrl_account_id, house_gid, start_date, end_date}`
- Schedule: Immediate on session start, periodic updates
- Retry: 3 attempts with exponential backoff

### 2.4 Store Administrator Access Control

#### 2.4.1 Middleware Enhancement
**File:** `pkg/plugin/middleware/auth.go`

**New Middleware:** `RequireStoreAdmin()`
- Validates user has store admin role
- Checks user is only admin for ONE store
- Injects store context into request

#### 2.4.2 API Endpoint Protection
Apply middleware to store management endpoints:
- `/game/shop/*` routes
- Ensure cross-store operations are blocked

---

## Phase 3: Frontend Implementation (battle-reusables)

### 3.1 User Registration Flow Enhancement

#### 3.1.1 Registration Form Component
**File:** `app/(auth)/register.tsx` (create/modify)

**Form Fields:**
- Username (text input)
- Password (password input)
- Game Login Mode (radio: account/mobile)
- Game Account (text input)
- Game Password (password input, will be MD5 hashed)

**Validation:**
- All fields required
- Real-time game account validation
- Password strength check

#### 3.1.2 Game Account Verification
**API Call:** `POST /game/accounts/verify`
- Show loading indicator during validation
- Display success/error message
- Disable submit until validation passes

### 3.2 Super Admin Dashboard

#### 3.2.1 Game Account Management Interface
**File:** `app/(admin)/game-accounts/page.tsx` (NEW)

**Features:**
- List all bound game accounts with store info
- Add new game account button
- Unbind account action
- View session status
- Trigger manual sync

#### 3.2.2 Add Game Account Modal
**Component:** `AddGameAccountModal`
- Game login mode selector
- Account/mobile input
- Password input
- Store selector
- Validation and submission

### 3.3 Store Admin Interface

#### 3.3.1 Single Store View
**File:** `app/(store)/dashboard/page.tsx`

**Features:**
- Display current store info
- View game records
- Player statistics
- Cannot switch stores (enforced by backend)

---

## Phase 4: Testing and Validation

### 4.1 Database Testing
- [ ] Verify unique constraints work correctly
- [ ] Test foreign key cascade behaviors
- [ ] Validate indexes improve query performance
- [ ] Test data migration scripts

### 4.2 Backend API Testing
- [ ] Test registration with game account binding
- [ ] Test super admin multi-account management
- [ ] Test automatic session activation
- [ ] Test game record synchronization
- [ ] Test store admin access restrictions

### 4.3 Frontend Testing
- [ ] Test registration flow end-to-end
- [ ] Test game account validation
- [ ] Test super admin dashboard
- [ ] Test store admin restrictions

### 4.4 Integration Testing
- [ ] Test full user journey from registration to game record viewing
- [ ] Test multi-store management by super admin
- [ ] Test concurrent session handling

---

## Phase 5: Deployment and Migration

### 5.1 Database Migration
1. Backup production database
2. Apply DDL scripts in transaction
3. Verify constraints and indexes
4. Run data migration scripts
5. Validate data integrity

### 5.2 Backend Deployment
1. Deploy new backend version
2. Monitor error logs
3. Verify API endpoints
4. Test critical paths

### 5.3 Frontend Deployment
1. Build and deploy frontend
2. Test registration flow
3. Test admin dashboards
4. Monitor user feedback

### 5.4 Legacy System Migration
1. Identify WeChat-dependent users
2. Notify users to bind game accounts
3. Provide migration assistance
4. Deprecate WeChat login
5. Remove WeChat dependencies

---

## Appendix A: API Specifications

### Registration API
```
POST /login/register
Request:
{
  "username": "string",
  "password": "string",
  "game_login_mode": "account|mobile",
  "game_account": "string",
  "game_password_md5": "string"
}
Response:
{
  "token": "string",
  "user": {
    "id": int,
    "username": "string",
    "nick_name": "string",
    "avatar": "string"
  }
}
```

### Bind Game Account API
```
POST /game/admin/ctrl-accounts/bind
Request:
{
  "login_mode": "account|mobile",
  "identifier": "string",
  "pwd_md5": "string",
  "house_gid": int
}
Response:
{
  "id": int,
  "game_user_id": "string",
  "house_gid": int,
  "session_state": "online"
}
```

---

## Appendix B: File Structure

### Backend Files to Create/Modify
```
battle-tiles/
├── internal/
│   ├── dal/
│   │   ├── model/game/
│   │   │   └── game_player_record.go (verify/create)
│   │   └── req/
│   │       └── basic_user.go (modify RegisterRequest)
│   ├── biz/
│   │   ├── basic/
│   │   │   └── basic_login.go (modify Register method)
│   │   └── game/
│   │       ├── ctrl_session.go (add AutoStartSessionOnBind)
│   │       └── game_record_sync.go (NEW)
│   └── service/
│       └── game/
│           ├── game_ctrl_account.go (extend)
│           └── game_record_sync.go (NEW)
```

### Frontend Files to Create/Modify
```
battle-reusables/
├── app/
│   ├── (auth)/
│   │   └── register.tsx (modify)
│   ├── (admin)/
│   │   └── game-accounts/
│   │       └── page.tsx (NEW)
│   └── (store)/
│       └── dashboard/
│           └── page.tsx (modify)
└── components/
    └── game-account/
        ├── AddGameAccountModal.tsx (NEW)
        └── GameAccountList.tsx (NEW)
```

---

## Success Criteria

1. ✅ All new users must bind a game account during registration
2. ✅ Game nicknames automatically populate user profiles
3. ✅ Super admins can manage multiple game accounts
4. ✅ Each game account can only bind to one store
5. ✅ Store sessions auto-activate on account binding
6. ✅ Game records auto-sync when sessions start
7. ✅ Store admins restricted to single store management
8. ✅ WeChat dependencies removed from registration flow
9. ✅ All constraints enforced at database level
10. ✅ Comprehensive error handling and validation

