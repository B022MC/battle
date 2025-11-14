package game

import (
	"battle-tiles/internal/biz/game"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type MemberService struct {
	memberUC *game.MemberUseCase
	log      *log.Helper
}

func NewMemberService(memberUC *game.MemberUseCase, logger log.Logger) *MemberService {
	return &MemberService{
		memberUC: memberUC,
		log:      log.NewHelper(log.With(logger, "module", "service/member")),
	}
}

// RegisterRouter 注册路由
func (s *MemberService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/members").Use(middleware.JWTAuth())

	g.POST("/list", s.ListAllUsers)          // 查看所有用户
	g.POST("/get", s.GetUser)                // 获取用户信息
	g.POST("/shop-admins", s.ListShopAdmins) // 获取店铺管理员列表
}

// ListAllUsersReq 查看所有用户请求
type ListAllUsersReq struct {
	Page    int32  `json:"page"`
	Size    int32  `json:"size"`
	Keyword string `json:"keyword"`
}

// ListAllUsers 查看所有用户（超级管理员和店铺管理员都可以查看）
// POST /api/members/list
func (s *MemberService) ListAllUsers(c *gin.Context) {
	var req ListAllUsersReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, nil)
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	users, total, err := s.memberUC.ListAllUsers(c.Request.Context(), req.Page, req.Size, req.Keyword)
	if err != nil {
		s.log.Errorf("list all users failed: %v", err)
		response.Fail(c, ecode.Failed, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"items": users,
		"total": total,
		"page":  req.Page,
		"size":  req.Size,
	})
}

// GetUserReq 获取用户信息请求
type GetUserReq struct {
	UserID int32 `json:"user_id" binding:"required"`
}

// GetUser 获取用户信息
// POST /api/members/get
func (s *MemberService) GetUser(c *gin.Context) {
	var req GetUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, nil)
		return
	}

	user, err := s.memberUC.GetUserByID(c.Request.Context(), req.UserID)
	if err != nil {
		s.log.Errorf("get user failed: %v", err)
		response.Fail(c, ecode.Failed, err.Error())
		return
	}

	response.Success(c, user)
}

// ListShopAdminsReq 获取店铺管理员列表请求
type ListShopAdminsReq struct {
	HouseGID int32 `json:"house_gid" binding:"required"`
}

// ListShopAdmins 获取店铺的所有管理员
// POST /api/members/shop-admins
func (s *MemberService) ListShopAdmins(c *gin.Context) {
	var req ListShopAdminsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, nil)
		return
	}

	users, err := s.memberUC.ListShopAdmins(c.Request.Context(), req.HouseGID)
	if err != nil {
		s.log.Errorf("list shop admins failed: %v", err)
		response.Fail(c, ecode.Failed, err.Error())
		return
	}

	response.Success(c, users)
}
