package biz

import (
	"context"
	"testing"
	"time"

	"battle-tiles/internal/dal/model/basic"
	"battle-tiles/internal/dal/model/game"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate tables
	err = db.AutoMigrate(
		&basic.BasicUser{},
		&game.GameAccount{},
		&game.GameAccountStoreBinding{},
		&game.GameSession{},
		&game.GameSyncLog{},
	)
	require.NoError(t, err)

	return db
}

// createTestUser creates a test user
func createTestUser(t *testing.T, db *gorm.DB, role string) *basic.BasicUser {
	user := &basic.BasicUser{
		Role:         role,
		GameNickname: "TestUser",
	}
	err := db.Create(user).Error
	require.NoError(t, err)
	return user
}

// createTestGameAccount creates a test game account
func createTestGameAccount(t *testing.T, db *gorm.DB, userID int32) *game.GameAccount {
	account := &game.GameAccount{
		UserID:             userID,
		Account:            "test_account",
		PwdMD5:             "test_pwd",
		Nickname:           "TestNickname",
		GameUserID:         "12345",
		VerificationStatus: game.VerificationStatusVerified,
		VerifiedAt:         timePtr(time.Now()),
	}
	err := db.Create(account).Error
	require.NoError(t, err)
	return account
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestAccountBindingService_BindGameAccountToStore(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create test data
	superAdmin := createTestUser(t, db, basic.UserRoleSuperAdmin)
	gameAccount := createTestGameAccount(t, db, superAdmin.Id)

	// Create service
	service := NewAccountBindingService(db)

	// Test binding
	binding, err := service.BindGameAccountToStore(ctx, gameAccount.Id, 100, superAdmin.Id)
	require.NoError(t, err)
	assert.NotNil(t, binding)
	assert.Equal(t, gameAccount.Id, binding.GameAccountID)
	assert.Equal(t, int32(100), binding.HouseGID)
	assert.Equal(t, game.BindingStatusActive, binding.Status)

	// Verify session was created
	var session game.GameSession
	err = db.Where("game_account_id = ? AND house_gid = ?", gameAccount.Id, 100).First(&session).Error
	require.NoError(t, err)
	assert.Equal(t, game.SessionStateActive, session.State)
	assert.True(t, session.AutoSyncEnabled)
	assert.Equal(t, game.SyncStatusIdle, session.SyncStatus)
}

func TestAccountBindingService_BindGameAccountToStore_AlreadyBound(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create test data
	superAdmin := createTestUser(t, db, basic.UserRoleSuperAdmin)
	gameAccount := createTestGameAccount(t, db, superAdmin.Id)

	// Create service
	service := NewAccountBindingService(db)

	// First binding
	_, err := service.BindGameAccountToStore(ctx, gameAccount.Id, 100, superAdmin.Id)
	require.NoError(t, err)

	// Try to bind to another store (should fail)
	_, err = service.BindGameAccountToStore(ctx, gameAccount.Id, 200, superAdmin.Id)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already bound")
}

func TestAccountBindingService_UnbindGameAccountFromStore(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create test data
	superAdmin := createTestUser(t, db, basic.UserRoleSuperAdmin)
	gameAccount := createTestGameAccount(t, db, superAdmin.Id)

	// Create service
	service := NewAccountBindingService(db)

	// Create binding
	_, err := service.BindGameAccountToStore(ctx, gameAccount.Id, 100, superAdmin.Id)
	require.NoError(t, err)

	// Unbind
	err = service.UnbindGameAccountFromStore(ctx, gameAccount.Id)
	require.NoError(t, err)

	// Verify binding is inactive
	var binding game.GameAccountStoreBinding
	err = db.Where("game_account_id = ?", gameAccount.Id).First(&binding).Error
	require.NoError(t, err)
	assert.Equal(t, game.BindingStatusInactive, binding.Status)

	// Verify session is inactive
	var session game.GameSession
	err = db.Where("game_account_id = ?", gameAccount.Id).First(&session).Error
	require.NoError(t, err)
	assert.Equal(t, game.SessionStateInactive, session.State)
	assert.False(t, session.AutoSyncEnabled)
}

func TestAccountBindingService_GetBindingWithSession(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create test data
	superAdmin := createTestUser(t, db, basic.UserRoleSuperAdmin)
	gameAccount := createTestGameAccount(t, db, superAdmin.Id)

	// Create service
	service := NewAccountBindingService(db)

	// Create binding
	_, err := service.BindGameAccountToStore(ctx, gameAccount.Id, 100, superAdmin.Id)
	require.NoError(t, err)

	// Get binding with session
	result, err := service.GetBindingWithSession(ctx, gameAccount.Id)
	require.NoError(t, err)
	assert.NotNil(t, result.Binding)
	assert.NotNil(t, result.Session)
	assert.Equal(t, gameAccount.Id, result.Binding.GameAccountID)
	assert.Equal(t, int32(100), result.Binding.HouseGID)
	assert.True(t, result.Session.AutoSyncEnabled)
}

func TestAccountBindingService_RestartSessionForBinding(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create test data
	superAdmin := createTestUser(t, db, basic.UserRoleSuperAdmin)
	gameAccount := createTestGameAccount(t, db, superAdmin.Id)

	// Create service
	service := NewAccountBindingService(db)

	// Create binding
	_, err := service.BindGameAccountToStore(ctx, gameAccount.Id, 100, superAdmin.Id)
	require.NoError(t, err)

	// Deactivate session manually
	err = db.Model(&game.GameSession{}).
		Where("game_account_id = ?", gameAccount.Id).
		Updates(map[string]interface{}{
			"state":             game.SessionStateInactive,
			"auto_sync_enabled": false,
		}).Error
	require.NoError(t, err)

	// Restart session
	err = service.RestartSessionForBinding(ctx, gameAccount.Id)
	require.NoError(t, err)

	// Verify session is active again
	var session game.GameSession
	err = db.Where("game_account_id = ?", gameAccount.Id).First(&session).Error
	require.NoError(t, err)
	assert.Equal(t, game.SessionStateActive, session.State)
	assert.True(t, session.AutoSyncEnabled)
}

func TestSessionManager_CreateSessionForBinding(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create test data
	superAdmin := createTestUser(t, db, basic.UserRoleSuperAdmin)
	gameAccount := createTestGameAccount(t, db, superAdmin.Id)

	// Create binding
	binding := &game.GameAccountStoreBinding{
		GameAccountID: gameAccount.Id,
		HouseGID:      100,
		BoundByUserID: superAdmin.Id,
		Status:        game.BindingStatusActive,
	}
	err := db.Create(binding).Error
	require.NoError(t, err)

	// Create session manager
	manager := NewSessionManager(db)

	// Create session
	session, err := manager.CreateSessionForBinding(ctx, gameAccount.Id, 100, superAdmin.Id)
	require.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, game.SessionStateActive, session.State)
	assert.True(t, session.AutoSyncEnabled)
}

