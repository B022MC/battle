package game

import (
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/dal/vo/game"
	"battle-tiles/internal/infra"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GameCtrlAccountRepo interface {
	Get(ctx context.Context, id int32) (*model.GameCtrlAccount, error)
	GetByIdentifier(ctx context.Context, loginMode int32, identifier string) (*model.GameCtrlAccount, error)
	Upsert(ctx context.Context, in *model.GameCtrlAccount) (*model.GameCtrlAccount, error) // by (login_mode, identifier)
	Update(ctx context.Context, in *model.GameCtrlAccount) error
	UpdateGameUserID(ctx context.Context, id int32, gameUserID int32) error
	UpdateStatus(ctx context.Context, id int32, status int32) error // 更新状态
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, cond game.CtrlAccountListCond) ([]*model.GameCtrlAccount, int64, error)
}

type gameCtrlAccountRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewCtrlAccountRepo(data *infra.Data, logger log.Logger) GameCtrlAccountRepo {
	return &gameCtrlAccountRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/game_ctrl_account")),
	}
}

func (r *gameCtrlAccountRepo) db(ctx context.Context) *gorm.DB { return r.data.GetDBWithContext(ctx) }

func (r *gameCtrlAccountRepo) Get(ctx context.Context, id int32) (*model.GameCtrlAccount, error) {
	var m model.GameCtrlAccount
	if err := r.db(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *gameCtrlAccountRepo) GetByIdentifier(ctx context.Context, loginMode int32, identifier string) (*model.GameCtrlAccount, error) {
	var m model.GameCtrlAccount
	if err := r.db(ctx).Where("login_mode = ? AND identifier = ?", loginMode, identifier).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *gameCtrlAccountRepo) Upsert(ctx context.Context, in *model.GameCtrlAccount) (*model.GameCtrlAccount, error) {
	db := r.db(ctx)
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "login_mode"}, {Name: "identifier"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"pwd_md5", "status", "last_verify_at", "updated_at",
			"game_player_id", "game_id", // 改为 game_player_id
		}),
	}).Create(in).Error
	if err != nil {
		return nil, err
	}
	// 取回（兼容插入/更新）
	return r.GetByIdentifier(ctx, in.LoginMode, in.Identifier)
}

func (r *gameCtrlAccountRepo) Update(ctx context.Context, in *model.GameCtrlAccount) error {
	return r.db(ctx).Save(in).Error
}

func (r *gameCtrlAccountRepo) UpdateGameUserID(ctx context.Context, id int32, gameUserID int32) error {
	return r.db(ctx).Model(&model.GameCtrlAccount{}).Where("id = ?", id).Update("game_player_id", gameUserID).Error
}

func (r *gameCtrlAccountRepo) UpdateStatus(ctx context.Context, id int32, status int32) error {
	return r.db(ctx).Model(&model.GameCtrlAccount{}).Where("id = ?", id).Update("status", status).Error
}

func (r *gameCtrlAccountRepo) Delete(ctx context.Context, id int32) error {
	return r.db(ctx).Delete(&model.GameCtrlAccount{}, id).Error
}

func (r *gameCtrlAccountRepo) List(ctx context.Context, cond game.CtrlAccountListCond) ([]*model.GameCtrlAccount, int64, error) {
	db := r.db(ctx).Model(&model.GameCtrlAccount{})

	if cond.LoginMode != nil {
		db = db.Where("login_mode = ?", *cond.LoginMode)
	}
	if cond.Status != nil {
		db = db.Where("status = ?", *cond.Status)
	}
	if cond.Keyword != "" {
		k := "%" + cond.Keyword + "%"
		db = db.Where("identifier ILIKE ?", k)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page, size := cond.Page, cond.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 200 {
		size = 20
	}
	offset := (page - 1) * size

	var list []*model.GameCtrlAccount
	if err := db.Order("id DESC").Limit(size).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
