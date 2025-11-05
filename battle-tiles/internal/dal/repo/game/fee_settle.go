package game

import (
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type FeeSettleRepo interface {
	Insert(ctx context.Context, in *model.GameFeeSettle) error
	Sum(ctx context.Context, houseGID int32, group string, start, end time.Time) (int64, error)
	ListGroupSums(ctx context.Context, houseGID int32, start, end time.Time) ([]GroupSum, error)
}

type feeSettleRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewFeeSettleRepo(data *infra.Data, logger log.Logger) FeeSettleRepo {
	return &feeSettleRepo{data: data, log: log.NewHelper(log.With(logger, "module", "repo/fee_settle"))}
}

func (r *feeSettleRepo) db(ctx context.Context) *gorm.DB { return r.data.GetDBWithContext(ctx) }

func (r *feeSettleRepo) Insert(ctx context.Context, in *model.GameFeeSettle) error {
	return r.db(ctx).Create(in).Error
}

func (r *feeSettleRepo) Sum(ctx context.Context, houseGID int32, group string, start, end time.Time) (int64, error) {
	var amt int64
	err := r.db(ctx).Model(&model.GameFeeSettle{}).
		Where("house_gid=? AND play_group=? AND feed_at>=? AND feed_at<?", houseGID, group, start, end).
		Select("COALESCE(SUM(amount),0)").
		Scan(&amt).Error
	return amt, err
}

type GroupSum struct {
	PlayGroup string `gorm:"column:play_group" json:"play_group"`
	Sum       int64  `gorm:"column:sum" json:"sum"`
}

func (r *feeSettleRepo) ListGroupSums(ctx context.Context, houseGID int32, start, end time.Time) ([]GroupSum, error) {
	var rows []GroupSum
	err := r.db(ctx).Model(&model.GameFeeSettle{}).
		Where("house_gid=? AND feed_at>=? AND feed_at<?", houseGID, start, end).
		Select("play_group, COALESCE(SUM(amount),0) as sum").
		Group("play_group").
		Scan(&rows).Error
	return rows, err
}
