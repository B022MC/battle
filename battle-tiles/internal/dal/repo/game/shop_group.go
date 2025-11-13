package game

import (
	"context"

	model "battle-tiles/internal/dal/model/game"
	"battle-tiles/internal/infra"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type ShopGroupRepo interface {
	// Create 创建圈子
	Create(ctx context.Context, m *model.GameShopGroup) error
	// GetByID 根据ID获取圈子
	GetByID(ctx context.Context, id int32) (*model.GameShopGroup, error)
	// GetByAdmin 根据管理员获取圈子
	GetByAdmin(ctx context.Context, houseGID int32, adminUserID int32) (*model.GameShopGroup, error)
	// ListByAdmin 获取管理员的所有圈子
	ListByAdmin(ctx context.Context, adminUserID int32) ([]*model.GameShopGroup, error)
	// ListByHouse 获取店铺下的所有圈子
	ListByHouse(ctx context.Context, houseGID int32) ([]*model.GameShopGroup, error)
	// Update 更新圈子信息
	Update(ctx context.Context, m *model.GameShopGroup) error
	// Delete 删除圈子
	Delete(ctx context.Context, id int32) error
	// Deactivate 停用圈子
	Deactivate(ctx context.Context, id int32) error
}

type shopGroupRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewShopGroupRepo(data *infra.Data, logger log.Logger) ShopGroupRepo {
	return &shopGroupRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/shop_group")),
	}
}

func (r *shopGroupRepo) db(ctx context.Context) *gorm.DB {
	return r.data.GetDBWithContext(ctx)
}

func (r *shopGroupRepo) Create(ctx context.Context, m *model.GameShopGroup) error {
	return r.db(ctx).Create(m).Error
}

func (r *shopGroupRepo) GetByID(ctx context.Context, id int32) (*model.GameShopGroup, error) {
	var out model.GameShopGroup
	err := r.db(ctx).
		Where("id = ? AND is_active = ?", id, true).
		First(&out).Error
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *shopGroupRepo) GetByAdmin(ctx context.Context, houseGID int32, adminUserID int32) (*model.GameShopGroup, error) {
	var out model.GameShopGroup
	err := r.db(ctx).
		Where("house_gid = ? AND admin_user_id = ? AND is_active = ?", houseGID, adminUserID, true).
		First(&out).Error
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *shopGroupRepo) ListByAdmin(ctx context.Context, adminUserID int32) ([]*model.GameShopGroup, error) {
	var out []*model.GameShopGroup
	err := r.db(ctx).
		Where("admin_user_id = ? AND is_active = ?", adminUserID, true).
		Order("id DESC").
		Find(&out).Error
	return out, err
}

func (r *shopGroupRepo) ListByHouse(ctx context.Context, houseGID int32) ([]*model.GameShopGroup, error) {
	var out []*model.GameShopGroup
	err := r.db(ctx).
		Where("house_gid = ? AND is_active = ?", houseGID, true).
		Order("id DESC").
		Find(&out).Error
	return out, err
}

func (r *shopGroupRepo) Update(ctx context.Context, m *model.GameShopGroup) error {
	return r.db(ctx).
		Model(&model.GameShopGroup{}).
		Where("id = ?", m.Id).
		Updates(map[string]interface{}{
			"group_name":  m.GroupName,
			"description": m.Description,
		}).Error
}

func (r *shopGroupRepo) Delete(ctx context.Context, id int32) error {
	// 物理删除圈子
	return r.db(ctx).Delete(&model.GameShopGroup{}, id).Error
}

func (r *shopGroupRepo) Deactivate(ctx context.Context, id int32) error {
	return r.db(ctx).
		Model(&model.GameShopGroup{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}
