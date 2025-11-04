package game

import (
	repo "battle-tiles/internal/dal/repo/game"
	resp "battle-tiles/internal/dal/resp"
	plazaHTTP "battle-tiles/internal/utils/plaza"
	"context"
	"encoding/json"
	"time"

	model "battle-tiles/internal/dal/model/game"

	"github.com/go-kratos/kratos/v2/log"
)

type BattleRecordUseCase struct {
	repo repo.BattleRecordRepo
	log  *log.Helper
}

func NewBattleRecordUseCase(r repo.BattleRecordRepo, logger log.Logger) *BattleRecordUseCase {
	return &BattleRecordUseCase{repo: r, log: log.NewHelper(log.With(logger, "module", "usecase/battle_record"))}
}

// PullAndSave 拉取 foxuc 战绩并入库
func (uc *BattleRecordUseCase) PullAndSave(ctx context.Context, httpc plazaHTTP.HTTPDoer, base string, houseGID, groupID, typeid int) (int, error) {
	list, err := plazaHTTP.GetGroupBattleInfoCtx(ctx, httpc, base, groupID, typeid)
	if err != nil {
		return 0, err
	}
	var batch []*model.GameBattleRecord
	now := time.Now()
	for _, b := range list {
		pbytes, _ := json.Marshal(b.Players)
		rec := &model.GameBattleRecord{
			HouseGID:    int32(houseGID),
			GroupID:     int32(groupID),
			RoomUID:     int32(b.RoomID),
			KindID:      int32(b.KindID),
			BaseScore:   int32(b.BaseScore),
			BattleAt:    time.Unix(int64(b.CreateTime), 0),
			PlayersJSON: string(pbytes),
			CreatedAt:   now,
		}
		batch = append(batch, rec)
	}
	if err := uc.repo.SaveBatch(ctx, batch); err != nil {
		return 0, err
	}
	return len(batch), nil
}

// List 本地战绩查询
func (uc *BattleRecordUseCase) List(ctx context.Context, houseGID int32, groupID, gameID *int32, start, end *time.Time, page, size int32) ([]resp.BattleRecordVO, int64, error) {
	list, total, err := uc.repo.List(ctx, houseGID, groupID, gameID, start, end, page, size)
	if err != nil {
		return nil, 0, err
	}
	out := make([]resp.BattleRecordVO, 0, len(list))
	for _, r := range list {
		// 简化：不展开 players_json，前端可透传或服务端可再解析
		out = append(out, resp.BattleRecordVO{RoomID: int(r.RoomUID), KindID: int(r.KindID), BaseScore: int(r.BaseScore), Time: int(r.BattleAt.Unix())})
	}
	return out, total, nil
}
