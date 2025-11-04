package basic

import (
	"battle-tiles/internal/consts"
	"battle-tiles/internal/dal/common"
	basicModel "battle-tiles/internal/dal/model/basic"
	"battle-tiles/internal/infra"
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// BasicMenuRepo is a BasicMenu repo.
type BasicMenuRepo interface {
	ListByPage(ctx context.Context, offset, limit int, funcs ...func(*gorm.DB) *gorm.DB) ([]*basicModel.BasicMenu, int64, error)
	OptionByPage(ctx context.Context, offset, limit int, query interface{}, funcs ...func(*gorm.DB) *gorm.DB) ([]*common.Option, int64, error)
	Save(ctx context.Context, m *basicModel.BasicMenu) error
	BatchSave(ctx context.Context, list []*basicModel.BasicMenu) error
	UpdateFieldByCond(ctx context.Context, fields map[string]interface{}, query interface{}, args ...interface{}) error
	UpdateByField(ctx context.Context, fields map[string]interface{}, value *basicModel.BasicMenu) (*basicModel.BasicMenu, error)
	DeleteById(ctx context.Context, ids ...interface{}) error
	LogicDelete(ctx context.Context, value *basicModel.BasicMenu) error
	LogicDeleteByConds(ctx context.Context, query interface{}, args ...interface{}) error
	FindOneByScope(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) (*basicModel.BasicMenu, error)
	FindByScope(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) ([]*basicModel.BasicMenu, error)
	FindByKeyword(ctx context.Context, keyword string) ([]*basicModel.BasicMenu, error)
}

type basicMenuRepo struct {
	data *infra.Data
	log  *log.Helper
}

// NewBaseMenuRepo .
func NewBaseMenuRepo(data *infra.Data, logger log.Logger) BasicMenuRepo {
	return &basicMenuRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/BasicMenu")),
	}
}

func (rp *basicMenuRepo) ListByPage(ctx context.Context, offset, limit int, funcs ...func(*gorm.DB) *gorm.DB) ([]*basicModel.BasicMenu, int64, error) {
	xdb := rp.data.GetDBWithContext(ctx).Model(&basicModel.BasicMenu{}).Scopes(funcs...)
	var (
		result []*basicModel.BasicMenu
		total  int64
	)
	if err := xdb.Offset(offset).Limit(limit).Find(&result).Error; err != nil {
		return nil, 0, err
	}
	if size := len(result); 0 < limit && 0 < size && size < limit {
		total = int64(size + offset)
		return result, total, nil
	}
	xdb.Offset(-1).Limit(-1).Count(&total)
	return result, total, nil
}

func (rp *basicMenuRepo) OptionByPage(ctx context.Context, offset, limit int, query interface{}, funcs ...func(*gorm.DB) *gorm.DB) ([]*common.Option, int64, error) {
	xdb := rp.data.GetDBWithContext(ctx).Model(&basicModel.BasicMenu{}).Select(query).Scopes(funcs...)
	var (
		result []*common.Option
		total  int64
	)
	if err := xdb.Offset(offset).Limit(limit).Find(&result).Error; err != nil {
		return nil, 0, err
	}
	if size := len(result); 0 < limit && 0 < size && size < limit {
		total = int64(size + offset)
		return result, total, nil
	}
	xdb.Offset(-1).Limit(-1).Count(&total)
	return result, total, nil
}

func (rp *basicMenuRepo) Save(ctx context.Context, m *basicModel.BasicMenu) error {
	return rp.data.GetDBWithContext(ctx).Create(m).Error
}

func (rp *basicMenuRepo) UpdateFieldByCond(ctx context.Context, fields map[string]interface{}, query interface{}, args ...interface{}) error {
	return rp.data.GetDBWithContext(ctx).Model(&basicModel.BasicMenu{}).Where(query, args...).Updates(fields).Error
}

func (rp *basicMenuRepo) UpdateByField(ctx context.Context, fields map[string]interface{}, value *basicModel.BasicMenu) (*basicModel.BasicMenu, error) {
	if err := rp.data.GetDBWithContext(ctx).Model(value).Updates(fields).Error; err != nil {
		return value, err
	}
	return value, nil
}

func (rp *basicMenuRepo) DeleteById(ctx context.Context, ids ...interface{}) error {
	// 软删除（deleted_at 设置为当前时间）
	return rp.data.GetDBWithContext(ctx).Delete(&basicModel.BasicMenu{}, ids).Error
}

