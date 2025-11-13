package validation

import (
	"context"
	"errors"
	"fmt"

	"battle-tiles/internal/dal/model/basic"
	"battle-tiles/internal/dal/model/game"

	"gorm.io/gorm"
)

var (
	// ErrGameAccountAlreadyBound is returned when game account is already bound to a store
	ErrGameAccountAlreadyBound = errors.New("game account is already bound to a store")

	// ErrStoreAdminAlreadyBound is returned when user is already admin of another store
	ErrStoreAdminAlreadyBound = errors.New("user is already administrator of another store")

	// ErrRegularUserMustHaveGameAccount is returned when regular user has no game account
	ErrRegularUserMustHaveGameAccount = errors.New("regular user must have at least one verified game account")

	// ErrGameAccountNotVerified is returned when game account is not verified
	ErrGameAccountNotVerified = errors.New("game account must be verified before binding")

	// ErrInvalidUserRole is returned when user role is invalid
	ErrInvalidUserRole = errors.New("invalid user role")
)

// AccountBindingValidator validates account binding business rules
type AccountBindingValidator struct {
	db *gorm.DB
}

// NewAccountBindingValidator creates a new validator
func NewAccountBindingValidator(db *gorm.DB) *AccountBindingValidator {
	return &AccountBindingValidator{db: db}
}

// ValidateGameAccountStoreBinding validates that a game account can be bound to a store
// Business Rule: One game account can only bind to ONE store
func (v *AccountBindingValidator) ValidateGameAccountStoreBinding(ctx context.Context, gameAccountID int32) error {
	var existingBinding game.GameAccountStoreBinding
	err := v.db.WithContext(ctx).
		Where("game_account_id = ? AND status = ?", gameAccountID, game.BindingStatusActive).
		First(&existingBinding).Error

	if err == nil {
		// Binding already exists
		return fmt.Errorf("%w: game account %d is already bound to store %d",
			ErrGameAccountAlreadyBound, gameAccountID, existingBinding.HouseGID)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Database error
		return fmt.Errorf("failed to check existing binding: %w", err)
	}

	// No existing binding found, validation passed
	return nil
}

// ValidateStoreAdminExclusiveBinding validates that a user can be assigned as store admin
// Business Rule: Store admin can only be admin of ONE store under ONE game account at a time
func (v *AccountBindingValidator) ValidateStoreAdminExclusiveBinding(ctx context.Context, userID int32) error {
	var count int64
	err := v.db.WithContext(ctx).
		Model(&game.GameShopAdmin{}).
		Where("user_id = ? AND is_exclusive = ? AND deleted_at IS NULL", userID, true).
		Count(&count).Error

	if err != nil {
		return fmt.Errorf("failed to check existing admin bindings: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("%w: user %d is already administrator of %d store(s)",
			ErrStoreAdminAlreadyBound, userID, count)
	}

	return nil
}

// ValidateGameAccountVerified validates that a game account is verified
func (v *AccountBindingValidator) ValidateGameAccountVerified(ctx context.Context, gameAccountID int32) error {
	var account game.GameAccount
	err := v.db.WithContext(ctx).
		Where("id = ?", gameAccountID).
		First(&account).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("game account %d not found", gameAccountID)
		}
		return fmt.Errorf("failed to check game account: %w", err)
	}

	if !account.IsVerified() {
		return fmt.Errorf("%w: game account %d has status %s",
			ErrGameAccountNotVerified, gameAccountID, account.VerificationStatus)
	}

	return nil
}

// ValidateRegularUserGameAccount validates that a regular user has at least one verified game account
func (v *AccountBindingValidator) ValidateRegularUserGameAccount(ctx context.Context, userID int32) error {
	var count int64
	err := v.db.WithContext(ctx).
		Model(&game.GameAccount{}).
		Where("user_id = ? AND verification_status = ?", userID, game.VerificationStatusVerified).
		Count(&count).Error

	if err != nil {
		return fmt.Errorf("failed to check user game accounts: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("%w: user %d has no verified game accounts",
			ErrRegularUserMustHaveGameAccount, userID)
	}

	return nil
}

