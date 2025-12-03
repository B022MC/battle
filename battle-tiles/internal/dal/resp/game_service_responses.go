package resp

type SessionStateResponse struct {
	// 会话状态字符串，如 online/offline/error
	State string `json:"state"`
}

// BattlePlayerVO 战绩中的单个玩家
// @description 战绩玩家及分数
// @example {"game_id":123456,"score":-20,"name":"玩家A"}
type BattlePlayerVO struct {
	// 用户游戏ID（GameID）
	GameID int `json:"game_id"`
	// 本局分数（正负）
	Score int `json:"score"`
	// 玩家昵称（可选）
	Name string `json:"name,omitempty"`
}

// BattleRecordVO 战绩记录
// @example {"room_id":7890,"kind_id":201,"base_score":5,"time":1731234567,"players":[{"game_id":1,"score":10}]}
type BattleRecordVO struct {
	// 房间号（RoomID）
	RoomID int `json:"room_id"`
	// 游戏KindID
	KindID int `json:"kind_id"`
	// 底分
	BaseScore int `json:"base_score"`
	// 创建时间（秒）
	Time int `json:"time"`
	// 参与玩家
	Players []BattlePlayerVO `json:"players"`
}

type FundsBalanceResponse struct {
	// 余额（分）
	Balance int32 `json:"balance"`
}

type FundsLimitResponse struct {
	// 余额（分）
	Balance int32 `json:"balance"`
	// 是否禁用
	Forbid bool `json:"forbid"`
	// 个性额度下限（分）
	LimitMin int32 `json:"limit_min"`
}

type VerifyAccountResponse struct {
	// 探活是否成功
	Ok bool `json:"ok"`
}

type ShopMemberListItem struct {
	UserID     uint32 `json:"user_id"`
	UserStatus int    `json:"user_status"`
	GameID     uint32 `json:"game_id"`
	MemberID   uint32 `json:"member_id"`
	MemberType int    `json:"member_type"`
	NickName   string `json:"nick_name"`
	// GroupID 若协议侧提供则填充；当前为 0（未知/未分组）
	GroupID int `json:"group_id"`
	// GroupName 圈子名称（圈主昵称）
	GroupName *string `json:"group_name,omitempty"`

	// 平台用户信息（关联信息）
	PlatformUser   *UserInfo `json:"platform_user,omitempty"`   // 绑定的平台用户信息
	GameAccountID  *uint32   `json:"game_account_id,omitempty"` // 游戏账号ID
	IsBindPlatform bool      `json:"is_bind_platform"`          // 是否已绑定平台用户

	// 拉圈功能字段
	GamePlayerID     string `json:"game_player_id,omitempty"`     // 游戏玩家ID（用于拉圈，等同于 game_id）
	CurrentGroupID   *int32 `json:"current_group_id,omitempty"`   // 当前所在圈子ID
	CurrentGroupName string `json:"current_group_name,omitempty"` // 当前所在圈子名称

	// 置顶和排序字段
	IsPinned bool  `json:"is_pinned"` // 是否置顶
	PinOrder int32 `json:"pin_order"` // 置顶排序（数字越小越靠前）

	// 管理员备注
	Remark string `json:"remark"` // 管理员对该成员的备注

	// 禁用状态
	Forbid bool `json:"forbid"` // 是否禁用
}

// UserInfo 用户信息响应（过滤敏感字段）
type UserInfo struct {
	ID           int32  `json:"id"`
	Username     string `json:"username"`
	NickName     string `json:"nick_name"`
	Avatar       string `json:"avatar"`
	Role         string `json:"role"`
	Introduction string `json:"introduction"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Items []UserInfo `json:"items"`
	Total int64      `json:"total"`
	Page  int32      `json:"page"`
	Size  int32      `json:"size"`
}

type ShopMemberListResponse struct {
	// 成员列表
	Items []ShopMemberListItem `json:"items"`
}

type DiamondQueryResponse struct {
	// 是否已触发查询
	Triggered bool `json:"triggered"`
}

type ShopTableListResponse struct {
	// 房间列表
	Items []TableInfoVO `json:"items"`
}

type ShopTableCheckResponse struct {
	// 是否触发了底层查询
	Triggered bool `json:"triggered"`
	// 是否在缓存中存在
	ExistsInCache bool `json:"exists_in_cache"`
	// 桌信息
	Table *TableInfoVO `json:"table,omitempty"`
}

type ShopTableDetailResponse struct {
	// 桌信息
	Table *TableInfoVO `json:"table"`
	// 是否触发刷新
	Triggered bool `json:"triggered"`
}

// —— 内联 TableInfoVO，避免外部包引用导致 swagger 解析失败 ——
type TableInfoVO struct {
	// 桌ID
	TableID int `json:"table_id"`
	// 映射号
	MappedNum int `json:"mapped_num"`
	// 圈ID
	GroupID int `json:"group_id"`
	// KindID
	KindID int `json:"kind_id"`
	// 底分
	BaseScore int `json:"base_score"`
}

type WalletListResponse struct {
	List     []WalletVO `json:"list"`
	Total    int64      `json:"total"`
	Page     int32      `json:"page"`
	PageSize int32      `json:"page_size"`
}

type LedgerListResponse struct {
	List     []LedgerVO `json:"list"`
	Total    int64      `json:"total"`
	Page     int32      `json:"page"`
	PageSize int32      `json:"page_size"`
}
