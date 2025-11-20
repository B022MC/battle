package game

import (
	repo "battle-tiles/internal/dal/repo/game"
	gameVO "battle-tiles/internal/dal/vo/game"
	plazaHTTP "battle-tiles/internal/utils/plaza"
	"context"
	"encoding/json"
	"fmt"
	"time"

	model "battle-tiles/internal/dal/model/game"

	"github.com/go-kratos/kratos/v2/log"
)

type BattleRecordUseCase struct {
	repo         repo.BattleRecordRepo
	ctrlRepo     repo.GameCtrlAccountRepo
	linkRepo     repo.GameCtrlAccountHouseRepo
	accountRepo  repo.GameAccountRepo
	memberRepo   repo.GameMemberRepo
	settingsRepo repo.HouseSettingsRepo
	feeRepo      repo.FeeSettleRepo
	log          *log.Helper
}

func NewBattleRecordUseCase(
	r repo.BattleRecordRepo,
	ctrlRepo repo.GameCtrlAccountRepo,
	linkRepo repo.GameCtrlAccountHouseRepo,
	accountRepo repo.GameAccountRepo,
	memberRepo repo.GameMemberRepo,
	settingsRepo repo.HouseSettingsRepo,
	feeRepo repo.FeeSettleRepo,
	logger log.Logger,
) *BattleRecordUseCase {
	return &BattleRecordUseCase{
		repo:         r,
		ctrlRepo:     ctrlRepo,
		linkRepo:     linkRepo,
		accountRepo:  accountRepo,
		memberRepo:   memberRepo,
		settingsRepo: settingsRepo,
		feeRepo:      feeRepo,
		log:          log.NewHelper(log.With(logger, "module", "usecase/battle_record")),
	}
}

// PullAndSave 拉取 foxuc 战绩并入库
func (uc *BattleRecordUseCase) PullAndSave(ctx context.Context, httpc plazaHTTP.HTTPDoer, base string, houseGID, groupID, typeid int) (int, error) {
	list, err := plazaHTTP.GetGroupBattleInfoCtx(ctx, httpc, base, groupID, typeid)
	if err != nil {
		return 0, err
	}

	// 加载店铺费用配置
	settings, err := uc.settingsRepo.Get(ctx, int32(houseGID))
	if err != nil {
		uc.log.Warnf("Failed to load house settings for %d: %v, fee calculation disabled", houseGID, err)
		settings = nil
	}

	var batch []*model.GameBattleRecord
	now := time.Now()

	// 处理每局战绩
	for bIdx, b := range list {
		// 序列化玩家数据
		pbytes, err := json.Marshal(b.Players)
		if err != nil {
			uc.log.Errorf("Failed to marshal players for battle %d: %v", bIdx, err)
			pbytes = []byte("[]")
		}

		// 1. 构建玩家圈子映射并验证数据有效性
		playerGroups, playerAccounts, validPlayers := uc.buildPlayerGroupMapping(ctx, int32(houseGID), b.Players)
		if len(validPlayers) == 0 {
			uc.log.Warnf("Battle room %d has no valid players, skipping", b.RoomID)
			continue
		}

		// 2. 过滤出有效玩家列表（避免重复过滤）
		validPlayerList := make([]*gameVO.BattleSettle, 0, len(validPlayers))
		for _, p := range b.Players {
			if validPlayers[p.UserGameID] {
				validPlayerList = append(validPlayerList, p)
			}
		}

		// 3. 计算运费（只基于有效玩家）
		totalFee := int32(0)
		feeDistribution := make(map[int32]*GroupInfo)
		if settings != nil {
			totalFee = CalculateFee(settings.FeesJSON, b)
			if totalFee > 0 {
				// 只基于有效玩家计算费用分配
				feeDistribution = CalculateFeeDistribution(validPlayerList, playerGroups, totalFee, settings.ShareFee)
			}
		}

		// 4. 为有效玩家创建战绩记录
		for _, player := range b.Players {
			// 跳过无效玩家（不在系统中或没有圈子）
			if !validPlayers[player.UserGameID] {
				uc.log.Debugf("Skipping invalid player %d in room %d", player.UserGameID, b.RoomID)
				continue
			}

			userGameID := int32(player.UserGameID)
			playerGroupID := playerGroups[player.UserGameID]
			playerAccountName := playerAccounts[player.UserGameID]

			// 计算玩家应付的费用
			playerFee := int32(0)
			if groupInfo, ok := feeDistribution[playerGroupID]; ok {
				playerFee = groupInfo.PlayerFee
			}

			rec := &model.GameBattleRecord{
				HouseGID:       int32(houseGID),
				GroupID:        playerGroupID,
				RoomUID:        int32(b.RoomID),
				KindID:         int32(b.KindID),
				BaseScore:      int32(b.BaseScore),
				BattleAt:       time.Unix(int64(b.CreateTime), 0),
				PlayerGameID:   &userGameID,       // 游戏玩家ID (22953243)
				PlayerGameName: playerAccountName, // 游戏账号名称 ("1106162940")
				PlayersJSON:    string(pbytes),
				Score:          int32(player.Score),
				Fee:            playerFee,
				CreatedAt:      now,
			}
			batch = append(batch, rec)
		}

		// 5. 保存费用结算记录（失败不影响战绩保存）
		if settings != nil && totalFee > 0 && len(feeDistribution) > 0 {
			if err := uc.saveFeeSettlements(ctx, int32(houseGID), int32(b.RoomID), feeDistribution, settings.ShareFee, validPlayerList, playerGroups, totalFee, now); err != nil {
				uc.log.Errorf("Failed to save fee settlements for room %d: %v", b.RoomID, err)
				// 继续处理，不中断战绩保存
			}
		}
	}

	// 批量保存战绩记录
	if len(batch) == 0 {
		return 0, nil
	}

	if err := uc.repo.SaveBatch(ctx, batch); err != nil {
		return 0, fmt.Errorf("保存战绩失败: %w", err)
	}

	uc.log.Infof("Successfully saved %d battle records for house %d (processed %d battles)", len(batch), houseGID, len(list))
	return len(batch), nil
}

