// internal/dal/repo/game/game_account.go
package game

import (
	"context"
	"strings"

	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"
	"battle-tiles/pkg/plugin/gormx/repo"

	"github.com/go-kratos/kratos/v2/log"
)

type GameAccountRepo interface {
	Create(ctx context.Context, a *model.GameAccount) error
	GetOneByUser(ctx context.Context, userID int32) (*model.GameAccount, error) // 若普通用户仅1条，可直接用
	ListByUser(ctx context.Context, userID int32) ([]*model.GameAccount, error)
	GetByID(ctx context.Context, id int32) (*model.GameAccount, error) // 通过ID查询游戏账号
	GetByIDForUser(ctx context.Context, id int32, userID int32) (*model.GameAccount, error)
	DeleteByUser(ctx context.Context, userID int32) error
	CountByCtrl(ctx context.Context, ctrlID int32) (int64, error)
	ListByGameUserIDEmpty(ctx context.Context) ([]*model.GameAccount, error) // 查询 game_user_id 为空的记录
	Update(ctx context.Context, a *model.GameAccount) error                  // 更新记录
	GetByAccount(ctx context.Context, account string) (*model.GameAccount, error)
	GetByGamePlayerID(ctx context.Context, gamePlayerID string) (*model.GameAccount, error) // 通过game_player_id查询
}

type gameAccountRepo struct {
	repo.CORMImpl[model.GameAccount]
	data *infra.Data
	log  *log.Helper
}

func NewGameAccountRepo(data *infra.Data, logger log.Logger) GameAccountRepo {
	return &gameAccountRepo{
		CORMImpl: repo.NewCORMImplRepo[model.GameAccount](data),
		data:     data,
		log:      log.NewHelper(log.With(logger, "module", "repo/gameAccount")),
	}
}

func (r *gameAccountRepo) Create(ctx context.Context, a *model.GameAccount) error {
	a.PwdMD5 = strings.ToUpper(strings.TrimSpace(a.PwdMD5))
	return r.data.GetDBWithContext(ctx).Create(a).Error
}

func (r *gameAccountRepo) GetOneByUser(ctx context.Context, userID int32) (*model.GameAccount, error) {
	var a model.GameAccount
	err := r.data.GetDBWithContext(ctx).
		Where("user_id = ? ", userID).
		Order("id DESC").
		First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *gameAccountRepo) ListByUser(ctx context.Context, userID int32) ([]*model.GameAccount, error) {
	var out []*model.GameAccount
	err := r.data.GetDBWithContext(ctx).
		Where("user_id = ? ", userID).
		Order("id DESC").
		Find(&out).Error
	return out, err
}

func (r *gameAccountRepo) GetByID(ctx context.Context, id int32) (*model.GameAccount, error) {
	var a model.GameAccount
	err := r.data.GetDBWithContext(ctx).
		Where("id = ?", id).
		First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *gameAccountRepo) GetByIDForUser(ctx context.Context, id int32, userID int32) (*model.GameAccount, error) {
	var a model.GameAccount
	err := r.data.GetDBWithContext(ctx).
		Where("id = ? AND user_id = ? ", id, userID).
		First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *gameAccountRepo) DeleteByUser(ctx context.Context, userID int32) error {
	return r.data.GetDBWithContext(ctx).
		Where("user_id = ? ", userID).
		Delete(&model.GameAccount{}).Error
}

func (r *gameAccountRepo) CountByCtrl(ctx context.Context, ctrlID int32) (int64, error) {
	var cnt int64
	err := r.data.GetDBWithContext(ctx).
		Model(&model.GameAccount{}).
		Where("ctrl_account_id = ?", ctrlID).
		Count(&cnt).Error
	return cnt, err
}

func (r *gameAccountRepo) ListByGameUserIDEmpty(ctx context.Context) ([]*model.GameAccount, error) {
	var out []*model.GameAccount
	err := r.data.GetDBWithContext(ctx).
		Where("(game_user_id = '' OR game_user_id IS NULL) AND is_del = 0").
		Order("id ASC").
		Find(&out).Error
	return out, err
}

func (r *gameAccountRepo) Update(ctx context.Context, a *model.GameAccount) error {
	return r.data.GetDBWithContext(ctx).
		Model(a).
		Updates(a).Error
}

func (r *gameAccountRepo) GetByAccount(ctx context.Context, account string) (*model.GameAccount, error) {
	var a model.GameAccount
	err := r.data.GetDBWithContext(ctx).Where("account = ? AND is_del = 0", account).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *gameAccountRepo) GetByGamePlayerID(ctx context.Context, gamePlayerID string) (*model.GameAccount, error) {
	var a model.GameAccount
	err := r.data.GetDBWithContext(ctx).Where("game_player_id = ? AND is_del = 0", gamePlayerID).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}
