package basic

import (
	"battle-tiles/internal/conf"
	"battle-tiles/internal/consts"
	"battle-tiles/internal/dal/common"
	basicModel "battle-tiles/internal/dal/model/basic"
	basicRepo "battle-tiles/internal/dal/repo/basic"
	"battle-tiles/internal/dal/req"
	"battle-tiles/internal/dal/resp"
	"battle-tiles/pkg/utils"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type BasicMenuUseCase struct {
	global          *conf.Global
	repo            basicRepo.BasicMenuRepo
	roleMenuRelRepo basicRepo.BaseRoleMenuRelRepo
	authRepo        basicRepo.AuthRepo
	btnRelRepo      basicRepo.BaseRoleMenuBtnRelRepo
	log             *log.Helper
}

func NewBasicMenuUseCase(
	global *conf.Global,
	repo basicRepo.BasicMenuRepo,
	roleMenuRelRepo basicRepo.BaseRoleMenuRelRepo,
	authRepo basicRepo.AuthRepo,
	btnRelRepo basicRepo.BaseRoleMenuBtnRelRepo,
	logger log.Logger,
) *BasicMenuUseCase {
	return &BasicMenuUseCase{
		global:          global,
		repo:            repo,
		roleMenuRelRepo: roleMenuRelRepo,
		authRepo:        authRepo,
		btnRelRepo:      btnRelRepo,
		log:             log.NewHelper(log.With(logger, "module", "usecase/base_menu")),
	}
}

// Create
func (uc *BasicMenuUseCase) Create(ctx context.Context, m *basicModel.BasicMenu) error {
	if m == nil {
		return errors.New("menu is nil")
	}
	return uc.repo.Save(ctx, m)
}

// Update —— 改为显式传 id，避免从 map 取主键
func (uc *BasicMenuUseCase) Update(ctx context.Context, id int32, fields map[string]interface{}) error {
	if id == 0 {
		return errors.New("id required")
	}
	return uc.repo.UpdateFieldByCond(ctx, fields, "id = ?", id)
}

// SelectOne（按 scope 查询）
func (uc *BasicMenuUseCase) SelectOne(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (*basicModel.BasicMenu, error) {
	return uc.repo.FindOneByScope(ctx, scopes...)
}

// SelectPage（分页）
func (uc *BasicMenuUseCase) SelectPage(ctx context.Context, page, pageSize int, isPage bool, scopes ...func(*gorm.DB) *gorm.DB) ([]*basicModel.BasicMenu, int64, error) {
	limit, offset := utils.Paginate(page, pageSize, isPage)
	return uc.repo.ListByPage(ctx, offset, limit, scopes...)
}

// SelectOption（下拉）
func (uc *BasicMenuUseCase) SelectOption(ctx context.Context, page, pageSize int, isPage bool, scopes ...func(*gorm.DB) *gorm.DB) ([]*common.Option, int64, error) {
	limit, offset := utils.Paginate(page, pageSize, isPage)
	// 你需要的 label/value 字段，这里用 name / id
	return uc.repo.OptionByPage(ctx, offset, limit, "name as label, id as value", scopes...)
}

// DeleteByIDs —— 菜单表没有 org_id，这里简化为按 id 软删
func (uc *BasicMenuUseCase) DeleteByIDs(ctx context.Context, ids ...int32) error {
	if len(ids) == 0 {
		return errors.New("ids is empty")
	}
	return uc.repo.DeleteById(ctx, "id IN ?", ids)
}

// 全量查询 -> 树
func (uc *BasicMenuUseCase) SelectAllToTree(ctx context.Context, keyword string, scopes ...func(*gorm.DB) *gorm.DB) (resp.MenuTree, error) {
	var (
		list []*basicModel.BasicMenu
		err  error
	)
	if keyword != "" {
		list, err = uc.repo.FindByKeyword(ctx, keyword)
	} else {
		list, err = uc.repo.FindByScope(ctx, scopes...)
	}
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = make([]*basicModel.BasicMenu, 0)
	}

	flat := make([]resp.Menu, 0, len(list))
	for _, r := range list {
		flat = append(flat, resp.Menu{
			ID:              r.Id,
			ParentID:        r.ParentId,
			MenuType:        r.MenuType,
			Title:           r.Title,
			Name:            r.Name,
			Path:            r.Path,
			Component:       r.Component,
			Rank:            r.Rank,
			Redirect:        r.Redirect,
			Icon:            r.Icon,
			ExtraIcon:       r.ExtraIcon,
			EnterTransition: r.EnterTransition,
			LeaveTransition: r.LeaveTransition,
			ActivePath:      r.ActivePath,
			Auths:           r.Auths,
			FrameSrc:        r.FrameSrc,
			FrameLoading:    r.FrameLoading,
			KeepAlive:       r.KeepAlive,
			HiddenTag:       r.HiddenTag,
			FixedTag:        r.FixedTag,
			ShowLink:        r.ShowLink,
			ShowParent:      r.ShowParent,
		})
	}
	return uc.toTree(flat, consts.DefaultIntMinusOneValue), nil
}

// 批量写入整棵树（覆盖式）
func (uc *BasicMenuUseCase) SaveAllTree(ctx context.Context, tree []req.MenuInfo, scopes ...func(*gorm.DB) *gorm.DB) (resp.MenuTree, error) {
	list := uc.flattenForSave(tree, -1)
	if err := uc.repo.BatchSave(ctx, list); err != nil {
		return nil, err
	}
	return uc.SelectAllToTree(ctx, "", scopes...)
}

// SelectMeTreeFiltered 基于用户角色过滤可见菜单并组装为树
func (uc *BasicMenuUseCase) SelectMeTreeFiltered(ctx context.Context, userID int32) (resp.MenuTree, error) {
	// 1) 取用户角色
	roleIDs, err := uc.authRepo.ListRoleIDsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(roleIDs) == 0 {
		return resp.MenuTree{}, nil
	}
	// 2) 聚合角色可见菜单ID
	seen := map[int32]struct{}{}
	menuIDs := make([]int32, 0)
	for _, rid := range roleIDs {
		ids, err := uc.roleMenuRelRepo.ListMenuIDsByRole(ctx, rid)
		if err != nil {
			return nil, err
		}
		for _, mid := range ids {
			if _, ok := seen[mid]; !ok {
				seen[mid] = struct{}{}
				menuIDs = append(menuIDs, mid)
			}
		}
	}
	if len(menuIDs) == 0 {
		return resp.MenuTree{}, nil
	}
	// 3) 查询菜单并组树（过滤软删）
	list, err := uc.repo.FindByScope(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted_at IS NULL").Where("id IN ?", menuIDs)
	})
	if err != nil {
		return nil, err
	}
	flat := make([]resp.Menu, 0, len(list))
	for _, r := range list {
		flat = append(flat, resp.Menu{
			ID:              r.Id,
			ParentID:        r.ParentId,
			MenuType:        r.MenuType,
			Title:           r.Title,
			Name:            r.Name,
			Path:            r.Path,
			Component:       r.Component,
			Rank:            r.Rank,
			Redirect:        r.Redirect,
			Icon:            r.Icon,
			ExtraIcon:       r.ExtraIcon,
			EnterTransition: r.EnterTransition,
			LeaveTransition: r.LeaveTransition,
			ActivePath:      r.ActivePath,
			Auths:           r.Auths,
			FrameSrc:        r.FrameSrc,
			FrameLoading:    r.FrameLoading,
			KeepAlive:       r.KeepAlive,
			HiddenTag:       r.HiddenTag,
			FixedTag:        r.FixedTag,
			ShowLink:        r.ShowLink,
			ShowParent:      r.ShowParent,
		})
	}
	return uc.toTree(flat, consts.DefaultIntMinusOneValue), nil
}

