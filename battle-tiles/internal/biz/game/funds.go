package game

import (
	"battle-tiles/internal/dal/resp"
	"context"
	"errors"
	"fmt"
	"time"

	model "battle-tiles/internal/dal/model/game"
	repo "battle-tiles/internal/dal/repo/game"

	"gorm.io/gorm"
)

const (
	LedgerTypeDeposit       int32 = 1
	LedgerTypeWithdraw      int32 = 2
	LedgerTypeForceWithdraw int32 = 3
	LedgerTypeAdjust        int32 = 4
)

type FundsUseCase struct {
	wallet     repo.WalletRepo
	walletRead repo.WalletReadRepo
	memberRepo repo.GameMemberRepo // 直接更新 game_member.balance
}

func NewFundsUseCase(w repo.WalletRepo, r repo.WalletReadRepo, m repo.GameMemberRepo) *FundsUseCase {
	return &FundsUseCase{wallet: w, walletRead: r, memberRepo: m}
}

func (uc *FundsUseCase) Deposit(ctx context.Context, opUser int32, houseGID, memberID int32, amount int32, bizNo, reason string) (*model.GameMemberWallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be > 0")
	}
	// 简单去重：检查 bizNo 是否已存在
	if ok, err := uc.wallet.ExistsLedgerBiz(ctx, houseGID, memberID, bizNo); err != nil {
		return nil, err
	} else if ok {
		// 重复请求，返回当前余额
		m, _ := uc.memberRepo.GetByMemberID(ctx, houseGID, memberID)
		if m != nil {
			return &model.GameMemberWallet{HouseGID: houseGID, MemberID: memberID, Balance: m.Balance}, nil
		}
		return nil, errors.New("member not found")
	}

	// 直接更新 game_member.balance
	before, after, err := uc.memberRepo.UpdateBalance(ctx, houseGID, memberID, amount)
	if err != nil {
		return nil, fmt.Errorf("上分失败: %w", err)
	}

	// 记录流水
	_ = uc.wallet.AppendLedger(ctx, nil, &model.GameWalletLedger{
		HouseGID:       houseGID,
		MemberID:       memberID,
		ChangeAmount:   amount,
		BalanceBefore:  before,
		BalanceAfter:   after,
		Type:           LedgerTypeDeposit,
		Reason:         reason,
		OperatorUserID: opUser,
		BizNo:          bizNo,
	})

	return &model.GameMemberWallet{HouseGID: houseGID, MemberID: memberID, Balance: after}, nil
}

func (uc *FundsUseCase) Withdraw(ctx context.Context, opUser int32, houseGID, memberID int32, amount int32, bizNo, reason string, force bool) (*model.GameMemberWallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be > 0")
	}
	// 简单去重
	if ok, err := uc.wallet.ExistsLedgerBiz(ctx, houseGID, memberID, bizNo); err != nil {
		return nil, err
	} else if ok {
		m, _ := uc.memberRepo.GetByMemberID(ctx, houseGID, memberID)
		if m != nil {
			return &model.GameMemberWallet{HouseGID: houseGID, MemberID: memberID, Balance: m.Balance}, nil
		}
		return nil, errors.New("member not found")
	}

	// 先查询当前余额，检查是否可以下分
	m, err := uc.memberRepo.GetByMemberID(ctx, houseGID, memberID)
	if err != nil {
		return nil, errors.New("member not found")
	}
	if !force {
		if m.Forbid {
			return nil, errors.New("member forbidden")
		}
		if m.Balance-amount < 0 {
			return nil, fmt.Errorf("余额不足")
		}
	}

	// 直接更新 game_member.balance（负数代表减少）
	before, after, err := uc.memberRepo.UpdateBalance(ctx, houseGID, memberID, -amount)
	if err != nil {
		return nil, fmt.Errorf("下分失败: %w", err)
	}

	// 记录流水
	tp := LedgerTypeWithdraw
	if force {
		tp = LedgerTypeForceWithdraw
	}
	_ = uc.wallet.AppendLedger(ctx, nil, &model.GameWalletLedger{
		HouseGID:       houseGID,
		MemberID:       memberID,
		ChangeAmount:   -amount,
		BalanceBefore:  before,
		BalanceAfter:   after,
		Type:           tp,
		Reason:         reason,
		OperatorUserID: opUser,
		BizNo:          bizNo,
	})

	return &model.GameMemberWallet{HouseGID: houseGID, MemberID: memberID, Balance: after}, nil
}

