package basic

import (
	"context"

	basicModel "battle-tiles/internal/dal/model/basic"
	"battle-tiles/internal/infra"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type BaseRoleMenuRelRepo interface {
	// 基本：替换某角色的全部菜单
	BatchReplaceByRole(ctx context.Context, roleID int32, menuIDs []int32) error
	// 查询
	ListMenuIDsByRole(ctx context.Context, roleID int32) ([]int32, error)
	ListRoleIDsByMenu(ctx context.Context, menuID int32) ([]int32, error)
	// 单条增删（可选）
	Add(ctx context.Context, roleID, menuID int32) error
	Remove(ctx context.Context, roleID, menuID int32) error
	// 批量保存（若你已有组装好的 []*BasicRoleMenuRel）
	BatchSave(ctx context.Context, roleID int32, list []*basicModel.BasicRoleMenuRel) error
}

type baseRoleMenuRelRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewBaseRoleMenuRelRepo(data *infra.Data, logger log.Logger) BaseRoleMenuRelRepo {
	return &baseRoleMenuRelRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/basicRoleMenuRel")),
	}
}

// 替换：先删旧，再插新（事务）
func (rp *baseRoleMenuRelRepo) BatchReplaceByRole(ctx context.Context, roleID int32, menuIDs []int32) error {
	return rp.data.GetDBWithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除该角色的所有旧关系
		if err := tx.
			Where("role_id = ?", roleID).
			Delete(&basicModel.BasicRoleMenuRel{}).Error; err != nil {
			return err
		}

		// 空集合则完成
		if len(menuIDs) == 0 {
			return nil
		}

		// 去重/组装
		rels := make([]*basicModel.BasicRoleMenuRel, 0, len(menuIDs))
		seen := map[int32]struct{}{}
		for _, mid := range menuIDs {
			if _, ok := seen[mid]; ok {
				continue
			}
			seen[mid] = struct{}{}
			rels = append(rels, &basicModel.BasicRoleMenuRel{
				RoleID: roleID,
				MenuID: mid,
			})
		}

		// 批量插入
		return tx.Create(&rels).Error
	})
}

func (rp *baseRoleMenuRelRepo) BatchSave(ctx context.Context, roleID int32, list []*basicModel.BasicRoleMenuRel) error {
	return rp.data.GetDBWithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).
			Delete(&basicModel.BasicRoleMenuRel{}).Error; err != nil {
			return err
		}
		if len(list) == 0 {
			return nil
		}
		return tx.Create(&list).Error
	})
}

func (rp *baseRoleMenuRelRepo) Add(ctx context.Context, roleID, menuID int32) error {
	rel := &basicModel.BasicRoleMenuRel{RoleID: roleID, MenuID: menuID}
	return rp.data.GetDBWithContext(ctx).Create(rel).Error
}

func (rp *baseRoleMenuRelRepo) Remove(ctx context.Context, roleID, menuID int32) error {
	return rp.data.GetDBWithContext(ctx).
		Where("role_id = ? AND menu_id = ?", roleID, menuID).
		Delete(&basicModel.BasicRoleMenuRel{}).Error
}

func (rp *baseRoleMenuRelRepo) ListMenuIDsByRole(ctx context.Context, roleID int32) ([]int32, error) {
	var rows []basicModel.BasicRoleMenuRel
	if err := rp.data.GetDBWithContext(ctx).
		Where("role_id = ?", roleID).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]int32, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.MenuID)
	}
	return out, nil
}

func (rp *baseRoleMenuRelRepo) ListRoleIDsByMenu(ctx context.Context, menuID int32) ([]int32, error) {
	var rows []basicModel.BasicRoleMenuRel
	if err := rp.data.GetDBWithContext(ctx).
		Where("menu_id = ?", menuID).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]int32, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.RoleID)
	}
	return out, nil
}