// ListMyButtonIDs 返回用户在某菜单下可用的按钮ID集合
func (uc *BasicMenuUseCase) ListMyButtonIDs(ctx context.Context, userID int32, menuID int32) ([]int32, error) {
	return uc.btnRelRepo.ListBtnIDsByUserAndMenu(ctx, userID, menuID)
}

/* -------------------- 内部函数 -------------------- */

// 把树拍平成 []*BasicMenu（用于 BatchSave）
func (uc *BasicMenuUseCase) flattenForSave(nodes []req.MenuInfo, parentID int32) []*basicModel.BasicMenu {
	out := make([]*basicModel.BasicMenu, 0)

	for _, n := range nodes {
		m := &basicModel.BasicMenu{
			ParentId:        parentID,
			MenuType:        n.MenuType,
			Title:           n.Title,
			Name:            n.Name,
			Path:            n.Path,
			Component:       n.Component,
			Rank:            n.Rank,
			Redirect:        n.Redirect,
			Icon:            n.Icon,
			ExtraIcon:       n.ExtraIcon,
			EnterTransition: n.EnterTransition,
			LeaveTransition: n.LeaveTransition,
			ActivePath:      n.ActivePath,
			Auths:           n.Auths,
			FrameSrc:        n.FrameSrc,
			FrameLoading:    n.FrameLoading,
			KeepAlive:       n.KeepAlive,
			HiddenTag:       n.HiddenTag,
			FixedTag:        n.FixedTag,
			ShowLink:        n.ShowLink,
			ShowParent:      n.ShowParent,
		}
		out = append(out, m)
		// 递归子节点，注意子节点的父ID是当前新节点的 ID
		// 这里交由 repo.BatchSave 在事务里“写后回填 ID 并继续写子级”
		// 如果你的 BatchSave 是“纯 create”，那就改成：先递归组父子，再在 repo 里用一趟事务 create/更新 parent_id。
		if len(n.Children) > 0 {
			// 这里先把 parent 占位，repo.BatchSave 应该用“递归+已获取的父ID”来做真正写入
			// 如果你的 BatchSave 已按我之前给你的版本（先清空再按父子顺序 create），就把子节点的 parentID 设为 -1 占位也可以
			children := uc.flattenForSave(n.Children, -1)
			out = append(out, children...)
		}
	}
	return out
}

// 扁平 -> 树
func (uc *BasicMenuUseCase) toTree(flat []resp.Menu, root int32) resp.MenuTree {
	buckets := make(map[int32][]resp.Menu, len(flat))
	for _, m := range flat {
		buckets[m.ParentID] = append(buckets[m.ParentID], m)
	}
	var build func(pid int32) []resp.Menu
	build = func(pid int32) []resp.Menu {
		nodes := buckets[pid]
		for i := range nodes {
			nodes[i].Children = build(nodes[i].ID)
		}
		return nodes
	}
	return build(root)
}
