// internal/dal/repo/game/game_account_house.go
package game

import (
	"context"

	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 用例需要的方法
type GameAccountHouseRepo interface {
	// 插入或更新一条绑定（唯一键：game_account_id + house_gid）
	Upsert(ctx context.Context, m *model.GameAccountHouse) error
	// 将该账号在该店铺设为默认（会清除该账号其它店铺的默认；若不存在绑定，则创建为默认）
	SetDefault(ctx context.Context, gameAccountID int32, houseGID int32) error
	// 校验该账号在该店铺的绑定是否为启用状态
	ExistsActive(ctx context.Context, gameAccountID int32, houseGID int32) (bool, error)
	// 列出某账号绑定到的所有店铺
	ListHousesByAccount(ctx context.Context, gameAccountID int32) ([]*model.GameAccountHouse, error)
	// 列出某店铺下的所有账号绑定
	ListByHouse(ctx context.Context, houseGID int32) ([]*model.GameAccountHouse, error)
	// 确保存在一条 house 记录；若不存在则插入（默认 auto_reconnect=true）
	Ensure(ctx context.Context, houseGID int32, name string) error
}

type gameAccountHouseRepo struct {
	data *infra.Data
	log  *log.Helper
}

func (r *gameAccountHouseRepo) Ensure(ctx context.Context, houseGID int32, name string) error {
	m := &model.GameAccountHouse{
		HouseGID: (houseGID),
		//Name:          name, // 没有就传 ""，后续可再改名
	}
	return r.db(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "house_gid"}},
			DoNothing: true, // 已存在就跳过
		}).
		Create(m).Error
}

func NewGameAccountHouseRepo(data *infra.Data, logger log.Logger) GameAccountHouseRepo {
	return &gameAccountHouseRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/game_account_house")),
	}
}

func (r *gameAccountHouseRepo) db(ctx context.Context) *gorm.DB { return r.data.GetDBWithContext(ctx) }

func (r *gameAccountHouseRepo) Upsert(ctx context.Context, m *model.GameAccountHouse) error {
	return r.db(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "game_account_id"}, {Name: "house_gid"}},
			DoUpdates: clause.AssignmentColumns([]string{"is_default", "status", "updated_at"}),
		}).
		Create(m).Error
}

func (r *gameAccountHouseRepo) SetDefault(ctx context.Context, gameAccountID int32, houseGID int32) error {
	db := r.db(ctx)

	// 清掉该账号下其它默认
	if err := db.
		Model(&model.GameAccountHouse{}).
		Where("game_account_id = ? AND is_default = TRUE", gameAccountID).
		Update("is_default", false).Error; err != nil {
		return err
	}

	// 目标置为默认（若不存在则插入）
	return db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "game_account_id"}, {Name: "house_gid"}},
			DoUpdates: clause.Assignments(map[string]any{"is_default": true, "status": 1}),
		}).
		Create(&model.GameAccountHouse{
			GameAccountID: gameAccountID,
			HouseGID:      houseGID,
			IsDefault:     true,
			Status:        1,
		}).Error
}

func (r *gameAccountHouseRepo) ExistsActive(ctx context.Context, gameAccountID int32, houseGID int32) (bool, error) {
	var cnt int64
	err := r.db(ctx).
		Model(&model.GameAccountHouse{}).
		Where("game_account_id = ? AND house_gid = ? AND status = 1", gameAccountID, houseGID).
		Count(&cnt).Error
	return cnt > 0, err
}

func (r *gameAccountHouseRepo) ListHousesByAccount(ctx context.Context, gameAccountID int32) ([]*model.GameAccountHouse, error) {
	var out []*model.GameAccountHouse
	err := r.db(ctx).
		Where("game_account_id = ?", gameAccountID).
		Order("is_default DESC, id DESC").
		Find(&out).Error
	return out, err
}

func (r *gameAccountHouseRepo) ListByHouse(ctx context.Context, houseGID int32) ([]*model.GameAccountHouse, error) {
	var out []*model.GameAccountHouse
	err := r.db(ctx).
		Where("house_gid = ?", houseGID).
		Order("is_default DESC, id DESC").
		Find(&out).Error
	return out, err
}
