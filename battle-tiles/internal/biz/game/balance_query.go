package game

import (
	model "battle-tiles/internal/dal/model/game"
	repo "battle-tiles/internal/dal/repo/game"
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
)

// MemberBalance 鎴愬憳浣欓淇℃伅
type MemberBalance struct {
	MemberID    int32   `json:"member_id"`
	GameID      int32   `json:"game_id"`
	GameName    string  `json:"game_name"`
	GroupID     *int32  `json:"group_id"`
	GroupName   string  `json:"group_name"`
	Balance     int32   `json:"balance"`      // 浣欓(鍒?
	BalanceYuan float64 `json:"balance_yuan"` // 浣欓(鍏?
	UpdatedAt   string  `json:"updated_at"`
}

// BalanceQueryUseCase
type BalanceQueryUseCase struct {
	memberRepo repo.GameMemberRepo
	walletRepo repo.WalletReadRepo
	groupRepo  repo.ShopGroupRepo
	log        *log.Helper
}

func NewBalanceQueryUseCase(
	memberRepo repo.GameMemberRepo,
	walletRepo repo.WalletReadRepo,
	groupRepo repo.ShopGroupRepo,
	logger log.Logger,
) *BalanceQueryUseCase {
	return &BalanceQueryUseCase{
		memberRepo: memberRepo,
		walletRepo: walletRepo,
		groupRepo:  groupRepo,
		log:        log.NewHelper(log.With(logger, "module", "usecase/balance_query")),
	}
}

// GetMyBalances 鐢ㄦ埛鏌ヨ鑷繁鐨勪綑棰?
func (uc *BalanceQueryUseCase) GetMyBalances(
	ctx context.Context,
	gameID int32,
	houseGID int32,
	groupID *int32,
) ([]*MemberBalance, error) {
	var balances []*MemberBalance

	if groupID != nil {
		// 鏌ヨ鎸囧畾鍦堝瓙鐨勪綑棰?
		member, err := uc.memberRepo.GetByGameIDAndGroup(ctx, houseGID, gameID, groupID)
		if err != nil {
			uc.log.Errorf("get member by game_id and group failed: %v", err)
			return nil, err
		}

		// 先按 game_id+group 查钱包，找不到再按 member_id+group
		wallet, err := uc.walletRepo.GetByGameID(ctx, houseGID, member.GameID, groupID)
		if err != nil {
			uc.log.Warnf("get wallet by game_id failed: %v, fallback member_id", err)
			wallet, err = uc.walletRepo.Get(ctx, houseGID, member.Id, groupID)
		}
		if err != nil || wallet == nil {
			wallet = &model.GameMemberWallet{Balance: member.Balance}
		}

		balance := &MemberBalance{
			MemberID:    member.Id,
			GameID:      member.GameID,
			GameName:    member.GameName,
			GroupID:     member.GroupID,
			GroupName:   member.GroupName,
			Balance:     wallet.Balance,
			BalanceYuan: float64(wallet.Balance) / 100.0,
			UpdatedAt:   wallet.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		balances = append(balances, balance)
	} else {
		// 鏌ヨ鎵€鏈夊湀瀛愮殑浣欓
		members, err := uc.memberRepo.ListByGameID(ctx, houseGID, gameID)
		if err != nil {
			uc.log.Errorf("list members by game_id failed: %v", err)
			return nil, err
		}

		for _, member := range members {
			// 先按 game_id+group 查钱包，找不到再按 member_id+group
			wallet, err := uc.walletRepo.GetByGameID(ctx, houseGID, member.GameID, member.GroupID)
			if err != nil {
				uc.log.Warnf("get wallet for member %d by game_id failed: %v, fallback member_id", member.Id, err)
				wallet, err = uc.walletRepo.Get(ctx, houseGID, member.Id, member.GroupID)
			}
			if err != nil || wallet == nil {
				wallet = &model.GameMemberWallet{Balance: member.Balance}
			}

			balance := &MemberBalance{
				MemberID:    member.Id,
				GameID:      member.GameID,
				GameName:    member.GameName,
				GroupID:     member.GroupID,
				GroupName:   member.GroupName,
				Balance:     wallet.Balance,
				BalanceYuan: float64(wallet.Balance) / 100.0,
			}

			if wallet.UpdatedAt.IsZero() {
				balance.UpdatedAt = member.UpdatedAt.Format("2006-01-02 15:04:05")
			} else {
				balance.UpdatedAt = wallet.UpdatedAt.Format("2006-01-02 15:04:05")
			}

			balances = append(balances, balance)
		}
	}

	return balances, nil
}

// ListGroupMemberBalances 绠＄悊鍛樻煡璇㈠湀瀛愭垚鍛樹綑棰?
func (uc *BalanceQueryUseCase) ListGroupMemberBalances(
	ctx context.Context,
	adminUserID int32,
	houseGID int32,
	groupID int32,
	min, max *int32,
	page, size int32,
) ([]*MemberBalance, int64, error) {
	// 楠岃瘉绠＄悊鍛樻潈闄?
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		uc.log.Errorf("get group failed: %v", err)
		return nil, 0, fmt.Errorf("鍦堝瓙涓嶅瓨鍦?")
	}

	if group.AdminUserID != adminUserID {
		uc.log.Warnf("admin %d has no permission to access group %d", adminUserID, groupID)
		return nil, 0, fmt.Errorf("鏃犳潈闄愯闂鍦堝瓙")
	}

	// 鏌ヨ鍦堝瓙鎴愬憳
	members, total, err := uc.memberRepo.ListByGroup(ctx, houseGID, groupID, page, size)
	if err != nil {
		uc.log.Errorf("list group members failed: %v", err)
		return nil, 0, err
	}

	var balances []*MemberBalance
	for _, member := range members {
		// 鏌ヨ閽卞寘淇℃伅
		wallet, err := uc.walletRepo.Get(ctx, houseGID, member.Id, member.GroupID)
		if err != nil {
			uc.log.Warnf("get wallet for member %d failed: %v, use member balance", member.Id, err)
			// 濡傛灉閽卞寘涓嶅瓨鍦?浣跨敤鎴愬憳琛ㄧ殑浣欓
			wallet = &model.GameMemberWallet{
				Balance: member.Balance,
			}
		}

		// 搴旂敤浣欓绛涢€?
		if min != nil && wallet.Balance < *min {
			continue
		}
		if max != nil && wallet.Balance > *max {
			continue
		}

		balance := &MemberBalance{
			MemberID:    member.Id,
			GameID:      member.GameID,
			GameName:    member.GameName,
			GroupID:     member.GroupID,
			GroupName:   member.GroupName,
			Balance:     wallet.Balance,
			BalanceYuan: float64(wallet.Balance) / 100.0,
		}

		if wallet.UpdatedAt.IsZero() {
			balance.UpdatedAt = member.UpdatedAt.Format("2006-01-02 15:04:05")
		} else {
			balance.UpdatedAt = wallet.UpdatedAt.Format("2006-01-02 15:04:05")
		}

		balances = append(balances, balance)
	}

	return balances, total, nil
}
