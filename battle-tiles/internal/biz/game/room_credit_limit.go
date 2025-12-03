package game

import (
	"context"
	"fmt"

	model "battle-tiles/internal/dal/model/game"
	repo "battle-tiles/internal/dal/repo/game"
	"battle-tiles/internal/infra/plaza"

	"github.com/go-kratos/kratos/v2/log"
)

type RoomCreditLimitUseCase struct {
	repo         repo.RoomCreditLimitRepo
	memberRepo   repo.GameMemberRepo
	settingsRepo repo.HouseSettingsRepo
	plazaMgr     plaza.Manager
	log          *log.Helper
}

func NewRoomCreditLimitUseCase(
	r repo.RoomCreditLimitRepo,
	memberRepo repo.GameMemberRepo,
	settingsRepo repo.HouseSettingsRepo,
	plazaMgr plaza.Manager,
	logger log.Logger,
) *RoomCreditLimitUseCase {
	return &RoomCreditLimitUseCase{
		repo:         r,
		memberRepo:   memberRepo,
		settingsRepo: settingsRepo,
		plazaMgr:     plazaMgr,
		log:          log.NewHelper(log.With(logger, "module", "usecase/room_credit_limit")),
	}
}

// SetCreditLimit 设置房间额度限制
func (uc *RoomCreditLimitUseCase) SetCreditLimit(ctx context.Context, opUser int32, houseGID int32, groupName string, gameKind int32, baseScore int32, creditLimit int32) error {
	limit := &model.GameRoomCreditLimit{
		HouseGID:    houseGID,
		GroupName:   groupName,
		GameKind:    gameKind,
		BaseScore:   baseScore,
		CreditLimit: creditLimit,
		UpdatedBy:   opUser,
	}
	return uc.repo.Upsert(ctx, limit)
}

// GetCreditLimit 获取特定的房间额度限制
func (uc *RoomCreditLimitUseCase) GetCreditLimit(ctx context.Context, houseGID int32, groupName string, gameKind int32, baseScore int32) (*model.GameRoomCreditLimit, error) {
	return uc.repo.Get(ctx, houseGID, groupName, gameKind, baseScore)
}

// ListCreditLimits 获取店铺的所有额度限制
func (uc *RoomCreditLimitUseCase) ListCreditLimits(ctx context.Context, houseGID int32, groupName string) ([]*model.GameRoomCreditLimit, error) {
	if groupName != "" {
		return uc.repo.ListByGroup(ctx, houseGID, groupName)
	}
	return uc.repo.List(ctx, houseGID)
}

// DeleteCreditLimit 删除房间额度限制
func (uc *RoomCreditLimitUseCase) DeleteCreditLimit(ctx context.Context, houseGID int32, groupName string, gameKind int32, baseScore int32) error {
	return uc.repo.Delete(ctx, houseGID, groupName, gameKind, baseScore)
}

