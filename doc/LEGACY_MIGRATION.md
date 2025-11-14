# Legacy System Migration Guide - passing-dragonfly/waiter

## Overview

This document provides detailed analysis and migration strategy for moving functionality from the legacy `passing-dragonfly/waiter` system to the new `battle-tiles` backend and `battle-reusables` frontend.

---

## Legacy System Analysis

### Source Location
**Path:** `passing-dragonfly/waiter`

### Technology Stack
- **Language:** Go
- **Database:** MySQL (legacy) → PostgreSQL (new)
- **Framework:** Custom framework (has/core)
- **WeChat Integration:** Heavy dependency on WeChat for authentication and notifications

### Key Components

#### 1. Data Models (`model/stoage.go`)

**Legacy Tables:**
```
TPlayer          - Player/user information (WeChat-based)
TGamer           - Game account bindings
THouse           - Store/house management
TManager         - Store manager information
TBattleRecord    - Game battle records
TRechargeRecord  - Recharge/balance records
TFeeRecord       - Fee collection records
TFeeSettle       - Fee settlement records
TPushStat        - Push notification statistics
TNotice          - Notification records
```

#### 2. Session Management (`plaza/session.go`)

**Legacy Session Features:**
- TCP connection to game servers (ports 82 and 87)
- Auto-reconnect functionality
- Command queue management
- Real-time game event handling
- Member list synchronization
- Room/table monitoring

#### 3. Service Layer (`robotsvs/service.go`)

**Legacy Service Features:**
- House (store) lifecycle management
- Robot/bot management for automation
- WeChat notification integration
- Background job scheduling
- Command caching

---

## Migration Mapping

### 1. Data Model Migration

#### TPlayer → basic_user + game_account

**Legacy Structure:**
```go
type TPlayer struct {
    ID           int
    WxKey        string  // WeChat key - TO BE REMOVED
    WxName       string  // WeChat name - TO BE REMOVED
    Balance      int
    Credit       int
    HouseGid     int
    PlayGroup    string
    Recommender  string
    UseMultiGids bool
    ActiveGid    int
    RecordedAt   int64
    UpdatedAt    int64
}
```

**New Structure:**
```sql
-- basic_user table
CREATE TABLE basic_user (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255),
    salt VARCHAR(50),
    nick_name VARCHAR(50),  -- From game profile
    avatar VARCHAR(255),    -- From game profile
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- game_account table
CREATE TABLE game_account (
    id SERIAL PRIMARY KEY,
    user_id INT4 NOT NULL REFERENCES basic_user(id),
    account VARCHAR(64) NOT NULL,
    pwd_md5 VARCHAR(64) NOT NULL,
    nickname VARCHAR(64),
    is_default BOOLEAN DEFAULT false,
    status INT4 DEFAULT 1,
    ctrl_account_id INT4,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
```

**Migration Strategy:**
1. Extract non-WeChat user data from TPlayer
2. Create basic_user records with game account as username
3. Create game_account records linked to basic_user
4. Migrate balance and credit to separate wallet tables
5. Map PlayGroup to store/house associations

---

#### TGamer → game_account + game_account_house

**Legacy Structure:**
```go
type TGamer struct {
    ID         int
    HouseGid   int
    PlayerId   int
    GameId     int
    GameName   string
    Forbid     bool
    RecordedAt int64
}
```

**New Structure:**
```sql
-- game_account_house table
CREATE TABLE game_account_house (
    id SERIAL PRIMARY KEY,
    game_account_id INT4 NOT NULL UNIQUE,  -- One store per account
    house_gid INT4 NOT NULL,
    is_default BOOLEAN DEFAULT false,
    status INT4 DEFAULT 1,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
```

**Migration Strategy:**
1. Map TGamer.GameId to game_account
2. Map TGamer.HouseGid to game_account_house.house_gid
3. Handle Forbid flag → status field
4. Enforce one-to-one constraint (game_account_id UNIQUE)

---

#### THouse → game_ctrl_account + game_session

