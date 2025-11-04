package basic

import (
	"context"
	"strings"

	basicModel "battle-tiles/internal/dal/model/basic"
	"battle-tiles/internal/infra"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthRepo interface {
	ListRoleIDsByUser(ctx context.Context, userID int32) ([]int32, error)
	ListPermsByUser(ctx context.Context, userID int32) ([]string, error)
	// EnsureUserHasRoleByCode 确保用户拥有指定 code 的角色（不存在则插入，幂等）
	EnsureUserHasRoleByCode(ctx context.Context, userID int32, roleCode string) error
	// EnsureUserHasOnlyRoleByCode 确保用户仅拥有指定 code 的单一角色（会清理该用户其它角色）
	EnsureUserHasOnlyRoleByCode(ctx context.Context, userID int32, roleCode string) error
}

type authRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewAuthRepo(data *infra.Data, logger log.Logger) AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/auth")),
	}
}

func (r *authRepo) ListRoleIDsByUser(ctx context.Context, userID int32) ([]int32, error) {
	db := r.data.GetDBWithContext(ctx)
	var rows []basicModel.BasicUserRoleRel
	if err := db.Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]int32, 0, len(rows))
	for _, v := range rows {
		out = append(out, v.RoleID)
	}
	return out, nil
}

// 根据用户 -> 角色 -> 菜单，汇总菜单上的 auths（逗号分隔），去重返回
func (r *authRepo) ListPermsByUser(ctx context.Context, userID int32) ([]string, error) {
	db := r.data.GetDBWithContext(ctx)

	// 取角色
	roleIDs, err := r.ListRoleIDsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(roleIDs) == 0 {
		return []string{}, nil
	}

	// 角色 -> 菜单 -> 菜单.auths
	type row struct {
		Auths string
	}
	var rows []row
	err = db.Table("basic_role_menu_rel AS rmr").
		Select("m.auths").
		Joins("JOIN basic_menu AS m ON m.id = rmr.menu_id").
		Where("rmr.role_id IN ?", roleIDs).
		Where("COALESCE(m.auths,'') <> ''").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	// 拆分/去重
	set := map[string]struct{}{}
	for _, one := range rows {
		for _, p := range strings.Split(one.Auths, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				set[p] = struct{}{}
			}
		}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	return out, nil
}

// EnsureUserHasRoleByCode 查询角色ID，并在 basic_user_role_rel 中为用户绑定该角色（幂等）
func (r *authRepo) EnsureUserHasRoleByCode(ctx context.Context, userID int32, roleCode string) error {
	db := r.data.GetDBWithContext(ctx)

	// 查角色ID
	type roleRow struct{ ID int32 }
	var rr roleRow
	if err := db.Table("basic_role").Select("id").
		Where("code = ? AND is_deleted = false AND enable = true", roleCode).
		First(&rr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 角色不存在时返回特定错误，便于上层记录
			return err
		}
		return err
	}

	// 绑定关系（user_id, role_id 唯一，冲突忽略）
	rel := basicModel.BasicUserRoleRel{UserID: userID, RoleID: rr.ID}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&rel).Error; err != nil {
		return err
	}
	return nil
}

// EnsureUserHasOnlyRoleByCode 先查询角色ID，然后将该用户的其它角色清理，仅保留（或插入）该角色
func (r *authRepo) EnsureUserHasOnlyRoleByCode(ctx context.Context, userID int32, roleCode string) error {
	db := r.data.GetDBWithContext(ctx)

	// 查角色ID
	type roleRow struct{ ID int32 }
	var rr roleRow
	if err := db.Table("basic_role").Select("id").
		Where("code = ? AND is_deleted = false AND enable = true", roleCode).
		First(&rr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return err
		}
		return err
	}

	// 事务：删除所有角色映射 -> 插入目标角色
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&basicModel.BasicUserRoleRel{}).Error; err != nil {
			return err
		}
		rel := basicModel.BasicUserRoleRel{UserID: userID, RoleID: rr.ID}
		if err := tx.Create(&rel).Error; err != nil {
			return err
		}
		return nil
	})
}
