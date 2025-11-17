package basic

import (
	"context"
	"strings"

	basicModel "battle-tiles/internal/dal/model/basic"
	"battle-tiles/internal/infra"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type PermissionRepo interface {
	// 权限CRUD
	Create(ctx context.Context, permission *basicModel.BasicPermission) error
	Update(ctx context.Context, permission *basicModel.BasicPermission) error
	Delete(ctx context.Context, id int32) error
	GetByID(ctx context.Context, id int32) (*basicModel.BasicPermission, error)
	GetByCode(ctx context.Context, code string) (*basicModel.BasicPermission, error)
	List(ctx context.Context, category string) ([]*basicModel.BasicPermission, error)
	ListAll(ctx context.Context) ([]*basicModel.BasicPermission, error)

	// 角色权限关联
	AssignPermissionsToRole(ctx context.Context, roleID int32, permissionIDs []int32) error
	RemovePermissionsFromRole(ctx context.Context, roleID int32, permissionIDs []int32) error
	GetRolePermissions(ctx context.Context, roleID int32) ([]*basicModel.BasicPermission, error)
	GetRolePermissionCodes(ctx context.Context, roleID int32) ([]string, error)

	// 用户权限查询（从权限表获取，区别于菜单权限）
	GetUserPermissionsFromTable(ctx context.Context, userID int32) ([]string, error)
}

type permissionRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewPermissionRepo(data *infra.Data, logger log.Logger) PermissionRepo {
	return &permissionRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/permission")),
	}
}

func (r *permissionRepo) Create(ctx context.Context, permission *basicModel.BasicPermission) error {
	db := r.data.GetDBWithContext(ctx)
	return db.Create(permission).Error
}

func (r *permissionRepo) Update(ctx context.Context, permission *basicModel.BasicPermission) error {
	db := r.data.GetDBWithContext(ctx)
	return db.Updates(permission).Error
}

func (r *permissionRepo) Delete(ctx context.Context, id int32) error {
	db := r.data.GetDBWithContext(ctx)
	return db.Model(&basicModel.BasicPermission{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

func (r *permissionRepo) GetByID(ctx context.Context, id int32) (*basicModel.BasicPermission, error) {
	db := r.data.GetDBWithContext(ctx)
	var permission basicModel.BasicPermission
	if err := db.Where("id = ? AND is_deleted = false", id).First(&permission).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepo) GetByCode(ctx context.Context, code string) (*basicModel.BasicPermission, error) {
	db := r.data.GetDBWithContext(ctx)
	var permission basicModel.BasicPermission
	if err := db.Where("code = ? AND is_deleted = false", code).First(&permission).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepo) List(ctx context.Context, category string) ([]*basicModel.BasicPermission, error) {
	db := r.data.GetDBWithContext(ctx)
	var permissions []*basicModel.BasicPermission
	query := db.Where("is_deleted = false")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if err := query.Order("category, code").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *permissionRepo) ListAll(ctx context.Context) ([]*basicModel.BasicPermission, error) {
	return r.List(ctx, "")
}

func (r *permissionRepo) AssignPermissionsToRole(ctx context.Context, roleID int32, permissionIDs []int32) error {
	db := r.data.GetDBWithContext(ctx)

	// 构建批量插入数据
	rels := make([]*basicModel.BasicRolePermissionRel, 0, len(permissionIDs))
	for _, permID := range permissionIDs {
		rels = append(rels, &basicModel.BasicRolePermissionRel{
			RoleID:       roleID,
			PermissionID: permID,
		})
	}

	// 使用 ON CONFLICT DO NOTHING 避免重复
	return db.Clauses().Create(&rels).Error
}

func (r *permissionRepo) RemovePermissionsFromRole(ctx context.Context, roleID int32, permissionIDs []int32) error {
	db := r.data.GetDBWithContext(ctx)
	return db.Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).
		Delete(&basicModel.BasicRolePermissionRel{}).Error
}

func (r *permissionRepo) GetRolePermissions(ctx context.Context, roleID int32) ([]*basicModel.BasicPermission, error) {
	db := r.data.GetDBWithContext(ctx)
	var permissions []*basicModel.BasicPermission

	err := db.Table(basicModel.TableNameBasicPermission+" AS p").
		Select("p.*").
		Joins("JOIN "+basicModel.TableNameBasicRolePermissionRel+" AS rpr ON rpr.permission_id = p.id").
		Where("rpr.role_id = ? AND p.is_deleted = false", roleID).
		Order("p.category, p.code").
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *permissionRepo) GetRolePermissionCodes(ctx context.Context, roleID int32) ([]string, error) {
	db := r.data.GetDBWithContext(ctx)

	type row struct {
		Code string
	}
	var rows []row

	err := db.Table(basicModel.TableNameBasicPermission+" AS p").
		Select("p.code").
		Joins("JOIN "+basicModel.TableNameBasicRolePermissionRel+" AS rpr ON rpr.permission_id = p.id").
		Where("rpr.role_id = ? AND p.is_deleted = false", roleID).
		Find(&rows).Error

	if err != nil {
		return nil, err
	}

	codes := make([]string, 0, len(rows))
	for _, r := range rows {
		codes = append(codes, r.Code)
	}
	return codes, nil
}

func (r *permissionRepo) GetUserPermissionsFromTable(ctx context.Context, userID int32) ([]string, error) {
	db := r.data.GetDBWithContext(ctx)

	type row struct {
		Code string
	}
	var rows []row

	// 用户 -> 角色 -> 权限
	err := db.Table(basicModel.TableNameBasicUserRoleRel+" AS urr").
		Select("DISTINCT p.code").
		Joins("JOIN "+basicModel.TableNameBasicRolePermissionRel+" AS rpr ON rpr.role_id = urr.role_id").
		Joins("JOIN "+basicModel.TableNameBasicPermission+" AS p ON p.id = rpr.permission_id").
		Where("urr.user_id = ? AND p.is_deleted = false", userID).
		Find(&rows).Error

	if err != nil {
		return nil, err
	}

	// 去重
	set := make(map[string]struct{})
	for _, r := range rows {
		code := strings.ToLower(strings.TrimSpace(r.Code))
		if code != "" {
			set[code] = struct{}{}
		}
	}

	codes := make([]string, 0, len(set))
	for code := range set {
		codes = append(codes, code)
	}
	return codes, nil
}
