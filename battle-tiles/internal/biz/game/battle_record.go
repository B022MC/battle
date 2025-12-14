package game

import (
	"battle-tiles/internal/conf"
	repo "battle-tiles/internal/dal/repo/game"
	gameVO "battle-tiles/internal/dal/vo/game"
	"battle-tiles/internal/infra"
	plazaHTTP "battle-tiles/internal/utils/plaza"
	"context"
	"encoding/json"
	"fmt"
	"time"

	model "battle-tiles/internal/dal/model/game"

	"github.com/go-kratos/kratos/v2/log"
)

type BattleRecordUseCase struct {
	repo             repo.BattleRecordRepo
	ctrlRepo         repo.GameCtrlAccountRepo
	linkRepo         repo.GameCtrlAccountHouseRepo
	accountRepo      repo.GameAccountRepo
	memberRepo       repo.GameMemberRepo
	accountGroupRepo repo.GameAccountGroupRepo // 用于查询活跃圈子
	settingsRepo     repo.HouseSettingsRepo
	feeRepo          repo.FeeSettleRepo
	walletRepo       repo.WalletReadRepo
	walletWriteRepo  repo.WalletRepo // 用于更新余额和记录流水
	log              *log.Helper
	verboseTaskLog   bool // 是否显示详细的任务日志
}

func NewBattleRecordUseCase(
	r repo.BattleRecordRepo,
	ctrlRepo repo.GameCtrlAccountRepo,
	linkRepo repo.GameCtrlAccountHouseRepo,
	accountRepo repo.GameAccountRepo,
	memberRepo repo.GameMemberRepo,
	accountGroupRepo repo.GameAccountGroupRepo,
	settingsRepo repo.HouseSettingsRepo,
	feeRepo repo.FeeSettleRepo,
	walletRepo repo.WalletReadRepo,
	walletWriteRepo repo.WalletRepo,
	logConf *conf.Log,
	logger log.Logger,
) *BattleRecordUseCase {
	verboseTaskLog := false
	if logConf != nil && logConf.VerboseTaskLog {
		verboseTaskLog = true
	}
	return &BattleRecordUseCase{
		repo:             r,
		ctrlRepo:         ctrlRepo,
		linkRepo:         linkRepo,
		accountRepo:      accountRepo,
		memberRepo:       memberRepo,
		accountGroupRepo: accountGroupRepo,
		settingsRepo:     settingsRepo,
		feeRepo:          feeRepo,
		walletRepo:       walletRepo,
		walletWriteRepo:  walletWriteRepo,
		log:              log.NewHelper(log.With(logger, "module", "usecase/battle_record")),
		verboseTaskLog:   verboseTaskLog,
	}
}

