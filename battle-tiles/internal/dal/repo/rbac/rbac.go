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

// 直接从数据库聚合权限码
func (s *Store) queryUserPermCodesDB(ctx context.Context, userID int32) ([]string, error) {
	db := s.data.GetDBWithContext(ctx)
	if db.Error != nil {
		return nil, db.Error
	}

	// user -> roles -> menus -> menus.auths
	type row struct{ Auths *string }
	var list []row

	// 这里使用 gorm 链式 join；你也可以写 Raw SQL
	err := db.
		Table("basic_user_role_rel AS ur").
		Select("m.auths").
		Joins("JOIN basic_role_menu_rel AS rm ON rm.role_id = ur.role_id").
		Joins("JOIN basic_menu AS m ON m.id = rm.menu_id AND m.deleted_at IS NULL").
		Where("ur.user_id = ?", userID).
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	// 展开、去重
	seen := map[string]struct{}{}
	out := make([]string, 0, len(list))
	for _, r := range list {
		if r.Auths == nil || *r.Auths == "" {
			continue
		}
		for _, code := range strings.Split(*r.Auths, ",") {
			code = strings.TrimSpace(code)
			if code == "" {
				continue
			}
			code = strings.ToLower(code) // 统一小写比较
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
