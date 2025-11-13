package biz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"battle-tiles/internal/biz/validation"
	"battle-tiles/internal/dal/model/game"

	"gorm.io/gorm"
)

var (
	// ErrBindingNotFound is returned when binding is not found
	ErrBindingNotFound = errors.New("binding not found")

	// ErrBindingAlreadyExists is returned when binding already exists
	ErrBindingAlreadyExists = errors.New("binding already exists")
)

// AccountBindingService manages game account to store bindings
// and triggers auto-sync session creation
type AccountBindingService struct {
	db             *gorm.DB
	validator      *validation.AccountBindingValidator
	sessionManager *SessionManager
}

// NewAccountBindingService creates a new account binding service
func NewAccountBindingService(db *gorm.DB) *AccountBindingService {
	return &AccountBindingService{
		db:             db,
		validator:      validation.NewAccountBindingValidator(db),
		sessionManager: NewSessionManager(db),
	}
}

// BindGameAccountToStore binds a game account to a store
// Business Rule: One game account can only bind to ONE store
// This automatically creates a session and starts auto-sync
func (s *AccountBindingService) BindGameAccountToStore(ctx context.Context, gameAccountID, houseGID, userID int32) (*game.GameAccountStoreBinding, error) {
	// Validate the binding using business rules
	if err := s.validator.ValidateStoreBinding(ctx, gameAccountID, houseGID, userID); err != nil {
		return nil, err
	}

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create binding
	binding := &game.GameAccountStoreBinding{
		GameAccountID: gameAccountID,
		HouseGID:      houseGID,
		BoundByUserID: userID,
		Status:        game.BindingStatusActive,
	}

	if err := tx.Create(binding).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create binding: %w", err)
	}

	// Automatically create session for this binding
	session, err := s.createSessionInTx(ctx, tx, gameAccountID, houseGID, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log the session creation
	fmt.Printf("Created session %d for binding: game_account=%d, store=%d\n", session.Id, gameAccountID, houseGID)

	return binding, nil
}

// createSessionInTx creates a session within a transaction
func (s *AccountBindingService) createSessionInTx(ctx context.Context, tx *gorm.DB, gameAccountID, houseGID, userID int32) (*game.GameSession, error) {
	// Check if session already exists
	var existingSession game.GameSession
	err := tx.Where("game_account_id = ? AND house_gid = ? AND state = ?",
		gameAccountID, houseGID, game.SessionStateActive).
		First(&existingSession).Error

	if err == nil {
		// Session already exists, return it
		return &existingSession, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing session: %w", err)
	}

	// Create new session with auto-sync enabled
	session := &game.GameSession{
		GameAccountID:   &gameAccountID,
		UserID:          userID,
		HouseGID:        houseGID,
		State:           game.SessionStateActive,
		AutoSyncEnabled: true,
		SyncStatus:      game.SyncStatusIdle,
	}

	if err := tx.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// UnbindGameAccountFromStore unbinds a game account from a store
// This also deactivates the associated session
func (s *AccountBindingService) UnbindGameAccountFromStore(ctx context.Context, gameAccountID int32) error {
	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get the binding
	var binding game.GameAccountStoreBinding
	err := tx.Where("game_account_id = ? AND status = ?", gameAccountID, game.BindingStatusActive).
		First(&binding).Error

	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: game account %d", ErrBindingNotFound, gameAccountID)
		}
		return fmt.Errorf("failed to get binding: %w", err)
	}

	// Update binding status to inactive
	if err := tx.Model(&binding).Update("status", game.BindingStatusInactive).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update binding: %w", err)
	}

	// Deactivate associated sessions
	err = tx.Model(&game.GameSession{}).
		Where("game_account_id = ? AND house_gid = ? AND state = ?",
			gameAccountID, binding.HouseGID, game.SessionStateActive).
		Updates(map[string]interface{}{
			"state":             game.SessionStateInactive,
			"auto_sync_enabled": false,
			"sync_status":       game.SyncStatusIdle,
			"end_at":            time.Now(),
			"updated_at":        time.Now(),
		}).Error

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to deactivate sessions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetBindingByGameAccount retrieves the binding for a game account
func (s *AccountBindingService) GetBindingByGameAccount(ctx context.Context, gameAccountID int32) (*game.GameAccountStoreBinding, error) {
	var binding game.GameAccountStoreBinding
	err := s.db.WithContext(ctx).
		Where("game_account_id = ? AND status = ?", gameAccountID, game.BindingStatusActive).
		First(&binding).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: game account %d", ErrBindingNotFound, gameAccountID)
		}
		return nil, fmt.Errorf("failed to get binding: %w", err)
	}

	return &binding, nil
}

