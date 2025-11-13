package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"battle-tiles/internal/biz"
	"battle-tiles/internal/dal/model/game"

	"gorm.io/gorm"
)

const (
	// Sync intervals
	SyncBattleRecordsInterval = 5 * time.Second
	SyncMemberListInterval    = 30 * time.Second
	SyncWalletInterval        = 10 * time.Second

	// Worker check interval
	WorkerCheckInterval = 10 * time.Second
)

// SyncWorker manages background synchronization tasks for game sessions
type SyncWorker struct {
	db             *gorm.DB
	sessionManager *biz.SessionManager
	sessions       map[int32]*SessionSyncTask // sessionID -> task
	mu             sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
}

// SessionSyncTask represents a sync task for a session
type SessionSyncTask struct {
	SessionID     int32
	GameAccountID int32
	HouseGID      int32
	Cancel        context.CancelFunc
	LastBattleSync time.Time
	LastMemberSync time.Time
	LastWalletSync time.Time
}

// NewSyncWorker creates a new sync worker
func NewSyncWorker(db *gorm.DB) *SyncWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &SyncWorker{
		db:             db,
		sessionManager: biz.NewSessionManager(db),
		sessions:       make(map[int32]*SessionSyncTask),
		ctx:            ctx,
		cancel:         cancel,
	}
}

// Start starts the sync worker
func (w *SyncWorker) Start() error {
	fmt.Println("Starting sync worker...")

	// Start the main worker loop
	w.wg.Add(1)
	go w.workerLoop()

	return nil
}

// Stop stops the sync worker
func (w *SyncWorker) Stop() error {
	fmt.Println("Stopping sync worker...")

	// Cancel context
	w.cancel()

	// Stop all session tasks
	w.mu.Lock()
	for sessionID, task := range w.sessions {
		task.Cancel()
		delete(w.sessions, sessionID)
	}
	w.mu.Unlock()

	// Wait for all goroutines to finish
	w.wg.Wait()

	fmt.Println("Sync worker stopped")
	return nil
}

// workerLoop is the main worker loop that checks for active sessions
func (w *SyncWorker) workerLoop() {
	defer w.wg.Done()

	ticker := time.NewTicker(WorkerCheckInterval)
	defer ticker.Stop()

	// Initial check
	w.checkAndStartSessions()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			w.checkAndStartSessions()
		}
	}
}

// checkAndStartSessions checks for active sessions and starts sync tasks
func (w *SyncWorker) checkAndStartSessions() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all active sessions with auto-sync enabled
	sessions, err := w.sessionManager.GetActiveSessionsForSync(ctx)
	if err != nil {
		fmt.Printf("Failed to get active sessions: %v\n", err)
		return
	}

	// Start sync tasks for new sessions
	for _, session := range sessions {
		if session.GameAccountID == nil {
			continue
		}

		w.mu.RLock()
		_, exists := w.sessions[session.Id]
		w.mu.RUnlock()

		if !exists {
			w.StartSyncForSession(session.Id, *session.GameAccountID, session.HouseGID)
		}
	}

	// Stop sync tasks for sessions that are no longer active
	w.mu.Lock()
	activeSessionIDs := make(map[int32]bool)
	for _, session := range sessions {
		activeSessionIDs[session.Id] = true
	}

	for sessionID, task := range w.sessions {
		if !activeSessionIDs[sessionID] {
			task.Cancel()
			delete(w.sessions, sessionID)
			fmt.Printf("Stopped sync task for session %d\n", sessionID)
		}
	}
	w.mu.Unlock()
}

// StartSyncForSession starts sync task for a specific session
func (w *SyncWorker) StartSyncForSession(sessionID, gameAccountID, houseGID int32) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Check if task already exists
	if _, exists := w.sessions[sessionID]; exists {
		return fmt.Errorf("sync task already exists for session %d", sessionID)
	}

	// Create task context
	taskCtx, taskCancel := context.WithCancel(w.ctx)

	task := &SessionSyncTask{
		SessionID:     sessionID,
		GameAccountID: gameAccountID,
		HouseGID:      houseGID,
		Cancel:        taskCancel,
	}

	w.sessions[sessionID] = task

	// Start sync goroutines
	w.wg.Add(3)
	go w.syncBattleRecordsLoop(taskCtx, task)
	go w.syncMemberListLoop(taskCtx, task)
	go w.syncWalletLoop(taskCtx, task)

	fmt.Printf("Started sync task for session %d (game_account=%d, house=%d)\n",
		sessionID, gameAccountID, houseGID)

	return nil
}

// StopSyncForSession stops sync task for a specific session
func (w *SyncWorker) StopSyncForSession(sessionID int32) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	task, exists := w.sessions[sessionID]
	if !exists {
		return fmt.Errorf("sync task not found for session %d", sessionID)
	}

	task.Cancel()
	delete(w.sessions, sessionID)

	fmt.Printf("Stopped sync task for session %d\n", sessionID)
	return nil
}

