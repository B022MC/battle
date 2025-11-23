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
	repo       repo.RoomCreditLimitRepo
	memberRepo repo.GameMemberRepo
	plazaMgr   plaza.Manager
	log        *log.Helper
}

func NewRoomCreditLimitUseCase(
	r repo.RoomCreditLimitRepo,
	memberRepo repo.GameMemberRepo,
	plazaMgr plaza.Manager,
	logger log.Logger,
) *RoomCreditLimitUseCase {
	return &RoomCreditLimitUseCase{
		repo:       r,
		memberRepo: memberRepo,
		plazaMgr:   plazaMgr,
		log:        log.NewHelper(log.With(logger, "module", "usecase/room_credit_limit")),
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
// 这是核心逻辑：当玩家坐下时，检查余额是否满足额度要求，不足则解散桌子
func (uc *RoomCreditLimitUseCase) HandlePlayerSitDown(ctx context.Context, userID int, houseGID int32, gameID int32, kindID int32, baseScore int32, tableNum int32) error {
	// 1. 获取玩家信息（余额、圈子、个人额度调整）
	member, err := uc.memberRepo.GetByGameID(ctx, houseGID, gameID)
	if err != nil {
		uc.log.Errorf("HandlePlayerSitDown: failed to get member, houseGID=%d, gameID=%d, err=%v", houseGID, gameID, err)
		// 玩家没有录入，解散桌子
		return uc.dismissTable(ctx, userID, houseGID, kindID, tableNum, "player not found")
	}

	// 2. 获取房间额度要求
	roomCreditLimit, found := uc.repo.GetCreditLimit(ctx, houseGID, member.GroupName, kindID, baseScore)
	if !found {
		uc.log.Warnf("HandlePlayerSitDown: no credit limit found, using default, houseGID=%d, group=%s, kind=%d, base=%d",
			houseGID, member.GroupName, kindID, baseScore)
	}

	// 3. 计算有效额度要求（房间额度 + 玩家个人额度调整）
	effectiveCredit := roomCreditLimit + member.Credit

	uc.log.Debugf("HandlePlayerSitDown: gameID=%d, balance=%d, roomCredit=%d, playerCredit=%d, effectiveCredit=%d",
		gameID, member.Balance, roomCreditLimit, member.Credit, effectiveCredit)

	// 4. 检查余额是否满足额度要求
	if member.Balance < effectiveCredit {
		// 余额不足，解散桌子
		uc.log.Warnf("HandlePlayerSitDown: insufficient balance, gameID=%d, balance=%d < credit=%d, dismissing table",
			gameID, member.Balance, effectiveCredit)

		// TODO: 可以在这里记录违规次数，连续3次自动禁用玩家（参考 passing-dragonfly 的 thiefCache 逻辑）

		return uc.dismissTable(ctx, userID, houseGID, kindID, tableNum, "insufficient balance")
	}

	uc.log.Debugf("HandlePlayerSitDown: balance sufficient, gameID=%d, balance=%d >= credit=%d, keep table",
		gameID, member.Balance, effectiveCredit)

	return nil
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
