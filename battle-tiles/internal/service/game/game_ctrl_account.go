// internal/service/game/game_ctrl_account.go
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
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CtrlAccountService struct{ uc *gameBiz.CtrlAccountUseCase }

func NewCtrlAccountService(uc *gameBiz.CtrlAccountUseCase) *CtrlAccountService {
	return &CtrlAccountService{uc: uc}
}

func (s *CtrlAccountService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/shops").Use(middleware.JWTAuth())

	// 仅创建/更新中控（不绑定店铺）
	g.POST("/ctrlAccounts", middleware.RequirePerm("game:ctrl:create"), s.CreateOrUpdate)
	// 更新状态（启用/停用）
	g.POST("/ctrlAccounts/updateStatus", middleware.RequirePerm("game:ctrl:update"), s.UpdateStatus)
	// 删除中控账号
	g.DELETE("/ctrlAccounts/:id", middleware.RequirePerm("game:ctrl:delete"), s.Delete)
	// 绑定/解绑
	g.POST("/ctrlAccounts/bind", middleware.RequirePerm("game:ctrl:update"), s.Bind)
	g.DELETE("/ctrlAccounts/bind", middleware.RequirePerm("game:ctrl:update"), s.Unbind)
	// 列表（按店铺）
	g.POST("/ctrlAccounts/list", middleware.RequirePerm("game:ctrl:view"), s.List)
	g.POST("/ctrlAccounts/listAll", middleware.RequirePerm("game:ctrl:view"), s.ListAll)
	// 店铺号下拉：基于 game_account_house 的去重列表
	g.GET("/houses/options", s.HouseOptions)
}

// ------- 入参结构（如果你已有同名 struct，可忽略这段，保证字段一致即可） -------

type createCtrlAccountOnlyReq struct {
	LoginMode  string `json:"login_mode"  binding:"required"`        // account|mobile
	Identifier string `json:"identifier"  binding:"required"`        // 账号或手机号
	PwdMD5     string `json:"pwd_md5"     binding:"required,len=32"` // 大写32
	Status     *int32 `json:"status"`                                // 可选，默认1
	// 兼容旧接口：如果传 house_gid，则等同“创建+绑定”
	HouseGID *int32 `json:"house_gid"` // optional
}

type bindCtrlReq struct {
	CtrlID   int32  `json:"ctrl_id"   binding:"required"`
	HouseGID int32  `json:"house_gid" binding:"required"`
	Status   *int32 `json:"status"` // 默认1
}

type unbindCtrlReq struct {
	CtrlID   int32 `json:"ctrl_id"   binding:"required"`
	HouseGID int32 `json:"house_gid" binding:"required"`
}

type updateStatusReq struct {
	CtrlID int32  `json:"ctrl_id" binding:"required"`
	Status *int32 `json:"status"  binding:"required,oneof=0 1"` // 0=停用, 1=启用
}

// ------------------------------------------------------------------

func parseModeStr(v string) (consts.GameLoginMode, bool) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "mobile":
		return consts.GameLoginModeMobile, true
	case "account":
		return consts.GameLoginModeAccount, true
	default:
		return consts.GameLoginMode(-1), false
	}
}

// CreateOrUpdate
// @Summary  创建/更新中控账号（不绑定店铺）；兼容：若携带 house_gid，则创建后立即绑定
// @Tags     游戏/中控
// @Accept   json
// @Produce  json
// @Param    in body createCtrlAccountOnlyReq true "入参在body"
// @Success  200 {object} response.Body{data=resp.CtrlAccountVO}
// @Router   /shops/ctrlAccounts [post]
func (s *CtrlAccountService) CreateOrUpdate(c *gin.Context) {
	var in createCtrlAccountOnlyReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	mode, ok := parseModeStr(in.LoginMode)
	if !ok {
		response.Fail(c, ecode.ParamsFailed, "invalid login_mode")
		return
	}
	status := int32(1)
	if in.Status != nil && (*in.Status == 0 || *in.Status == 1) {
		status = *in.Status
	}

	// 1) 仅创建/更新 ctrl
	m, err := s.uc.CreateOrUpdateCtrl(c.Request.Context(), mode, in.Identifier, in.PwdMD5, status, claims.BaseClaims.UserID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	// 2) 兼容：如果携带 house_gid => 立即绑定
	if in.HouseGID != nil && *in.HouseGID > 0 {
		if err := s.uc.BindCtrlToHouse(c.Request.Context(), m.Id, *in.HouseGID, status, "", claims.BaseClaims.UserID); err != nil {
			response.Fail(c, ecode.Failed, err)
			return
		}
	}

	out := &resp.CtrlAccountVO{
		ID:         m.Id,
		LoginMode:  strings.ToLower(in.LoginMode),
		Identifier: m.Identifier,
		Status:     m.Status,
	}
	response.Success(c, out)
}

// UpdateStatus
// @Summary  更新中控账号状态（启用/停用）
// @Tags     游戏/中控
// @Accept   json
// @Produce  json
// @Param    in body updateStatusReq true "ctrl_id, status(0=停用, 1=启用)"
// @Success  200 {object} response.Body
// @Router   /shops/ctrlAccounts/updateStatus [post]
func (s *CtrlAccountService) UpdateStatus(c *gin.Context) {
	var in updateStatusReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	// 更新状态
	if err := s.uc.UpdateStatus(c.Request.Context(), in.CtrlID, *in.Status); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.SuccessWithOK(c)
}

// Delete
// @Summary  删除中控账号
// @Tags     游戏/中控
// @Accept   json
// @Produce  json
// @Param    id path int true "中控账号ID"
// @Success  200 {object} response.Body
// @Router   /shops/ctrlAccounts/{id} [delete]
func (s *CtrlAccountService) Delete(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Fail(c, ecode.ParamsFailed, "id required")
		return
	}

	var id int32
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		response.Fail(c, ecode.ParamsFailed, "invalid id")
		return
	}

	if err := s.uc.Delete(c.Request.Context(), id); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.SuccessWithOK(c)
}

