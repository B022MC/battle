package req

// 踢出成员
type KickMemberRequest struct {
	HouseGID int `json:"house_gid" binding:"required"`
	MemberID int `json:"member_id" binding:"required"`
}
