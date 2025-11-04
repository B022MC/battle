package req

// ListApplicationsRequest 请求：按店铺号列出管理员申请
type ListApplicationsRequest struct {
	HouseGID int  `json:"house_gid" binding:"required"`      // 店铺号（HouseGID）
	Type     *int `json:"type" binding:"omitempty"`          // 可选：过滤类型（admin/join 的枚举码）
	AdminUID *int `json:"admin_user_id" binding:"omitempty"` // 可选：过滤指定圈主/管理员
}

// DecideApplicationRequest 请求：审批指定申请（同意/拒绝）
type DecideApplicationRequest struct {
	ID int `json:"id" binding:"required"` // 申请消息ID
}

// ApplyAdminRequest 发起管理员申请
type ApplyAdminRequest struct {
	HouseGID int32  `json:"house_gid" binding:"required"`
	Note     string `json:"note"`
}

// ApplyJoinGroupRequest 发起加入圈子申请
type ApplyJoinGroupRequest struct {
	HouseGID    int32  `json:"house_gid" binding:"required"`
	AdminUserID int32  `json:"admin_user_id" binding:"required"`
	Note        string `json:"note"`
}