func (rp *basicMenuRepo) LogicDeleteByConds(ctx context.Context, query interface{}, args ...interface{}) error {
	return rp.data.GetDBWithContext(ctx).
		Where(query, args...).
		Delete(&basicModel.BasicMenu{}).Error
}

func (rp *basicMenuRepo) LogicDelete(ctx context.Context, value *basicModel.BasicMenu) error {
	// 软删除（按主键）
	return rp.data.GetDBWithContext(ctx).Delete(value).Error
}

func (rp *basicMenuRepo) FindOneByScope(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) (*basicModel.BasicMenu, error) {
	var dest basicModel.BasicMenu
	if err := rp.data.GetDBWithContext(ctx).Model(&basicModel.BasicMenu{}).Scopes(funcs...).First(&dest).Error; err != nil {
		return nil, err
	}
	return &dest, nil
}

func (rp *basicMenuRepo) FindByScope(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) ([]*basicModel.BasicMenu, error) {
	var dest []*basicModel.BasicMenu
	if err := rp.data.GetDBWithContext(ctx).Model(&basicModel.BasicMenu{}).Scopes(funcs...).Find(&dest).Error; err != nil {
		return nil, err
	}
	return dest, nil
}

func (rp *basicMenuRepo) FindByKeyword(ctx context.Context, keyword string) ([]*basicModel.BasicMenu, error) {
	var dest []*basicModel.BasicMenu

	// 参数化 + 表名修正 + "rank" 加引号
	sqlStr := fmt.Sprintf(`
WITH RECURSIVE cte AS (
  SELECT * FROM %s WHERE name ILIKE ?
  UNION ALL
  SELECT p.* FROM %s p
  INNER JOIN cte ON p.id = cte.parent_id
)
SELECT DISTINCT ON (cte.id) cte.*
FROM cte
ORDER BY cte.id, cte."rank" ASC
`, basicModel.TableNameBasicMenu, basicModel.TableNameBasicMenu)

	if err := rp.data.GetDBWithContext(ctx).Raw(sqlStr, "%"+keyword+"%").Find(&dest).Error; err != nil {
		return nil, err
	}
	return dest, nil
}

// 递归写入/更新（按 parent_id + path 判重）
func (rp *basicMenuRepo) createChildren(ctx context.Context, tx *gorm.DB, parentID int32, list []*basicModel.BasicMenu) error {
	for i := range list {
		menu := list[i]
		menu.ParentId = parentID

		var existing basicModel.BasicMenu
		err := tx.Table(basicModel.TableNameBasicMenu).
			Where("parent_id = ? AND path = ?", menu.ParentId, menu.Path).
			First(&existing).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// 不存在 -> 新建
				if e := tx.Table(basicModel.TableNameBasicMenu).Create(menu).Error; e != nil {
					return e
				}
			} else {
				return err
			}
		} else {
			// 已存在 -> 更新字段
			updates := map[string]interface{}{
				"menu_type":        menu.MenuType,
				"title":            menu.Title,
				"name":             menu.Name,
				"component":        menu.Component,
				"rank":             menu.Rank,
				"redirect":         menu.Redirect,
				"icon":             menu.Icon,
				"extra_icon":       menu.ExtraIcon,
				"enter_transition": menu.EnterTransition,
				"leave_transition": menu.LeaveTransition,
				"active_path":      menu.ActivePath,
				"auths":            menu.Auths,
				"frame_src":        menu.FrameSrc,
				"frame_loading":    menu.FrameLoading,
				"keep_alive":       menu.KeepAlive,
				"hidden_tag":       menu.HiddenTag,
				"fixed_tag":        menu.FixedTag,
				"show_link":        menu.ShowLink,
				"show_parent":      menu.ShowParent,
			}
			if e := tx.Table(basicModel.TableNameBasicMenu).Where("id = ?", existing.Id).Updates(updates).Error; e != nil {
				return e
			}
			menu.Id = existing.Id
		}

		// 递归子节点
		if len(menu.Children) > 0 {
			if err := rp.createChildren(ctx, tx, menu.Id, menu.Children); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rp *basicMenuRepo) BatchSave(ctx context.Context, list []*basicModel.BasicMenu) error {
	return rp.data.GetDBWithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return rp.createChildren(ctx, tx, consts.DefaultIntMinusOneValue, list)
	})
}