// ValidateUserRoleChange validates that a user role change is allowed
func (v *AccountBindingValidator) ValidateUserRoleChange(ctx context.Context, userID int32, newRole string) error {
	// Validate role value
	validRoles := map[string]bool{
		basic.UserRoleSuperAdmin:  true,
		basic.UserRoleStoreAdmin:  true,
		basic.UserRoleRegularUser: true,
	}

	if !validRoles[newRole] {
		return fmt.Errorf("%w: %s", ErrInvalidUserRole, newRole)
	}

	// If changing to regular user, ensure they have a verified game account
	if newRole == basic.UserRoleRegularUser {
		if err := v.ValidateRegularUserGameAccount(ctx, userID); err != nil {
			return fmt.Errorf("cannot change to regular user role: %w", err)
		}
	}

	// If changing to store admin, ensure they don't have exclusive bindings
	if newRole == basic.UserRoleStoreAdmin {
		// Check if user is already admin of multiple stores
		var count int64
		err := v.db.WithContext(ctx).
			Model(&game.GameShopAdmin{}).
			Where("user_id = ? AND deleted_at IS NULL", userID).
			Count(&count).Error

		if err != nil {
			return fmt.Errorf("failed to check admin bindings: %w", err)
		}

		if count > 1 {
			return fmt.Errorf("user is admin of %d stores, must be admin of at most 1 store to have store_admin role", count)
		}
	}

	return nil
}

// ValidateGameAccountDeletion validates that a game account can be deleted
func (v *AccountBindingValidator) ValidateGameAccountDeletion(ctx context.Context, gameAccountID int32) error {
	// Check if game account is bound to any store
	var bindingCount int64
	err := v.db.WithContext(ctx).
		Model(&game.GameAccountStoreBinding{}).
		Where("game_account_id = ? AND status = ?", gameAccountID, game.BindingStatusActive).
		Count(&bindingCount).Error

	if err != nil {
		return fmt.Errorf("failed to check store bindings: %w", err)
	}

	if bindingCount > 0 {
		return fmt.Errorf("cannot delete game account: it is bound to %d store(s)", bindingCount)
	}

	// Check if this is the user's last verified game account
	var account game.GameAccount
	err = v.db.WithContext(ctx).
		Where("id = ?", gameAccountID).
		First(&account).Error

	if err != nil {
		return fmt.Errorf("failed to get game account: %w", err)
	}

	// Check user role
	var user basic.BasicUser
	err = v.db.WithContext(ctx).
		Where("id = ?", account.UserID).
		First(&user).Error

	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// If user is regular user, ensure they have at least one other verified game account
	if user.IsRegularUser() {
		var verifiedCount int64
		err = v.db.WithContext(ctx).
			Model(&game.GameAccount{}).
			Where("user_id = ? AND id != ? AND verification_status = ?",
				account.UserID, gameAccountID, game.VerificationStatusVerified).
			Count(&verifiedCount).Error

		if err != nil {
			return fmt.Errorf("failed to check other game accounts: %w", err)
		}

		if verifiedCount == 0 {
			return fmt.Errorf("cannot delete last verified game account for regular user")
		}
	}

	return nil
}

// ValidateStoreBinding validates all requirements for binding a game account to a store
func (v *AccountBindingValidator) ValidateStoreBinding(ctx context.Context, gameAccountID, houseGID int32, boundByUserID int32) error {
	// 1. Validate game account is verified
	if err := v.ValidateGameAccountVerified(ctx, gameAccountID); err != nil {
		return err
	}

	// 2. Validate game account is not already bound
	if err := v.ValidateGameAccountStoreBinding(ctx, gameAccountID); err != nil {
		return err
	}

	// 3. Validate user has permission (must be super admin)
	var user basic.BasicUser
	err := v.db.WithContext(ctx).
		Where("id = ?", boundByUserID).
		First(&user).Error

	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsSuperAdmin() {
		return fmt.Errorf("only super administrators can bind game accounts to stores")
	}

	return nil
}