**Legacy Structure:**
```go
type THouse struct {
    ID           int
    GameId       int     // Store/house game ID
    Groups       string  // Store groups
    UserLogon    string  // Admin game account (mobile)
    UserGid      int     // Admin game ID
    UserPwd      string  // Admin password
    Magic        string  // Magic code
    Price        int     // Push price
    UseMultiGids bool
    GameFees     string  // JSON
    GameCredits  string  // JSON
    Balance      int
    ManualFreeze bool
    Stopped      bool
    PushCredit   int
    Mode         int
    OwnerWxKey   string  // WeChat key - TO BE REMOVED
    ShareFee     bool
}
```

**New Structure:**
```sql
-- game_ctrl_account table
CREATE TABLE game_ctrl_account (
    id SERIAL PRIMARY KEY,
    login_mode INT4 NOT NULL,  -- 1=account, 2=mobile
    identifier VARCHAR(64) NOT NULL,
    pwd_md5 VARCHAR(64) NOT NULL,
    game_user_id VARCHAR(32),
    game_id VARCHAR(32),
    status INT4 DEFAULT 1,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- game_session table
CREATE TABLE game_session (
    id SERIAL PRIMARY KEY,
    game_ctrl_account_id INT4 NOT NULL,
    user_id INT4 NOT NULL,
    house_gid INT4 NOT NULL,
    state VARCHAR(20) NOT NULL,  -- 'online', 'offline', 'error'
    device_ip VARCHAR(64),
    error_msg VARCHAR(255),
    end_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
```

**Migration Strategy:**
1. Create game_ctrl_account from THouse.UserLogon/UserPwd
2. Map THouse.GameId to house_gid
3. Create game_session records for active houses
4. Migrate GameFees/GameCredits to separate configuration tables
5. Remove WeChat-specific fields (OwnerWxKey)

---

#### TBattleRecord → game_battle_record + game_player_record

**Legacy Structure:**
```go
type TBattleRecord struct {
    ID            int
    PlayerId      int
    CreateTime    int
    HouseGid      int
    RoomId        int
    RoomUid       int64
    PlayerGid     int
    PlayerGname   string
    PlayGroup     string
    GameKind      int
    Score         int
    Base          int
    Fee           int
    Factor        float64
    PlayerBalance int
}
```

**New Structure:**
```sql
-- game_battle_record table (one row per game)
CREATE TABLE game_battle_record (
    id SERIAL PRIMARY KEY,
    house_gid INT4 NOT NULL,
    group_id INT4 NOT NULL,
    room_uid INT4 NOT NULL,
    kind_id INT4 NOT NULL,
    base_score INT4 NOT NULL,
    battle_at TIMESTAMPTZ NOT NULL,
    players_json TEXT NOT NULL,  -- JSON array of players
    created_at TIMESTAMPTZ
);

-- game_player_record table (one row per player per game)
CREATE TABLE game_player_record (
    id SERIAL PRIMARY KEY,
    battle_record_id INT4 NOT NULL REFERENCES game_battle_record(id),
    house_gid INT4 NOT NULL,
    player_gid INT8 NOT NULL,
    game_account_id INT4 REFERENCES game_account(id),
    group_id INT4 NOT NULL,
    room_uid INT4 NOT NULL,
    kind_id INT4 NOT NULL,
    base_score INT4 NOT NULL,
    score_delta INT4 NOT NULL,
    is_winner BOOLEAN NOT NULL,
    battle_at TIMESTAMPTZ NOT NULL,
    meta_json JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
```

**Migration Strategy:**
1. Group TBattleRecord by RoomUid to create game_battle_record
2. Create one game_player_record per TBattleRecord row
3. Link game_account_id if player is registered
4. Store fee and factor in meta_json
5. Calculate is_winner from Score
6. Create indexes for efficient querying

---

#### TManager → game_shop_admin

**Legacy Structure:**
```go
type TManager struct {
    ID          int
    DeviceIP    string
    DeviceToken string  // WeChat login ID - TO BE REMOVED
    HouseGid    int
    PlayGroup   string
    WxKey       string  // WeChat key - TO BE REMOVED
    WxName      string  // WeChat name - TO BE REMOVED
}
```

