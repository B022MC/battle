package middleware

import "context"

// PermissionStore：抽象权限来源
type PermissionStore interface {
	// 返回用户拥有的权限码集合（去重后的 code set）
	GetUserPermCodes(ctx context.Context, userID int32) (map[string]struct{}, error)
}

// 全局绑定（通过 Init 在应用启动时注入）
var globalStore PermissionStore

func BindPermissionStore(s PermissionStore) { globalStore = s }
func store() PermissionStore                { return globalStore }
