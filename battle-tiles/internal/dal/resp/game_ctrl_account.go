package resp

import (
	model "battle-tiles/internal/dal/model/game"
	"time"
)

type CtrlAccountVO struct {
	ID         int32  `json:"id"`
	HouseGID   int32  `json:"house_gid"`
	LoginMode  string `json:"login_mode"`
	Identifier string `json:"identifier"`
	Status     int32  `json:"status"`
	//LastVerify *time.Time `json:"last_verify_at,omitempty"`
}

//type CtrlAccountVO struct {
//	ID         int32      `json:"id"`
//	LoginMode  string     `json:"login_mode"`
//	Identifier string     `json:"identifier"`
//	Status     int32      `json:"status"`
//	LastVerify *time.Time `json:"last_verify_at,omitempty"`
//}

type CtrlAccountBindVO struct {
	CtrlID   int32  `json:"ctrl_id"`
	HouseGID int32  `json:"house_gid"`
	Status   int32  `json:"status"`
	Alias    string `json:"alias"`
}
type CtrlAccountWithHouses struct {
	*model.GameCtrlAccount
	Houses []int32
}
type CtrlAccountAllVO struct {
	ID           int32      `json:"id"`
	LoginMode    string     `json:"login_mode"` // account|mobile
	Identifier   string     `json:"identifier"`
	Status       int32      `json:"status"`
	LastVerifyAt *time.Time `json:"last_verify_at,omitempty"`
	Houses       []int32    `json:"houses"` // 绑定的 house_gid 列表
}
