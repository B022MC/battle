package game

import (
	"fmt"

	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

// UpdateRemark 更新成员备注
// POST /api/shops/members/update-remark
func (s *GameShopMemberService) UpdateRemark(c *gin.Context) {
	var in struct {
		HouseGID     int32  `json:"house_gid" binding:"required"`      // 店铺ID
		GamePlayerID string `json:"game_player_id" binding:"required"` // 游戏玩家ID
		Remark       string `json:"remark"`                            // 备注内容（可为空）
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	// 转换 GamePlayerID 为 int32
	var gameID int32
	if _, err := fmt.Sscanf(in.GamePlayerID, "%d", &gameID); err != nil {
		response.Fail(c, ecode.ParamsFailed, "invalid game_player_id")
		return
	}

	// 更新备注
	if err := s.gameMember.UpdateRemark(c.Request.Context(), in.HouseGID, gameID, in.Remark); err != nil {
		response.Fail(c, ecode.Failed, fmt.Sprintf("更新备注失败: %v", err))
		return
	}

	response.SuccessWithOK(c)
}
