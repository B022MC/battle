package game

import (
	gameBiz "battle-tiles/internal/biz/game"
	"battle-tiles/internal/infra/plaza"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ShopGroupAdminService 管理“我的圈子”与圈子列表
type ShopGroupAdminService struct {
	uc  *gameBiz.ShopGroupAdminUseCase
	mgr plaza.Manager
}

func NewShopGroupAdminService(uc *gameBiz.ShopGroupAdminUseCase, mgr plaza.Manager) *ShopGroupAdminService {
	return &ShopGroupAdminService{uc: uc, mgr: mgr}
}

func (s *ShopGroupAdminService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/shops").Use(middleware.JWTAuth())
	g.GET("/my/groups", s.MyGroups)
	g.POST("/groups/list", s.ListGroupsByHouse)
}

// MyGroups 返回当前用户作为圈管理员的 (house_gid, group_id) 列表
func (s *ShopGroupAdminService) MyGroups(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	arr, err := s.uc.ListByUser(c.Request.Context(), claims.BaseClaims.UserID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	type item struct{ HouseGID, GroupID int32 }
	out := make([]item, 0, len(arr))
	for _, it := range arr {
		out = append(out, item{HouseGID: it[0], GroupID: it[1]})
	}
	response.Success(c, out)
}

// ListGroupsByHouse 基于在线会话快照，返回指定店铺下当前可见的 group_id 列表
func (s *ShopGroupAdminService) ListGroupsByHouse(c *gin.Context) {
	var in struct {
		HouseGID int `json:"house_gid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid house_gid")
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
			c.JSON(http.StatusConflict, response.Body{Code: ecode.Failed, Msg: "session not found or not online for this house"})
			return
		}
	}
	// 从桌台快照推断可见的 group_id 集合
	set := map[int32]struct{}{}
	if tables := sess.ListTables(); len(tables) > 0 {
		for _, t := range tables {
			if t != nil && t.GroupID > 0 {
				set[int32(t.GroupID)] = struct{}{}
			}
		}
	}
	out := make([]int32, 0, len(set))
	for gid := range set {
		out = append(out, gid)
	}
	response.Success(c, out)
}
