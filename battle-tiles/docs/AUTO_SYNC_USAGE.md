# Auto-Sync Store Session Feature - Usage Guide

## Overview

The auto-sync store session feature automatically creates game sessions and starts synchronization when a super admin binds a game account to a store. This eliminates the need for manual session management and ensures continuous data synchronization.

## Architecture

```
Super Admin Binds Account → Create Binding → Create Session → Start Sync Worker
                                    ↓              ↓              ↓
                            game_account_store_binding  game_session  Background Sync
```

## Components

### 1. SessionManager (`internal/biz/session_manager.go`)
Manages game session lifecycle and sync status.

**Key Functions:**
- `CreateSessionForBinding()` - Creates session when binding is created
- `StartAutoSync()` - Enables auto-sync for a session
- `StopAutoSync()` - Disables auto-sync for a session
- `GetActiveSessionsForSync()` - Gets all sessions that need syncing
- `UpdateSyncStatus()` - Updates sync status
- `DeactivateSession()` - Deactivates a session
- `ReactivateSession()` - Reactivates an inactive session

### 2. AccountBindingService (`internal/biz/account_binding_service.go`)
Manages game account to store bindings and triggers session creation.

**Key Functions:**
- `BindGameAccountToStore()` - Creates binding and automatically creates session
- `UnbindGameAccountFromStore()` - Removes binding and deactivates session
- `GetBindingWithSession()` - Gets binding with its associated session
- `RestartSessionForBinding()` - Restarts session for a binding

### 3. SyncWorker (`internal/worker/sync_worker.go`)
Background worker that performs actual synchronization tasks.

**Key Functions:**
- `Start()` - Starts the worker
- `Stop()` - Stops the worker
- `StartSyncForSession()` - Starts sync for a specific session
- `StopSyncForSession()` - Stops sync for a specific session

## Usage Examples

### Example 1: Bind Game Account to Store (Triggers Auto-Sync)

```go
package main

import (
    "context"
    "fmt"
    "battle-tiles/internal/biz"
    "gorm.io/gorm"
)

func bindAccountAndStartSync(db *gorm.DB) {
    ctx := context.Background()
    
    // Create binding service
    bindingService := biz.NewAccountBindingService(db)
    
    // Super admin binds game account to store
    // This automatically:
    // 1. Validates the binding (business rules)
    // 2. Creates the binding record
    // 3. Creates a game session
    // 4. Enables auto-sync
    binding, err := bindingService.BindGameAccountToStore(
        ctx,
        gameAccountID: 123,  // Game account ID
        houseGID: 456,       // Store ID
        userID: 1,           // Super admin user ID
    )
    
    if err != nil {
        fmt.Printf("Failed to bind account: %v\n", err)
        return
    }
    
    fmt.Printf("Binding created: %+v\n", binding)
    
    // Get binding with session info
    bindingWithSession, err := bindingService.GetBindingWithSession(ctx, 123)
    if err != nil {
        fmt.Printf("Failed to get binding: %v\n", err)
        return
    }
    
    fmt.Printf("Session created: %+v\n", bindingWithSession.Session)
    fmt.Printf("Auto-sync enabled: %v\n", bindingWithSession.Session.AutoSyncEnabled)
}
```

### Example 2: Start Sync Worker

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "battle-tiles/internal/worker"
    "gorm.io/gorm"
)

func startSyncWorker(db *gorm.DB) {
    // Create sync worker
    syncWorker := worker.NewSyncWorker(db)
    
    // Start the worker
    if err := syncWorker.Start(); err != nil {
        fmt.Printf("Failed to start sync worker: %v\n", err)
        return
    }
    
    fmt.Println("Sync worker started successfully")
    
    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan
    
    // Stop the worker gracefully
    fmt.Println("Shutting down sync worker...")
    if err := syncWorker.Stop(); err != nil {
        fmt.Printf("Failed to stop sync worker: %v\n", err)
    }
}
```

### Example 3: Check Sync Status

```go
package main

import (
    "context"
    "fmt"
    "battle-tiles/internal/biz"
    "gorm.io/gorm"
)

func checkSyncStatus(db *gorm.DB, sessionID int32) {
    ctx := context.Background()
    
    // Create session manager
    sessionManager := biz.NewSessionManager(db)
    
    // Get sync status
    status, err := sessionManager.GetSyncStatus(ctx, sessionID)
    if err != nil {
        fmt.Printf("Failed to get sync status: %v\n", err)
        return
    }
    
    fmt.Printf("Session ID: %d\n", status.SessionID)
    fmt.Printf("State: %s\n", status.State)
    fmt.Printf("Auto-sync enabled: %v\n", status.AutoSyncEnabled)
    fmt.Printf("Sync status: %s\n", status.SyncStatus)
    fmt.Printf("Last sync at: %v\n", status.LastSyncAt)
    
    if status.ErrorMsg != "" {
        fmt.Printf("Error: %s\n", status.ErrorMsg)
    }
    
    if status.LastSyncLog != nil {
        fmt.Printf("Last sync type: %s\n", status.LastSyncLog.SyncType)
        fmt.Printf("Records synced: %d\n", status.LastSyncLog.RecordsSynced)
        fmt.Printf("Duration: %v\n", status.LastSyncLog.Duration())
    }
}
```

### Example 4: Manually Control Sync

```go
package main

import (
    "context"
    "fmt"
    "battle-tiles/internal/biz"
    "gorm.io/gorm"
)

