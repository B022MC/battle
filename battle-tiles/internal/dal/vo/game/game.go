package game

// —— DTO —— //
type BattleSettle struct {
	UserGameID int
	Score      int
	NickName   string `json:",omitempty"` // 玩家昵称（可选）
	Balance    int32  `json:"-"`          // 玩家余额（不序列化到JSON，仅内部使用）
	Credit     int32  `json:"-"`          // 玩家额度（不序列化到JSON，仅内部使用）
}
type BattleInfo struct {
	RoomID     int
	KindID     int
	CreateTime int
	BaseScore  int
	Players    []*BattleSettle
}

type HouseSessionCred struct {
	UserLogon string
	UserPwd   string
	UserID    int
	HouseGID  int
}
type UserLogonInfo struct {
	UserID uint32
	GameID uint32
}