// Bind
// @Summary  绑定中控到店铺
// @Tags     游戏/中控
// @Accept   json
// @Produce  json
// @Param    in body bindCtrlReq true "ctrl_id, house_gid, status(可选), alias(可选)"
// @Success  200 {object} response.Body
// @Router   /shops/ctrlAccounts/bind [post]
func (s *CtrlAccountService) Bind(c *gin.Context) {
	var in bindCtrlReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	status := int32(1)
	if in.Status != nil && (*in.Status == 0 || *in.Status == 1) {
		status = *in.Status
	}
	if err := s.uc.BindCtrlToHouse(c.Request.Context(), in.CtrlID, in.HouseGID, status, "", claims.BaseClaims.UserID); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}

// Unbind
// @Summary  解绑中控与店铺
// @Tags     游戏/中控
// @Accept   json
// @Produce  json
// @Param    in body unbindCtrlReq true "ctrl_id, house_gid"
// @Success  200 {object} response.Body
// @Router   /shops/ctrlAccounts/bind [delete]
func (s *CtrlAccountService) Unbind(c *gin.Context) {
	var in unbindCtrlReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	if err := s.uc.UnbindCtrlFromHouse(c.Request.Context(), in.CtrlID, in.HouseGID); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}

// List
// @Summary  按店铺列出中控账号
// @Tags     游戏/中控
// @Accept   json
// @Produce  json
// @Param    in  body   req.CtrlAccountListRequest  true  "house_gid 必填"
// @Success  200 {object} response.Body{data=[]resp.CtrlAccountVO}
// @Failure  400 {object} response.Body
// @Failure  401 {object} response.Body
// @Router   /shops/ctrlAccounts/list [post]
func (s *CtrlAccountService) List(c *gin.Context) {
	var in req.CtrlAccountListRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	if in.HouseGID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "house_gid required")
		return
	}

	list, err := s.uc.ListCtrlByHouse(c.Request.Context(), in.HouseGID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Success(c, make([]resp.CtrlAccountVO, 0))
			return
		}
		response.Fail(c, ecode.Failed, err)
		return
	}

	out := make([]resp.CtrlAccountVO, 0, len(list))
	for _, m := range list {
		modeStr := map[int32]string{1: "account", 2: "mobile"}[m.LoginMode]
		if modeStr == "" {
			modeStr = "account"
		}
		out = append(out, resp.CtrlAccountVO{
			ID:         m.Id,
			HouseGID:   in.HouseGID, // 这里是“通过绑定”得来的 house 维度
			LoginMode:  modeStr,
			Identifier: m.Identifier,
			Status:     m.Status,
		})
	}
	response.Success(c, out)
}

// ListAll
// @Summary     中控账号全量列表（可筛选分页）
// @Tags        游戏/中控
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       in body req.ListAllCtrlAccountsRequest true "筛选：login_mode(1/2) status(0/1) keyword 模糊 identifier"
// @Success     200 {object} response.Body{data=[]resp.CtrlAccountAllVO} "附带 houses（绑定店铺列表）"
// @Failure     400 {object} response.Body
// @Failure     401 {object} response.Body
// @Router      /shops/ctrlAccounts/listAll [post]
func (s *CtrlAccountService) ListAll(c *gin.Context) {
	var in req.ListAllCtrlAccountsRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	// 统一字符串模式 -> 内部枚举过滤值
	var lm *int32
	if v := strings.ToLower(strings.TrimSpace(in.LoginMode)); v != "" {
		if m, ok := parseModeStr(v); ok {
			t := int32(m)
			lm = &t
		} else {
			response.Fail(c, ecode.ParamsFailed, "invalid login_mode")
			return
		}
	}
	items, _, err := s.uc.ListAll(c.Request.Context(), req.CtrlListFilter{
		LoginMode: lm,
		Status:    in.Status,
		Keyword:   in.Keyword,
	}, in.Page, in.Size)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	out := make([]resp.CtrlAccountAllVO, 0, len(items))
	for _, it := range items {
		modeStr := map[int32]string{1: "account", 2: "mobile"}[it.LoginMode]
		if modeStr == "" {
			modeStr = "account"
		}
		out = append(out, resp.CtrlAccountAllVO{
			ID:           it.Id,
			LoginMode:    modeStr,
			Identifier:   it.Identifier,
			Status:       it.Status,
			LastVerifyAt: it.LastVerifyAt,
			Houses:       it.Houses,
		})
	}
	response.Success(c, out)
}

// HouseOptions
// @Summary  店铺号下拉选项（基于 game_account_house）
// @Tags     游戏/中控
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Success  200 {object} response.Body{data=[]int32}
// @Router   /shops/houses/options [get]
func (s *CtrlAccountService) HouseOptions(c *gin.Context) {
	rows, err := s.uc.ListDistinctHouses(c.Request.Context())
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	if rows == nil {
		rows = []int32{}
	}
	response.Success(c, rows)
}
