// internal/service/game/account_service.go
package game

import (
	gameBiz "battle-tiles/internal/biz/game"
	"battle-tiles/internal/consts"
	"battle-tiles/internal/dal/req"
	"battle-tiles/internal/dal/resp"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

type AccountService struct{ uc *gameBiz.GameAccountUseCase }

func NewAccountService(uc *gameBiz.GameAccountUseCase) *AccountService {
	return &AccountService{uc: uc}
}

func (s *AccountService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/game").Use(middleware.JWTAuth())
	g.POST("/accounts/verify", s.VerifyAccount) // 探活 82
	g.POST("/accounts", s.BindMyAccount)        // 仅 1 条（只建 game_account）
	g.GET("/accounts/me", s.GetMyAccount)       // 查询我的账号
	g.DELETE("/accounts/me", s.DeleteMyAccount) // 解绑我的账号
}

// VerifyAccount
// @Summary     校验游戏账号是否可用（探活82）
// @Description 只做登录探测，不写库、不建立会话
// @Tags        游戏/我的账号
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       in  body     req.VerifyAccountRequest  true  "mode: account|mobile；account: 账号或手机号"
// @Success     200 {object} response.Body{data=resp.VerifyAccountResponse} "data: { ok: true }"
// @Failure     400 {object} response.Body
// @Failure     401 {object} response.Body
// @Failure     500 {object} response.Body
// @Router      /game/accounts/verify [post]
func (s *AccountService) VerifyAccount(c *gin.Context) {
	var in req.VerifyAccountRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	if _, err := utils.GetClaims(c); err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	var mode consts.GameLoginMode
	switch in.Mode {
	case "account":
		mode = consts.GameLoginModeAccount
	case "mobile":
		mode = consts.GameLoginModeMobile
	default:
		response.Fail(c, ecode.ParamsFailed, "invalid mode")
		return
	}
	if err := s.uc.Verify(c.Request.Context(), mode, in.Account, in.PwdMD5); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, resp.VerifyAccountResponse{Ok: true})
}

// BindMyAccount
// @Summary     绑定“我的”游戏账号（仅允许 1 条）
// @Description 普通用户仅能绑定1个游戏账号（DB 触发器兜底）；管理员不受限
// @Tags        游戏/我的账号
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       in  body     req.BindMyAccountRequest  true  "入参在 body"
// @Success     200 {object} response.Body{data=resp.AccountVO}
// @Failure     400 {object} response.Body
// @Failure     401 {object} response.Body
// @Failure     409 {object} response.Body "you have already bound a game account"
// @Failure     500 {object} response.Body
// @Router      /game/accounts [post]
func (s *AccountService) BindMyAccount(c *gin.Context) {
	var in req.BindMyAccountRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	var mode consts.GameLoginMode
	switch in.Mode {
	case "account":
		mode = consts.GameLoginModeAccount
	case "mobile":
		mode = consts.GameLoginModeMobile
	default:
		response.Fail(c, ecode.ParamsFailed, "invalid mode")
		return
	}
	acc, err := s.uc.BindSingle(c.Request.Context(), claims.BaseClaims.UserID, mode, in.Account, in.PwdMD5, in.Nickname)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, &resp.AccountVO{
		ID:        acc.Id,
		Account:   acc.Account,
		Nickname:  acc.Nickname,
		IsDefault: acc.IsDefault,
		Status:    acc.Status,
		LoginMode: acc.LoginMode,
	})
}

// GetMyAccount
// @Summary  我的游戏账号
// @Tags     游戏/我的账号
// @Produce  json
// @Security BearerAuth
// @Success  200 {object} response.Body{data=resp.AccountVO} "若未绑定，data 为 null"
// @Failure  401 {object} response.Body
// @Failure  500 {object} response.Body
// @Router   /game/accounts/me [get]
func (s *AccountService) GetMyAccount(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	acc, err := s.uc.GetMine(c.Request.Context(), claims.BaseClaims.UserID)
	if err != nil || acc == nil {
		response.Success(c, nil)
		return
	}
	response.Success(c, &resp.AccountVO{
		ID:        acc.Id,
		Account:   acc.Account,
		Nickname:  acc.Nickname,
		IsDefault: acc.IsDefault,
		Status:    acc.Status,
		LoginMode: acc.LoginMode,
	})
}

// DeleteMyAccount
// @Summary  解绑我的游戏账号
// @Tags     游戏/我的账号
// @Produce  json
// @Security BearerAuth
// @Success  200 {object} response.Body
// @Failure  401 {object} response.Body
// @Failure  500 {object} response.Body
// @Router   /game/accounts/me [delete]
func (s *AccountService) DeleteMyAccount(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	if err := s.uc.DeleteMine(c.Request.Context(), claims.BaseClaims.UserID); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}
