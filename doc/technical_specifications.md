# Technical Specifications - User Management System Enhancement

## Table of Contents
1. [Database Layer Specifications](#database-layer-specifications)
2. [Backend API Specifications](#backend-api-specifications)
3. [Frontend Component Specifications](#frontend-component-specifications)
4. [Integration Specifications](#integration-specifications)
5. [Security Specifications](#security-specifications)

---

## Database Layer Specifications

### 1.1 Table Modifications

#### basic_user Table
**Current Schema:**
```sql
CREATE TABLE basic_user (
  id SERIAL PRIMARY KEY,
  username VARCHAR(50) NOT NULL UNIQUE,
  password VARCHAR(255),
  salt VARCHAR(50),
  wechat_id VARCHAR(64),  -- TO BE DEPRECATED
  avatar VARCHAR(255) DEFAULT '',
  nick_name VARCHAR(50) DEFAULT '',
  introduction TEXT,
  pinyin_code VARCHAR(100),
  first_letter VARCHAR(50),
  last_login_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ,
  is_del SMALLINT
);
```

**Migration Strategy:**
- Phase 1: Make `wechat_id` nullable (already is)
- Phase 2: Remove `wechat_id` from registration requirements
- Phase 3: Add comment marking field as deprecated
- Phase 4: Drop column after full migration (future)

#### game_account Table
**Current Schema:**
```sql
CREATE TABLE game_account (
  id SERIAL PRIMARY KEY,
  user_id INT4 NOT NULL,
  account VARCHAR(64) NOT NULL,
  pwd_md5 VARCHAR(64) NOT NULL,
  nickname VARCHAR(64) DEFAULT '',
  is_default BOOLEAN DEFAULT false,
  status INT4 DEFAULT 1,
  last_login_at TIMESTAMPTZ,
  login_mode VARCHAR(10) DEFAULT 'account',
  ctrl_account_id INT4,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ,
  is_del SMALLINT
);
```

**New Constraint:**
```sql
ALTER TABLE game_account
  ADD CONSTRAINT fk_game_account_user
  FOREIGN KEY (user_id) REFERENCES basic_user(id)
  ON UPDATE CASCADE ON DELETE RESTRICT;
```

#### game_account_house Table
**Current Schema:**
```sql
CREATE TABLE game_account_house (
  id SERIAL PRIMARY KEY,
  game_account_id INT4 NOT NULL,
  house_gid INT4 NOT NULL,
  is_default BOOLEAN DEFAULT false,
  status INT4 DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
```

**New Constraints:**
```sql
-- Enforce one store per game account
ALTER TABLE game_account_house
  ADD CONSTRAINT uk_game_account_house_account
  UNIQUE (game_account_id);

-- Also prevent duplicate (account, house) combinations
ALTER TABLE game_account_house
  ADD CONSTRAINT uk_game_account_house_account_house
  UNIQUE (game_account_id, house_gid);
```

#### game_shop_admin Table
**Current Schema:**
```sql
CREATE TABLE game_shop_admin (
  id SERIAL PRIMARY KEY,
  house_gid INT4 NOT NULL,
  user_id INT4 NOT NULL,
  role VARCHAR(20) NOT NULL,  -- 'admin' | 'operator'
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
```

**New Constraint:**
```sql
-- Enforce one active store per user
CREATE UNIQUE INDEX uk_game_shop_admin_user_active
  ON game_shop_admin(user_id)
  WHERE deleted_at IS NULL;
```

#### game_player_record Table (NEW)
**Purpose:** Store player-dimension game records for efficient querying

**Schema:**
```sql
CREATE TABLE game_player_record (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  
  -- Links
  battle_record_id INT4 NOT NULL,
  house_gid INT4 NOT NULL,
  player_gid INT8 NOT NULL,
  game_account_id INT4,  -- NULL for unregistered players
  
  -- Game info
  group_id INT4 NOT NULL,
  room_uid INT4 NOT NULL,
  kind_id INT4 NOT NULL,
  base_score INT4 NOT NULL,
  
  -- Player stats
  score_delta INT4 NOT NULL DEFAULT 0,
  is_winner BOOLEAN NOT NULL DEFAULT false,
  battle_at TIMESTAMPTZ NOT NULL,
  
  -- Flexible metadata
  meta_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  
  -- Foreign keys
  CONSTRAINT fk_gpr_battle
    FOREIGN KEY (battle_record_id)
    REFERENCES game_battle_record(id)
    ON UPDATE CASCADE ON DELETE CASCADE,
  
  CONSTRAINT fk_gpr_game_account
    FOREIGN KEY (game_account_id)
    REFERENCES game_account(id)
    ON UPDATE CASCADE ON DELETE SET NULL
);

-- Indexes for common queries
CREATE INDEX idx_gpr_house_battleat ON game_player_record(house_gid, battle_at DESC);
CREATE INDEX idx_gpr_player_battleat ON game_player_record(player_gid, battle_at DESC);
CREATE INDEX idx_gpr_account_battleat ON game_player_record(game_account_id, battle_at DESC);
CREATE INDEX idx_gpr_group_room ON game_player_record(group_id, room_uid);
```

### 1.2 Data Migration Scripts

#### Script 1: Verify Existing Data Integrity
```sql
-- Check for users without game accounts
SELECT u.id, u.username, u.nick_name
FROM basic_user u
LEFT JOIN game_account ga ON ga.user_id = u.id
WHERE ga.id IS NULL
  AND u.deleted_at IS NULL;

-- Check for game accounts with multiple store bindings
SELECT game_account_id, COUNT(*) as binding_count
FROM game_account_house
GROUP BY game_account_id
HAVING COUNT(*) > 1;

-- Check for users managing multiple stores
SELECT user_id, COUNT(*) as store_count
FROM game_shop_admin
WHERE deleted_at IS NULL
GROUP BY user_id
HAVING COUNT(*) > 1;
```

#### Script 2: Handle Constraint Violations Before Migration
```sql
-- Archive duplicate game_account_house bindings
-- Keep the most recent binding, mark others as inactive
WITH ranked_bindings AS (
  SELECT id, game_account_id,
         ROW_NUMBER() OVER (PARTITION BY game_account_id ORDER BY created_at DESC) as rn
  FROM game_account_house
)
UPDATE game_account_house
SET status = 0
WHERE id IN (
  SELECT id FROM ranked_bindings WHERE rn > 1
);

-- Archive duplicate game_shop_admin assignments
-- Keep the most recent assignment
WITH ranked_admins AS (
  SELECT id, user_id,
         ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at DESC) as rn
  FROM game_shop_admin
  WHERE deleted_at IS NULL
)
UPDATE game_shop_admin
SET deleted_at = NOW()
WHERE id IN (
  SELECT id FROM ranked_admins WHERE rn > 1
);
```

---

## Backend API Specifications

### 2.1 Registration API Enhancement

#### Endpoint: POST /login/register

**Current Request:**
```json
{
  "username": "string",
  "password": "string",
  "nick_name": "string",
  "avatar": "string",
  "wechat_id": "string"
}
```

**New Request:**
```json
{
  "username": "string",
  "password": "string",
  "game_login_mode": "account|mobile",
  "game_account": "string",
  "game_password_md5": "string"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "token": "string",
    "user": {
      "id": 1,
      "username": "string",
      "nick_name": "string",  // From game profile
      "avatar": "string",      // From game profile
      "roles": [1, 2],
      "perms": ["user:read", "user:write"]
    }
  }
}
```

**Implementation Steps:**

1. **Update Request Model** (`internal/dal/req/basic_user.go`):
```go
type RegisterRequest struct {
    Username        string `json:"username" binding:"required,min=3,max=50"`
    Password        string `json:"password" binding:"required,min=6,max=30"`
    GameLoginMode   string `json:"game_login_mode" binding:"required,oneof=account mobile"`
    GameAccount     string `json:"game_account" binding:"required"`
    GamePasswordMD5 string `json:"game_password_md5" binding:"required,len=32"`
}
```

2. **Update Registration Logic** (`internal/biz/basic/basic_login.go`):
```go
func (uc *BasicLoginUseCase) Register(ctx context.Context, c *gin.Context, req *req.RegisterRequest) (*resp.LoginResponse, error) {
    // 1. Check username uniqueness
    if existUser, _ := uc.repo.FindByUsername(ctx, req.Username); existUser != nil {
        return nil, errors.New("username already exists")
    }
    
    // 2. Validate game account credentials
    var loginMode consts.GameLoginMode
    if req.GameLoginMode == "account" {
        loginMode = consts.GameLoginModeAccount
    } else {
        loginMode = consts.GameLoginModeMobile
    }
    
    gameProfile, err := uc.gameAccountUC.VerifyAndFetchProfile(ctx, loginMode, req.GameAccount, req.GamePasswordMD5)
    if err != nil {
        return nil, errors.Wrap(err, "invalid game account credentials")
    }
    
    // 3. Hash platform password
    passwordPlain, err := uc.BeforeValidatorPwd(ctx, req.Password)
    if err != nil {
        return nil, err
    }
    salt, err := utils.GenerateSalt()
    if err != nil {
        return nil, errors.New("failed to generate salt")
    }
    passwordHash := utils.BcryptHash(passwordPlain, salt)
    
    // 4. Create basic_user with game nickname
    user := &basicModel.BasicUser{
        Username: req.Username,
        Password: passwordHash,
        Salt:     salt,
        NickName: gameProfile.Nickname,
        Avatar:   gameProfile.Avatar,
    }
    
    userID, err := uc.repo.Create(ctx, user)
    if err != nil {
        return nil, errors.Wrap(err, "failed to create user")
    }
    user.Id = userID
    
    // 5. Create game_ctrl_account
    ctrlAccount := &gameModel.GameCtrlAccount{
        LoginMode:  int32(loginMode),
        Identifier: req.GameAccount,
        PwdMD5:     req.GamePasswordMD5,
        GameUserID: gameProfile.GameUserID,
        GameID:     gameProfile.GameID,
        Status:     1,
    }
    ctrlAccountID, err := uc.gameCtrlAccountRepo.Create(ctx, ctrlAccount)
    if err != nil {
        // Rollback user creation
        uc.repo.Delete(ctx, userID)
        return nil, errors.Wrap(err, "failed to create game ctrl account")
    }
    
    // 6. Create game_account
    gameAccount := &gameModel.GameAccount{
        UserID:        userID,
        Account:       req.GameAccount,
        PwdMD5:        req.GamePasswordMD5,
        Nickname:      gameProfile.Nickname,
        IsDefault:     true,
        Status:        1,
        LoginMode:     req.GameLoginMode,
        CtrlAccountID: &ctrlAccountID,
    }
    _, err = uc.gameAccountRepo.Create(ctx, gameAccount)
    if err != nil {
        // Rollback previous creations
        uc.gameCtrlAccountRepo.Delete(ctx, ctrlAccountID)
        uc.repo.Delete(ctx, userID)
        return nil, errors.Wrap(err, "failed to create game account")
    }
    
    // 7. Assign default role
    if uc.authRepo != nil {
        if err := uc.authRepo.EnsureUserHasOnlyRoleByCode(ctx, userID, "ordinary-user"); err != nil {
            uc.log.Errorf("failed to assign role: %v", err)
        }
    }
    
    // 8. Build and return login response
    return uc.buildLoginResponse(c, user)
}
```

3. **Add Game Profile Fetching Method** (`internal/biz/game/game_account.go`):
```go
type GameProfile struct {
    GameUserID string
    GameID     string
    Nickname   string
    Avatar     string
}

func (uc *GameAccountUseCase) VerifyAndFetchProfile(ctx context.Context, mode consts.GameLoginMode, account, pwdMD5 string) (*GameProfile, error) {
    // Use existing plaza client to verify and fetch profile
    client := uc.plazaClient
    
    loginResp, err := client.Login(ctx, mode, account, pwdMD5)
    if err != nil {
        return nil, errors.Wrap(err, "game login failed")
    }
    
    profile := &GameProfile{
        GameUserID: loginResp.UserID,
        GameID:     loginResp.GameID,
        Nickname:   loginResp.Nickname,
        Avatar:     loginResp.Avatar,
    }
    
    return profile, nil
}
```

### 2.2 Super Admin Game Account Management APIs

#### Endpoint: POST /game/admin/ctrl-accounts/bind

**Purpose:** Bind a new game account to a store for super admin

**Request:**
```json
{
  "login_mode": 1,  // 1=account, 2=mobile
  "identifier": "string",
  "pwd_md5": "string",
  "house_gid": 123
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "game_user_id": "string",
    "house_gid": 123,
    "session_state": "online",
    "sync_status": "in_progress"
  }
}
```

**Implementation** (`internal/service/game/game_ctrl_account.go`):
```go
func (s *CtrlAccountService) BindToHouse(c *gin.Context) {
    var in req.BindCtrlAccountRequest
    if err := c.ShouldBindJSON(&in); err != nil {
        response.Fail(c, ecode.ParamsFailed, err)
        return
    }
    
    claims, err := utils.GetClaims(c)
    if err != nil {
        response.Fail(c, ecode.TokenValidateFailed, err)
        return
    }
    
    // Verify super admin role
    if !utils.HasRole(claims, "super-admin") {
        response.Fail(c, ecode.PermissionDenied, "super admin role required")
        return
    }
    
    // Bind account and auto-start session
    result, err := s.uc.BindAccountToHouseWithAutoStart(
        c.Request.Context(),
        claims.UserID,
        in.LoginMode,
        in.Identifier,
        in.PwdMD5,
        in.HouseGID,
    )
    if err != nil {
        response.Fail(c, ecode.Failed, err)
        return
    }
    
    response.Success(c, result)
}
```

#### Endpoint: GET /game/admin/ctrl-accounts

**Purpose:** List all game accounts bound by super admin

**Response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "identifier": "string",
        "game_user_id": "string",
        "house_gid": 123,
        "house_name": "Store A",
        "session_state": "online",
        "last_sync_at": "2025-11-11T10:00:00Z"
      }
    ],
    "total": 1
  }
}
```

### 2.3 Game Record Synchronization Service

#### Background Job: TaskSyncGameRecords

**File:** `internal/biz/asynq.go`

```go
const (
    TypeSyncGameRecords = "sync:game_records"
)

type SyncGameRecordsPayload struct {
    CtrlAccountID int32     `json:"ctrl_account_id"`
    HouseGID      int32     `json:"house_gid"`
    StartDate     time.Time `json:"start_date"`
    EndDate       time.Time `json:"end_date"`
}

func (b *BizAsynq) HandleSyncGameRecords(ctx context.Context, t *asynq.Task) error {
    var p SyncGameRecordsPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return fmt.Errorf("json.Unmarshal failed: %v", err)
    }
    
    return b.gameRecordSyncUC.SyncRecords(ctx, p.CtrlAccountID, p.HouseGID, p.StartDate, p.EndDate)
}
```

**Sync Logic** (`internal/biz/game/game_record_sync.go`):
```go
func (uc *GameRecordSyncUseCase) SyncRecords(ctx context.Context, ctrlAccountID, houseGID int32, startDate, endDate time.Time) error {
    // 1. Fetch battle records from game API
    records, err := uc.plazaClient.FetchBattleRecords(ctx, houseGID, startDate, endDate)
    if err != nil {
        return errors.Wrap(err, "failed to fetch battle records")
    }
    
    // 2. Process each record
    for _, record := range records {
        // Save to game_battle_record if not exists
        battleRecord, err := uc.saveBattleRecord(ctx, record)
        if err != nil {
            uc.log.Errorf("failed to save battle record: %v", err)
            continue
        }
        
        // Parse players and create player records
        players, err := uc.parsePlayers(record.PlayersJSON)
        if err != nil {
            uc.log.Errorf("failed to parse players: %v", err)
            continue
        }
        
        for _, player := range players {
            playerRecord := &gameModel.GamePlayerRecord{
                BattleRecordID: battleRecord.Id,
                HouseGID:       houseGID,
                PlayerGID:      player.PlayerGID,
                GameAccountID:  uc.findGameAccountID(ctx, player.PlayerGID),
                GroupID:        record.GroupID,
                RoomUID:        record.RoomUID,
                KindID:         record.KindID,
                BaseScore:      record.BaseScore,
                ScoreDelta:     player.ScoreDelta,
                IsWinner:       player.IsWinner,
                BattleAt:       record.BattleAt,
                MetaJSON:       player.MetaJSON,
            }
            
            if err := uc.repo.CreatePlayerRecord(ctx, playerRecord); err != nil {
                uc.log.Errorf("failed to create player record: %v", err)
            }
        }
    }
    
    return nil
}
```

---

## Frontend Component Specifications

### 3.1 Registration Form Component

**File:** `app/(auth)/register.tsx`

```typescript
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Label } from '@/components/ui/label';
import { useToast } from '@/components/ui/use-toast';
import { registerUser, verifyGameAccount } from '@/services/auth';
import md5 from 'md5';

