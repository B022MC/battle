package game

import (
	model "battle-tiles/internal/dal/model/game"
	repo "battle-tiles/internal/dal/repo/game"
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// BattleStats 鎴樼哗缁熻鏁版嵁
type BattleStats struct {
	TotalGames int64   `json:"total_games"`
	TotalScore int     `json:"total_score"`
	TotalFee   int     `json:"total_fee"`
	AvgScore   float64 `json:"avg_score"`
	GroupID    *int32  `json:"group_id,omitempty"`
	GroupName  string  `json:"group_name,omitempty"`
}

// GroupStats 鍦堝瓙缁熻鏁版嵁
type GroupStats struct {
	GroupID       int32  `json:"group_id"`
	GroupName     string `json:"group_name"`
	TotalGames    int64  `json:"total_games"`
	TotalScore    int    `json:"total_score"`
	TotalFee      int    `json:"total_fee"`
	ActiveMembers int64  `json:"active_members"`
}

// HouseStats 搴楅摵缁熻鏁版嵁
type HouseStats struct {
	HouseGID   int32 `json:"house_gid"`
	TotalGames int64 `json:"total_games"`
	TotalScore int   `json:"total_score"`
	TotalFee   int   `json:"total_fee"`
}

// BattleQueryUseCase 鎴樼哗鏌ヨ涓氬姟閫昏緫
type BattleQueryUseCase struct {
	battleRepo repo.BattleRecordRepo
	memberRepo repo.GameMemberRepo
	groupRepo  repo.ShopGroupRepo
	log        *log.Helper
}

func NewBattleQueryUseCase(
	battleRepo repo.BattleRecordRepo,
	memberRepo repo.GameMemberRepo,
	groupRepo repo.ShopGroupRepo,
	logger log.Logger,
) *BattleQueryUseCase {
	return &BattleQueryUseCase{
		battleRepo: battleRepo,
		memberRepo: memberRepo,
		groupRepo:  groupRepo,
		log:        log.NewHelper(log.With(logger, "module", "usecase/battle_query")),
	}
}

func (uc *BattleQueryUseCase) ListMyBattles(
	ctx context.Context,
	gameID int32,
	houseGID int32,
	groupID *int32,
	start, end *time.Time,
	page, size int32,
) ([]*model.GameBattleRecord, int64, error) {
	records, total, err := uc.battleRepo.ListByPlayer(ctx, houseGID, gameID, groupID, start, end, page, size)
	if err != nil {
		uc.log.Errorf("list my battles failed: %v", err)
		return nil, 0, err
	}

	return records, total, nil
}

// GetMyStats 鐢ㄦ埛鏌ヨ鑷繁鐨勭粺璁℃暟鎹?
func (uc *BattleQueryUseCase) GetMyStats(
	ctx context.Context,
	gameID int32,
	houseGID int32,
	groupID *int32,
	start, end *time.Time,
) (*BattleStats, error) {
	// 鏌ヨ缁熻鏁版嵁
	totalGames, totalScore, totalFee, err := uc.battleRepo.GetPlayerStats(ctx, houseGID, gameID, groupID, start, end)
	if err != nil {
		uc.log.Errorf("get my stats failed: %v", err)
		return nil, err
	}

	stats := &BattleStats{
		TotalGames: totalGames,
		TotalScore: totalScore,
		TotalFee:   totalFee,
		GroupID:    groupID,
	}

	if totalGames > 0 {
		stats.AvgScore = float64(totalScore) / float64(totalGames)
	}

	// 濡傛灉鎸囧畾浜嗗湀瀛?鏌ヨ鍦堝瓙鍚嶇О
	if groupID != nil {
		group, err := uc.groupRepo.GetByID(ctx, *groupID)
		if err == nil {
			stats.GroupName = group.GroupName
		}
	}

	return stats, nil
}

// ListGroupBattles 绠＄悊鍛樻煡璇㈠湀瀛愭垬缁?
func (uc *BattleQueryUseCase) ListGroupBattles(
	ctx context.Context,
	adminUserID int32,
	houseGID int32,
	groupID int32,
	playerGameID *int32,
	start, end *time.Time,
	page, size int32,
) ([]*model.GameBattleRecord, int64, error) {
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		uc.log.Errorf("get group failed: %v", err)
		return nil, 0, fmt.Errorf("")
	}

	if group.AdminUserID != adminUserID {
		uc.log.Warnf("admin %d has no permission to access group %d", adminUserID, groupID)
		return nil, 0, fmt.Errorf("鏃犳潈闄愯闂鍦堝瓙")
	}

	// 鏌ヨ鎴樼哗
	var records []*model.GameBattleRecord
	var total int64

	if playerGameID != nil {
		// 鏌ヨ鎸囧畾鐜╁鐨勬垬缁?
		records, total, err = uc.battleRepo.ListByPlayer(ctx, houseGID, *playerGameID, &groupID, start, end, page, size)
	} else {
		// 鏌ヨ鍦堝瓙鎵€鏈夋垬缁?
		records, total, err = uc.battleRepo.List(ctx, houseGID, &groupID, nil, start, end, page, size)
	}

	if err != nil {
		uc.log.Errorf("list group battles failed: %v", err)
		return nil, 0, err
	}

	return records, total, nil
}

// GetGroupStats 绠＄悊鍛樻煡璇㈠湀瀛愮粺璁?
func (uc *BattleQueryUseCase) GetGroupStats(
	ctx context.Context,
	adminUserID int32,
	houseGID int32,
	groupID int32,
	start, end *time.Time,
) (*GroupStats, error) {
	// 楠岃瘉绠＄悊鍛樻潈闄?
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		uc.log.Errorf("get group failed: %v", err)
		return nil, fmt.Errorf("鍦堝瓙涓嶅瓨鍦?")
	}

	if group.AdminUserID != adminUserID {
		uc.log.Warnf("admin %d has no permission to access group %d", adminUserID, groupID)
		return nil, fmt.Errorf("鏃犳潈闄愯闂鍦堝瓙")
	}

	// 鏌ヨ缁熻鏁版嵁
	totalGames, totalScore, totalFee, activeMembers, err := uc.battleRepo.GetGroupStats(ctx, houseGID, groupID, start, end)
	if err != nil {
		uc.log.Errorf("get group stats failed: %v", err)
		return nil, err
	}

	stats := &GroupStats{
		GroupID:       groupID,
		GroupName:     group.GroupName,
		TotalGames:    totalGames,
		TotalScore:    totalScore,
		TotalFee:      totalFee,
		ActiveMembers: activeMembers,
	}

	return stats, nil
}

// GetHouseStats 瓒呯骇绠＄悊鍛樻煡璇㈠簵閾虹粺璁?
func (uc *BattleQueryUseCase) GetHouseStats(
	ctx context.Context,
	superAdminUserID int32,
	houseGID int32,
	start, end *time.Time,
) (*HouseStats, error) {
	// TODO: 妫€鏌ヨ秴绾х鐞嗗憳鏉冮檺
	// 杩欓噷搴旇楠岃瘉 superAdminUserID 鏄惁鏄秴绾х鐞嗗憳

	// 鏌ヨ缁熻鏁版嵁
	totalGames, totalScore, totalFee, err := uc.battleRepo.GetHouseStats(ctx, houseGID, start, end)
	if err != nil {
		uc.log.Errorf("get house stats failed: %v", err)
		return nil, err
	}

	stats := &HouseStats{
		HouseGID:   houseGID,
		TotalGames: totalGames,
		TotalScore: totalScore,
		TotalFee:   totalFee,
	}

	return stats, nil
}