// GetBindingsByStore retrieves all bindings for a store
func (s *AccountBindingService) GetBindingsByStore(ctx context.Context, houseGID int32) ([]*game.GameAccountStoreBinding, error) {
	var bindings []*game.GameAccountStoreBinding
	err := s.db.WithContext(ctx).
		Where("house_gid = ? AND status = ?", houseGID, game.BindingStatusActive).
		Find(&bindings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get bindings: %w", err)
	}

	return bindings, nil
}

// GetBindingsByUser retrieves all bindings created by a user
func (s *AccountBindingService) GetBindingsByUser(ctx context.Context, userID int32) ([]*game.GameAccountStoreBinding, error) {
	var bindings []*game.GameAccountStoreBinding
	err := s.db.WithContext(ctx).
		Where("bound_by_user_id = ?", userID).
		Order("created_at DESC").
		Find(&bindings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get bindings: %w", err)
	}

	return bindings, nil
}

// GetBindingWithSession retrieves a binding along with its associated session
func (s *AccountBindingService) GetBindingWithSession(ctx context.Context, gameAccountID int32) (*BindingWithSession, error) {
	// Get binding
	binding, err := s.GetBindingByGameAccount(ctx, gameAccountID)
	if err != nil {
		return nil, err
	}

	// Get associated session
	var session game.GameSession
	err = s.db.WithContext(ctx).
		Where("game_account_id = ? AND house_gid = ?", gameAccountID, binding.HouseGID).
		Order("created_at DESC").
		First(&session).Error

	var sessionPtr *game.GameSession
	if err == nil {
		sessionPtr = &session
	}

	return &BindingWithSession{
		Binding: binding,
		Session: sessionPtr,
	}, nil
}

// BindingWithSession represents a binding with its associated session
type BindingWithSession struct {
	Binding *game.GameAccountStoreBinding `json:"binding"`
	Session *game.GameSession             `json:"session,omitempty"`
}

// GetAllActiveBindingsWithSessions retrieves all active bindings with their sessions
func (s *AccountBindingService) GetAllActiveBindingsWithSessions(ctx context.Context) ([]*BindingWithSession, error) {
	var bindings []*game.GameAccountStoreBinding
	err := s.db.WithContext(ctx).
		Where("status = ?", game.BindingStatusActive).
		Find(&bindings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get bindings: %w", err)
	}

	result := make([]*BindingWithSession, 0, len(bindings))
	for _, binding := range bindings {
		// Get associated session
		var session game.GameSession
		err = s.db.WithContext(ctx).
			Where("game_account_id = ? AND house_gid = ?", binding.GameAccountID, binding.HouseGID).
			Order("created_at DESC").
			First(&session).Error

		var sessionPtr *game.GameSession
		if err == nil {
			sessionPtr = &session
		}

		result = append(result, &BindingWithSession{
			Binding: binding,
			Session: sessionPtr,
		})
	}

	return result, nil
}

// RestartSessionForBinding restarts the session for a binding
// This is useful when a session encounters an error and needs to be restarted
func (s *AccountBindingService) RestartSessionForBinding(ctx context.Context, gameAccountID int32) error {
	// Get binding
	binding, err := s.GetBindingByGameAccount(ctx, gameAccountID)
	if err != nil {
		return err
	}

	// Get existing session
	var session game.GameSession
	err = s.db.WithContext(ctx).
		Where("game_account_id = ? AND house_gid = ?", gameAccountID, binding.HouseGID).
		Order("created_at DESC").
		First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No session exists, create a new one
			_, err = s.sessionManager.CreateSessionForBinding(ctx, gameAccountID, binding.HouseGID, binding.BoundByUserID)
			return err
		}
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Reactivate the session
	if !session.IsActive() {
		return s.sessionManager.ReactivateSession(ctx, session.Id)
	}

	// Session is already active, just restart auto-sync
	return s.sessionManager.StartAutoSync(ctx, session.Id)
}

