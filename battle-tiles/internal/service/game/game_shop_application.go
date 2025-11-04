package game

import (
	gameModel "battle-tiles/internal/dal/model/game"
	gameRepo "battle-tiles/internal/dal/repo/game"
	"battle-tiles/internal/dal/req"
	"battle-tiles/internal/dal/resp"
	"battle-tiles/internal/infra/plaza"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShopApplicationService struct {
	mgr      plaza.Manager
	userRepo gameRepo.UserApplicationRepo
}

func NewShopApplicationService(mgr plaza.Manager, userRepo gameRepo.UserApplicationRepo) *ShopApplicationService {
	return &ShopApplicationService{mgr: mgr, userRepo: userRepo}
}

func (s *ShopApplicationService) RegisterRouter(r *gin.RouterGroup) {
	shops := r.Group("/shops").Use(middleware.JWTAuth())
	apps := r.Group("/applications").Use(middleware.JWTAuth())

	// 路由不带参数，全部从 body 取
	shops.POST("/applications/list", middleware.RequirePerm("shop:apply:view"), s.List)
	shops.POST("/applications/applyAdmin", s.ApplyAdmin)
	shops.POST("/applications/applyJoin", s.ApplyJoin)
	apps.POST("/approve", middleware.RequirePerm("shop:apply:approve"), s.Approve)
	apps.POST("/reject", middleware.RequirePerm("shop:apply:reject"), s.Reject)
}

// List
// @Summary      店铺管理员申请列表
// @Description  按店铺号列出最近收到的管理员申请（依赖在线会话缓存）
// @Tags         店铺/申请
// @Accept       json
// @Produce      json
// @Param        in  body  req.ListApplicationsRequest  true  "入参在body：{ house_gid }"
// @Success      200 {object} response.Body{data=resp.ApplicationsVO}
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Failure      409 {object} response.Body "session not online"
// @Router       /shops/applications/list [post]
func (s *ShopApplicationService) List(c *gin.Context) {
	var in req.ListApplicationsRequest
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid house_gid")
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	// 需要该店铺的在线会话（只读：支持共享会话回退）
	sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
	if !ok || sess == nil {
		if shared, ok2 := s.mgr.GetAnyByHouse(in.HouseGID); ok2 && shared != nil {
			sess = shared
		} else {
			c.JSON(http.StatusConflict, response.Body{Code: ecode.Failed, Msg: "session not found or not online for this house"})
			return
		}
	}

	apps := sess.ListApplications(in.HouseGID)
	out := make([]*resp.ApplicationItemVO, 0, len(apps))
	for _, a := range apps {
		// 过滤：类型 / 指定圈主（若传）
		if in.Type != nil && a.ApplyType != *in.Type {
			continue
		}
		if in.AdminUID != nil && a.AdminUserID != *in.AdminUID {
			continue
		}
		out = append(out, &resp.ApplicationItemVO{
			ID:          a.MessageId,
			Status:      a.MessageStatus,
			ApplierID:   a.AplierId,
			ApplierGID:  a.ApplierGid,
			ApplierName: a.ApplierGName,
			HouseGID:    a.HouseGid,
			Type:        a.ApplyType,
			AdminUserID: a.AdminUserID,
			CreatedAt:   a.CreatedAt,
		})
	}
	response.Success(c, &resp.ApplicationsVO{Items: out})
}

// Approve
// @Summary      审批通过
// @Description  通过指定申请（按消息ID在在线会话缓存中查找）
// @Tags         店铺/申请
// @Accept       json
// @Produce      json
// @Param        in  body  req.DecideApplicationRequest  true  "入参在body：{ id }"
// @Success      200 {object} response.Body{data=resp.AckVO}
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Failure      404 {object} response.Body "application not found in online sessions"
// @Failure      409 {object} response.Body "no online session"
// @Router       /applications/approve [post]
func (s *ShopApplicationService) Approve(c *gin.Context) {
	s.respond(c, true)
}

// Reject
// @Summary      审批拒绝
// @Description  拒绝指定申请（按消息ID在在线会话缓存中查找）
// @Tags         店铺/申请
// @Accept       json
// @Produce      json
// @Param        in  body  req.DecideApplicationRequest  true  "入参在body：{ id }"
// @Success      200 {object} response.Body{data=resp.AckVO}
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Failure      404 {object} response.Body "application not found in online sessions"
// @Failure      409 {object} response.Body "no online session"
// @Router       /applications/reject [post]
func (s *ShopApplicationService) Reject(c *gin.Context) {
	s.respond(c, false)
}

func (s *ShopApplicationService) respond(c *gin.Context, agree bool) {
	var in req.DecideApplicationRequest
	if err := c.ShouldBindJSON(&in); err != nil || in.ID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid application id")
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	// 遍历该用户的所有在线会话（不同店铺），找到这条申请
	sessions := s.mgr.GetByUser(int(claims.BaseClaims.UserID))
	if len(sessions) == 0 {
		c.JSON(http.StatusConflict, response.Body{Code: ecode.Failed, Msg: "no online sessions"})
		return
	}
	for _, sess := range sessions {
		if sess == nil {
			continue
		}
		if ai, ok := sess.FindApplicationByID(in.ID); ok && ai != nil {
			// 下发审批指令（同意/拒绝）
			sess.RespondApplication(ai, agree)

			// 通过之后：申请人成为管理员，还需要主动对某个店铺“值班”（启动会话）后，才能做后续操作。
			response.Success(c, &resp.AckVO{OK: true})
			return
		}
	}

	c.JSON(http.StatusNotFound, response.Body{
		Code: ecode.Failed,
		Msg:  "application not found in online sessions cache; ensure target shop session is online and has recently received apply messages",
	})
}

// ApplyAdmin 发起管理员申请（平台侧持久化）
func (s *ShopApplicationService) ApplyAdmin(c *gin.Context) {
	var in req.ApplyAdminRequest
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid house_gid")
		return
	}
	uid := utils.GetUserID(c)
	app := &gameModel.UserApplication{HouseGID: in.HouseGID, Applicant: uid, Type: 1, AdminUID: 0, Note: in.Note, Status: 0}
	if err := s.userRepo.Insert(c.Request.Context(), app); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}

// ApplyJoin 发起入圈申请（平台侧持久化）
func (s *ShopApplicationService) ApplyJoin(c *gin.Context) {
	var in req.ApplyJoinGroupRequest
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 || in.AdminUserID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid params")
		return
	}
	uid := utils.GetUserID(c)
	app := &gameModel.UserApplication{HouseGID: in.HouseGID, Applicant: uid, Type: 2, AdminUID: in.AdminUserID, Note: in.Note, Status: 0}
	if err := s.userRepo.Insert(c.Request.Context(), app); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}
