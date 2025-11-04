// internal/dal/req/game_shop_admin.go
package req

// 分配管理员（role 可选，缺省按 admin 处理）
type AssignShopAdminRequest struct {
	HouseGID int32  `json:"house_gid" binding:"required"`                      // 店铺(茶馆)号
	UserID   int32  `json:"user_id"   binding:"required"`                      // 被授权的用户ID
	Role     string `json:"role"     binding:"omitempty,oneof=admin operator"` // 仅 admin|operator，留空=admin
}

// 撤销管理员
type RevokeShopAdminRequest struct {
	HouseGID int32 `json:"house_gid" binding:"required"`
	UserID   int32 `json:"user_id"   binding:"required"`
}
type ListShopAdminsRequest struct {
	HouseGID int32 `json:"house_gid" binding:"required"` // 茶馆号
}