**New Structure:**
```sql
-- game_shop_admin table
CREATE TABLE game_shop_admin (
    id SERIAL PRIMARY KEY,
    house_gid INT4 NOT NULL,
    user_id INT4 NOT NULL,
    role VARCHAR(20) NOT NULL,  -- 'admin' | 'operator'
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- Unique constraint: one active store per user
CREATE UNIQUE INDEX uk_game_shop_admin_user_active
    ON game_shop_admin(user_id)
    WHERE deleted_at IS NULL;
```

**Migration Strategy:**
1. Map TManager to game_shop_admin
2. Remove WeChat-specific fields
3. Link to basic_user via user_id
4. Set role based on permissions
5. Enforce one-store-per-user constraint

---

### 2. Session Management Migration

#### Legacy Session (`plaza/session.go`)

**Key Features:**
- TCP connections to game servers (ports 82, 87)
- Auto-reconnect with exponential backoff
- Command queue for game operations
- Real-time event handling
- Member list synchronization
- Room/table monitoring

**Migration to battle-tiles:**

**New Implementation:**
```go
// internal/biz/game/ctrl_session.go

type CtrlSessionUseCase struct {
    repo          CtrlSessionRepo
    plazaClient   plaza.IPlazaClient
    syncUseCase   *GameRecordSyncUseCase
    sessionMonitor *SessionMonitor
}

func (uc *CtrlSessionUseCase) StartSession(ctx context.Context, userID, ctrlAccountID, houseGID int32) error {
    // 1. Verify ctrl account exists
    ctrlAccount, err := uc.repo.GetCtrlAccount(ctx, ctrlAccountID)
    if err != nil {
        return err
    }
    
    // 2. Create session record
    session := &model.GameSession{
        GameCtrlAccountID: ctrlAccountID,
        UserID:            userID,
        HouseGID:          houseGID,
        State:             "connecting",
        DeviceIP:          uc.getLocalIP(),
    }
    
    sessionID, err := uc.repo.CreateSession(ctx, session)
    if err != nil {
        return err
    }
    
    // 3. Connect to plaza (game server)
    plazaSession, err := uc.plazaClient.Connect(
        ctrlAccount.Identifier,
        ctrlAccount.PwdMD5,
        int(houseGID),
    )
    if err != nil {
        uc.repo.UpdateSessionState(ctx, sessionID, "error", err.Error())
        return err
    }
    
    // 4. Update session state
    uc.repo.UpdateSessionState(ctx, sessionID, "online", "")
    
    // 5. Register session monitor
    uc.sessionMonitor.Register(sessionID, plazaSession)
    
    // 6. Trigger game record sync
    go uc.syncUseCase.SyncRecords(ctx, ctrlAccountID, houseGID, time.Now().AddDate(0, -1, 0), time.Now())
    
    return nil
}

func (uc *CtrlSessionUseCase) AutoStartSessionOnBind(ctx context.Context, userID, ctrlAccountID, houseGID int32) error {
    return uc.StartSession(ctx, userID, ctrlAccountID, houseGID)
}
```

**Session Monitor:**
```go
// internal/server/session_monitor.go

type SessionMonitor struct {
    sessions sync.Map  // sessionID -> *PlazaSession
    repo     CtrlSessionRepo
}

func (m *SessionMonitor) Register(sessionID int32, plazaSession plaza.IPlazaSession) {
    m.sessions.Store(sessionID, plazaSession)
    
    // Start monitoring goroutine
    go m.monitor(sessionID, plazaSession)
}

func (m *SessionMonitor) monitor(sessionID int32, plazaSession plaza.IPlazaSession) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if !plazaSession.IsConnected() {
                // Attempt reconnect
                if err := plazaSession.Reconnect(); err != nil {
                    m.repo.UpdateSessionState(context.Background(), sessionID, "error", err.Error())
                } else {
                    m.repo.UpdateSessionState(context.Background(), sessionID, "online", "")
                }
            }
        case <-plazaSession.Done():
            // Session ended
            m.sessions.Delete(sessionID)
            m.repo.UpdateSessionState(context.Background(), sessionID, "offline", "")
            return
        }
    }
}
```

