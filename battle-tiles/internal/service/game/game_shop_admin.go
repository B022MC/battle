// internal/service/game/game_shop_admin.go
package game

import (
	gameBiz "battle-tiles/internal/biz/game"
	basicRepo "battle-tiles/internal/dal/repo/basic"
	"battle-tiles/internal/dal/req"
	"battle-tiles/internal/dal/resp"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

type ShopAdminService struct {
	uc    *gameBiz.ShopAdminUseCase
	users basicRepo.BasicUserRepo
}

func NewShopAdminService(uc *gameBiz.ShopAdminUseCase, users basicRepo.BasicUserRepo) *ShopAdminService {
	return &ShopAdminService{uc: uc, users: users}
}

func (s *ShopAdminService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/shops").Use(middleware.JWTAuth())

	g.POST("/admins", middleware.RequirePerm("shop:admin:assign"), s.Assign)
	g.DELETE("/admins", middleware.RequirePerm("shop:admin:revoke"), s.Revoke)
	g.POST("/admins/list", middleware.RequirePerm("shop:admin:view"), s.List)
	g.GET("/admins/me", s.GetMyAdminInfo) // 获取当前用户的店铺管理员信息
}

// 提供店铺设置相关的 HTTP 处理
type HouseSettingsService struct{ uc *gameBiz.HouseSettingsUseCase }

func NewHouseSettingsService(uc *gameBiz.HouseSettingsUseCase) *HouseSettingsService {
	return &HouseSettingsService{uc: uc}
}

// RegisterRouter
func (s *HouseSettingsService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/shops").Use(middleware.JWTAuth())
	g.POST("/fees/get", middleware.RequirePerm("shop:fees:view"), s.Get)
	g.POST("/fees/set", middleware.RequirePerm("shop:fees:update"), s.SetFees)
	g.POST("/sharefee/set", middleware.RequirePerm("shop:sharefee:write"), s.SetShare)
	g.POST("/pushcredit/get", middleware.RequirePerm("shop:pushcredit:view"), s.Get)
	g.POST("/pushcredit/set", middleware.RequirePerm("shop:pushcredit:write"), s.SetPushCredit)
	// 费用结算（基础）
	g.POST("/fees/settle/insert", middleware.RequirePerm("shop:fees:update"), s.InsertFeeSettle)
	g.POST("/fees/settle/sum", middleware.RequirePerm("shop:fees:view"), s.SumFeeSettle)
	g.POST("/fees/settle/payoffs", middleware.RequirePerm("shop:fees:view"), s.ListGroupPayoffs)
}

// Get
// @Summary      查询店铺设置（运费/分运/推送额度）
// @Tags         店铺
// @Accept       json
// @Produce      json
// @Param        in body req.GetWalletRequest true "house_gid"
// @Success      200 {object} response.Body{data=resp.HouseSettingsVO}
// @Router       /shops/fees/get [post]
func (s *HouseSettingsService) Get(c *gin.Context) {
	var in struct {
		HouseGID int32 `json:"house_gid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	m, err := s.uc.Get(c.Request.Context(), in.HouseGID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, resp.HouseSettingsVO{HouseGID: m.HouseGID, FeesJSON: m.FeesJSON, ShareFee: m.ShareFee, PushCredit: m.PushCredit})
}

// SetFees
func (s *HouseSettingsService) SetFees(c *gin.Context) {
	var in req.SetFeesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	if err := validateFeesJSON(in.FeesJSON); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	uid := utils.GetUserID(c)
	if err := s.uc.SetFees(c.Request.Context(), uid, in.HouseGID, in.FeesJSON); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, nil)
}

// SetShare
func (s *HouseSettingsService) SetShare(c *gin.Context) {
	var in req.SetShareFeeRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	uid := utils.GetUserID(c)
	if err := s.uc.SetShareFee(c.Request.Context(), uid, in.HouseGID, in.Share); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, nil)
}

// SetPushCredit
func (s *HouseSettingsService) SetPushCredit(c *gin.Context) {
	var in req.SetPushCreditRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	uid := utils.GetUserID(c)
	if err := s.uc.SetPushCredit(c.Request.Context(), uid, in.HouseGID, in.Credit); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, nil)
}

// InsertFeeSettle 新增费用结算
func (s *HouseSettingsService) InsertFeeSettle(c *gin.Context) {
	var in req.InsertFeeSettleRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	t, err := time.Parse(time.RFC3339, in.FeedAt)
	if err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	if err := s.uc.InsertFeeSettle(c.Request.Context(), in.HouseGID, in.PlayGroup, in.Amount, t); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, gin.H{"ok": true})
}

// SumFeeSettle 汇总费用
func (s *HouseSettingsService) SumFeeSettle(c *gin.Context) {
	var in req.SumFeeSettleRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	start, err := time.Parse(time.RFC3339, in.StartAt)
	if err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	end, err := time.Parse(time.RFC3339, in.EndAt)
	if err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	sum, err := s.uc.SumFeeSettle(c.Request.Context(), in.HouseGID, in.PlayGroup, start, end)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, gin.H{"sum": sum})
}

// ListGroupPayoffs 汇总时间范围内各圈费用并计算圈间结转（正数组出借给负数组）
func (s *HouseSettingsService) ListGroupPayoffs(c *gin.Context) {
	var in req.ListGroupPayoffsRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	start, err := time.Parse(time.RFC3339, in.StartAt)
	if err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	end, err := time.Parse(time.RFC3339, in.EndAt)
	if err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	sums, err := s.uc.ListGroupSums(c.Request.Context(), in.HouseGID, start, end)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	type item struct {
		Group string
		Value int64
	}
	var positives []item
	var negatives []item
	for _, gs := range sums {
		if gs.Sum > 0 {
			positives = append(positives, item{Group: gs.PlayGroup, Value: gs.Sum})
		} else if gs.Sum < 0 {
			negatives = append(negatives, item{Group: gs.PlayGroup, Value: gs.Sum})
		}
	}
	// 结转
	type payoff struct {
		From  string `json:"from_group"`
		To    string `json:"to_group"`
		Value int64  `json:"value"`
	}
	var payoffs []payoff
	pi, ni := 0, 0
	for pi < len(positives) && ni < len(negatives) {
		p := &positives[pi]
		n := &negatives[ni]
		need := p.Value
		have := -n.Value
		v := need
		if have < v {
			v = have
		}
		if v > 0 {
			payoffs = append(payoffs, payoff{From: n.Group, To: p.Group, Value: v})
			p.Value -= v
			n.Value += int64(v) // n.Value 是负数，+v 使其趋近 0
		}
		if p.Value == 0 {
			pi++
		}
		if n.Value == 0 {
			ni++
		}
		if v == 0 {
			break
		}
	}
	response.Success(c, gin.H{"group_sums": sums, "payoffs": payoffs})
}

// --- fees_json 校验 ---
type feeRule struct {
	Threshold int    `json:"threshold"`
	Fee       int    `json:"fee"`
	Kind      string `json:"kind,omitempty"`
	Base      int    `json:"base,omitempty"`
}
type feePayload struct {
	Rules []feeRule `json:"rules"`
}

func validateFeesJSON(s string) error {
	if len(s) == 0 {
		return nil
	}
	var p feePayload
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return err
	}
	if len(p.Rules) == 0 {
		return nil
	}
	if len(p.Rules) > 200 {
		return errors.New("fees rules too many")
	}
	for _, r := range p.Rules {
		if r.Threshold <= 0 || r.Fee <= 0 {
			return errors.New("invalid threshold/fee")
		}
		if r.Kind != "" && r.Base <= 0 {
			return errors.New("invalid base for kind rule")
		}
	}
	return nil
}

// Assign
// @Summary      分配店铺管理员
// @Tags         店铺/管理员
// @Accept       json
// @Produce      json
// @Param        in body req.AssignShopAdminRequest true "house_gid, user_id, role(admin|operator, 默认admin)"
// @Success      200 {object} response.Body
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Failure      500 {object} response.Body
// @Router       /shops/admins [post]
func (s *ShopAdminService) Assign(c *gin.Context) {
	var in req.AssignShopAdminRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	// 仅校验登录
	if _, err := utils.GetClaims(c); err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	if err := s.uc.Assign(c.Request.Context(), in.HouseGID, in.UserID, in.Role); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}

// Revoke
// @Summary      撤销店铺管理员
// @Tags         店铺/管理员
// @Accept       json
// @Produce      json
// @Param        in body req.RevokeShopAdminRequest true "house_gid, user_id"
// @Success      200 {object} response.Body
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Failure      500 {object} response.Body
// @Router       /shops/admins [delete]
func (s *ShopAdminService) Revoke(c *gin.Context) {
	var in req.RevokeShopAdminRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	if _, err := utils.GetClaims(c); err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	if err := s.uc.Revoke(c.Request.Context(), in.HouseGID, in.UserID); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}

// List
// @Summary      店铺管理员列表
// @Tags         店铺/管理员
// @Accept       json
// @Produce      json
// @Param        in body req.ListShopAdminsRequest true "house_gid 必填"
// @Success      200 {object} response.Body{data=[]resp.ShopAdminVO}
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Router       /shops/admins/list [post]
func (s *ShopAdminService) List(c *gin.Context) {
	var in req.ListShopAdminsRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	if _, err := utils.GetClaims(c); err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	// 优先使用平台侧圈管理员映射表（game_shop_group_admin）作为“店铺管理员”来源；若未配置则回退到旧表
	// 输出构建辅助
	buildOut := func(ids []struct {
		ID, HouseGID, UserID int32
		Role                 string
	}) []resp.ShopAdminVO {
		// 批量查询昵称
		uidSet := make(map[int32]struct{}, len(ids))
		for _, m := range ids {
			uidSet[m.UserID] = struct{}{}
		}
		uids := make([]int32, 0, len(uidSet))
		for uid := range uidSet {
			uids = append(uids, uid)
		}
		nameMap := map[int32]string{}
		if len(uids) > 0 && s.users != nil {
			if users, e := s.users.SelectByPK(c.Request.Context(), uids); e == nil {
				for _, u := range users {
					if u != nil {
						nameMap[u.Id] = u.NickName
					}
				}
			}
		}
		out := make([]resp.ShopAdminVO, 0, len(ids))
		for _, m := range ids {
			// 过滤平台超级管理员（约定 user_id=1）
			if m.UserID == 1 {
				continue
			}
			out = append(out, resp.ShopAdminVO{
				ID:       m.ID,
				HouseGID: m.HouseGID,
				UserID:   m.UserID,
				Role:     m.Role,
				NickName: nameMap[m.UserID],
			})
		}
		return out
	}

	// 使用 game_shop_admin 表
	oldList, err := s.uc.List(c.Request.Context(), in.HouseGID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	ids := make([]struct {
		ID, HouseGID, UserID int32
		Role                 string
	}, 0, len(oldList))
	for _, m := range oldList {
		if m == nil {
			continue
		}
		ids = append(ids, struct {
			ID, HouseGID, UserID int32
			Role                 string
		}{ID: m.Id, HouseGID: m.HouseGID, UserID: m.UserID, Role: m.Role})
	}
	out := buildOut(ids)
	response.Success(c, out)
}

// GetMyAdminInfo
// @Summary      获取当前用户的店铺管理员信息
// @Tags         店铺/管理员
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Body{data=resp.ShopAdminVO}
// @Failure      401 {object} response.Body
// @Failure      404 {object} response.Body
// @Router       /shops/admins/me [get]
func (s *ShopAdminService) GetMyAdminInfo(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	// 查询当前用户的店铺管理员记录
	admins, err := s.uc.ListByUser(c.Request.Context(), claims.UserID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	if len(admins) == 0 {
		// 不是店铺管理员，返回 null
		response.Success(c, nil)
		return
	}

	// 返回第一个管理员记录（假设一个用户只能是一个店铺的管理员）
	admin := admins[0]

	// 获取用户昵称
	nickName := ""
	if s.users != nil {
		if user, e := s.users.SelectOneByPK(c.Request.Context(), admin.UserID); e == nil && user != nil {
			nickName = user.NickName
		}
	}

	response.Success(c, resp.ShopAdminVO{
		ID:       admin.Id,
		HouseGID: admin.HouseGID,
		UserID:   admin.UserID,
		Role:     admin.Role,
		NickName: nickName,
	})
}