func (uc *FundsUseCase) UpdateLimit(ctx context.Context, opUser int32, houseGID, memberID int32, limitMin *int32, forbid *bool, reason string) (*model.GameMemberWallet, error) {
	tx, txErr := uc.wallet.BeginTx(ctx)
	if txErr != nil {
		return nil, txErr
	}
	defer func() { _ = tx.Rollback() }()

	w, err := uc.wallet.GetForUpdate(ctx, tx, houseGID, memberID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w = &model.GameMemberWallet{
			HouseGID: houseGID,
			MemberID: memberID,
		}
	}
	if limitMin != nil {
		w.LimitMin = *limitMin
	}
	if forbid != nil {
		w.Forbid = *forbid
	}

	if err = uc.wallet.Upsert(ctx, tx, w); err != nil {
		return nil, err
	}
	// 审计流水：调整，变动额=0
	if err = uc.wallet.AppendLedger(ctx, tx, &model.GameWalletLedger{
		HouseGID:       houseGID,
		MemberID:       memberID,
		ChangeAmount:   0,
		BalanceBefore:  w.Balance,
		BalanceAfter:   w.Balance,
		Type:           LedgerTypeAdjust,
		Reason:         reason,
		OperatorUserID: opUser, // int32
		BizNo:          fmt.Sprintf("limit-%d-%d-%d", houseGID, memberID, opUser),
	}); err != nil {
		return nil, err
	}

	if commit := tx.Commit(); commit.Error != nil {
		return nil, commit.Error
	}
	return w, nil
}

// —— 单人钱包 ——
// 这里你的 model 就是 time.Time，直接塞出去即可
func (uc *FundsUseCase) GetWallet(ctx context.Context, houseGID, memberID int32) (*resp.WalletVO, error) {
	m, err := uc.walletRead.Get(ctx, houseGID, memberID, nil)
	if err != nil {
		return nil, err
	}
	return &resp.WalletVO{
		HouseGID:  m.HouseGID,
		MemberID:  m.MemberID,
		Balance:   m.Balance,
		Forbid:    m.Forbid,
		LimitMin:  m.LimitMin,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

// —— 钱包列表 ——
// 直接映射，无需额外转换
func (uc *FundsUseCase) ListWallets(ctx context.Context, houseGID int32, min, max *int32, hasCustomLimit *bool, page, size int32) ([]resp.WalletVO, int64, error) {
	list, total, err := uc.walletRead.ListWallets(ctx, houseGID, nil, min, max, hasCustomLimit, page, size)
	if err != nil {
		return nil, 0, err
	}
	out := make([]resp.WalletVO, 0, len(list))
	for _, m := range list {
		out = append(out, resp.WalletVO{
			HouseGID:  m.HouseGID,
			MemberID:  m.MemberID,
			Balance:   m.Balance,
			Forbid:    m.Forbid,
			LimitMin:  m.LimitMin,
			UpdatedAt: m.UpdatedAt,
		})
	}
	return out, total, nil
}

// 按成员集合过滤的钱包列表
func (uc *FundsUseCase) ListWalletsByMembers(ctx context.Context, houseGID int32, memberIDs []int32, min, max *int32, page, size int32) ([]resp.WalletVO, int64, error) {
	list, total, err := uc.walletRead.ListWalletsByMembers(ctx, houseGID, memberIDs, min, max, page, size)
	if err != nil {
		return nil, 0, err
	}
	out := make([]resp.WalletVO, 0, len(list))
	for _, m := range list {
		out = append(out, resp.WalletVO{
			HouseGID:  m.HouseGID,
			MemberID:  m.MemberID,
			Balance:   m.Balance,
			Forbid:    m.Forbid,
			LimitMin:  m.LimitMin,
			UpdatedAt: m.UpdatedAt,
		})
	}
	return out, total, nil
}

// —— 流水列表 ——
// start/end 为空时默认近7天；CreatedAt 已是 time.Time，直接赋值
func (uc *FundsUseCase) ListLedger(ctx context.Context, houseGID int32, memberID *int32, tp *int32, start, end *time.Time, page, size int32) ([]resp.LedgerVO, int64, error) {
	s, e := timeRangeDefault(start, end)
	list, total, err := uc.walletRead.ListLedger(ctx, houseGID, memberID, tp, s, e, page, size)
	if err != nil {
		return nil, 0, err
	}
	out := make([]resp.LedgerVO, 0, len(list))
	for _, l := range list {
		out = append(out, resp.LedgerVO{
			ID:             l.Id,
			HouseGID:       l.HouseGID,
			MemberID:       l.MemberID,
			ChangeAmount:   l.ChangeAmount,
			BalanceBefore:  l.BalanceBefore,
			BalanceAfter:   l.BalanceAfter,
			Type:           l.Type,
			Reason:         l.Reason,
			OperatorUserID: l.OperatorUserID,
			BizNo:          l.BizNo,
			CreatedAt:      l.CreatedAt, // 直接用 time.Time
		})
	}
	return out, total, nil
}

func timeRangeDefault(start, end *time.Time) (time.Time, time.Time) {
	now := time.Now()
	todayStart, tomorrowStart := dayRange(now)
	if start == nil && end == nil {
		return todayStart.AddDate(0, 0, -6), tomorrowStart
	}
	var s time.Time
	if start != nil {
		s = *start
	} else {
		s = todayStart.AddDate(0, 0, -6)
	}
	var e time.Time
	if end != nil {
		e = *end
	} else {
		e = tomorrowStart
	}
	return s, e
}

func dayRange(now time.Time) (start, end time.Time) {
	loc := now.Location()
	y, m, d := now.In(loc).Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, loc)
	end = start.AddDate(0, 0, 1)
	return
}
