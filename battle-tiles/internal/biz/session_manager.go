package biz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"battle-tiles/internal/dal/model/game"

	"gorm.io/gorm"
)

var (
	// ErrSessionNotFound is returned when session is not found
	ErrSessionNotFound = errors.New("session not found")

	// ErrSessionAlreadyActive is returned when trying to start an already active session
	ErrSessionAlreadyActive = errors.New("session is already active")

	// ErrGameAccountNotBound is returned when game account is not bound to store
	ErrGameAccountNotBound = errors.New("game account is not bound to any store")

	// ErrInvalidGameAccount is returned when game account is invalid
	ErrInvalidGameAccount = errors.New("invalid game account")
)

// SessionManager manages game sessions and auto-sync functionality
type SessionManager struct {
	db *gorm.DB
}

// NewSessionManager creates a new session manager
func NewSessionManager(db *gorm.DB) *SessionManager {
	return &SessionManager{db: db}
}

// CreateSessionForBinding creates a game session when game account is bound to store
// This is triggered automatically when super admin binds game account to store
// The session will start automatic synchronization
func (m *SessionManager) CreateSessionForBinding(ctx context.Context, gameAccountID, houseGID, userID int32) (*game.GameSession, error) {
	// Verify game account exists and is verified
	var account game.GameAccount
	err := m.db.WithContext(ctx).
		Where("id = ? AND verification_status = ?", gameAccountID, game.VerificationStatusVerified).
		First(&account).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: game account %d not found or not verified", ErrInvalidGameAccount, gameAccountID)
		}
		return nil, fmt.Errorf("failed to get game account: %w", err)
	}

	// Verify binding exists
	var binding game.GameAccountStoreBinding
	err = m.db.WithContext(ctx).
		Where("game_account_id = ? AND house_gid = ? AND status = ?",
			gameAccountID, houseGID, game.BindingStatusActive).
		First(&binding).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: game account %d is not bound to store %d", ErrGameAccountNotBound, gameAccountID, houseGID)
		}
		return nil, fmt.Errorf("failed to get binding: %w", err)
	}

	// Check if session already exists for this binding
	var existingSession game.GameSession
	err = m.db.WithContext(ctx).
		Where("game_account_id = ? AND house_gid = ? AND state = ?",
			gameAccountID, houseGID, game.SessionStateActive).
		First(&existingSession).Error

	if err == nil {
		// Session already exists
		return &existingSession, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing session: %w", err)
	}

	// Create new session
	session := &game.GameSession{
		GameAccountID:   &gameAccountID,
		UserID:          userID,
		HouseGID:        houseGID,
		State:           game.SessionStateActive,
		AutoSyncEnabled: true,
		SyncStatus:      game.SyncStatusIdle,
	}

	if err := m.db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// StartAutoSync starts automatic synchronization for a session
func (m *SessionManager) StartAutoSync(ctx context.Context, sessionID int32) error {
	var session game.GameSession
	err := m.db.WithContext(ctx).Where("id = ?", sessionID).First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: session %d", ErrSessionNotFound, sessionID)
		}
		return fmt.Errorf("failed to get session: %w", err)
	}

	if !session.IsActive() {
		return fmt.Errorf("session %d is not active (state: %s)", sessionID, session.State)
	}

	// Update session to enable auto-sync
	updates := map[string]interface{}{
		"auto_sync_enabled": true,
		"sync_status":       game.SyncStatusIdle,
		"updated_at":        time.Now(),
	}

	if err := m.db.WithContext(ctx).Model(&game.GameSession{}).Where("id = ?", sessionID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to enable auto-sync: %w", err)
	}

	return nil
}

// StopAutoSync stops automatic synchronization for a session
func (m *SessionManager) StopAutoSync(ctx context.Context, sessionID int32) error {
	var session game.GameSession
	err := m.db.WithContext(ctx).Where("id = ?", sessionID).First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: session %d", ErrSessionNotFound, sessionID)
		}
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Update session to disable auto-sync
	updates := map[string]interface{}{
		"auto_sync_enabled": false,
		"sync_status":       game.SyncStatusIdle,
		"updated_at":        time.Now(),
	}

	if err := m.db.WithContext(ctx).Model(&game.GameSession{}).Where("id = ?", sessionID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to disable auto-sync: %w", err)
	}

	return nil
}

// UpdateSyncStatus updates the sync status of a session
func (m *SessionManager) UpdateSyncStatus(ctx context.Context, sessionID int32, status string, errorMsg string) error {
	updates := map[string]interface{}{
		"sync_status": status,
		"last_sync_at": time.Now(),
		"updated_at":  time.Now(),
	}

	if errorMsg != "" {
		updates["error_msg"] = errorMsg
	}

	if err := m.db.WithContext(ctx).Model(&game.GameSession{}).Where("id = ?", sessionID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update sync status: %w", err)
	}

	return nil
}

