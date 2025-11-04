// internal/biz/game/game_shop_admin.go
package game

import (
	model "battle-tiles/internal/dal/model/game"
	repo "battle-tiles/internal/dal/repo/game"
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

type ShopAdminUseCase struct {
	repo repo.GameShopAdminRepo
	log  *log.Helper
}

func NewShopAdminUseCase(r repo.GameShopAdminRepo, logger log.Logger) *ShopAdminUseCase {
	return &ShopAdminUseCase{
		repo: r,
		log:  log.NewHelper(log.With(logger, "module", "usecase/shop_admin")),
	}
}

func validateRole(role string) (string, error) {
	role = strings.ToLower(strings.TrimSpace(role))
	if role == "" {
		return "admin", nil
	}
	switch role {
	case "admin", "operator":
		return role, nil
	default:
		return "", errors.New("invalid role (must be admin|operator)")
	}
}

func (uc *ShopAdminUseCase) Assign(ctx context.Context, houseGID int32, targetUserID int32, role string) error {
	if houseGID <= 0 || targetUserID <= 0 {
		return errors.New("invalid house_gid or user_id")
	}
	r, err := validateRole(role)
	if err != nil {
		return err
	}
	m := &model.GameShopAdmin{
		HouseGID: int32(houseGID),
		UserID:   targetUserID,
		Role:     r,
	}
	return uc.repo.Assign(ctx, m)
}

func (uc *ShopAdminUseCase) Revoke(ctx context.Context, houseGID int32, targetUserID int32) error {
	if houseGID <= 0 || targetUserID <= 0 {
		return errors.New("invalid house_gid or user_id")
	}
	return uc.repo.Revoke(ctx, houseGID, targetUserID)
}

func (uc *ShopAdminUseCase) IsAdmin(ctx context.Context, houseGID int32, userID int32) (bool, error) {
	return uc.repo.Exists(ctx, houseGID, userID)
}

func (uc *ShopAdminUseCase) List(ctx context.Context, houseGID int32) ([]*model.GameShopAdmin, error) {
	if houseGID <= 0 {
		return nil, errors.New("invalid house_gid")
	}
	return uc.repo.ListByHouse(ctx, houseGID)
}
