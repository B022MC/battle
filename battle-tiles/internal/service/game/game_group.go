// internal/service/game/game_group.go
package game

import (
	"battle-tiles/internal/dal/req"
	"battle-tiles/internal/infra/plaza"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

type GameGroupService struct{ mgr plaza.Manager }

func NewGameGroupService(mgr plaza.Manager) *GameGroupService { return &GameGroupService{mgr: mgr} }

func (s *GameGroupService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/shops").Use(middleware.JWTAuth())

	g.POST("/groups/forbid", middleware.RequirePerm("group:forbid"), s.Forbid)
	g.POST("/groups/unforbid", middleware.RequirePerm("group:forbid"), s.Unforbid)
	g.POST("/groups/delete", middleware.RequirePerm("group:delete"), s.Delete)

	g.POST("/groups/bind", middleware.RequirePerm("group:bind"), s.Bind)
	g.POST("/groups/unbind", middleware.RequirePerm("group:unbind"), s.Unbind)
}

// Forbid 禁圈（可选指定成员；默认对当前圈全部成员执行，排除自身）
// @Summary      禁圈
// @Tags         圈/群
// @Accept       json
// @Produce      json
// @Param        in body req.GroupForbidRequest true "house_gid, key, member_ids 可选"
// @Success      200 {object} response.Body
// @Router       /shops/groups/forbid [post]
func (s *GameGroupService) Forbid(c *gin.Context) {
	s.doForbid(c, true)
}

// Unforbid 解圈
// @Summary      解圈
// @Tags         圈/群
// @Accept       json
// @Produce      json
// @Param        in body req.GroupForbidRequest true "house_gid, key, member_ids 可选"
// @Success      200 {object} response.Body
// @Router       /shops/groups/unforbid [post]
func (s *GameGroupService) Unforbid(c *gin.Context) { s.doForbid(c, false) }

func (s *GameGroupService) doForbid(c *gin.Context, forbid bool) {
	var in req.GroupForbidRequest
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 || in.Key == "" {
		response.Fail(c, ecode.ParamsFailed, "invalid house_gid or key")
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	actor := int(claims.BaseClaims.UserID)

	// 收集成员
	members := in.MemberIDs
	if len(members) == 0 {
		// 从快照读取全量成员
		if sess, ok := s.mgr.Get(actor, in.HouseGID); ok && sess != nil {
			if ms := sess.ListMembers(); len(ms) > 0 {
				for _, m := range ms {
					if int(m.UserID) == actor { // 跳过自身
						continue
					}
					members = append(members, int(m.MemberID))
				}
			}
		}
	}
	if len(members) == 0 {
		response.SuccessWithOK(c)
		return
	}

	if err := s.mgr.ForbidMembers(actor, in.HouseGID, in.Key, members, forbid); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}

// Delete 删圈（等价：批量踢出所有成员）
// @Summary      删圈
// @Description  若无直通删除圈协议，则逐个踢出成员以达到删除效果
// @Tags         圈/群
// @Accept       json
// @Produce      json
// @Param        in body req.GroupBaseRequest true "house_gid"
// @Success      200 {object} response.Body
// @Router       /shops/groups/delete [post]
func (s *GameGroupService) Delete(c *gin.Context) {
	var in req.GroupBaseRequest
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid house_gid")
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	actor := int(claims.BaseClaims.UserID)
	sess, ok := s.mgr.Get(actor, in.HouseGID)
	if !ok || sess == nil {
		response.Fail(c, ecode.Failed, "session not found or not online for this house")
		return
	}
	ms := sess.ListMembers()
	for _, m := range ms {
		_ = s.mgr.KickMember(actor, in.HouseGID, int(m.MemberID))
	}
	response.SuccessWithOK(c)
}

// Bind 绑圈（同意申请）
// @Summary      绑圈（同意申请）
// @Tags         圈/群
// @Accept       json
// @Produce      json
// @Param        in body req.GroupBindRequest true "house_gid, message_id"
// @Success      200 {object} response.Body
// @Router       /shops/groups/bind [post]
func (s *GameGroupService) Bind(c *gin.Context) {
	var in req.GroupBindRequest
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 || in.MessageID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid house_gid or message_id")
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	actor := int(claims.BaseClaims.UserID)
	sess, ok := s.mgr.Get(actor, in.HouseGID)
	if !ok || sess == nil {
		response.Fail(c, ecode.Failed, "session not found or not online for this house")
		return
	}
	ai, ok2 := sess.FindApplicationByID(in.MessageID)
	if !ok2 || ai == nil {
		response.Fail(c, ecode.Failed, "application message not found")
		return
	}
	sess.RespondApplication(ai, true)
	response.SuccessWithOK(c)
}

// Unbind 退圈（停止该店铺会话）
// @Summary      退圈
// @Description  无直通退圈协议时，等价：停止当前店铺会话
// @Tags         圈/群
// @Accept       json
// @Produce      json
// @Param        in body req.GroupBaseRequest true "house_gid"
// @Success      200 {object} response.Body
// @Router       /shops/groups/unbind [post]
func (s *GameGroupService) Unbind(c *gin.Context) {
	var in req.GroupBaseRequest
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid house_gid")
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	actor := int(claims.BaseClaims.UserID)
	s.mgr.StopUser(actor, in.HouseGID)
	response.SuccessWithOK(c)
}
