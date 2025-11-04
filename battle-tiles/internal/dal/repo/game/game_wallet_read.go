// internal/dal/repo/game/game_wallet_read.go
package game

import (
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type WalletReadRepo interface {
	// 读：非加锁
	Get(ctx context.Context, houseGID, memberID int32) (*model.GameMemberWallet, error)
	// 列表 + 计数（余额区间、是否个性额度、分页）
	ListWallets(ctx context.Context, houseGID int32, min, max *int32, hasCustomLimit *bool, page, size int32) ([]*model.GameMemberWallet, int64, error)
	// 流水列表 + 计数（可按成员/类型/时间范围筛选，分页）
	ListLedger(ctx context.Context, houseGID int32, memberID *int32, tp *int32, start, end time.Time, page, size int32) ([]*model.GameWalletLedger, int64, error)
}

type walletReadRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewWalletReadRepo(data *infra.Data, logger log.Logger) WalletReadRepo {
	return &walletReadRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/wallet_read")),
	}
}

func (r *walletReadRepo) db(ctx context.Context) *gorm.DB { return r.data.GetDBWithContext(ctx) }

func (r *walletReadRepo) Get(ctx context.Context, houseGID, memberID int32) (*model.GameMemberWallet, error) {
	var m model.GameMemberWallet
	if err := r.db(ctx).
		Where("house_gid = ? AND member_id = ?", houseGID, memberID).
		First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *walletReadRepo) ListWallets(ctx context.Context, houseGID int32, min, max *int32, hasCustomLimit *bool, page, size int32) ([]*model.GameMemberWallet, int64, error) {
	db := r.db(ctx).Model(&model.GameMemberWallet{}).Where("house_gid = ?", houseGID)

	if min != nil {
		db = db.Where("balance >= ?", *min)
	}
	if max != nil {
		db = db.Where("balance <= ?", *max)
	}
	if hasCustomLimit != nil {
		if *hasCustomLimit {
			db = db.Where("limit_min > 0")
		} else {
			db = db.Where("limit_min = 0")
		}
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 200 {
		size = 20
	}
	offset := (page - 1) * size

	var list []*model.GameMemberWallet
	if err := db.
		Order("updated_at DESC").
		Limit(int(size)).Offset(int(offset)).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *walletReadRepo) ListLedger(ctx context.Context, houseGID int32, memberID *int32, tp *int32, start, end time.Time, page, size int32) ([]*model.GameWalletLedger, int64, error) {
	db := r.db(ctx).Model(&model.GameWalletLedger{}).
		Where("house_gid = ? AND created_at >= ? AND created_at < ?", houseGID, start, end)

	if memberID != nil {
		db = db.Where("member_id = ?", *memberID)
	}
	if tp != nil {
		db = db.Where("type = ?", *tp)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 200 {
		size = 20
	}
	offset := (page - 1) * size

	var list []*model.GameWalletLedger
	if err := db.
		Order("created_at DESC"). // 更直观
		Limit(int(size)).Offset(int(offset)).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
