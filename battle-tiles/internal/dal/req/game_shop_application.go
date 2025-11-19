package req

// ============ 游戏内申请功能（从 Plaza Session 内存读取）============

// ListGameApplicationsRequest 请求：列出游戏内待处理申请（从内存读取）
type ListGameApplicationsRequest struct {
	HouseGID int32 `json:"house_gid" binding:"required"` // 店铺游戏ID
}

// RespondGameApplicationRequest 请求：处理游戏内申请（通过/拒绝）
type RespondGameApplicationRequest struct {
	HouseGID  int32 `json:"house_gid" binding:"required"`  // 店铺游戏ID
	MessageID int   `json:"message_id" binding:"required"` // 游戏消息ID
}

// ============ 旧的管理员申请功能（已废弃，保留用于兼容）============

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