func controlSync(db *gorm.DB, sessionID int32) {
    ctx := context.Background()
    sessionManager := biz.NewSessionManager(db)
    
    // Stop auto-sync
    if err := sessionManager.StopAutoSync(ctx, sessionID); err != nil {
        fmt.Printf("Failed to stop auto-sync: %v\n", err)
        return
    }
    fmt.Println("Auto-sync stopped")
    
    // Do some maintenance work...
    
    // Restart auto-sync
    if err := sessionManager.StartAutoSync(ctx, sessionID); err != nil {
        fmt.Printf("Failed to start auto-sync: %v\n", err)
        return
    }
    fmt.Println("Auto-sync restarted")
}
```

### Example 5: Unbind Account (Stops Sync)

```go
package main

import (
    "context"
    "fmt"
    "battle-tiles/internal/biz"
    "gorm.io/gorm"
)

func unbindAccount(db *gorm.DB, gameAccountID int32) {
    ctx := context.Background()
    
    // Create binding service
    bindingService := biz.NewAccountBindingService(db)
    
    // Unbind account from store
    // This automatically:
    // 1. Marks binding as inactive
    // 2. Deactivates the session
    // 3. Stops auto-sync
    err := bindingService.UnbindGameAccountFromStore(ctx, gameAccountID)
    if err != nil {
        fmt.Printf("Failed to unbind account: %v\n", err)
        return
    }
    
    fmt.Println("Account unbound and sync stopped")
}
```

### Example 6: Restart Failed Session

```go
package main

import (
    "context"
    "fmt"
    "battle-tiles/internal/biz"
    "gorm.io/gorm"
)

func restartFailedSession(db *gorm.DB, gameAccountID int32) {
    ctx := context.Background()
    
    // Create binding service
    bindingService := biz.NewAccountBindingService(db)
    
    // Restart session for binding
    // This will:
    // 1. Check if session exists
    // 2. Reactivate if inactive
    // 3. Restart auto-sync
    err := bindingService.RestartSessionForBinding(ctx, gameAccountID)
    if err != nil {
        fmt.Printf("Failed to restart session: %v\n", err)
        return
    }
    
    fmt.Println("Session restarted successfully")
}
```

## Integration with Application

### In main.go or server initialization:

```go
package main

import (
    "battle-tiles/internal/worker"
    "gorm.io/gorm"
)

var syncWorker *worker.SyncWorker

func initSyncWorker(db *gorm.DB) error {
    syncWorker = worker.NewSyncWorker(db)
    return syncWorker.Start()
}

func shutdownSyncWorker() error {
    if syncWorker != nil {
        return syncWorker.Stop()
    }
    return nil
}

func main() {
    // ... initialize database ...
    
    // Start sync worker
    if err := initSyncWorker(db); err != nil {
        panic(err)
    }
    defer shutdownSyncWorker()
    
    // ... start HTTP server ...
}
```

## Business Rules Enforced

1. **One Game Account → One Store**
   - Enforced by unique constraint in database
   - Validated before binding creation
   - Prevents multiple bindings for same account

2. **Automatic Session Creation**
   - Session created automatically when binding is created
   - No manual session management needed
   - Session linked to both game account and store

3. **Auto-Sync by Default**
   - All new sessions have auto-sync enabled
   - Sync starts immediately after session creation
   - Can be manually controlled if needed

4. **Graceful Shutdown**
   - Sessions deactivated when binding is removed
   - Sync workers stopped gracefully
   - No data loss during shutdown

## Sync Operations

The sync worker performs three types of synchronization:

### 1. Battle Records Sync (Every 5 seconds)
- Fetches new battle records from game server
- Stores records with player dimension data
- Updates player statistics

### 2. Member List Sync (Every 30 seconds)
- Fetches current member list
- Updates member status
- Tracks member changes

### 3. Wallet Sync (Every 10 seconds)
- Fetches wallet updates
- Records balance changes
- Tracks recharge/withdrawal operations

## Monitoring

### Check Active Sessions:
```go
activeSessions := syncWorker.GetActiveSessions()
fmt.Printf("Active sync sessions: %v\n", activeSessions)
```

### Check Sync Logs:
```sql
SELECT * FROM game_sync_log 
WHERE session_id = ? 
ORDER BY started_at DESC 
LIMIT 10;
```

### Check Session Status:
```sql
SELECT id, state, auto_sync_enabled, sync_status, last_sync_at, error_msg
FROM game_session
WHERE state = 'active';
```

## Error Handling

The system handles errors gracefully:

1. **Connection Errors**: Session marked as error state, auto-retry
2. **Sync Errors**: Logged in sync_log table, session continues
3. **Validation Errors**: Binding creation fails, no session created
4. **Database Errors**: Transaction rolled back, no partial state

## Next Steps

1. **Implement Plaza Session Integration**
   - Port session logic from legacy project
   - Implement protocol handlers
   - Add connection management

2. **Add Monitoring Dashboard**
   - Real-time sync status
   - Error alerts
   - Performance metrics

3. **Optimize Sync Performance**
   - Batch operations
   - Connection pooling
   - Caching strategies

## Notes

- The sync worker automatically discovers new sessions
- Sessions are checked every 10 seconds
- Failed sessions are automatically retried
- All operations are logged for debugging
- Transaction safety ensures data consistency

---

**Status**: Core Implementation Complete ✅  
**Next**: Plaza Session Integration  
**Version**: 1.0  
**Last Updated**: 2025-11-12