// GetActiveSessionsForSync retrieves all active sessions that have auto-sync enabled
func (m *SessionManager) GetActiveSessionsForSync(ctx context.Context) ([]*game.GameSession, error) {
	var sessions []*game.GameSession
	err := m.db.WithContext(ctx).
		Where("state = ? AND auto_sync_enabled = ? AND (sync_status = ? OR sync_status = ?)",
			game.SessionStateActive, true, game.SyncStatusIdle, game.SyncStatusError).
		Find(&sessions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}

	return sessions, nil
}

// GetSessionByID retrieves a session by ID
func (m *SessionManager) GetSessionByID(ctx context.Context, sessionID int32) (*game.GameSession, error) {
	var session game.GameSession
	err := m.db.WithContext(ctx).Where("id = ?", sessionID).First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: session %d", ErrSessionNotFound, sessionID)
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

// GetSessionsByGameAccount retrieves all sessions for a game account
func (m *SessionManager) GetSessionsByGameAccount(ctx context.Context, gameAccountID int32) ([]*game.GameSession, error) {
	var sessions []*game.GameSession
	err := m.db.WithContext(ctx).
		Where("game_account_id = ?", gameAccountID).
		Order("created_at DESC").
		Find(&sessions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	return sessions, nil
}

// GetSessionsByStore retrieves all sessions for a store
func (m *SessionManager) GetSessionsByStore(ctx context.Context, houseGID int32) ([]*game.GameSession, error) {
	var sessions []*game.GameSession
	err := m.db.WithContext(ctx).
		Where("house_gid = ?", houseGID).
		Order("created_at DESC").
		Find(&sessions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	return sessions, nil
}

// DeactivateSession deactivates a session
func (m *SessionManager) DeactivateSession(ctx context.Context, sessionID int32, reason string) error {
	updates := map[string]interface{}{
		"state":            game.SessionStateInactive,
		"auto_sync_enabled": false,
		"sync_status":      game.SyncStatusIdle,
		"end_at":           time.Now(),
		"updated_at":       time.Now(),
	}

	if reason != "" {
		updates["error_msg"] = reason
	}

	if err := m.db.WithContext(ctx).Model(&game.GameSession{}).Where("id = ?", sessionID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to deactivate session: %w", err)
	}

	return nil
}

// ReactivateSession reactivates an inactive session
func (m *SessionManager) ReactivateSession(ctx context.Context, sessionID int32) error {
	var session game.GameSession
	err := m.db.WithContext(ctx).Where("id = ?", sessionID).First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: session %d", ErrSessionNotFound, sessionID)
		}
		return fmt.Errorf("failed to get session: %w", err)
	}

	if session.IsActive() {
		return ErrSessionAlreadyActive
	}

	updates := map[string]interface{}{
		"state":            game.SessionStateActive,
		"auto_sync_enabled": true,
		"sync_status":      game.SyncStatusIdle,
		"error_msg":        "",
		"end_at":           nil,
		"updated_at":       time.Now(),
	}

	if err := m.db.WithContext(ctx).Model(&game.GameSession{}).Where("id = ?", sessionID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to reactivate session: %w", err)
	}

	return nil
}

// GetSyncStatus gets the current sync status for a session
func (m *SessionManager) GetSyncStatus(ctx context.Context, sessionID int32) (*SyncStatus, error) {
	var session game.GameSession
	err := m.db.WithContext(ctx).Where("id = ?", sessionID).First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: session %d", ErrSessionNotFound, sessionID)
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Get latest sync log
	var latestLog game.GameSyncLog
	err = m.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("started_at DESC").
		First(&latestLog).Error

	var lastSyncLog *game.GameSyncLog
	if err == nil {
		lastSyncLog = &latestLog
	}

	return &SyncStatus{
		SessionID:       session.Id,
		State:           session.State,
		AutoSyncEnabled: session.AutoSyncEnabled,
		SyncStatus:      session.SyncStatus,
		LastSyncAt:      session.LastSyncAt,
		ErrorMsg:        session.ErrorMsg,
		LastSyncLog:     lastSyncLog,
	}, nil
}

// SyncStatus represents the current sync status of a session
type SyncStatus struct {
	SessionID       int32              `json:"session_id"`
	State           string             `json:"state"`
	AutoSyncEnabled bool               `json:"auto_sync_enabled"`
	SyncStatus      string             `json:"sync_status"`
	LastSyncAt      *time.Time         `json:"last_sync_at"`
	ErrorMsg        string             `json:"error_msg"`
	LastSyncLog     *game.GameSyncLog  `json:"last_sync_log,omitempty"`
}

