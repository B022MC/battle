package basic

import (
	"context"

	"battle-tiles/internal/infra"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type BaseRoleMenuBtnRelRepo interface {
	// 返回用户在指定菜单下可用的按钮ID集合（聚合其角色）
	ListBtnIDsByUserAndMenu(ctx context.Context, userID int32, menuID int32) ([]int32, error)
}

type baseRoleMenuBtnRelRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewBaseRoleMenuBtnRelRepo(data *infra.Data, logger log.Logger) BaseRoleMenuBtnRelRepo {
	return &baseRoleMenuBtnRelRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/basicRoleMenuBtnRel")),
	}
}

func (r *baseRoleMenuBtnRelRepo) db(ctx context.Context) *gorm.DB {
	return r.data.GetDBWithContext(ctx)
}

func (r *baseRoleMenuBtnRelRepo) ListBtnIDsByUserAndMenu(ctx context.Context, userID int32, menuID int32) ([]int32, error) {
	// basic_user_role_rel (ur) -> basic_role_menu_btn_rel (rmb)
	type row struct{ BtnID int32 }
	var rows []row
	if err := r.db(ctx).
		Table("basic_user_role_rel AS ur").
		Select("rmb.btn_id AS btn_id").
		Joins("JOIN basic_role_menu_btn_rel AS rmb ON rmb.role_id = ur.role_id").
		Where("ur.user_id = ? AND rmb.menu_id = ?", userID, menuID).
		Distinct().
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]int32, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.BtnID)
	}
	return out, nil
}