func TestSessionManager_StartStopAutoSync(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create test data
	superAdmin := createTestUser(t, db, basic.UserRoleSuperAdmin)
	gameAccount := createTestGameAccount(t, db, superAdmin.Id)

	// Create session
	session := &game.GameSession{
		GameAccountID:   &gameAccount.Id,
		UserID:          superAdmin.Id,
		HouseGID:        100,
		State:           game.SessionStateActive,
		AutoSyncEnabled: true,
		SyncStatus:      game.SyncStatusIdle,
	}
	err := db.Create(session).Error
	require.NoError(t, err)

	// Create session manager
	manager := NewSessionManager(db)

	// Stop auto-sync
	err = manager.StopAutoSync(ctx, session.Id)
	require.NoError(t, err)

	// Verify
	var updated game.GameSession
	err = db.First(&updated, session.Id).Error
	require.NoError(t, err)
	assert.False(t, updated.AutoSyncEnabled)

	// Start auto-sync
	err = manager.StartAutoSync(ctx, session.Id)
	require.NoError(t, err)

	// Verify
	err = db.First(&updated, session.Id).Error
	require.NoError(t, err)
	assert.True(t, updated.AutoSyncEnabled)
}

func TestSessionManager_GetActiveSessionsForSync(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create test data
	superAdmin := createTestUser(t, db, basic.UserRoleSuperAdmin)
	gameAccount1 := createTestGameAccount(t, db, superAdmin.Id)
	gameAccount2 := createTestGameAccount(t, db, superAdmin.Id)

	// Create active session
	session1 := &game.GameSession{
		GameAccountID:   &gameAccount1.Id,
		UserID:          superAdmin.Id,
		HouseGID:        100,
		State:           game.SessionStateActive,
		AutoSyncEnabled: true,
		SyncStatus:      game.SyncStatusIdle,
	}
	err := db.Create(session1).Error
	require.NoError(t, err)

	// Create inactive session
	session2 := &game.GameSession{
		GameAccountID:   &gameAccount2.Id,
		UserID:          superAdmin.Id,
		HouseGID:        200,
		State:           game.SessionStateInactive,
		AutoSyncEnabled: false,
		SyncStatus:      game.SyncStatusIdle,
	}
	err = db.Create(session2).Error
	require.NoError(t, err)

	// Create session manager
	manager := NewSessionManager(db)

	// Get active sessions
	sessions, err := manager.GetActiveSessionsForSync(ctx)
	require.NoError(t, err)
	assert.Len(t, sessions, 1)
	assert.Equal(t, session1.Id, sessions[0].Id)
}

