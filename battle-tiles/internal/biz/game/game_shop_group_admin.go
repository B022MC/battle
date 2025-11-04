package game

import (
	repo "battle-tiles/internal/dal/repo/game"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type ShopGroupAdminUseCase struct {
	repo repo.GameShopGroupAdminRepo
	log  *log.Helper
}

func NewShopGroupAdminUseCase(r repo.GameShopGroupAdminRepo, logger log.Logger) *ShopGroupAdminUseCase {
	return &ShopGroupAdminUseCase{repo: r, log: log.NewHelper(log.With(logger, "module", "usecase/shop_group_admin"))}
}

func (uc *ShopGroupAdminUseCase) IsGroupAdmin(ctx context.Context, houseGID, groupID, userID int32) (bool, error) {
	return uc.repo.Exists(ctx, houseGID, groupID, userID)
}

func (uc *ShopGroupAdminUseCase) ListByUser(ctx context.Context, userID int32) ([][2]int32, error) {
	rows, err := uc.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([][2]int32, 0, len(rows))
	for _, r := range rows {
		out = append(out, [2]int32{r.HouseGID, r.GroupID})
	}
	return out, nil
}
