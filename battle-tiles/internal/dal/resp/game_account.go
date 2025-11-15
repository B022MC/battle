package resp

type AccountVO struct {
	ID        int32  `json:"id"`
	Account   string `json:"account"`
	Nickname  string `json:"nickname"`
	IsDefault bool   `json:"is_default"`
	Status    int32  `json:"status"`
	LoginMode string `json:"login_mode"`
}

type GameAccountHouseVO struct {
	ID        int32 `json:"id"`
	HouseGID  int32 `json:"house_gid"`
	IsDefault bool  `json:"is_default"`
	Status    int32 `json:"status"`
}
