package service

import (
	"battle-tiles/internal/biz"
	"battle-tiles/internal/consts"
	"battle-tiles/internal/infra/plaza"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

// PublicService is a public service.
type PublicService struct {
	uc *biz.PublicUseCase
}

// NewPublicService new a public service.
func NewPublicService(uc *biz.PublicUseCase) *PublicService {
	return &PublicService{uc: uc}
}

// OpsService 提供只读运维端点
type OpsService struct{ mgr plaza.Manager }

func NewOpsService(mgr plaza.Manager) *OpsService { return &OpsService{mgr: mgr} }

func (s *OpsService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/ops").Use(middleware.JWTAuth())
	// 如需更细权限，可加 RequirePerm("ops:view")
	g.GET("/plaza/metrics", s.Metrics)
	g.GET("/plaza/health", s.Health)
	g.POST("/plaza/housesByLogin", s.HousesByLogin)
	// plaza 控制/数据
	g.POST("/plaza/forbidMembers", s.ForbidMembers)
	g.POST("/plaza/members/pull", s.PullMembers)
	g.POST("/plaza/table/query", s.QueryTable)
	// 申请列表与处理
	g.POST("/plaza/applications/list", s.ListApplications)
	g.POST("/plaza/applications/respond", s.RespondApplication)
}

// Metrics 返回 plaza.Manager 的指标快照
// @Summary      获取广场指标
// @Description  返回 plaza.Manager 的指标快照
// @Tags         Ops
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} response.Body
// @Router       /ops/plaza/metrics [get]
func (s *OpsService) Metrics(c *gin.Context) {
	m := s.mgr.Metrics()
	response.Success(c, m)
}

// Health 返回 plaza.Manager 的健康状态
// @Summary      获取广场健康状态
// @Description  返回 plaza.Manager 的健康状态
// @Tags         Ops
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} response.Body
// @Router       /ops/plaza/health [get]
func (s *OpsService) Health(c *gin.Context) {
	h := s.mgr.Health()
	response.Success(c, h)
}

type housesByLoginReq struct {
	LoginMode  string `json:"login_mode" binding:"omitempty,oneof=account mobile"`
	Identifier string `json:"identifier" binding:"required"`
	PwdMD5     string `json:"pwd_md5" binding:"required,len=32"`
}

// HousesByLogin 通过游戏登录态尝试列出可见店铺（best-effort）
// @Summary  通过账号登录尝试列出店铺
// @Tags     Ops
// @Accept   json
// @Produce  json
// @Param    in body housesByLoginReq true "login_mode 默认 account"
// @Success  200 {object} response.Body{data=[]int}
// @Router   /ops/plaza/housesByLogin [post]
func (s *OpsService) HousesByLogin(c *gin.Context) {
	var in housesByLoginReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, err)
		return
	}
	mode := consts.GameLoginModeAccount
	if in.LoginMode == "mobile" {
		mode = consts.GameLoginModeMobile
	}
	hs, err := s.mgr.ListHousesByLogin(c.Request.Context(), mode, in.Identifier, in.PwdMD5)
	if err != nil {
		response.Fail(c, 500, err)
		return
	}
	response.Success(c, hs)
}

type forbidMembersReq struct {
	HouseGID int    `json:"house_gid" binding:"required"`
	Key      string `json:"key"      binding:"required"`
	Members  []int  `json:"members"  binding:"required,min=1"`
	Forbid   bool   `json:"forbid"`
}

// ForbidMembers 禁言/解禁一组成员（通过 87 指令）
func (s *OpsService) ForbidMembers(c *gin.Context) {
	var in forbidMembersReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, 401, err)
		return
	}
	if err := s.mgr.ForbidMembers(int(claims.BaseClaims.UserID), in.HouseGID, in.Key, in.Members, in.Forbid); err != nil {
		response.Fail(c, 500, err)
		return
	}
	response.SuccessWithOK(c)
}

type pullMembersReq struct {
	HouseGID int `json:"house_gid" binding:"required"`
}

// PullMembers 主动拉取成员列表（触发 87 推送）
func (s *OpsService) PullMembers(c *gin.Context) {
	var in pullMembersReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, 401, err)
		return
	}
	if err := s.mgr.GetGroupMembers(int(claims.BaseClaims.UserID), in.HouseGID); err != nil {
		response.Fail(c, 500, err)
		return
	}
	response.SuccessWithOK(c)
}

type queryTableReq struct {
	HouseGID  int `json:"house_gid"  binding:"required"`
	MappedNum int `json:"mapped_num" binding:"required"`
}

// QueryTable 查询单桌（向 87 队列发指令，快照随后下发）
func (s *OpsService) QueryTable(c *gin.Context) {
	var in queryTableReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, 401, err)
		return
	}
	sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
	if !ok || sess == nil {
		response.Fail(c, 400, "session not found; start it first")
		return
	}
	sess.QueryTable(in.MappedNum)
	response.SuccessWithOK(c)
}

type listApplicationsReq struct {
	HouseGID int `json:"house_gid"`
}

// ListApplications 获取申请消息快照（保存在会话缓存中）
func (s *OpsService) ListApplications(c *gin.Context) {
	var in listApplicationsReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, 401, err)
		return
	}
	// 若指定 house 取该会话，否则聚合该用户所有会话里的申请
	if in.HouseGID > 0 {
		sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
		if !ok || sess == nil {
			response.Success(c, []any{})
			return
		}
		response.Success(c, sess.ListApplications(in.HouseGID))
		return
	}
	// 聚合所有 house 的申请
	m := s.mgr.GetByUser(int(claims.BaseClaims.UserID))
	out := make([]any, 0, 16)
	for _, sess := range m {
		for _, ai := range sess.ListApplications(0) {
			out = append(out, ai)
		}
	}
	response.Success(c, out)
}

type respondApplicationReq struct {
	MessageID int  `json:"message_id" binding:"required"`
	Agree     bool `json:"agree"`
}

// RespondApplication 同意/拒绝申请（在任意活跃会话中查找到消息即处理）
func (s *OpsService) RespondApplication(c *gin.Context) {
	var in respondApplicationReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, 401, err)
		return
	}
	m := s.mgr.GetByUser(int(claims.BaseClaims.UserID))
	for _, sess := range m {
		if ai, ok := sess.FindApplicationByID(in.MessageID); ok {
			sess.RespondApplication(ai, in.Agree)
			response.SuccessWithOK(c)
			return
		}
	}
	response.Fail(c, 404, "application not found in active sessions")
}
