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
	// 璇伙細闈炲姞閿?	// 淇敼: 澧炲姞 groupID 鍙傛暟
	Get(ctx context.Context, houseGID, memberID int32, groupID *int32) (*model.GameMemberWallet, error)
	// 新增：按 game_id 查询
	GetByGameID(ctx context.Context, houseGID, gameID int32, groupID *int32) (*model.GameMemberWallet, error)

	// 鍒楄〃 + 璁℃暟锛堜綑棰濆尯闂淬€佹槸鍚︿釜鎬ч搴︺€佸垎椤碉級
	// 淇敼: 澧炲姞 groupID 鍙傛暟
	ListWallets(ctx context.Context, houseGID int32, groupID *int32, min, max *int32, hasCustomLimit *bool, page, size int32) ([]*model.GameMemberWallet, int64, error)

	// 鏂板: 鏌ヨ鎴愬憳鍦ㄦ墍鏈夊湀瀛愮殑浣欓
	ListMemberBalances(ctx context.Context, houseGID int32, memberID int32) ([]*model.GameMemberWallet, error)

	// 娴佹按鍒楄〃 + 璁℃暟锛堝彲鎸夋垚鍛?绫诲瀷/鏃堕棿鑼冨洿绛涢€夛紝鍒嗛〉锛?
	ListLedger(ctx context.Context, houseGID int32, memberID *int32, tp *int32, start, end time.Time, page, size int32) ([]*model.GameWalletLedger, int64, error)

	// 鎸夋垚鍛橀泦鍚堣繃婊ょ殑閽卞寘鍒楄〃
	ListWalletsByMembers(ctx context.Context, houseGID int32, memberIDs []int32, min, max *int32, page, size int32) ([]*model.GameMemberWallet, int64, error)
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

// Get 鑾峰彇閽卞寘淇℃伅
func (r *walletReadRepo) Get(ctx context.Context, houseGID, memberID int32, groupID *int32) (*model.GameMemberWallet, error) {
	var m model.GameMemberWallet
	db := r.db(ctx).Where("house_gid = ? AND member_id = ?", houseGID, memberID)

	// 鏂板: 鏀寔鎸夊湀瀛愮瓫閫?
	if groupID != nil {
		db = db.Where("group_id = ?", *groupID)
	}

	if err := db.First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// GetByGameID 按 game_id 查询（可选 group_id）
func (r *walletReadRepo) GetByGameID(ctx context.Context, houseGID, gameID int32, groupID *int32) (*model.GameMemberWallet, error) {
	var m model.GameMemberWallet
	db := r.db(ctx).Where("house_gid = ? AND game_id = ?", houseGID, gameID)

	if groupID != nil {
		db = db.Where("group_id = ?", *groupID)
	}

	if err := db.First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// ListWallets 鏌ヨ閽卞寘鍒楄〃
func (r *walletReadRepo) ListWallets(ctx context.Context, houseGID int32, groupID *int32, min, max *int32, hasCustomLimit *bool, page, size int32) ([]*model.GameMemberWallet, int64, error) {
	db := r.db(ctx).Model(&model.GameMemberWallet{}).Where("house_gid = ?", houseGID)

	// 鏂板: 鏀寔鎸夊湀瀛愮瓫閫?
	if groupID != nil {
		db = db.Where("group_id = ?", *groupID)
	}

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

// ListMemberBalances 鏌ヨ鎴愬憳鍦ㄦ墍鏈夊湀瀛愮殑浣欓
func (r *walletReadRepo) ListMemberBalances(ctx context.Context, houseGID int32, memberID int32) ([]*model.GameMemberWallet, error) {
	var wallets []*model.GameMemberWallet
	if err := r.db(ctx).
		Where("house_gid = ? AND member_id = ?", houseGID, memberID).
		Order("group_id").
		Find(&wallets).Error; err != nil {
		return nil, err
	}

	return wallets, nil
}

// ListLedger 鏌ヨ娴佹按鍒楄〃
func (r *walletReadRepo) ListLedger(ctx context.Context, houseGID int32, memberID *int32, tp *int32, start, end time.Time, page, size int32) ([]*model.GameWalletLedger, int64, error) {
	db := r.db(ctx).Model(&model.GameWalletLedger{}).Where("house_gid = ?", houseGID)
	if memberID != nil {
		db = db.Where("member_id = ?", *memberID)
	}
	if tp != nil {
		db = db.Where("type = ?", *tp)
	}
	db = db.Where("created_at >= ? AND created_at < ?", start, end)

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
	if err := db.Order("created_at DESC").Limit(int(size)).Offset(int(offset)).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// ListWalletsByMembers 鎸夋垚鍛橀泦鍚堣繃婊ょ殑閽卞寘鍒楄〃
func (r *walletReadRepo) ListWalletsByMembers(ctx context.Context, houseGID int32, memberIDs []int32, min, max *int32, page, size int32) ([]*model.GameMemberWallet, int64, error) {
	if len(memberIDs) == 0 {
		return []*model.GameMemberWallet{}, 0, nil
	}
	db := r.db(ctx).Model(&model.GameMemberWallet{}).
		Where("house_gid = ? AND member_id IN ?", houseGID, memberIDs)

	if min != nil {
		db = db.Where("balance >= ?", *min)
	}
	if max != nil {
		db = db.Where("balance <= ?", *max)
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
	if err := db.Order("updated_at DESC").Limit(int(size)).Offset(int(offset)).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
