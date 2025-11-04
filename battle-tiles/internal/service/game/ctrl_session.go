// internal/service/game/session_service.go
package game

import (
	biz "battle-tiles/internal/biz/game"
	"battle-tiles/internal/dal/req"
	resp "battle-tiles/internal/dal/resp"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

type SessionService struct{ uc *biz.CtrlSessionUseCase }

func NewSessionService(uc *biz.CtrlSessionUseCase) *SessionService {
	return &SessionService{uc: uc}
}

func (s *SessionService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/game").Use(middleware.JWTAuth())
	g.POST("/accounts/sessionStart", middleware.RequirePerm("game:ctrl:create"), s.Start)
	g.POST("/accounts/sessionStop", middleware.RequirePerm("game:ctrl:update"), s.Stop)
}

// Start
// @Summary      启动会话（中控账号）
// @Description  使用中控账号在指定店铺建立会话
// @Tags         游戏/会话
// @Accept       json
// @Produce      json
// @Param        in body req.StartSessionRequest true "id=game_ctrl_account主键ID; house_gid=店铺号"
// @Success      200 {object} response.Body{data=resp.SessionStateResponse} "state=online"
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Failure      500 {object} response.Body
// @Router       /game/accounts/sessionStart [post]
func (s *SessionService) Start(c *gin.Context) {
	var in req.StartSessionRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	if err := s.uc.StartSession(c.Request.Context(), claims.BaseClaims.UserID, in.Id, in.HouseGID); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, resp.SessionStateResponse{State: "online"})
}

// Stop
// @Summary      停止会话（中控账号）
// @Description  停止中控账号与店铺的会话
// @Tags         游戏/会话
// @Accept       json
// @Produce      json
// @Param        in body req.StopSessionRequest true "id=game_ctrl_account主键ID; house_gid=店铺号"
// @Success      200 {object} response.Body{data=resp.SessionStateResponse} "state=offline"
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Failure      500 {object} response.Body
// @Router       /game/accounts/sessionStop [post]
func (s *SessionService) Stop(c *gin.Context) {
	var in req.StopSessionRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	if err := s.uc.StopSession(c.Request.Context(), claims.BaseClaims.UserID, in.Id, in.HouseGID); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, resp.SessionStateResponse{State: "offline"})
}
