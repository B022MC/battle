package resp

type ShopAdminVO struct {
	ID       int32  `json:"id"`
	HouseGID int32  `json:"house_gid"`
	UserID   int32  `json:"user_id"`
	Role     string `json:"role"     binding:"omitempty,oneof=admin operator"` // 仅 admin|operator，留空=admin
	NickName string `json:"nick_name,omitempty"`
}
