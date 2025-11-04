package req

// 启动会话入参
type StartSessionRequest struct {
	Id       int32 `json:"id"        binding:"required"` // game_ctrl_account 主键ID（int8）
	HouseGID int32 `json:"house_gid" binding:"required"` // 店铺(茶馆)号
}

// 停止会话入参
type StopSessionRequest struct {
	Id       int32 `json:"id"        binding:"required"` // game_ctrl_account 主键ID
	HouseGID int32 `json:"house_gid" binding:"required"` // 店铺(茶馆)号
}
