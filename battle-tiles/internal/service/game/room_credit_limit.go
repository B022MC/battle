package game

import (
	biz "battle-tiles/internal/biz/game"
	"battle-tiles/internal/dal/req"
	"battle-tiles/internal/dal/resp"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type RoomCreditLimitService struct {
	uc  *biz.RoomCreditLimitUseCase
	log *log.Helper
}

func NewRoomCreditLimitService(uc *biz.RoomCreditLimitUseCase, logger log.Logger) *RoomCreditLimitService {
	return &RoomCreditLimitService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "service/room_credit_limit")),
	}
}

func (s *RoomCreditLimitService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/room-credit").Use(middleware.JWTAuth())

	g.POST("/set", middleware.RequirePerm("room:credit:set"), s.SetCreditLimit)
	g.POST("/get", middleware.RequirePerm("room:credit:view"), s.GetCreditLimit)
	g.POST("/list", middleware.RequirePerm("room:credit:view"), s.ListCreditLimits)
	g.POST("/delete", middleware.RequirePerm("room:credit:delete"), s.DeleteCreditLimit)
	g.POST("/check", middleware.RequirePerm("room:credit:check"), s.CheckPlayerCredit)
}

// SetCreditLimit
// @Summary      设置房间额度限制
// @Description  设置不同游戏类型、底分、圈子对应的进入房间所需的最低余额
// @Tags         房间额度
// @Accept       json
// @Produce      json
// @Param        in body req.SetRoomCreditLimitRequest true "入参在body"
// @Success      200 {object} response.Body{data=resp.RoomCreditLimitResponse}
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Router       /room-credit/set [post]
func (s *RoomCreditLimitService) SetCreditLimit(c *gin.Context) {
	var in req.SetRoomCreditLimitRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	if err := s.uc.SetCreditLimit(c.Request.Context(), claims.BaseClaims.UserID, in.HouseGID, in.GroupName, in.GameKind, in.BaseScore, in.CreditLimit); err != nil {
		s.log.Errorf("SetCreditLimit failed: %v", err)
		response.Fail(c, ecode.Failed, err)
		return
	}

	// 查询返回设置结果
	limit, err := s.uc.GetCreditLimit(c.Request.Context(), in.HouseGID, in.GroupName, in.GameKind, in.BaseScore)
	if err != nil {
		response.Success(c, gin.H{"message": "设置成功"})
		return
	}

	response.Success(c, &resp.RoomCreditLimitResponse{
		RoomCreditLimitItem: &resp.RoomCreditLimitItem{
			Id:          limit.Id,
			HouseGID:    limit.HouseGID,
			GroupName:   limit.GroupName,
			GameKind:    limit.GameKind,
			BaseScore:   limit.BaseScore,
			CreditLimit: limit.CreditLimit,
			CreditYuan:  float64(limit.CreditLimit) / 100.0,
			CreatedAt:   limit.CreatedAt,
			UpdatedAt:   limit.UpdatedAt,
			UpdatedBy:   limit.UpdatedBy,
		},
	})
}

// GetCreditLimit
// @Summary      查询房间额度限制
// @Description  查询特定游戏类型、底分、圈子的额度限制
// @Tags         房间额度
// @Accept       json
// @Produce      json
// @Param        in body req.GetRoomCreditLimitRequest true "入参在body"
// @Success      200 {object} response.Body{data=resp.RoomCreditLimitResponse}
// @Failure      400 {object} response.Body
// @Router       /room-credit/get [post]
func (s *RoomCreditLimitService) GetCreditLimit(c *gin.Context) {
	var in req.GetRoomCreditLimitRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	limit, err := s.uc.GetCreditLimit(c.Request.Context(), in.HouseGID, in.GroupName, in.GameKind, in.BaseScore)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, &resp.RoomCreditLimitResponse{
		RoomCreditLimitItem: &resp.RoomCreditLimitItem{
			Id:          limit.Id,
			HouseGID:    limit.HouseGID,
			GroupName:   limit.GroupName,
			GameKind:    limit.GameKind,
			BaseScore:   limit.BaseScore,
			CreditLimit: limit.CreditLimit,
			CreditYuan:  float64(limit.CreditLimit) / 100.0,
			CreatedAt:   limit.CreatedAt,
			UpdatedAt:   limit.UpdatedAt,
			UpdatedBy:   limit.UpdatedBy,
		},
	})
}

