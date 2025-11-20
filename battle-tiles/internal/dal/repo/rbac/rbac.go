package rbac

import (
	"battle-tiles/internal/infra"
	pdb "battle-tiles/pkg/plugin/dbx"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type Store struct {
	data    *infra.Data
	logger  *log.Helper
	ttl     time.Duration
	redisNS string
}

func NewStore(data *infra.Data, logger log.Logger) *Store {
	return &Store{
		data:    data,
		logger:  log.NewHelper(log.With(logger, "module", "infra/rbac/store")),
		ttl:     10 * time.Minute,
		redisNS: "rbac:perm",
	}
}

func (s *Store) cacheKey(ctx context.Context, userID int32) string {
	dbKey := ""
	if v := ctx.Value(pdb.CtxDBKey); v != nil {
		if str, ok := v.(string); ok {
			dbKey = str
		}
	}
	return fmt.Sprintf("%s:%s:%d", s.redisNS, dbKey, userID)
}

// GetUserPermCodes 从缓存读取，缓存 miss 则查库并回填
func (s *Store) GetUserPermCodes(ctx context.Context, userID int32) (map[string]struct{}, error) {
	if s.data.RDB != nil {
		key := s.cacheKey(ctx, userID)
		if bs, err := s.data.RDB.Get(ctx, key).Bytes(); err == nil && len(bs) > 0 {
			var arr []string
			if json.Unmarshal(bs, &arr) == nil {
				return toSet(arr), nil
			}
		}
	}

	codes, err := s.queryUserPermCodesDB(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 回填缓存
	if s.data.RDB != nil {
		key := s.cacheKey(ctx, userID)
		_ = s.data.RDB.Set(ctx, key, mustJSON(codes), s.ttl).Err()
	}
	return toSet(codes), nil
}

// ClearUserCache 清除用户权限缓存
func (s *Store) ClearUserCache(ctx context.Context, userID int32) error {
	if s.data.RDB != nil {
		key := s.cacheKey(ctx, userID)
		return s.data.RDB.Del(ctx, key).Err()
	}
	return nil
}

// 直接从数据库聚合权限码（同时从菜单auths和权限表获取）
func (s *Store) queryUserPermCodesDB(ctx context.Context, userID int32) ([]string, error) {
	db := s.data.GetDBWithContext(ctx)
	if db.Error != nil {
		return nil, db.Error
	}

	seen := map[string]struct{}{}
	out := make([]string, 0)

	// 1. 从菜单的 auths 字段获取权限（兼容旧系统）
	type menuRow struct{ Auths *string }
	var menuList []menuRow
	err := db.
		Table("basic_user_role_rel AS ur").
		Select("m.auths").
		Joins("JOIN basic_role_menu_rel AS rm ON rm.role_id = ur.role_id").
		Joins("JOIN basic_menu AS m ON m.id = rm.menu_id AND m.deleted_at IS NULL").
		Where("ur.user_id = ?", userID).
		Find(&menuList).Error
	if err != nil {
		return nil, err
	}

	for _, r := range menuList {
		if r.Auths == nil || *r.Auths == "" {
			continue
		}
		for _, code := range strings.Split(*r.Auths, ",") {
			code = strings.TrimSpace(code)
			if code == "" {
				continue
			}
			code = strings.ToLower(code)
			if _, ok := seen[code]; !ok {
				seen[code] = struct{}{}
				out = append(out, code)
			}
		}
	}

	// 2. 从权限表获取权限（新系统）
	type permRow struct{ Code string }
	var permList []permRow
	err = db.
		Table("basic_user_role_rel AS ur").
		Select("DISTINCT p.code").
		Joins("JOIN basic_role_permission_rel AS rpr ON rpr.role_id = ur.role_id").
		Joins("JOIN basic_permission AS p ON p.id = rpr.permission_id AND p.is_deleted = false").
		Where("ur.user_id = ?", userID).
		Find(&permList).Error
	if err != nil {
		// 如果权限表不存在（旧版本数据库），忽略错误
		s.logger.Warnf("query permission table failed: %v", err)
	} else {
		for _, r := range permList {
			code := strings.ToLower(strings.TrimSpace(r.Code))
			if code == "" {
				continue
			}
			if _, ok := seen[code]; !ok {
				seen[code] = struct{}{}
				out = append(out, code)
			}
		}
	}

	return out, nil
}

// 小工具
func toSet(arr []string) map[string]struct{} {
	s := make(map[string]struct{}, len(arr))
	for _, v := range arr {
		s[strings.ToLower(strings.TrimSpace(v))] = struct{}{}
	}
	return s
}
func mustJSON(v any) []byte { b, _ := json.Marshal(v); return b }
