// internal/service/game/game_shop_member.go
package game

import (
	biz "battle-tiles/internal/biz/game"
	"battle-tiles/internal/dal/req"
	resp "battle-tiles/internal/dal/resp"
	"battle-tiles/internal/infra/plaza"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type GameShopMemberService struct {
	mgr  plaza.Manager
	rule *biz.MemberRuleUseCase
}

func NewGameShopMemberService(mgr plaza.Manager, rule *biz.MemberRuleUseCase) *GameShopMemberService {
	return &GameShopMemberService{mgr: mgr, rule: rule}
}

func (s *GameShopMemberService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/shops").Use(middleware.JWTAuth())
	g.POST("/members/kick", middleware.RequirePerm("shop:member:kick"), s.Kick)
	g.POST("/members/list", middleware.RequirePerm("shop:member:view"), s.List)
	g.POST("/members/logout", middleware.RequirePerm("shop:member:logout"), s.Logout)
	g.POST("/diamond/query", middleware.RequirePerm("shop:member:view"), s.QueryDiamond)
	g.POST("/members/pull", middleware.RequirePerm("shop:member:view"), s.PullMembers)
	// 成员规则
	g.POST("/members/rules/vip", middleware.RequirePerm("shop:member:update"), s.SetVIP)
	g.POST("/members/rules/multi", middleware.RequirePerm("shop:member:update"), s.SetMulti)
	g.POST("/members/rules/temp_release", middleware.RequirePerm("shop:member:update"), s.SetTempRelease)
}

// Kick
// @Summary      踢出成员
// @Tags         店铺/成员
// @Accept       json
// @Produce      json
// @Param        in body req.KickMemberRequest true "house_gid, member_id"
// @Success      200 {object} response.Body
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Failure      409 {object} response.Body "no online session"
// @Router       /shops/members/kick [post]
func (s *GameShopMemberService) Kick(c *gin.Context) {
	var in req.KickMemberRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	actorUID := int(claims.BaseClaims.UserID)

	if err := s.mgr.KickMember(actorUID, in.HouseGID, in.MemberID); err != nil {
		if strings.Contains(err.Error(), "session not found") {
			response.Fail(c, ecode.Failed, "no online session")
			return
		}
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.SuccessWithOK(c)
}

// List
// @Summary      成员列表快照
// @Description  触发 GetGroupMembers 并返回最近一次的成员快照（若无则返回空列表）
// @Tags         店铺/成员
// @Accept       json
// @Produce      json
// @Param        in body req.ListTablesRequest true "house_gid"
// @Success      200 {object} response.Body{data=resp.ShopMemberListResponse} "data: { items: [] }"
// @Router       /shops/members/list [post]
func (s *GameShopMemberService) List(c *gin.Context) {
	var in req.ListTablesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
	if !ok || sess == nil {
		if shared, ok2 := s.mgr.GetAnyByHouse(in.HouseGID); ok2 && shared != nil {
			sess = shared
		} else {
			response.Success(c, resp.ShopMemberListResponse{Items: []resp.ShopMemberListItem{}})
			return
		}
	}
	// 触发拉取
	sess.GetGroupMembers()
	// 返回快照（若无为 nil -> 空数组）
	mems := sess.ListMembers()
	// 简单映射（直接透出字段）
	out := make([]resp.ShopMemberListItem, 0, len(mems))
	for _, m := range mems {
		out = append(out, resp.ShopMemberListItem{
			UserID:     m.UserID,
			UserStatus: m.UserStatus,
			GameID:     m.GameID,
			MemberID:   m.MemberID,
			MemberType: m.MemberType,
			NickName:   m.NickName,
		})
	}
	response.Success(c, resp.ShopMemberListResponse{Items: out})
}

// QueryDiamond
// @Summary      查询当前中控账号财富（钻石）
// @Description  调用会话 GetDiamond；接口返回触发状态。
// @Tags         店铺/成员
// @Accept       json
// @Produce      json
// @Param        in body req.ListTablesRequest true "house_gid"
// @Success      200 {object} response.Body{data=resp.DiamondQueryResponse} "data: { triggered: bool }"
// @Router       /shops/diamond/query [post]
func (s *GameShopMemberService) QueryDiamond(c *gin.Context) {
	var in req.ListTablesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
	if !ok || sess == nil {
		if shared, ok2 := s.mgr.GetAnyByHouse(in.HouseGID); ok2 && shared != nil {
			sess = shared
		} else {
			response.Fail(c, ecode.Failed, "session not found or not online for this house")
			return
		}
	}
	sess.GetDiamond()
	response.Success(c, resp.DiamondQueryResponse{Triggered: true})
}

// Logout
// @Summary      强制用户下线（区别于踢出）
// @Description  使用 KickOffMember 强制下线（实现侧等价于从圈成员删除），与“踢人”动作在结果上相同。
// @Tags         店铺/成员
// @Accept       json
// @Produce      json
// @Param        in body req.KickMemberRequest true "house_gid, member_id"
// @Success      200 {object} response.Body
// @Router       /shops/members/logout [post]
func (s *GameShopMemberService) Logout(c *gin.Context) {
	var in req.KickMemberRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	if err := s.mgr.KickMember(int(claims.BaseClaims.UserID), in.HouseGID, in.MemberID); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}

// PullMembers
// @Summary      手动刷新成员列表（触发拉取）
// @Tags         店铺/成员
// @Accept       json
// @Produce      json
// @Param        in body req.ListTablesRequest true "house_gid"
// @Success      200 {object} response.Body
// @Router       /shops/members/pull [post]
func (s *GameShopMemberService) PullMembers(c *gin.Context) {
	var in req.ListTablesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
	if !ok || sess == nil {
		if shared, ok2 := s.mgr.GetAnyByHouse(in.HouseGID); ok2 && shared != nil {
			sess = shared
		} else {
			response.Fail(c, ecode.Failed, "session not found or not online for this house")
			return
		}
	}
	sess.GetGroupMembers()
	response.SuccessWithOK(c)
}

// 规则设置
func (s *GameShopMemberService) SetVIP(c *gin.Context) {
	var in struct {
		HouseGID int32 `json:"house_gid" binding:"required"`
		MemberID int32 `json:"member_id" binding:"required"`
		VIP      bool  `json:"vip"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	uid := utils.GetUserID(c)
	if err := s.rule.SetVIP(c.Request.Context(), uid, in.HouseGID, in.MemberID, in.VIP); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}
func (s *GameShopMemberService) SetMulti(c *gin.Context) {
	var in struct {
		HouseGID int32 `json:"house_gid" binding:"required"`
		MemberID int32 `json:"member_id" binding:"required"`
		Allow    bool  `json:"allow"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	uid := utils.GetUserID(c)
	if err := s.rule.SetMultiGIDs(c.Request.Context(), uid, in.HouseGID, in.MemberID, in.Allow); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}
func (s *GameShopMemberService) SetTempRelease(c *gin.Context) {
	var in struct {
		HouseGID int32 `json:"house_gid" binding:"required"`
		MemberID int32 `json:"member_id" binding:"required"`
		Limit    int32 `json:"limit"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	uid := utils.GetUserID(c)
	if err := s.rule.SetTempRelease(c.Request.Context(), uid, in.HouseGID, in.MemberID, in.Limit); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}