// HandlePlayerSitDown 处理玩家坐下事件（类似 passing-dragonfly 的 _handlePlayerSitDown）
// 这是核心逻辑：当玩家坐下时，检查余额是否满足额度要求和费用规则，不足则解散桌子
func (uc *RoomCreditLimitUseCase) HandlePlayerSitDown(ctx context.Context, userID int, houseGID int32, gameID int32, kindID int32, baseScore int32, tableNum int32) error {
	// 1. 获取玩家信息（余额、圈子、个人额度调整）
	member, err := uc.memberRepo.GetByGameID(ctx, houseGID, gameID)
	if err != nil {
		uc.log.Warnf("HandlePlayerSitDown: player not found in member table, houseGID=%d, gameID=%d, err=%v", houseGID, gameID, err)

		// 玩家没有录入，先解散桌子
		dismissErr := uc.dismissTable(ctx, userID, houseGID, kindID, tableNum, "player not found")
		if dismissErr != nil {
			uc.log.Errorf("Failed to dismiss table for unregistered player: %v", dismissErr)
		}

		// 自动创建成员记录（forbid=true），方便管理员在列表中看到并解禁
		newMember := &model.GameMember{
			HouseGID:  houseGID,
			GameID:    gameID,
			GameName:  fmt.Sprintf("未录入玩家-%d", gameID), // 临时名称
			GroupName: "",                              // 无圈子
			Balance:   0,
			Credit:    0,
			Forbid:    true, // 禁用状态
		}
		if createErr := uc.memberRepo.Create(ctx, newMember); createErr != nil {
			uc.log.Errorf("Failed to create member record for unregistered player gameID=%d: %v", gameID, createErr)
		} else {
			uc.log.Infof("Created forbidden member record: gameID=%d, houseGID=%d", gameID, houseGID)
		}

		// 推送禁用指令到游戏端（使用 plaza.ForbidMembers）
		forbidErr := uc.plazaMgr.ForbidMembers(userID, int(houseGID), "", []int{int(gameID)}, true)
		if forbidErr != nil {
			uc.log.Errorf("Failed to forbid unregistered player gameID=%d: %v", gameID, forbidErr)
		} else {
			uc.log.Infof("Forbid unregistered player via plaza: gameID=%d, houseGID=%d", gameID, houseGID)
		}

		return err
	}

	// 2. 检查房间额度要求
	roomCreditLimit, found := uc.repo.GetCreditLimit(ctx, houseGID, member.GroupName, kindID, baseScore)
	if !found {
		uc.log.Warnf("HandlePlayerSitDown: no credit limit found, using default, houseGID=%d, group=%s, kind=%d, base=%d",
			houseGID, member.GroupName, kindID, baseScore)
	}

	// 3. 计算有效额度要求（房间额度 + 玩家个人额度调整）
	effectiveCredit := roomCreditLimit + member.Credit

	uc.log.Infof("[费用检查] 玩家信息: gameID=%d, balance=%d, roomCredit=%d, playerCredit=%d, effectiveCredit=%d",
		gameID, member.Balance, roomCreditLimit, member.Credit, effectiveCredit)

	// 4. 检查余额是否满足额度要求
	if member.Balance < effectiveCredit {
		// 余额不足，解散桌子
		uc.log.Warnf("HandlePlayerSitDown: insufficient balance for credit limit, gameID=%d, balance=%d < credit=%d, dismissing table",
			gameID, member.Balance, effectiveCredit)

		// 先解散桌子
		dismissErr := uc.dismissTable(ctx, userID, houseGID, kindID, tableNum, "insufficient balance for credit limit")
		if dismissErr != nil {
			uc.log.Errorf("Failed to dismiss table: %v", dismissErr)
		}

		// 禁用玩家（更新数据库 forbid 字段，类似 passing-dragonfly 的 _freezeMembers）
		// TODO: 可以实现 thiefCache 机制，连续3次才禁用，现在先直接禁用
		if err := uc.memberRepo.UpdateForbid(ctx, houseGID, gameID, true); err != nil {
			uc.log.Errorf("Failed to forbid player with insufficient balance gameID=%d: %v", gameID, err)
		} else {
			uc.log.Infof("Forbid player with insufficient balance: gameID=%d, balance=%d < credit=%d",
				gameID, member.Balance, effectiveCredit)
		}

		return fmt.Errorf("insufficient balance: %d < %d", member.Balance, effectiveCredit)
	}

	// 5. 检查店铺费用规则（与 passing-dragonfly 的逻辑一致）
	settings, err := uc.settingsRepo.Get(ctx, houseGID)
	if err != nil {
		uc.log.Warnf("HandlePlayerSitDown: failed to get shop settings, houseGID=%d, err=%v", houseGID, err)
		// 如果无法获取设置，继续执行（不因为配置错误而影响游戏）
	} else if settings != nil && settings.FeesJSON != "" {
		// 解析费用规则
		feesConfig, err := ParseFeesJSON(settings.FeesJSON)
		if err != nil {
			uc.log.Errorf("HandlePlayerSitDown: failed to parse fees config, houseGID=%d, err=%v", houseGID, err)
		} else if len(feesConfig.Rules) > 0 {
			// 计算需要的费用（使用当前房间参数作为模拟战绩）
			// 由于玩家刚坐下，还没有实际分数，我们检查规则中的最低阈值
			requiredFee := uc.calculateMinRequiredFee(feesConfig, int(kindID), int(baseScore))

			if requiredFee > 0 {
				uc.log.Infof("[费用检查] 检查费用资格: gameID=%d, balance=%d, requiredFee=%d",
					gameID, member.Balance, requiredFee)

				// 检查玩家余额是否足以支付费用
				if member.Balance < requiredFee {
					uc.log.Warnf("HandlePlayerSitDown: insufficient balance for fee, gameID=%d, balance=%d < fee=%d, dismissing table",
						gameID, member.Balance, requiredFee)

					// 先解散桌子
					dismissErr := uc.dismissTable(ctx, userID, houseGID, kindID, tableNum, "insufficient balance for shop fee")
					if dismissErr != nil {
						uc.log.Errorf("Failed to dismiss table: %v", dismissErr)
					}

					// 禁用玩家（更新数据库 forbid 字段）
					if err := uc.memberRepo.UpdateForbid(ctx, houseGID, gameID, true); err != nil {
						uc.log.Errorf("Failed to forbid player with insufficient fee gameID=%d: %v", gameID, err)
					} else {
						uc.log.Infof("Forbid player with insufficient fee: gameID=%d, balance=%d < fee=%d",
							gameID, member.Balance, requiredFee)
					}

					return fmt.Errorf("insufficient balance for fee: %d < %d", member.Balance, requiredFee)
				}

				uc.log.Infof("[费用检查] 费用检查通过: gameID=%d, balance=%d >= fee=%d",
					gameID, member.Balance, requiredFee)
			}
		}
	}

	uc.log.Infof("[费用检查] 所有检查通过: gameID=%d, balance=%d >= credit=%d, keep table",
		gameID, member.Balance, effectiveCredit)

	return nil
}

