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
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type FundsService struct {
	uc  *biz.FundsUseCase
	mgr plaza.Manager
}

func NewFundsService(uc *biz.FundsUseCase, mgr plaza.Manager) *FundsService {
	return &FundsService{uc: uc, mgr: mgr}
}

func (s *FundsService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/members").Use(middleware.JWTAuth())

	g.POST("/credit/deposit", middleware.RequirePerm("fund:deposit"), s.Deposit)
	g.POST("/credit/withdraw", middleware.RequirePerm("fund:withdraw"), s.Withdraw)
	g.POST("/credit/force_withdraw", middleware.RequirePerm("fund:force_withdraw"), s.ForceWithdraw)
	g.PATCH("/limit", middleware.RequirePerm("fund:limit:update"), s.UpdateLimit)
}

// ---- House settings APIs (fees/share_fee/push_credit) ----
// 为保持聚合度，将它们放在 /shops 下更合适，但当前文件已引入资金相关 req/resp，复用即可
func (s *FundsService) registerHouseSettings(r *gin.RouterGroup, hs *HouseSettingsService) {
	shops := r.Group("/shops").Use(middleware.JWTAuth())
	shops.POST("/fees/set", middleware.RequirePerm("shop:fees:update"), hs.SetFees)
	shops.POST("/fees/get", middleware.RequirePerm("shop:fees:view"), hs.Get)
	shops.POST("/sharefee/set", middleware.RequirePerm("shop:sharefee:write"), hs.SetShare)
	shops.POST("/pushcredit/set", middleware.RequirePerm("shop:pushcredit:write"), hs.SetPushCredit)
	shops.POST("/pushcredit/get", middleware.RequirePerm("shop:pushcredit:view"), hs.Get)
}

// Deposit
// @Summary      上分
// @Tags         资金/额度
// @Accept       json
// @Produce      json
// @Param        in body req.CreditDepositRequest true "入参在body"
// @Success      200 {object} response.Body{data=resp.FundsBalanceResponse} "data: { balance: 123 }"
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Router       /members/credit/deposit [post]
func (s *FundsService) Deposit(c *gin.Context) {
	var in req.CreditDepositRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	// 可选：biz_no 兜底
	if len(in.BizNo) == 0 || len(in.BizNo) > 64 {
		response.Fail(c, ecode.ParamsFailed, "invalid biz_no")
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	w, err := s.uc.Deposit(c.Request.Context(), claims.BaseClaims.UserID, in.HouseGID, in.MemberID, in.Amount, in.BizNo, in.Reason)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, resp.FundsBalanceResponse{Balance: w.Balance})
}

// Withdraw
// @Summary      下分
// @Description  受 forbid 与 limit_min 约束
// @Tags         资金/额度
// @Accept       json
// @Produce      json
// @Param        in body req.CreditWithdrawRequest true "入参在body"
// @Success      200 {object} response.Body{data=resp.FundsBalanceResponse} "data: { balance: 123 }"
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Failure      422 {object} response.Body "越过 limit_min 或已禁分"
// @Router       /members/credit/withdraw [post]
func (s *FundsService) Withdraw(c *gin.Context) {
	var in req.CreditWithdrawRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	// 可选：biz_no 兜底
	if len(in.BizNo) == 0 || len(in.BizNo) > 64 {
		response.Fail(c, ecode.ParamsFailed, "invalid biz_no")
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	w, err := s.uc.Withdraw(c.Request.Context(), claims.BaseClaims.UserID, in.HouseGID, in.MemberID, in.Amount, in.BizNo, in.Reason, false)
	if err != nil {
		// 兼容版：根据错误信息粗分；更理想是用 usecase 暴露的 ErrMemberForbidden / ErrCrossLimitMin 等哨兵错误再用 errors.Is 判断
		msg := err.Error()
		if strings.Contains(msg, "forbidden") || strings.Contains(msg, "limit_min") {
			c.JSON(http.StatusUnprocessableEntity, response.Body{Code: ecode.ParamsFailed, Msg: msg})
			return
		}
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, resp.FundsBalanceResponse{Balance: w.Balance})
}

// ForceWithdraw
// @Summary      强制下分
// @Description  忽略 forbid 与 limit_min，仅记录流水
// @Tags         资金/额度
// @Accept       json
// @Produce      json
// @Param        in body req.CreditForceWithdrawRequest true "入参在body"
// @Success      200 {object} response.Body{data=resp.FundsBalanceResponse} "data: { balance: 123 }"
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Router       /members/credit/force_withdraw [post]
func (s *FundsService) ForceWithdraw(c *gin.Context) {
	var in req.CreditForceWithdrawRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	// 可选：biz_no 兜底
	if len(in.BizNo) == 0 || len(in.BizNo) > 64 {
		response.Fail(c, ecode.ParamsFailed, "invalid biz_no")
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	w, err := s.uc.Withdraw(c.Request.Context(), claims.BaseClaims.UserID, in.HouseGID, in.MemberID, in.Amount, in.BizNo, in.Reason, true)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, resp.FundsBalanceResponse{Balance: w.Balance})
}

// UpdateLimit
// @Summary      设置禁分阈值/禁用状态
// @Description  forbid=true 会尝试调用会话进行禁分；false 解除禁分
// @Tags         资金/额度
// @Accept       json
// @Produce      json
// @Param        in body req.UpdateMemberLimitRequest true "入参在body（limit_min/forbid 任一可选）"
// @Success      200 {object} response.Body{data=resp.FundsLimitResponse} "data: { balance, forbid, limit_min }"
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Failure      409 {object} response.Body "会话不存在（需要下发禁分指令时）"
// @Router       /members/limit [patch]
func (s *FundsService) UpdateLimit(c *gin.Context) {
	var in req.UpdateMemberLimitRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	w, err := s.uc.UpdateLimit(c.Request.Context(), claims.BaseClaims.UserID, in.HouseGID, in.MemberID, in.LimitMin, in.Forbid, in.Reason)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	// 只有在传入 forbid 且「确实发生变化」时才下发游戏端禁/解指令
	if in.Forbid != nil && *in.Forbid != w.Forbid {
		if sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), int(in.HouseGID)); !ok || sess == nil {
			c.JSON(http.StatusConflict, response.Body{Code: ecode.Failed, Msg: "session not found or not online for this house"})
			return
		}
		key := fmt.Sprintf("limit-%d-%d-%d", in.HouseGID, in.MemberID, time.Now().UnixNano())
		if err := s.mgr.ForbidMembers(int(claims.BaseClaims.UserID), int(in.HouseGID), key, []int{int(in.MemberID)}, *in.Forbid); err != nil {
			// 失败按文档返回 409；此处不回滚 DB，保持“以库为准”
			c.JSON(http.StatusConflict, response.Body{Code: ecode.Failed, Msg: "push forbid to game failed: " + err.Error()})
			return
		}
	}

	response.Success(c, resp.FundsLimitResponse{
		Balance:  w.Balance,
		Forbid:   w.Forbid,
		LimitMin: w.LimitMin,
	})
}
