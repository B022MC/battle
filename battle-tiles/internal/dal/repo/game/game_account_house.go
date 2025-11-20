// internal/dal/repo/game/game_account_house.go
package game

import (
	"context"

	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// 用例需要的方法
type GameAccountHouseRepo interface {
	// 账号绑定的店铺（一个账号只能绑定一个店铺）
	HousesByAccount(ctx context.Context, gameAccountID int32) (*model.GameAccountHouse, error)
	// 创建账号店铺绑定
	Create(ctx context.Context, accountHouse *model.GameAccountHouse) error
	// 删除账号店铺绑定（通过 game_account_id 删除）
	DeleteByGameAccountID(ctx context.Context, gameAccountID int32) error
}

type gameAccountHouseRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewGameAccountHouseRepo(data *infra.Data, logger log.Logger) GameAccountHouseRepo {
	return &gameAccountHouseRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/game_account_house")),
	}
}

func (r *gameAccountHouseRepo) db(ctx context.Context) *gorm.DB { return r.data.GetDBWithContext(ctx) }

func (r *gameAccountHouseRepo) HousesByAccount(ctx context.Context, gameAccountID int32) (*model.GameAccountHouse, error) {
	var out model.GameAccountHouse
	err := r.db(ctx).
		Where("game_account_id = ?", gameAccountID).
		First(&out).Error
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Create 创建账号店铺绑定（使用 FirstOrCreate 实现幂等性）
func (r *gameAccountHouseRepo) Create(ctx context.Context, accountHouse *model.GameAccountHouse) error {
	return r.db(ctx).
		Where("game_account_id = ? AND house_gid = ?", accountHouse.GameAccountID, accountHouse.HouseGID).
		FirstOrCreate(accountHouse).Error
}

// DeleteByGameAccountID 删除账号店铺绑定（通过 game_account_id 删除）
func (r *gameAccountHouseRepo) DeleteByGameAccountID(ctx context.Context, gameAccountID int32) error {
	return r.db(ctx).
		Where("game_account_id = ?", gameAccountID).
		Delete(&model.GameAccountHouse{}).Error
}