// calculateMinRequiredFee 计算玩家需要的最低费用
// 与 passing-dragonfly 的逻辑一致：
// 1. 优先匹配全局规则（kind=0, base=0）
// 2. 如果没有全局规则，匹配特定游戏类型和底分的规则
func (uc *RoomCreditLimitUseCase) calculateMinRequiredFee(config *FeesConfig, kindID int, baseScore int) int32 {
	if len(config.Rules) == 0 {
		return 0
	}

	// 1. 优先查找全局通用规则 (kind=0 && base=0)
	for _, rule := range config.Rules {
		if rule.Kind == 0 && rule.Base == 0 {
			// 全局规则：返回费用（阈值为0时也需要支付费用）
			return int32(rule.Fee)
		}
	}

	// 2. 查找特定游戏类型和底分的规则
	for _, rule := range config.Rules {
		// 跳过全局规则（已经处理过）
		if rule.Kind == 0 && rule.Base == 0 {
			continue
		}

		// 匹配游戏类型
		kindMatches := (rule.Kind == 0 || rule.Kind == kindID)
		// 匹配底分
		baseMatches := (rule.Base == 0 || rule.Base == baseScore)

		// 同时匹配游戏类型和底分
		if kindMatches && baseMatches {
			return int32(rule.Fee)
		}
	}

	return 0 // 未匹配到任何规则
}

// dismissTable 解散桌子（调用 plaza 接口）
func (uc *RoomCreditLimitUseCase) dismissTable(ctx context.Context, userID int, houseGID int32, kindID int32, tableNum int32, reason string) error {
	uc.log.Infof("Dismissing table: userID=%d, houseGID=%d, kindID=%d, tableNum=%d, reason=%s",
		userID, houseGID, kindID, tableNum, reason)

	err := uc.plazaMgr.DismissTable(userID, int(houseGID), int(kindID), int(tableNum))
	if err != nil {
		uc.log.Errorf("Failed to dismiss table: %v", err)
		return fmt.Errorf("failed to dismiss table: %w", err)
	}

	return nil
}

// CheckPlayerCanEnterRoom 检查玩家是否可以进入房间（用于前端显示）
// 返回：是否可以进入、玩家余额、房间额度、玩家个人额度、有效额度
func (uc *RoomCreditLimitUseCase) CheckPlayerCanEnterRoom(ctx context.Context, houseGID int32, gameID int32, groupName string, gameKind int32, baseScore int32) (bool, int32, int32, int32, int32, error) {
	// 1. 获取玩家信息
	member, err := uc.memberRepo.GetByGameID(ctx, houseGID, gameID)
	if err != nil {
		uc.log.Errorf("CheckPlayerCanEnterRoom: failed to get member, houseGID=%d, gameID=%d, err=%v", houseGID, gameID, err)
		return false, 0, 0, 0, 0, err
	}

	// 2. 获取房间额度限制
	roomCreditLimit, _ := uc.repo.GetCreditLimit(ctx, houseGID, groupName, gameKind, baseScore)

	// 3. 计算有效额度
	effectiveCredit := roomCreditLimit + member.Credit

	// 4. 判断是否可以进入
	canEnter := member.Balance >= effectiveCredit

	return canEnter, member.Balance, roomCreditLimit, member.Credit, effectiveCredit, nil
}