export default function RegisterPage() {
  const router = useRouter();
  const { toast } = useToast();
  
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    confirmPassword: '',
    gameLoginMode: 'account',
    gameAccount: '',
    gamePassword: '',
  });
  
  const [gameAccountVerified, setGameAccountVerified] = useState(false);
  const [isVerifying, setIsVerifying] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  const handleVerifyGameAccount = async () => {
    if (!formData.gameAccount || !formData.gamePassword) {
      toast({
        title: 'Error',
        description: 'Please enter game account and password',
        variant: 'destructive',
      });
      return;
    }
    
    setIsVerifying(true);
    try {
      const result = await verifyGameAccount({
        mode: formData.gameLoginMode,
        account: formData.gameAccount,
        pwd_md5: md5(formData.gamePassword),
      });
      
      if (result.ok) {
        setGameAccountVerified(true);
        toast({
          title: 'Success',
          description: 'Game account verified successfully',
        });
      }
    } catch (error) {
      toast({
        title: 'Verification Failed',
        description: error.message || 'Invalid game account credentials',
        variant: 'destructive',
      });
    } finally {
      setIsVerifying(false);
    }
  };
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!gameAccountVerified) {
      toast({
        title: 'Error',
        description: 'Please verify your game account first',
        variant: 'destructive',
      });
      return;
    }
    
    if (formData.password !== formData.confirmPassword) {
      toast({
        title: 'Error',
        description: 'Passwords do not match',
        variant: 'destructive',
      });
      return;
    }
    
    setIsSubmitting(true);
    try {
      const result = await registerUser({
        username: formData.username,
        password: formData.password,
        game_login_mode: formData.gameLoginMode,
        game_account: formData.gameAccount,
        game_password_md5: md5(formData.gamePassword),
      });
      
      // Save token and redirect
      localStorage.setItem('token', result.token);
      toast({
        title: 'Success',
        description: 'Registration successful',
      });
      router.push('/dashboard');
    } catch (error) {
      toast({
        title: 'Registration Failed',
        description: error.message || 'Failed to register',
        variant: 'destructive',
      });
    } finally {
      setIsSubmitting(false);
    }
  };
  
  return (
    <div className="container max-w-md mx-auto py-10">
      <h1 className="text-2xl font-bold mb-6">Register</h1>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <Label htmlFor="username">Username</Label>
          <Input
            id="username"
            value={formData.username}
            onChange={(e) => setFormData({ ...formData, username: e.target.value })}
            required
          />
        </div>
        
        <div>
          <Label htmlFor="password">Password</Label>
          <Input
            id="password"
            type="password"
            value={formData.password}
            onChange={(e) => setFormData({ ...formData, password: e.target.value })}
            required
          />
        </div>
        
        <div>
          <Label htmlFor="confirmPassword">Confirm Password</Label>
          <Input
            id="confirmPassword"
            type="password"
            value={formData.confirmPassword}
            onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
            required
          />
        </div>
        
        <div>
          <Label>Game Login Mode</Label>
          <RadioGroup
            value={formData.gameLoginMode}
            onValueChange={(value) => setFormData({ ...formData, gameLoginMode: value })}
          >
            <div className="flex items-center space-x-2">
              <RadioGroupItem value="account" id="account" />
              <Label htmlFor="account">Account</Label>
            </div>
            <div className="flex items-center space-x-2">
              <RadioGroupItem value="mobile" id="mobile" />
              <Label htmlFor="mobile">Mobile</Label>
            </div>
          </RadioGroup>
        </div>
        
        <div>
          <Label htmlFor="gameAccount">Game Account</Label>
          <div className="flex gap-2">
            <Input
              id="gameAccount"
              value={formData.gameAccount}
              onChange={(e) => {
                setFormData({ ...formData, gameAccount: e.target.value });
                setGameAccountVerified(false);
              }}
              required
            />
            <Button
              type="button"
              onClick={handleVerifyGameAccount}
              disabled={isVerifying || gameAccountVerified}
            >
              {isVerifying ? 'Verifying...' : gameAccountVerified ? 'Verified' : 'Verify'}
            </Button>
          </div>
        </div>
        
        <div>
          <Label htmlFor="gamePassword">Game Password</Label>
          <Input
            id="gamePassword"
            type="password"
            value={formData.gamePassword}
            onChange={(e) => {
              setFormData({ ...formData, gamePassword: e.target.value });
              setGameAccountVerified(false);
            }}
            required
          />
        </div>
        
        <Button type="submit" className="w-full" disabled={isSubmitting || !gameAccountVerified}>
          {isSubmitting ? 'Registering...' : 'Register'}
        </Button>
      </form>
    </div>
  );
}
```

---

## Security Specifications

### 5.1 Authentication & Authorization

**JWT Token Structure:**
```json
{
  "user_id": 1,
  "username": "string",
  "nick_name": "string",
  "platform": "string",
  "roles": [1, 2],
  "perms": ["user:read", "game:ctrl:create"],
  "exp": 1699999999
}
```

**Role Hierarchy:**
- `super-admin`: Full system access, can manage multiple game accounts
- `store-admin`: Single store management
- `store-operator`: Store operations only
- `ordinary-user`: Basic user access

**Permission Checks:**
- Registration: Public (no auth required)
- Game account binding: Super admin only
- Store management: Store admin + ownership check
- Game record viewing: Store admin/operator + store ownership check

### 5.2 Data Validation

**Registration Validation:**
- Username: 3-50 characters, alphanumeric + underscore
- Password: 6-30 characters, must contain letter and number
- Game account: Required, must pass game API validation
- Game password: Required, MD5 hashed before transmission

**API Input Validation:**
- All inputs sanitized to prevent SQL injection
- CSRF tokens for state-changing operations
- Rate limiting on authentication endpoints
- Request size limits enforced

### 5.3 Error Handling

**Error Response Format:**
```json
{
  "code": 400,
  "msg": "error message",
  "data": null
}
```

**Common Error Codes:**
- 400: Bad request / validation failed
- 401: Unauthorized / token invalid
- 403: Forbidden / insufficient permissions
- 409: Conflict / constraint violation
- 500: Internal server error

---

## Appendix: Code Examples

### Example: Transaction Handling in Registration

```go
func (uc *BasicLoginUseCase) Register(ctx context.Context, c *gin.Context, req *req.RegisterRequest) (*resp.LoginResponse, error) {
    // Start transaction
    tx := uc.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // Create user
    user, err := uc.createUser(tx, req)
    if err != nil {
        tx.Rollback()
        return nil, err
    }
    
    // Create game accounts
    if err := uc.createGameAccounts(tx, user.Id, req); err != nil {
        tx.Rollback()
        return nil, err
    }
    
    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return nil, errors.Wrap(err, "failed to commit transaction")
    }
    
    return uc.buildLoginResponse(c, user)
}
```

### Example: Concurrent Sync Handling

```go
func (uc *GameRecordSyncUseCase) SyncRecords(ctx context.Context, ctrlAccountID, houseGID int32, startDate, endDate time.Time) error {
    // Use semaphore to limit concurrent syncs
    if err := uc.syncSemaphore.Acquire(ctx, 1); err != nil {
        return errors.Wrap(err, "failed to acquire sync semaphore")
    }
    defer uc.syncSemaphore.Release(1)
    
    // Check if sync already in progress
    key := fmt.Sprintf("sync:%d:%d", ctrlAccountID, houseGID)
    if exists, _ := uc.redis.Exists(ctx, key).Result(); exists > 0 {
        return errors.New("sync already in progress")
    }
    
    // Set sync lock with expiration
    uc.redis.SetEX(ctx, key, "1", 30*time.Minute)
    defer uc.redis.Del(ctx, key)
    
    // Perform sync...
    return uc.doSync(ctx, ctrlAccountID, houseGID, startDate, endDate)
}
```

