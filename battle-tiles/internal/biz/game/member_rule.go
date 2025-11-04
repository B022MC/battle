package game

import (
	model "battle-tiles/internal/dal/model/game"
	repo "battle-tiles/internal/dal/repo/game"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type MemberRuleUseCase struct {
	repo repo.MemberRuleRepo
	log  *log.Helper
}

func NewMemberRuleUseCase(r repo.MemberRuleRepo, logger log.Logger) *MemberRuleUseCase {
	return &MemberRuleUseCase{repo: r, log: log.NewHelper(log.With(logger, "module", "usecase/member_rule"))}
}

func (uc *MemberRuleUseCase) SetVIP(ctx context.Context, op int32, houseGID, memberID int32, vip bool) error {
	return uc.repo.Upsert(ctx, &model.GameMemberRule{HouseGID: houseGID, MemberID: memberID, VIP: vip, UpdatedBy: op})
}
func (uc *MemberRuleUseCase) SetMultiGIDs(ctx context.Context, op int32, houseGID, memberID int32, allow bool) error {
	return uc.repo.Upsert(ctx, &model.GameMemberRule{HouseGID: houseGID, MemberID: memberID, MultiGIDs: allow, UpdatedBy: op})
}
func (uc *MemberRuleUseCase) SetTempRelease(ctx context.Context, op int32, houseGID, memberID int32, limit int32) error {
	return uc.repo.Upsert(ctx, &model.GameMemberRule{HouseGID: houseGID, MemberID: memberID, TempRelease: limit, UpdatedBy: op})
}
