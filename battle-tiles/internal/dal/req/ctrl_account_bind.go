package req

// 创建/更新中控（不绑定店铺）
type CreateOrUpdateCtrlRequest struct {
	LoginMode  string `json:"login_mode"  binding:"required,oneof=account mobile"`
	Identifier string `json:"identifier"  binding:"required"`
	PwdMD5     string `json:"pwd_md5"     binding:"required,len=32"`
	Status     int32  `json:"status"      binding:"omitempty,oneof=0 1"`
}

// 绑定/解绑中控到店铺
type BindCtrlHouseRequest struct {
	CtrlID   int32 `json:"ctrl_id"  binding:"required"`
	HouseGID int32 `json:"house_gid" binding:"required"`
	Status   int32 `json:"status"    binding:"omitempty,oneof=0 1"` // 仅绑定时可传
}
type UnbindCtrlHouseRequest struct {
	CtrlID   int32 `json:"ctrl_id"  binding:"required"`
	HouseGID int32 `json:"house_gid" binding:"required"`
}

// 兼容：按店铺列中控
type CtrlAccountListRequest struct {
	HouseGID int32 `json:"house_gid" binding:"required"`
}