// ListCreditLimits
// @Summary      列出房间额度限制
// @Description  列出店铺的所有房间额度限制
// @Tags         房间额度
// @Accept       json
// @Produce      json
// @Param        in body req.ListRoomCreditLimitRequest true "入参在body"
// @Success      200 {object} response.Body{data=resp.RoomCreditLimitListResponse}
// @Failure      400 {object} response.Body
// @Router       /room-credit/list [post]
func (s *RoomCreditLimitService) ListCreditLimits(c *gin.Context) {
	var in req.ListRoomCreditLimitRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	limits, err := s.uc.ListCreditLimits(c.Request.Context(), in.HouseGID, in.GroupName)
	if err != nil {
		s.log.Errorf("ListCreditLimits failed: %v", err)
		response.Fail(c, ecode.Failed, err)
		return
	}

	items := make([]*resp.RoomCreditLimitItem, 0, len(limits))
	for _, limit := range limits {
		items = append(items, &resp.RoomCreditLimitItem{
			Id:          limit.Id,
			HouseGID:    limit.HouseGID,
			GroupName:   limit.GroupName,
			GameKind:    limit.GameKind,
			BaseScore:   limit.BaseScore,
			CreditLimit: limit.CreditLimit,
			CreditYuan:  float64(limit.CreditLimit) / 100.0,
			CreatedAt:   limit.CreatedAt,
			UpdatedAt:   limit.UpdatedAt,
			UpdatedBy:   limit.UpdatedBy,
		})
	}

	response.Success(c, &resp.RoomCreditLimitListResponse{
		Total: int32(len(items)),
		Items: items,
	})
}

// DeleteCreditLimit
// @Summary      删除房间额度限制
// @Description  删除特定的房间额度限制
// @Tags         房间额度
// @Accept       json
// @Produce      json
// @Param        in body req.DeleteRoomCreditLimitRequest true "入参在body"
// @Success      200 {object} response.Body
// @Failure      400 {object} response.Body
// @Router       /room-credit/delete [post]
func (s *RoomCreditLimitService) DeleteCreditLimit(c *gin.Context) {
	var in req.DeleteRoomCreditLimitRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	if err := s.uc.DeleteCreditLimit(c.Request.Context(), in.HouseGID, in.GroupName, in.GameKind, in.BaseScore); err != nil {
		s.log.Errorf("DeleteCreditLimit failed: %v", err)
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// CheckPlayerCredit
// @Summary      检查玩家是否满足房间额度要求
// @Description  检查玩家余额是否满足进入特定房间的额度要求
// @Tags         房间额度
// @Accept       json
// @Produce      json
// @Param        in body req.CheckPlayerCreditRequest true "入参在body"
// @Success      200 {object} response.Body{data=resp.CheckPlayerCreditResponse}
// @Failure      400 {object} response.Body
// @Router       /room-credit/check [post]
func (s *RoomCreditLimitService) CheckPlayerCredit(c *gin.Context) {
	var in req.CheckPlayerCreditRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	canEnter, playerBalance, roomCredit, playerCredit, effectiveCredit, err := s.uc.CheckPlayerCanEnterRoom(
		c.Request.Context(),
		in.HouseGID,
		in.GameID,
		in.GroupName,
		in.GameKind,
		in.BaseScore,
	)
	if err != nil {
		s.log.Errorf("CheckPlayerCredit failed: %v", err)
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, &resp.CheckPlayerCreditResponse{
		CanEnter:        canEnter,
		PlayerBalance:   playerBalance,
		RequiredCredit:  roomCredit,
		PlayerCredit:    playerCredit,
		EffectiveCredit: effectiveCredit,
		BalanceYuan:     float64(playerBalance) / 100.0,
		RequiredYuan:    float64(effectiveCredit) / 100.0,
	})
}