// syncBattleRecordsLoop syncs battle records periodically
func (w *SyncWorker) syncBattleRecordsLoop(ctx context.Context, task *SessionSyncTask) {
	defer w.wg.Done()

	ticker := time.NewTicker(SyncBattleRecordsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.syncBattleRecords(ctx, task); err != nil {
				fmt.Printf("Failed to sync battle records for session %d: %v\n", task.SessionID, err)
			}
		}
	}
}

// syncMemberListLoop syncs member list periodically
func (w *SyncWorker) syncMemberListLoop(ctx context.Context, task *SessionSyncTask) {
	defer w.wg.Done()

	ticker := time.NewTicker(SyncMemberListInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.syncMemberList(ctx, task); err != nil {
				fmt.Printf("Failed to sync member list for session %d: %v\n", task.SessionID, err)
			}
		}
	}
}

// syncWalletLoop syncs wallet updates periodically
func (w *SyncWorker) syncWalletLoop(ctx context.Context, task *SessionSyncTask) {
	defer w.wg.Done()

	ticker := time.NewTicker(SyncWalletInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.syncWallet(ctx, task); err != nil {
				fmt.Printf("Failed to sync wallet for session %d: %v\n", task.SessionID, err)
			}
		}
	}
}

// syncBattleRecords syncs battle records from game server
func (w *SyncWorker) syncBattleRecords(ctx context.Context, task *SessionSyncTask) error {
	// Update sync status
	if err := w.sessionManager.UpdateSyncStatus(ctx, task.SessionID, game.SyncStatusSyncing, ""); err != nil {
		return err
	}

	// TODO: Implement actual sync logic with plaza session
	// This should:
	// 1. Connect to game server using plaza session
	// 2. Fetch battle records
	// 3. Parse and store records in database
	// 4. Update player dimension data

	// Create sync log
	syncLog := &game.GameSyncLog{
		SessionID:     task.SessionID,
		SyncType:      game.SyncTypeBattleRecord,
		Status:        game.SyncStatusSuccess,
		RecordsSynced: 0, // TODO: actual count
		StartedAt:     time.Now(),
		CompletedAt:   timePtr(time.Now()),
	}

	if err := w.db.WithContext(ctx).Create(syncLog).Error; err != nil {
		fmt.Printf("Failed to create sync log: %v\n", err)
	}

	// Update sync status back to idle
	if err := w.sessionManager.UpdateSyncStatus(ctx, task.SessionID, game.SyncStatusIdle, ""); err != nil {
		return err
	}

	task.LastBattleSync = time.Now()
	return nil
}

// syncMemberList syncs member list from game server
func (w *SyncWorker) syncMemberList(ctx context.Context, task *SessionSyncTask) error {
	// TODO: Implement actual sync logic with plaza session
	// This should:
	// 1. Connect to game server
	// 2. Fetch member list
	// 3. Update member records in database

	// Create sync log
	syncLog := &game.GameSyncLog{
		SessionID:     task.SessionID,
		SyncType:      game.SyncTypeMemberList,
		Status:        game.SyncStatusSuccess,
		RecordsSynced: 0, // TODO: actual count
		StartedAt:     time.Now(),
		CompletedAt:   timePtr(time.Now()),
	}

	if err := w.db.WithContext(ctx).Create(syncLog).Error; err != nil {
		fmt.Printf("Failed to create sync log: %v\n", err)
	}

	task.LastMemberSync = time.Now()
	return nil
}

// syncWallet syncs wallet updates from game server
func (w *SyncWorker) syncWallet(ctx context.Context, task *SessionSyncTask) error {
	// TODO: Implement actual sync logic with plaza session
	// This should:
	// 1. Connect to game server
	// 2. Fetch wallet updates
	// 3. Update wallet records in database

	// Create sync log
	syncLog := &game.GameSyncLog{
		SessionID:     task.SessionID,
		SyncType:      game.SyncTypeWalletUpdate,
		Status:        game.SyncStatusSuccess,
		RecordsSynced: 0, // TODO: actual count
		StartedAt:     time.Now(),
		CompletedAt:   timePtr(time.Now()),
	}

	if err := w.db.WithContext(ctx).Create(syncLog).Error; err != nil {
		fmt.Printf("Failed to create sync log: %v\n", err)
	}

	task.LastWalletSync = time.Now()
	return nil
}

// GetActiveSessions returns all active sync tasks
func (w *SyncWorker) GetActiveSessions() []int32 {
	w.mu.RLock()
	defer w.mu.RUnlock()

	sessionIDs := make([]int32, 0, len(w.sessions))
	for sessionID := range w.sessions {
		sessionIDs = append(sessionIDs, sessionID)
	}

	return sessionIDs
}

// timePtr returns a pointer to a time.Time
func timePtr(t time.Time) *time.Time {
	return &t
}