---

### 3. Service Layer Migration

#### Legacy Service (`robotsvs/service.go`)

**Key Features:**
- House (store) lifecycle management
- Robot management for automation
- WeChat notification integration
- Background job scheduling
- Command caching

**Migration to battle-tiles:**

**New Implementation:**
```go
// internal/biz/game/game_record_sync.go

type GameRecordSyncUseCase struct {
    repo          GameRecordRepo
    plazaClient   plaza.IPlazaClient
    asynqClient   *asynq.Client
    log           *log.Helper
}

func (uc *GameRecordSyncUseCase) SyncRecords(ctx context.Context, ctrlAccountID, houseGID int32, startDate, endDate time.Time) error {
    // 1. Fetch battle records from game API
    records, err := uc.plazaClient.FetchBattleRecords(ctx, houseGID, startDate, endDate)
    if err != nil {
        return errors.Wrap(err, "failed to fetch battle records")
    }
    
    // 2. Process each record
    for _, record := range records {
        // Check if already synced
        exists, _ := uc.repo.BattleRecordExists(ctx, houseGID, record.RoomUID)
        if exists {
            continue
        }
        
        // Save battle record
        battleRecord := &model.GameBattleRecord{
            HouseGID:    houseGID,
            GroupID:     record.GroupID,
            RoomUID:     record.RoomUID,
            KindID:      record.KindID,
            BaseScore:   record.BaseScore,
            BattleAt:    record.BattleAt,
            PlayersJSON: record.PlayersJSON,
        }
        
        battleRecordID, err := uc.repo.CreateBattleRecord(ctx, battleRecord)
        if err != nil {
            uc.log.Errorf("failed to create battle record: %v", err)
            continue
        }
        
        // Parse players and create player records
        players, err := uc.parsePlayers(record.PlayersJSON)
        if err != nil {
            uc.log.Errorf("failed to parse players: %v", err)
            continue
        }
        
        for _, player := range players {
            // Find game account if registered
            gameAccountID := uc.findGameAccountID(ctx, player.PlayerGID)
            
            playerRecord := &model.GamePlayerRecord{
                BattleRecordID: battleRecordID,
                HouseGID:       houseGID,
                PlayerGID:      player.PlayerGID,
                GameAccountID:  gameAccountID,
                GroupID:        record.GroupID,
                RoomUID:        record.RoomUID,
                KindID:         record.KindID,
                BaseScore:      record.BaseScore,
                ScoreDelta:     player.ScoreDelta,
                IsWinner:       player.ScoreDelta > 0,
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

func (uc *GameRecordSyncUseCase) parsePlayers(playersJSON string) ([]*PlayerInfo, error) {
    var players []*PlayerInfo
    if err := json.Unmarshal([]byte(playersJSON), &players); err != nil {
        return nil, err
    }
    return players, nil
}

func (uc *GameRecordSyncUseCase) findGameAccountID(ctx context.Context, playerGID int64) *int32 {
    account, err := uc.repo.FindGameAccountByPlayerGID(ctx, playerGID)
    if err != nil || account == nil {
        return nil
    }
    return &account.ID
}
```

---

## Migration Checklist

### Phase 1: Data Migration

- [ ] **Export legacy data from MySQL**
  ```bash
  mysqldump -u user -p waiter_db > legacy_data.sql
  ```

- [ ] **Create data transformation scripts**
  - [ ] TPlayer → basic_user + game_account
  - [ ] TGamer → game_account_house
  - [ ] THouse → game_ctrl_account + game_session
  - [ ] TBattleRecord → game_battle_record + game_player_record
  - [ ] TManager → game_shop_admin
  - [ ] TRechargeRecord → game_member_wallet
  - [ ] TFeeRecord → game_fee_settle

