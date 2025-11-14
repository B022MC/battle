// internal/dal/repo/game/game_ctrl_account_house.go
package game

import (
	"context"

	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type GameCtrlAccountHouseRepo interface {
	// 绑定中控到店铺：通过 ctrl_id 解析其默认/最近的 game_account_id，写入 game_account_house（幂等）
	BindByCtrl(ctx context.Context, ctrlID int32, houseGID int32, status int32) error
	// 解绑中控与店铺：删除该 ctrl 旗下所有 game_account 在指定店铺的绑定
	UnbindByCtrl(ctx context.Context, ctrlID int32, houseGID int32) error
	// 按店铺列出所有中控账号
	ListByHouse(ctx context.Context, houseGID int32) ([]*model.GameCtrlAccount, error)
	// 按中控列出其绑定的店铺
	ListHousesByCtrl(ctx context.Context, ctrlID int32) ([]*model.GameCtrlAccountHouse, error)
	// 按中控ID列出其绑定的店铺
	ListByCtrlID(ctx context.Context, ctrlID int32) ([]*model.GameCtrlAccountHouse, error)
	ListHouseMapByCtrlIDs(ctx context.Context, ctrlIDs []int32) (map[int32][]int32, error)
	// 删除中控账号的所有绑定
	DeleteByCtrlID(ctx context.Context, ctrlID int32) error
	// 列出所有出现过的店铺号（去重）
	ListDistinctHouses(ctx context.Context) ([]int32, error)
}

type gameCtrlAccountHouseRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewCtrlAccountHouseRepo(data *infra.Data, logger log.Logger) GameCtrlAccountHouseRepo {
	return &gameCtrlAccountHouseRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/ctrlAccountHouse")),
	}
}

func (r *gameCtrlAccountHouseRepo) db(ctx context.Context) *gorm.DB {
	return r.data.GetDBWithContext(ctx)
}

func (r *gameCtrlAccountHouseRepo) BindByCtrl(ctx context.Context, ctrlID int32, houseGID int32, status int32) error {
	// 业务规则：每个中控账号只能绑定一个店铺
	// 先检查该中控账号是否已经绑定了其他店铺
	var existingHouses []int32
	err := r.db(ctx).
		Table("game_account_house AS l").
		Select("DISTINCT l.house_gid").
		Joins("JOIN game_account AS ga ON ga.id = l.game_account_id").
		Where("ga.ctrl_account_id = ? AND l.house_gid != ?", ctrlID, houseGID).
		Pluck("house_gid", &existingHouses).Error

	if err != nil {
		return err
	}

	if len(existingHouses) > 0 {
		return gorm.ErrInvalidData // 已经绑定了其他店铺，不允许再绑定
	}

	// 将 ctrlID 映射到某一个 game_account_id（优先 is_default=true，否则取最近一条）并做 UPSERT
	// 使用原生 SQL 以便一次性完成选择与插入/更新
	q := `
INSERT INTO game_account_house (game_account_id, house_gid, status)
SELECT id AS game_account_id, ?, ? FROM game_account
WHERE ctrl_account_id = ?
ORDER BY is_default DESC, id DESC
LIMIT 1
ON CONFLICT (game_account_id, house_gid)
DO UPDATE SET status = EXCLUDED.status, updated_at = NOW();`
	return r.db(ctx).Exec(q, houseGID, status, ctrlID).Error
}

func (r *gameCtrlAccountHouseRepo) UnbindByCtrl(ctx context.Context, ctrlID int32, houseGID int32) error {
	// 删除该 ctrl 旗下所有 game_account 在该店铺的绑定
	return r.db(ctx).
		Where("game_account_id IN (SELECT id FROM game_account WHERE ctrl_account_id = ?) AND house_gid = ?", ctrlID, houseGID).
		Delete(&model.GameCtrlAccountHouse{}).Error
}

func (r *gameCtrlAccountHouseRepo) ListByHouse(ctx context.Context, houseGID int32) ([]*model.GameCtrlAccount, error) {
	db := r.db(ctx)
	var out []*model.GameCtrlAccount
	err := db.Table("game_account_house AS l").
		Select("c.*").
		Joins("JOIN game_account AS ga ON ga.id = l.game_account_id").
		Joins("JOIN game_ctrl_account AS c ON c.id = ga.ctrl_account_id").
		Where("l.house_gid = ?", houseGID).
		Order("l.id DESC").
		Find(&out).Error
	return out, err
}

func (r *gameCtrlAccountHouseRepo) ListHousesByCtrl(ctx context.Context, ctrlID int32) ([]*model.GameCtrlAccountHouse, error) {
	var out []*model.GameCtrlAccountHouse
	err := r.db(ctx).
		Table("game_account_house AS l").
		Select("l.*").
		Joins("JOIN game_account AS ga ON ga.id = l.game_account_id").
		Where("ga.ctrl_account_id = ?", ctrlID).
		Order("l.id DESC").
		Find(&out).Error
	return out, err
}
func (r *gameCtrlAccountHouseRepo) ListHouseMapByCtrlIDs(ctx context.Context, ctrlIDs []int32) (map[int32][]int32, error) {
	if len(ctrlIDs) == 0 {
		return map[int32][]int32{}, nil
	}
	type row struct {
		CtrlID   int32 `gorm:"column:ctrl_id"`
		HouseGID int32 `gorm:"column:house_gid"`
	}
	var rows []row
	if err := r.db(ctx).
		Table("game_account_house AS l").
		Select("ga.ctrl_account_id AS ctrl_id, l.house_gid").
		Joins("JOIN game_account AS ga ON ga.id = l.game_account_id").
		Where("ga.ctrl_account_id IN ?", ctrlIDs).
		Order("ga.ctrl_account_id ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	m := make(map[int32][]int32, len(ctrlIDs))
	for _, id := range ctrlIDs {
		m[id] = []int32{}
	}
	for _, r := range rows {
		m[r.CtrlID] = append(m[r.CtrlID], r.HouseGID)
	}
	return m, nil
}

func (r *gameCtrlAccountHouseRepo) ListByCtrlID(ctx context.Context, ctrlID int32) ([]*model.GameCtrlAccountHouse, error) {
	var rows []*model.GameCtrlAccountHouse
	err := r.db(ctx).
		Table("game_account_house AS l").
		Select("l.*").
		Joins("JOIN game_account AS ga ON ga.id = l.game_account_id").
		Where("ga.ctrl_account_id = ?", ctrlID).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *gameCtrlAccountHouseRepo) DeleteByCtrlID(ctx context.Context, ctrlID int32) error {
	// 删除该中控账号下所有 game_account 的绑定关系
	return r.db(ctx).
		Exec(`DELETE FROM game_account_house
			  WHERE game_account_id IN (
				  SELECT id FROM game_account WHERE ctrl_account_id = ?
			  )`, ctrlID).Error
}

func (r *gameCtrlAccountHouseRepo) ListDistinctHouses(ctx context.Context) ([]int32, error) {
	var rows []int32
	if err := r.db(ctx).Table("game_account_house").Distinct("house_gid").Pluck("house_gid", &rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
