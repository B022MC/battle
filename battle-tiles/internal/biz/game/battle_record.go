package game

import (
	repo "battle-tiles/internal/dal/repo/game"
	resp "battle-tiles/internal/dal/resp"
	plazaHTTP "battle-tiles/internal/utils/plaza"
	"context"
	"encoding/json"
	"fmt"
	"time"

	model "battle-tiles/internal/dal/model/game"

	"github.com/go-kratos/kratos/v2/log"
)

type BattleRecordUseCase struct {
	repo        repo.BattleRecordRepo
	ctrlRepo    repo.GameCtrlAccountRepo
	linkRepo    repo.GameCtrlAccountHouseRepo
	accountRepo repo.GameAccountRepo
	memberRepo  repo.GameMemberRepo
	log         *log.Helper
}

func NewBattleRecordUseCase(
	r repo.BattleRecordRepo,
	ctrlRepo repo.GameCtrlAccountRepo,
	linkRepo repo.GameCtrlAccountHouseRepo,
	accountRepo repo.GameAccountRepo,
	memberRepo repo.GameMemberRepo,
	logger log.Logger,
) *BattleRecordUseCase {
	return &BattleRecordUseCase{
		repo:        r,
		ctrlRepo:    ctrlRepo,
		linkRepo:    linkRepo,
		accountRepo: accountRepo,
		memberRepo:  memberRepo,
		log:         log.NewHelper(log.With(logger, "module", "usecase/battle_record")),
	}
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

		// 为每个玩家创建战绩记录，并从 game_member 表中查询玩家的当前圈子
		for _, player := range b.Players {
			playerGroupID := int32(groupID) // 默认使用参数中的 groupID

			// 从 game_member 表中查询玩家的当前圈子
			member, err := uc.memberRepo.GetByGameID(ctx, int32(houseGID), int32(player.UserGameID))
			if err == nil && member != nil && member.GroupID != nil {
				playerGroupID = *member.GroupID
			} else if err != nil {
				uc.log.Warnf("Failed to get member info for gameID %d: %v, using default groupID %d", player.UserGameID, err, groupID)
			}
			userGameID := int32(player.UserGameID)
			rec := &model.GameBattleRecord{
				HouseGID:     int32(houseGID),
				GroupID:      playerGroupID, // 使用玩家当前所在的圈子
				RoomUID:      int32(b.RoomID),
				KindID:       int32(b.KindID),
				BaseScore:    int32(b.BaseScore),
				BattleAt:     time.Unix(int64(b.CreateTime), 0),
				PlayerGameID: &userGameID,
				PlayersJSON:  string(pbytes),
				CreatedAt:    now,
			}
			batch = append(batch, rec)
		}
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

// ListMyBattleRecords 用户查看自己的战绩（通过绑定的游戏账号）
func (uc *BattleRecordUseCase) ListMyBattleRecords(
	ctx context.Context,
	userID int32,
	houseGID int32,
	GroupID *int32,
	start, end *time.Time,
	page, size int32,
) ([]*model.GameBattleRecord, int64, error) {
	account, err := uc.accountRepo.GetOneByUser(ctx, userID)
	if err != nil {
		uc.log.Errorf("Failed to get game account for user %d: %v", userID, err)
		return nil, 0, fmt.Errorf("未找到绑定的游戏账号")
	}

	// 检查是否有游戏账号
	if account.Account == "" {
		uc.log.Warnf("User %d has no game account", userID)
		return []*model.GameBattleRecord{}, 0, nil
	}

	// 查询战绩（使用游戏账号 account.Account）
	return uc.repo.ListByPlayerGameName(ctx, houseGID, account.Account, GroupID, start, end, page, size)
}

// GetMyBattleStats 获取用户的战绩统计
func (uc *BattleRecordUseCase) GetMyBattleStats(
	ctx context.Context,
	userID int32,
	houseGID int32,
	groupID *int32,
	start, end *time.Time,
) (totalGames int64, totalScore int, totalFee int, err error) {
	// 1. 查询用户绑定的游戏账号
	account, err := uc.accountRepo.GetOneByUser(ctx, userID)
	if err != nil {
		uc.log.Errorf("Failed to get game account for user %d: %v", userID, err)
		return 0, 0, 0, fmt.Errorf("未找到绑定的游戏账号")
	}

	// 检查是否有游戏账号
	if account.Account == "" {
		uc.log.Warnf("User %d has no game account", userID)
		return 0, 0, 0, nil
	}

	// 查询统计（使用游戏账号 account.Account）
	return uc.repo.GetPlayerStatsByGameName(ctx, houseGID, account.Account, groupID, start, end)
}

// ListHouseBattleRecords 管理员查看店铺战绩
func (uc *BattleRecordUseCase) ListHouseBattleRecords(
	ctx context.Context,
	houseGID int32,
	groupID *int32,
	gameID *int32,
	start, end *time.Time,
	page, size int32,
) ([]*model.GameBattleRecord, int64, error) {
	// TODO: 添加 groupID 和 gameID 过滤
	// 目前 repo 层的 ListByPlayer 方法不支持这些过滤条件
	// 可以在后续扩展 repo 层方法来支持

	// 暂时返回所有店铺的战绩
	// 需要在 repo 层添加 ListByHouse 方法
	return nil, 0, fmt.Errorf("ListHouseBattleRecords not implemented yet")
}

// GetPlayerBattleStats 管理员查看玩家战绩统计
func (uc *BattleRecordUseCase) GetPlayerBattleStats(
	ctx context.Context,
	houseGID int32,
	playerGameID int32,
	start, end *time.Time,
) (totalGames int64, totalScore int, totalFee int, err error) {
	return uc.repo.GetPlayerStats(ctx, houseGID, string(playerGameID), nil, start, end)
}

// GetMyBalances 获取用户的余额
func (uc *BattleRecordUseCase) GetMyBalances(
	ctx context.Context,
	userID int32,
	houseGID int32,
	groupID *int32,
) (interface{}, error) {
	// 1. 查询用户绑定的游戏账号
	account, err := uc.accountRepo.GetOneByUser(ctx, userID)
	if err != nil {
		uc.log.Errorf("Failed to get game account for user %d: %v", userID, err)
		return nil, fmt.Errorf("未找到绑定的游戏账号")
	}

	// 2. 检查是否有游戏账号
	if account.Account == "" {
		uc.log.Warnf("User %d has no game account", userID)
		return []interface{}{}, nil
	}

	// 3. 返回空列表（暂时实现）
	// TODO: 实现实际的余额查询逻辑
	return []interface{}{}, nil
}

// ListGroupMemberBalances 查询圈子成员余额
func (uc *BattleRecordUseCase) ListGroupMemberBalances(
	ctx context.Context,
	houseGID int32,
	groupID int32,
	minYuan *int32,
	maxYuan *int32,
	page, size int32,
) (interface{}, int64, error) {
	// 暂时返回空列表
	return []interface{}{}, 0, nil
}

// GetGroupStats 查询圈子统计
func (uc *BattleRecordUseCase) GetGroupStats(
	ctx context.Context,
	houseGID int32,
	groupID int32,
	start, end *time.Time,
) (interface{}, error) {
	totalGames, totalScore, totalFee, activeMembers, err := uc.repo.GetGroupStats(ctx, houseGID, groupID, start, end)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"group_id":       groupID,
		"total_games":    totalGames,
		"total_score":    totalScore,
		"total_fee":      totalFee,
		"active_members": activeMembers,
	}, nil
}

// GetHouseStats 查询店铺统计
func (uc *BattleRecordUseCase) GetHouseStats(
	ctx context.Context,
	houseGID int32,
	start, end *time.Time,
) (interface{}, error) {
	totalGames, totalScore, totalFee, err := uc.repo.GetHouseStats(ctx, houseGID, start, end)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"house_gid":   houseGID,
		"total_games": totalGames,
		"total_score": totalScore,
		"total_fee":   totalFee,
	}, nil
}

// parseGameUserID 解析 game_user_id 字符串为整数
func parseGameUserID(gameUserIDStr string, out *int32) (bool, error) {
	// game_user_id 已经是字符串形式的数字，直接转换
	var id int
	if n, err := fmt.Sscanf(gameUserIDStr, "%d", &id); err != nil || n != 1 {
		return false, fmt.Errorf("invalid game_user_id format: %s", gameUserIDStr)
	}
	*out = int32(id)
	return true, nil
}
