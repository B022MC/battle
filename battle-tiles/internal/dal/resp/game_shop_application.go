package resp

// ApplicationItemVO 单条申请信息
type ApplicationItemVO struct {
	ID          int    `json:"id"`            // 申请消息ID
	Status      int    `json:"status"`        // 平台返回的消息状态码
	ApplierID   int    `json:"applier_id"`    // 申请者用户ID（游戏用户）
	ApplierGID  int    `json:"applier_gid"`   // 申请者圈ID
	ApplierName string `json:"applier_name"`  // 申请者昵称
	HouseGID    int    `json:"house_gid"`     // 目标店铺号
	Type        int    `json:"type"`          // 申请类型（解析自 applyType）
	AdminUserID int    `json:"admin_user_id"` // 圈主/馆主用户ID
	CreatedAt   int64  `json:"created_at"`    // 申请时间（unix ms）
}

// ApplicationsVO 申请列表响应
type ApplicationsVO struct {
	Items []*ApplicationItemVO `json:"items"`
}

// AckVO 简单确认响应
type AckVO struct {
	OK bool `json:"ok"`
}