// buildPlayerGroupMapping 构建玩家到圈子的映射
// 验证玩家是否在系统中且有有效圈子
// 关键流程：游戏玩家ID（UserGameID） -> game_account.game_player_id -> game_account.id -> game_member -> group_id
// 返回：
//   - playerGroups: 游戏玩家ID到圈子ID的映射
//   - playerAccounts: 游戏玩家ID到账号名称的映射
//   - validPlayers: 有效玩家的标记（只有在系统中且有圈子的玩家才是有效的）
func (uc *BattleRecordUseCase) buildPlayerGroupMapping(
	ctx context.Context,
	houseGID int32,
	players []*gameVO.BattleSettle,
) (map[int]int32, map[int]string, map[int]bool) {
	playerGroups := make(map[int]int32, len(players))
	playerAccounts := make(map[int]string, len(players))
	validPlayers := make(map[int]bool, len(players))

	for _, player := range players {
		gamePlayerID := fmt.Sprintf("%d", player.UserGameID)

		// 第1步：通过游戏玩家ID查询 game_account
		// player.UserGameID 是Plaza API返回的GameID，对应 game_account.game_player_id
		account, err := uc.accountRepo.GetByGamePlayerID(ctx, gamePlayerID)
		if err != nil {
			uc.log.Warnf("Game account with game_player_id=%s not found: %v", gamePlayerID, err)
			continue // 游戏账号不在系统中，跳过
		}

		if account == nil {
			uc.log.Warnf("Game account with game_player_id=%s exists but record is nil", gamePlayerID)
			continue // 数据异常，跳过
		}

		// 第2步：通过 game_account.id 查询 game_member
		// game_member.game_id 对应 game_account.id（不是游戏账号ID）
		member, err := uc.memberRepo.GetByGameID(ctx, houseGID, account.Id)
		if err != nil {
			uc.log.Warnf("Member not found for account %d (game_player_id=%s) in house %d: %v",
				account.Id, gamePlayerID, houseGID, err)
			continue // 玩家不在此店铺中，跳过
		}

		if member == nil {
			uc.log.Warnf("Member record is nil for account %d", account.Id)
			continue // 数据异常，跳过
		}

		// 第3步：验证圈子
		if member.GroupID == nil || *member.GroupID == 0 {
			uc.log.Warnf("Player %s (account_id=%d) has no valid group in house %d",
				gamePlayerID, account.Id, houseGID)
			continue // 没有圈子，费用无法计算，跳过
		}

		// 玩家有效：游戏玩家ID -> 游戏账号 -> 成员信息 -> 圈子
		playerGroups[player.UserGameID] = *member.GroupID
		playerAccounts[player.UserGameID] = account.Account // 保存账号名称
		validPlayers[player.UserGameID] = true
	}

	uc.log.Debugf("Mapped %d valid players out of %d total players", len(validPlayers), len(players))
	return playerGroups, playerAccounts, validPlayers
}

// saveFeeSettlements 保存费用结算记录
func (uc *BattleRecordUseCase) saveFeeSettlements(
	ctx context.Context,
	houseGID int32,
	roomID int32,
	feeDistribution map[int32]*GroupInfo,
	shareFee bool,
	players []*gameVO.BattleSettle,
	playerGroups map[int]int32,
	totalFee int32,
	feedAt time.Time,
) error {
	// 1. 保存每个圈子的费用记录
	for groupID, info := range feeDistribution {
		if info.TotalFee > 0 {
			groupName := fmt.Sprintf("group_%d", groupID)
			feeRecord := &model.GameFeeSettle{
				HouseGID:  houseGID,
				PlayGroup: groupName,
				Amount:    info.TotalFee,
				FeedAt:    feedAt,
				CreatedAt: time.Now(),
			}
			if err := uc.feeRepo.Insert(ctx, feeRecord); err != nil {
				uc.log.Errorf("Failed to save fee record for group %d: %v", groupID, err)
				// 继续处理其他圈子
			}
		}
	}

	// 2. 如果是分运费模式，保存结转记录
	if shareFee {
		settlements := CalculateFeeSettlements(players, playerGroups, totalFee, shareFee)
		for _, s := range settlements {
			if s.Amount != 0 && s.IsPayoff {
				groupName := fmt.Sprintf("group_%d", s.GroupID)
				settleRecord := &model.GameFeeSettle{
					HouseGID:  houseGID,
					PlayGroup: groupName,
					Amount:    s.Amount, // 正数=支出，负数=收入
					FeedAt:    feedAt,
					CreatedAt: time.Now(),
				}
				if err := uc.feeRepo.Insert(ctx, settleRecord); err != nil {
					uc.log.Errorf("Failed to save fee settlement for group %d: %v", s.GroupID, err)
				}
			}
		}
	}

	return nil
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
