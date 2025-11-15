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
	g := r.Group("/game")
	// 验证账号不需要 JWT（用于注册前验证）
	g.POST("/accounts/verify", s.VerifyAccount) // 探活 82

	// 以下接口需要 JWT 认证
	auth := g.Use(middleware.JWTAuth())
	auth.POST("/accounts", s.BindMyAccount)               // 仅 1 条（只建 game_account）
	auth.GET("/accounts/me", s.GetMyAccount)              // 查询我的账号
	auth.GET("/accounts/me/houses", s.GetMyAccountHouses) // 查询我的账号绑定的店铺游戏ID
	auth.DELETE("/accounts/me", s.DeleteMyAccount)        // 解绑我的账号

	// 管理员接口
	admin := g.Use(middleware.JWTAuth(), middleware.AdminOnly())
	admin.POST("/accounts/fix-empty-game-user-id", s.FixEmptyGameUserID) // 修复空的 game_user_id
}

// VerifyAccount
// @Summary     校验游戏账号是否可用（探活82）
// @Description 只做登录探测，不写库、不建立会话。此接口无需认证，用于注册前验证游戏账号。
// @Tags        游戏/我的账号
// @Accept      json
// @Produce     json
// @Param       in  body     req.VerifyAccountRequest  true  "mode: account|mobile；account: 账号或手机号"
// @Success     200 {object} response.Body{data=resp.VerifyAccountResponse} "data: { ok: true }"
// @Failure     400 {object} response.Body
// @Failure     500 {object} response.Body
// @Router      /game/accounts/verify [post]
func (s *AccountService) VerifyAccount(c *gin.Context) {
	var in req.VerifyAccountRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
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

// GetMyAccountHouses
// @Summary  查询我的游戏账号绑定的店铺游戏ID列表
// @Tags     游戏/我的账号
// @Produce  json
// @Security BearerAuth
// @Success  200 {object} response.Body{data=[]resp.GameAccountHouseVO} "返回绑定的店铺游戏ID列表"
// @Failure  401 {object} response.Body
// @Failure  500 {object} response.Body
// @Router   /game/accounts/me/houses [get]
func (s *AccountService) GetMyAccountHouses(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	houses, err := s.uc.GetMyHouses(c.Request.Context(), claims.BaseClaims.UserID)
	if err != nil {
		// 如果用户没有绑定游戏账号，返回空数组而不是错误
		response.Success(c, []resp.GameAccountHouseVO{})
		return
	}
	var result []resp.GameAccountHouseVO
	for _, h := range houses {
		result = append(result, resp.GameAccountHouseVO{
			ID:        h.Id,
			HouseGID:  h.HouseGID,
			IsDefault: h.IsDefault,
			Status:    h.Status,
		})
	}
	response.Success(c, result)
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

// FixEmptyGameUserID
// @Summary  修复空的 game_user_id 字段（管理员接口）
// @Description 修复在修复代码之前注册的用户的 game_user_id。此接口仅限管理员使用。
// @Tags     游戏/我的账号
// @Produce  json
// @Security BearerAuth
// @Success  200 {object} response.Body{data=map[string]int64} "data: { fixed: 10, failed: 2 }"
// @Failure  401 {object} response.Body
// @Failure  403 {object} response.Body
// @Failure  500 {object} response.Body
// @Router   /game/accounts/fix-empty-game-user-id [post]
func (s *AccountService) FixEmptyGameUserID(c *gin.Context) {
	fixed, failed, err := s.uc.FixEmptyGameUserID(c.Request.Context())
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, map[string]int64{
		"fixed":  fixed,
		"failed": failed,
	})
}
