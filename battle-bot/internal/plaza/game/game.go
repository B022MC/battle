package game

// —— DTO —— //
type BattleSettle struct {
	UserGameID int
	Score      int
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
