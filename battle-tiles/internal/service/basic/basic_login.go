package basic

import (
	basicBiz "battle-tiles/internal/biz/basic"
	"battle-tiles/internal/dal/req"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

type BasicLoginService struct {
	uc *basicBiz.BasicLoginUseCase
}

func NewBasicLoginService(uc *basicBiz.BasicLoginUseCase) *BasicLoginService {
	return &BasicLoginService{uc: uc}
}

func (s *BasicLoginService) RegisterRouter(router *gin.RouterGroup) {
	r := router.Group("/login").Use(middleware.SwitchingDB())
	r.POST("/username", s.LoginByUsernamePassword)
	r.POST("/register", s.Register)
}

// Register
// @Summary      用户注册（用户名 + 密码）
// @Description  使用用户名与密码注册新用户，成功后返回访问令牌与用户信息
// @Tags         基础管理/登录
// @Accept       json
// @Produce      json
// @Param        payload    body      req.RegisterRequest  true  "注册参数"
// @Success      200        {object}  response.Body{data=resp.LoginResponse,msg=string}
// @Failure      400        {object}  response.Body{msg=string}  "参数错误"
// @Failure      500        {object}  response.Body{msg=string}  "注册失败"
// @Router       /login/register [post]
func (s *BasicLoginService) Register(ctx *gin.Context) {
	var req req.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	res, err := s.uc.Register(ctx.Request.Context(), ctx, &req)
	if err != nil {
		response.Fail(ctx, ecode.RegisterFailed, err)
		return
	}
	response.Success(ctx, res)
}

// LoginByUsernamePassword
// @Summary      用户名密码登录
// @Description  使用用户名与密码登录，成功后返回访问令牌与用户信息
// @Tags         基础管理/登录
// @Accept       json
// @Produce      json
// @Param        payload    body      req.UsernamePasswordLoginRequest  true  "登录参数"
// @Success      200        {object}  response.Body{data=resp.LoginResponse,msg=string}
// @Failure      400        {object}  response.Body{msg=string}  "参数错误"
// @Failure      401        {object}  response.Body{msg=string}  "登录失败"
// @Router       /login/username [post]
func (s *BasicLoginService) LoginByUsernamePassword(ctx *gin.Context) {
	var req req.UsernamePasswordLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	res, err := s.uc.LoginByUsernamePassword(ctx.Request.Context(), ctx, &req)
	if err != nil {
		response.Fail(ctx, ecode.LoginFailed, err)
		return
	}
	response.Success(ctx, res)
}
