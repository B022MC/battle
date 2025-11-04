package game

import (
	"context"

	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"
	corm "battle-tiles/pkg/plugin/gormx/repo"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepo interface {
	// 事务：由 biz 控制提交/回滚
	BeginTx(ctx context.Context) (*gorm.DB, error)

	// 行级锁获取（FOR UPDATE），tx 允许为 nil（则内部自取 db）
	GetForUpdate(ctx context.Context, tx *gorm.DB, houseGID, memberID int32) (*model.GameMemberWallet, error)

	// UPSERT（(house_gid,member_id) 冲突更新）
	Upsert(ctx context.Context, tx *gorm.DB, w *model.GameMemberWallet) error

	// 幂等流水写入（(house_gid,member_id,biz_no) 唯一）
	AppendLedger(ctx context.Context, tx *gorm.DB, l *model.GameWalletLedger) error

	// 幂等检查
	ExistsLedgerBiz(ctx context.Context, houseGID, memberID int32, bizNo string) (bool, error)
}

type walletRepo struct {
	corm.CORMImpl[model.GameMemberWallet]
	data *infra.Data
	log  *log.Helper
}

func NewWalletRepo(data *infra.Data, logger log.Logger) WalletRepo {
	return &walletRepo{
		CORMImpl: corm.NewCORMImplRepo[model.GameMemberWallet](data),
		data:     data,
		log:      log.NewHelper(log.With(logger, "module", "repo/wallet")),
	}
}

func (r *walletRepo) db(ctx context.Context) *gorm.DB {
	return r.data.GetDBWithContext(ctx)
}

func (r *walletRepo) BeginTx(ctx context.Context) (*gorm.DB, error) {
	return r.db(ctx).Begin(), nil
}

func (r *walletRepo) GetForUpdate(ctx context.Context, tx *gorm.DB, houseGID, memberID int32) (*model.GameMemberWallet, error) {
	db := r.db(ctx)
	if tx != nil {
		db = tx
	}
	var m model.GameMemberWallet
	if err := db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("house_gid = ? AND member_id = ?", houseGID, memberID).
		First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *walletRepo) Upsert(ctx context.Context, tx *gorm.DB, w *model.GameMemberWallet) error {
	db := r.db(ctx)
	if tx != nil {
		db = tx
	}
	return db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "house_gid"}, {Name: "member_id"}},
			DoUpdates: clause.Assignments(map[string]any{
				"balance":   w.Balance,
				"forbid":    w.Forbid,
				"limit_min": w.LimitMin,
			}),
		}).
		Create(w).Error
}

func (r *walletRepo) AppendLedger(ctx context.Context, tx *gorm.DB, l *model.GameWalletLedger) error {
	db := r.db(ctx)
	if tx != nil {
		db = tx
	}
	return db.Create(l).Error
}

func (r *walletRepo) ExistsLedgerBiz(ctx context.Context, houseGID, memberID int32, bizNo string) (bool, error) {
	var cnt int64
	err := r.db(ctx).
		Model(&model.GameWalletLedger{}).
		Where("house_gid = ? AND member_id = ? AND biz_no = ?", houseGID, memberID, bizNo).
		Count(&cnt).Error
	return cnt > 0, err
}
