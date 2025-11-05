package game

import (
	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type UserApplicationRepo interface {
	Insert(ctx context.Context, in *model.UserApplication) error
	ListByHouse(ctx context.Context, houseGID int32, typ *int32, status *int32) ([]*model.UserApplication, error)
	ExistsPending(ctx context.Context, houseGID, applicant, typ, adminUID int32) (bool, error)
	GetByID(ctx context.Context, id int32) (*model.UserApplication, error)
	UpdateStatusByID(ctx context.Context, id int32, status int32) (int64, error)
	ListHistory(ctx context.Context, houseGID int32, applicant *int32, typ, status *int32, start, end *time.Time) ([]*model.UserApplication, error)
	// ListApprovedJoinsByAdmin 按圈主返回已通过的入圈申请（作为平台侧圈成员）
	ListApprovedJoinsByAdmin(ctx context.Context, houseGID int32, adminUID int32) ([]*model.UserApplication, error)
	// ListApprovedJoins 返回该店铺下所有已通过的入圈申请
	ListApprovedJoins(ctx context.Context, houseGID int32) ([]*model.UserApplication, error)
	// RemoveApprovedJoin 将指定用户在指定圈的已通过入圈记录标记为移除（status=3）
	RemoveApprovedJoin(ctx context.Context, houseGID int32, adminUID int32, applicant int32) (int64, error)
}

type userApplicationRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewUserApplicationRepo(data *infra.Data, logger log.Logger) UserApplicationRepo {
	return &userApplicationRepo{data: data, log: log.NewHelper(log.With(logger, "module", "repo/user_application"))}
}

func (r *userApplicationRepo) db(ctx context.Context) *gorm.DB { return r.data.GetDBWithContext(ctx) }

func (r *userApplicationRepo) Insert(ctx context.Context, in *model.UserApplication) error {
	return r.db(ctx).Create(in).Error
}

func (r *userApplicationRepo) ListByHouse(ctx context.Context, houseGID int32, typ *int32, status *int32) ([]*model.UserApplication, error) {
	db := r.db(ctx).Model(&model.UserApplication{}).Where("house_gid = ?", houseGID)
	if typ != nil {
		db = db.Where("type = ?", *typ)
	}
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	var list []*model.UserApplication
	if err := db.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *userApplicationRepo) ExistsPending(ctx context.Context, houseGID, applicant, typ, adminUID int32) (bool, error) {
	var cnt int64
	err := r.db(ctx).Model(&model.UserApplication{}).
		Where("house_gid = ? AND applicant = ? AND type = ? AND admin_user_id = ? AND status = 0", houseGID, applicant, typ, adminUID).
		Count(&cnt).Error
	return cnt > 0, err
}

func (r *userApplicationRepo) GetByID(ctx context.Context, id int32) (*model.UserApplication, error) {
	var m model.UserApplication
	if err := r.db(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *userApplicationRepo) UpdateStatusByID(ctx context.Context, id int32, status int32) (int64, error) {
	tx := r.db(ctx).Model(&model.UserApplication{}).Where("id = ?", id).Update("status", status)
	return tx.RowsAffected, tx.Error
}

func (r *userApplicationRepo) ListHistory(ctx context.Context, houseGID int32, applicant *int32, typ, status *int32, start, end *time.Time) ([]*model.UserApplication, error) {
	db := r.db(ctx).Model(&model.UserApplication{}).Where("house_gid = ?", houseGID)
	if applicant != nil {
		db = db.Where("applicant = ?", *applicant)
	}
	if typ != nil {
		db = db.Where("type = ?", *typ)
	}
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	if start != nil {
		db = db.Where("created_at >= ?", *start)
	}
	if end != nil {
		db = db.Where("created_at < ?", *end)
	}
	var list []*model.UserApplication
	if err := db.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *userApplicationRepo) ListApprovedJoinsByAdmin(ctx context.Context, houseGID int32, adminUID int32) ([]*model.UserApplication, error) {
	var list []*model.UserApplication
	err := r.db(ctx).Where("house_gid = ? AND type = 2 AND status = 1 AND admin_user_id = ?", houseGID, adminUID).
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

func (r *userApplicationRepo) ListApprovedJoins(ctx context.Context, houseGID int32) ([]*model.UserApplication, error) {
	var list []*model.UserApplication
	err := r.db(ctx).Where("house_gid = ? AND type = 2 AND status = 1", houseGID).
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

func (r *userApplicationRepo) RemoveApprovedJoin(ctx context.Context, houseGID int32, adminUID int32, applicant int32) (int64, error) {
	tx := r.db(ctx).Model(&model.UserApplication{}).
		Where("house_gid = ? AND type = 2 AND status = 1 AND admin_user_id = ? AND applicant = ?", houseGID, adminUID, applicant).
		Update("status", 3)
	return tx.RowsAffected, tx.Error
}
