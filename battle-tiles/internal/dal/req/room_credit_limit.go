package req

// SetRoomCreditLimitRequest 设置房间额度限制请求
type SetRoomCreditLimitRequest struct {
	HouseGID    int32  `json:"house_gid" binding:"required,gt=0"`     // 店铺GID
	GroupName   string `json:"group_name" binding:"omitempty,max=64"` // 圈子名称，空表示全局
	GameKind    int32  `json:"game_kind" binding:"omitempty,gte=0"`   // 游戏类型，0表示默认
	BaseScore   int32  `json:"base_score" binding:"omitempty,gte=0"`  // 底分，0表示默认
	CreditLimit int32  `json:"credit_limit" binding:"required,gte=0"` // 额度限制（分）
}

// GetRoomCreditLimitRequest 查询房间额度限制请求
type GetRoomCreditLimitRequest struct {
	HouseGID  int32  `json:"house_gid" binding:"required,gt=0"`     // 店铺GID
	GroupName string `json:"group_name" binding:"omitempty,max=64"` // 圈子名称
	GameKind  int32  `json:"game_kind" binding:"omitempty,gte=0"`   // 游戏类型
	BaseScore int32  `json:"base_score" binding:"omitempty,gte=0"`  // 底分
}

// ListRoomCreditLimitRequest 列出房间额度限制请求
type ListRoomCreditLimitRequest struct {
	HouseGID  int32  `json:"house_gid" binding:"required,gt=0"`     // 店铺GID
	GroupName string `json:"group_name" binding:"omitempty,max=64"` // 圈子名称，空表示查询所有
}

// DeleteRoomCreditLimitRequest 删除房间额度限制请求
type DeleteRoomCreditLimitRequest struct {
	HouseGID  int32  `json:"house_gid" binding:"required,gt=0"`     // 店铺GID
	GroupName string `json:"group_name" binding:"omitempty,max=64"` // 圈子名称
	GameKind  int32  `json:"game_kind" binding:"omitempty,gte=0"`   // 游戏类型
	BaseScore int32  `json:"base_score" binding:"omitempty,gte=0"`  // 底分
}

// CheckPlayerCreditRequest 检查玩家是否满足房间额度要求
type CheckPlayerCreditRequest struct {
	HouseGID  int32  `json:"house_gid" binding:"required,gt=0"`     // 店铺GID
	GameID    int32  `json:"game_id" binding:"required,gt=0"`       // 玩家游戏ID
	GroupName string `json:"group_name" binding:"omitempty,max=64"` // 圈子名称
	GameKind  int32  `json:"game_kind" binding:"required,gt=0"`     // 游戏类型
	BaseScore int32  `json:"base_score" binding:"required,gt=0"`    // 底分
}
