package req

type BindGameAccountHouseRequest struct {
	GameAccountID int64 `json:"game_account_id" binding:"required"`
	HouseGID      int   `json:"house_gid"       binding:"required"`
	IsDefault     *bool `json:"is_default"      binding:"omitempty"`
	Status        *int  `json:"status"          binding:"omitempty,oneof=0 1"`
}
type UnbindGameAccountHouseRequest struct {
	GameAccountID int64 `json:"game_account_id" binding:"required"`
	HouseGID      int   `json:"house_gid"       binding:"required"`
}
type ListHousesByGameAccountRequest struct {
	GameAccountID int64 `json:"game_account_id" binding:"required"`
}
type ListGameAccountsByHouseRequest struct {
	HouseGID int `json:"house_gid" binding:"required"`
}
