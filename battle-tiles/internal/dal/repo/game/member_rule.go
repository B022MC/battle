package game

import (
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type MemberRuleRepo interface {
	Upsert(ctx context.Context, in *model.GameMemberRule) error
	Get(ctx context.Context, houseGID, memberID int32) (*model.GameMemberRule, error)
}

type memberRuleRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewMemberRuleRepo(data *infra.Data, logger log.Logger) MemberRuleRepo {
	return &memberRuleRepo{data: data, log: log.NewHelper(log.With(logger, "module", "repo/member_rule"))}
}

func (r *memberRuleRepo) db(ctx context.Context) *gorm.DB { return r.data.GetDBWithContext(ctx) }

func (r *memberRuleRepo) Upsert(ctx context.Context, in *model.GameMemberRule) error {
	return r.db(ctx).Clauses().Save(in).Error
}

func (r *memberRuleRepo) Get(ctx context.Context, houseGID, memberID int32) (*model.GameMemberRule, error) {
	var m model.GameMemberRule
	if err := r.db(ctx).Where("house_gid = ? AND member_id = ?", houseGID, memberID).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
