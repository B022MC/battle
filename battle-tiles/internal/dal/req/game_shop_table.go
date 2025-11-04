package req

// ListTablesRequest 请求房间列表
// @Description 列出当前在线会话的房间快照
// @example {"house_gid":20001}
type ListTablesRequest struct {
	// 店铺号（圈ID/HouseGID）
	HouseGID int `json:"house_gid" binding:"required"`
}

// DismissTableRequest 解散房间（kind_id 可选；若会话缓存可推断则无需传）
// @example {"house_gid":20001, "mapped_num":123456, "kind_id":201}
type DismissTableRequest struct {
	// 店铺号（圈ID/HouseGID）
	HouseGID int `json:"house_gid" binding:"required"`
	// 桌子映射号（MappedNum）
	MappedNum int `json:"mapped_num" binding:"required"`
	// 游戏 KindID，可选；不传时从缓存推断
	KindID int `json:"kind_id"`
}

// QueryTableRequest 查询单桌（触发 QueryTable）
// @example {"house_gid":20001, "mapped_num":123456}
type QueryTableRequest struct {
	// 店铺号（圈ID/HouseGID）
	HouseGID int `json:"house_gid" binding:"required"`
	// 桌子映射号（MappedNum）
	MappedNum int `json:"mapped_num" binding:"required"`
}