- [ ] **Remove WeChat dependencies**
  - [ ] Strip WxKey fields
  - [ ] Strip WxName fields
  - [ ] Strip DeviceToken fields
  - [ ] Map to game account identifiers

- [ ] **Validate data integrity**
  - [ ] Check foreign key relationships
  - [ ] Verify unique constraints
  - [ ] Validate data completeness

### Phase 2: Session Management Migration

- [ ] **Implement plaza client in battle-tiles**
  - [ ] Port TCP connection logic
  - [ ] Implement auto-reconnect
  - [ ] Add command queue
  - [ ] Handle game events

- [ ] **Create session monitor service**
  - [ ] Health check mechanism
  - [ ] Auto-restart on failure
  - [ ] State tracking

- [ ] **Integrate with game_session table**
  - [ ] Create session on connect
  - [ ] Update state on events
  - [ ] Clean up on disconnect

### Phase 3: Service Migration

- [ ] **Remove WeChat notification service**
  - [ ] Replace with email/SMS notifications
  - [ ] Update notification templates
  - [ ] Test notification delivery

- [ ] **Migrate background jobs**
  - [ ] Port to asynq task queue
  - [ ] Implement sync jobs
  - [ ] Add monitoring

- [ ] **Update API endpoints**
  - [ ] Remove WeChat-specific endpoints
  - [ ] Add game account endpoints
  - [ ] Update authentication

### Phase 4: Testing

- [ ] **Unit tests for data transformation**
- [ ] **Integration tests for session management**
- [ ] **End-to-end tests for sync process**
- [ ] **Load tests for concurrent sessions**
- [ ] **Verify data accuracy**

---

## WeChat Dependency Removal

### Fields to Remove

| Legacy Table | Field | Replacement |
|--------------|-------|-------------|
| TPlayer | WxKey | game_account.account |
| TPlayer | WxName | basic_user.nick_name (from game) |
| THouse | OwnerWxKey | basic_user.id |
| TManager | WxKey | basic_user.id |
| TManager | WxName | basic_user.nick_name |
| TManager | DeviceToken | Removed (no replacement) |

### Notification Migration

**Legacy:** WeChat messages via wechaty  
**New:** Email/SMS/Push notifications

```go
// Old (WeChat)
func (s *Service) sendWeChatNotification(wxKey string, message string) error {
    // Send via wechaty
}

// New (Email/SMS)
func (s *NotificationService) sendNotification(userID int32, message string) error {
    user, _ := s.userRepo.Get(ctx, userID)
    
    // Send email
    s.emailClient.Send(user.Email, "Notification", message)
    
    // Send SMS (optional)
    if user.Mobile != "" {
        s.smsClient.Send(user.Mobile, message)
    }
    
    return nil
}
```

---

## Timeline

| Phase | Duration | Dependencies |
|-------|----------|--------------|
| Data Migration | 1 week | Database schema ready |
| Session Management | 2 weeks | Plaza client implementation |
| Service Migration | 1 week | Session management complete |
| Testing | 1 week | All migrations complete |
| **Total** | **5 weeks** | |

---

## Risks and Mitigation

### Risk 1: Data Loss During Migration
**Mitigation:**
- Full backup before migration
- Parallel run of old and new systems
- Data validation scripts
- Rollback plan

### Risk 2: Session Instability
**Mitigation:**
- Comprehensive testing
- Gradual rollout
- Monitoring and alerting
- Quick rollback capability

### Risk 3: WeChat User Resistance
**Mitigation:**
- Clear communication
- Migration support
- Gradual deprecation
- Alternative notification methods

---

## Success Criteria

- [ ] All data migrated successfully
- [ ] Zero data loss
- [ ] All sessions stable
- [ ] No WeChat dependencies remaining
- [ ] Performance equal or better than legacy
- [ ] All tests passing
- [ ] Documentation complete

---

## Support

For migration questions:
- **Data Migration**: [DBA contact]
- **Session Management**: [Backend lead]
- **WeChat Removal**: [Product owner]