// PullAndSave 拉取 foxuc 战绩并入库
func (uc *BattleRecordUseCase) PullAndSave(ctx context.Context, httpc plazaHTTP.HTTPDoer, base string, houseGID, groupID, typeid int) (int, error) {
	// 如果不显示详细日志，则使用静默模式
	if !uc.verboseTaskLog {
		ctx = infra.WithQuietDB(ctx)
	}

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
		// 1. 构建玩家圈子映射并验证数据有效性
		playerGroups, playerAccounts, validPlayers := uc.buildPlayerGroupMapping(ctx, int32(houseGID), b.Players)

		// 填充玩家昵称到 players 数据中
		for _, p := range b.Players {
			if name, ok := playerAccounts[p.UserGameID]; ok {
				p.NickName = name
			}
		}

		// 序列化玩家数据（包含昵称）
		pbytes, err := json.Marshal(b.Players)
		if err != nil {
			uc.log.Errorf("Failed to marshal players for battle %d: %v", bIdx, err)
			pbytes = []byte("[]")
		}
		if len(validPlayers) == 0 {
			if uc.verboseTaskLog {
				uc.log.Warnf("Battle room %d has no valid players, skipping", b.RoomID)
			}
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
			// 调试日志：显示费用计算结果
			uc.log.Infof("[费用计算] room=%d kind=%d base=%d maxScore=%d totalFee=%d validPlayers=%d",
				b.RoomID, b.KindID, b.BaseScore, getMaxScore(b.Players), totalFee, len(validPlayerList))
			if totalFee > 0 {
				// 只基于有效玩家计算费用分配
				feeDistribution = CalculateFeeDistribution(validPlayerList, playerGroups, totalFee, settings.ShareFee)
				// 调试日志：显示费用分配结果
				for gid, info := range feeDistribution {
					uc.log.Infof("[费用分配] room=%d group=%d isWinner=%v totalFee=%d playerFee=%d players=%v",
						b.RoomID, gid, info.IsWinner, info.TotalFee, info.PlayerFee, info.PlayerIDs)
				}
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
				PlayerBalance:  player.Balance, // 玩家当前余额
				PlayerCredit:   player.Credit,  // 玩家额度
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

	// 批量保存战绩记录（自动去重）
	if len(batch) == 0 {
		return 0, nil
	}

	savedRecords, err := uc.repo.SaveBatchWithDedup(ctx, batch)
	if err != nil {
		return 0, fmt.Errorf("保存战绩失败: %w", err)
	}

	// 根据新保存的战绩更新玩家余额
	if len(savedRecords) > 0 && uc.walletWriteRepo != nil {
		uc.updatePlayerBalances(ctx, int32(houseGID), savedRecords)
	}

	if len(savedRecords) > 0 {
		uc.log.Infof("Successfully saved %d new battle records for house %d (processed %d battles, %d duplicates skipped)",
			len(savedRecords), houseGID, len(list), len(batch)-len(savedRecords))
	}
	return len(savedRecords), nil
}

// updatePlayerBalances 根据战绩更新玩家余额并记录流水
// 只更新已存在且有余额（或曾经上过分）的成员，不自动创建新成员
// 余额计算公式（与 passing-dragonfly 对齐）：delta = score * factor - fee
func (uc *BattleRecordUseCase) updatePlayerBalances(ctx context.Context, houseGID int32, records []*model.GameBattleRecord) {
	const (
		LedgerTypeBattleSettle int32 = 5 // 战绩结算类型
	)

	for _, record := range records {
		if record.PlayerGameID == nil {
			continue // 跳过无效记录
		}

		gameID := *record.PlayerGameID

		// 计算余额变化：score * factor - fee
		// 与 passing-dragonfly 对齐：delta = int(float64(record.Score*100)*record.Factor) - record.Fee
		// 注意：battle-tiles 中分数和余额都是分，不需要乘以100
		delta := int32(float64(record.Score)*record.Factor) - record.Fee

		// 如果变化为0，跳过
		if delta == 0 {
			continue
		}

		// 查询成员信息（必须已存在）
		member, err := uc.memberRepo.GetByGameID(ctx, houseGID, gameID)
		if err != nil || member == nil {
			// 成员不存在，跳过（不自动创建）
			if uc.verboseTaskLog {
				uc.log.Debugf("Member not found for game_id=%d, skip balance update", gameID)
			}
			continue
		}

		// 检查成员是否有余额记录（只有上过分的玩家才同步战绩余额）
		// 判断条件：balance != 0 或者有过流水记录
		if member.Balance == 0 {
			// 检查是否有过流水记录（说明曾经上过分）
			hasLedger, _ := uc.walletWriteRepo.ExistsLedgerByMember(ctx, houseGID, gameID)
			if !hasLedger {
				// 没有余额也没有流水记录，说明从未上过分，跳过
				if uc.verboseTaskLog {
					uc.log.Debugf("Member game_id=%d has no balance and no ledger history, skip", gameID)
				}
				continue
			}
		}

		// 更新余额
		before, after, err := uc.updateMemberBalance(ctx, houseGID, member.Id, gameID, delta)
		if err != nil {
			uc.log.Errorf("Failed to update balance for game_id=%d: %v", gameID, err)
			continue
		}

		// 记录流水（member_id 使用 game_id，与上分接口保持一致）
		bizNo := fmt.Sprintf("battle-%d-%d-%d", houseGID, record.RoomUID, gameID)
		reason := fmt.Sprintf("战绩结算 房间:%d 分数:%d 运费:%d", record.RoomUID, record.Score, record.Fee)

		ledger := &model.GameWalletLedger{
			HouseGID:       houseGID,
			MemberID:       gameID, // 使用 game_id，与上分接口保持一致
			ChangeAmount:   delta,
			BalanceBefore:  before,
			BalanceAfter:   after,
			Type:           LedgerTypeBattleSettle,
			Reason:         reason,
			OperatorUserID: 0, // 系统自动结算
			BizNo:          bizNo,
		}

		if err := uc.walletWriteRepo.AppendLedger(ctx, nil, ledger); err != nil {
			uc.log.Errorf("Failed to append ledger for game_id=%d: %v", gameID, err)
		}

		if uc.verboseTaskLog {
			uc.log.Infof("Updated balance for game_id=%d: %d -> %d (delta: %d, score: %d, factor: %.4f, fee: %d)",
				gameID, before, after, delta, record.Score, record.Factor, record.Fee)
		}
	}
}

// updateMemberBalance 更新成员余额（不自动创建）
func (uc *BattleRecordUseCase) updateMemberBalance(ctx context.Context, houseGID int32, memberID int32, gameID int32, amount int32) (before int32, after int32, err error) {
	// 直接使用 memberRepo.UpdateBalance，但传入的是 member.Id 而不是 game_id
	// 这样可以避免自动创建新成员
	return uc.memberRepo.UpdateBalance(ctx, houseGID, gameID, amount)
}

// CompensateUnSettledBattles 补偿未结算的战绩（定时任务调用）
// 查找有战绩记录但没有对应流水的记录，进行补偿结算
// 只处理已存在且曾经上过分的成员（有余额或有流水记录）
// 余额计算公式（与 passing-dragonfly 对齐）：delta = score * factor - fee
func (uc *BattleRecordUseCase) CompensateUnSettledBattles(ctx context.Context, houseGID int32) (int, error) {
	const LedgerTypeBattleSettle int32 = 5

	// 查询最近24小时内的战绩记录
	now := time.Now()
	start := now.Add(-24 * time.Hour)

	records, _, err := uc.repo.List(ctx, houseGID, nil, nil, &start, &now, 1, 1000)
	if err != nil {
		return 0, fmt.Errorf("查询战绩失败: %w", err)
	}

	if len(records) == 0 {
		return 0, nil
	}

	compensated := 0
	for _, record := range records {
		if record.PlayerGameID == nil {
			continue
		}

		gameID := *record.PlayerGameID

		// 计算余额变化：score * factor - fee
		delta := int32(float64(record.Score)*record.Factor) - record.Fee
		if delta == 0 {
			continue // 变化为0，跳过
		}

		// 检查是否已有对应流水（通过 bizNo 判断）
		bizNo := fmt.Sprintf("battle-%d-%d-%d", houseGID, record.RoomUID, gameID)
		exists, err := uc.walletWriteRepo.ExistsLedgerBiz(ctx, houseGID, 0, bizNo)
		if err != nil {
			uc.log.Warnf("Check ledger exists failed for bizNo=%s: %v", bizNo, err)
			continue
		}
		if exists {
			continue // 已结算，跳过
		}

		// 查询成员信息（必须已存在）
		member, err := uc.memberRepo.GetByGameID(ctx, houseGID, gameID)
		if err != nil || member == nil {
			// 成员不存在，跳过（不自动创建）
			continue
		}

		// 检查成员是否曾经上过分（有余额或有流水记录）
		if member.Balance == 0 {
			hasLedger, _ := uc.walletWriteRepo.ExistsLedgerByMember(ctx, houseGID, gameID)
			if !hasLedger {
				// 没有余额也没有流水记录，说明从未上过分，跳过
				if uc.verboseTaskLog {
					uc.log.Debugf("Compensate: Member game_id=%d has no balance and no ledger, skip", gameID)
				}
				continue
			}
		}

		// 更新余额（传入 gameID，因为 UpdateBalance 是按 game_id 查询的）
		before, after, err := uc.memberRepo.UpdateBalance(ctx, houseGID, gameID, delta)
		if err != nil {
			uc.log.Errorf("Compensate: Failed to update balance for game_id=%d: %v", gameID, err)
			continue
		}

		// 记录流水（member_id 使用 game_id，与上分接口保持一致）
		reason := fmt.Sprintf("战绩补偿结算 房间:%d 分数:%d 运费:%d", record.RoomUID, record.Score, record.Fee)
		ledger := &model.GameWalletLedger{
			HouseGID:       houseGID,
			MemberID:       gameID, // 使用 game_id，与上分接口保持一致
			ChangeAmount:   delta,
			BalanceBefore:  before,
			BalanceAfter:   after,
			Type:           LedgerTypeBattleSettle,
			Reason:         reason,
			OperatorUserID: 0,
			BizNo:          bizNo,
		}

		if err := uc.walletWriteRepo.AppendLedger(ctx, nil, ledger); err != nil {
			uc.log.Errorf("Compensate: Failed to append ledger for game_id=%d: %v", gameID, err)
			continue
		}

		compensated++
		if uc.verboseTaskLog {
			uc.log.Infof("Compensated: game_id=%d, delta=%d (score=%d, factor=%.4f, fee=%d)",
				gameID, delta, record.Score, record.Factor, record.Fee)
		}
	}

	if compensated > 0 {
		uc.log.Infof("Compensated %d unsettled battles for house %d", compensated, houseGID)
	}

	return compensated, nil
}

// PlayerMemberInfo 玩家成员信息（用于战绩保存）
type PlayerMemberInfo struct {
	GroupID int32
	Name    string
	Balance int32
	Credit  int32
}

// buildPlayerGroupMapping 构建玩家到圈子的映射
// 验证玩家是否在系统中且有有效圈子
// 关键流程：游戏玩家ID（UserGameID） -> game_account.game_player_id -> game_account.id -> game_member -> group_id
// 返回：
//   - playerGroups: 游戏玩家ID到圈子ID的映射
//   - playerAccounts: 游戏玩家ID到账号名称的映射（已废弃，使用 playerMembers）
//   - validPlayers: 有效玩家的标记（只有在系统中且有圈子的玩家才是有效的）
//   - playerMembers: 游戏玩家ID到成员信息的映射（包含余额和额度）
func (uc *BattleRecordUseCase) buildPlayerGroupMapping(
	ctx context.Context,
	houseGID int32,
	players []*gameVO.BattleSettle,
) (map[int]int32, map[int]string, map[int]bool) {
	playerGroups := make(map[int]int32, len(players))
	playerAccounts := make(map[int]string, len(players))
	validPlayers := make(map[int]bool, len(players))

	for _, player := range players {
		// 直接查询 game_member 表获取玩家的圈子信息
		// ctx 已经包含静默标记，会自动使用静默模式
		member, err := uc.memberRepo.GetByGameID(ctx, houseGID, int32(player.UserGameID))
		if err != nil {
			if uc.verboseTaskLog {
				uc.log.Debugf("Member not found for game_id=%d in house %d: %v",
					player.UserGameID, houseGID, err)
			}
			continue // 玩家不在成员表中，跳过
		}

		if member == nil {
			uc.log.Warnf("Member record is nil for game_id=%d", player.UserGameID)
			continue // 数据异常，跳过
		}

		// 验证圈子ID（必须有圈子才参与计费）
		if member.GroupID == nil || *member.GroupID == 0 {
			if uc.verboseTaskLog {
				uc.log.Debugf("Player %d has no group in house %d, skipping",
					player.UserGameID, houseGID)
			}
			continue // 没有圈子，跳过
		}

		// 使用成员的游戏名称，如果没有则使用 game_id
		accountName := member.GameName
		if accountName == "" {
			accountName = fmt.Sprintf("%d", player.UserGameID)
		}

		// 玩家有效：game_member 表中有记录且有 group_id
		playerGroups[player.UserGameID] = *member.GroupID
		playerAccounts[player.UserGameID] = accountName
		validPlayers[player.UserGameID] = true

		// 保存余额和额度到 player 对象（用于后续保存战绩）
		player.Balance = member.Balance
		player.Credit = member.Credit
	}

	if uc.verboseTaskLog {
		uc.log.Debugf("Mapped %d valid players out of %d total players", len(validPlayers), len(players))
	}
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

// ListGroupBattles 查询圈子成员战绩（管理员）
func (uc *BattleRecordUseCase) ListGroupBattles(
	ctx context.Context,
	houseGID int32,
	groupID int32,
	playerGameID *int32,
	start, end *time.Time,
	page, size int32,
) ([]*model.GameBattleRecord, int64, error) {
	// 如果指定了玩家ID，查询该玩家的战绩
	if playerGameID != nil && *playerGameID > 0 {
		return uc.repo.ListByPlayer(ctx, houseGID, *playerGameID, &groupID, start, end, page, size)
	}

	// 否则查询整个圈子的战绩
	return uc.repo.List(ctx, houseGID, &groupID, nil, start, end, page, size)
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
	if account.Account == "" {
		uc.log.Warnf("User %d has no game account", userID)
		return []interface{}{}, nil
	}

	// 2. 查询成员（按 game_id）
	member, err := uc.memberRepo.GetByGameID(ctx, houseGID, account.GameIDInt())
	if err != nil || member == nil {
		uc.log.Warnf("GetMyBalances: member not found house=%d game_id=%d err=%v", houseGID, account.GameIDInt(), err)
		return []interface{}{}, nil
	}

	// 3. 选择圈子：优先入参 group_id，否则用成员自身 group_id
	targetGroupID := groupID
	if targetGroupID == nil && member.GroupID != nil {
		targetGroupID = member.GroupID
	}

	// 4. 查询钱包：先用 game_id+group 再回退 member_id+group
	var balance int32 = member.Balance
	if w, err := uc.walletRepo.GetByGameID(ctx, houseGID, member.GameID, targetGroupID); err == nil && w != nil {
		balance = w.Balance
	} else if w2, err2 := uc.walletRepo.Get(ctx, houseGID, member.Id, targetGroupID); err2 == nil && w2 != nil {
		balance = w2.Balance
	}

	// 5. 返回
	item := &MemberBalance{
		MemberID:    member.Id,
		GameID:      member.GameID,
		GameName:    member.GameName,
		GroupID:     targetGroupID,
		GroupName:   member.GroupName,
		Balance:     balance,
		BalanceYuan: float64(balance) / 100.0,
		UpdatedAt:   member.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return []interface{}{item}, nil
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

// GetHouseStats 查询店铺完整统计（包括充值、余额、圈子间转账等）
func (uc *BattleRecordUseCase) GetHouseStats(
	ctx context.Context,
	houseGID int32,
	start, end *time.Time,
) (interface{}, error) {
	stats, err := uc.repo.GetHouseStatsDetail(ctx, houseGID, start, end)
	if err != nil {
		return nil, err
	}
	return stats, nil
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

// getMaxScore 获取玩家列表中的最高分
func getMaxScore(players []*gameVO.BattleSettle) int {
	maxScore := 0
	for _, p := range players {
		if p.Score > maxScore {
			maxScore = p.Score
		}
	}
	return maxScore
}
