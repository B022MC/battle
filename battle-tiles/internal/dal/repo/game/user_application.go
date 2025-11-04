package game

import (
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type UserApplicationRepo interface {
	Insert(ctx context.Context, in *model.UserApplication) error
}

type userApplicationRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewUserApplicationRepo(data *infra.Data, logger log.Logger) UserApplicationRepo {
	return &userApplicationRepo{data: data, log: log.NewHelper(log.With(logger, "module", "repo/user_application"))}
}

func (r *userApplicationRepo) db(ctx context.Context) *gorm.DB { return r.data.GetDBWithContext(ctx) }

func (r *userApplicationRepo) Insert(ctx context.Context, in *model.UserApplication) error {
	return r.db(ctx).Create(in).Error
}
