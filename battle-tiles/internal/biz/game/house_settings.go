package game

import (
	model "battle-tiles/internal/dal/model/game"
	repo "battle-tiles/internal/dal/repo/game"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type HouseSettingsUseCase struct {
	repo    repo.HouseSettingsRepo
	feeRepo repo.FeeSettleRepo
	log     *log.Helper
}

func NewHouseSettingsUseCase(r repo.HouseSettingsRepo, feeRepo repo.FeeSettleRepo, logger log.Logger) *HouseSettingsUseCase {
	return &HouseSettingsUseCase{repo: r, feeRepo: feeRepo, log: log.NewHelper(log.With(logger, "module", "usecase/house_settings"))}
}

func (uc *HouseSettingsUseCase) Get(ctx context.Context, houseGID int32) (*model.GameHouseSettings, error) {
	return uc.repo.Get(ctx, houseGID)
}

func (uc *HouseSettingsUseCase) SetFees(ctx context.Context, opUser int32, houseGID int32, feesJSON string) error {
	return uc.repo.Upsert(ctx, &model.GameHouseSettings{HouseGID: houseGID, FeesJSON: feesJSON, UpdatedBy: opUser})
}

func (uc *HouseSettingsUseCase) SetShareFee(ctx context.Context, opUser int32, houseGID int32, share bool) error {
	return uc.repo.Upsert(ctx, &model.GameHouseSettings{HouseGID: houseGID, ShareFee: share, UpdatedBy: opUser})
}

func (uc *HouseSettingsUseCase) SetPushCredit(ctx context.Context, opUser int32, houseGID int32, credit int32) error {
	return uc.repo.Upsert(ctx, &model.GameHouseSettings{HouseGID: houseGID, PushCredit: credit, UpdatedBy: opUser})
}

// ---- Fee Settle ----
func (uc *HouseSettingsUseCase) InsertFeeSettle(ctx context.Context, houseGID int32, group string, amount int32, feedAt time.Time) error {
	return uc.feeRepo.Insert(ctx, &model.GameFeeSettle{HouseGID: houseGID, PlayGroup: group, Amount: amount, FeedAt: feedAt})
}

func (uc *HouseSettingsUseCase) SumFeeSettle(ctx context.Context, houseGID int32, group string, start, end time.Time) (int64, error) {
	return uc.feeRepo.Sum(ctx, houseGID, group, start, end)
}

func (uc *HouseSettingsUseCase) ListGroupSums(ctx context.Context, houseGID int32, start, end time.Time) ([]repo.GroupSum, error) {
	return uc.feeRepo.ListGroupSums(ctx, houseGID, start, end)
}
