package game

import (
	gameModel "battle-tiles/internal/dal/model/game"
	basicRepo "battle-tiles/internal/dal/repo/basic"
	gameRepo "battle-tiles/internal/dal/repo/game"
	"battle-tiles/internal/dal/req"
	"battle-tiles/internal/dal/resp"
	"battle-tiles/internal/infra/plaza"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ShopApplicationService struct {
	mgr      plaza.Manager
	userRepo gameRepo.UserApplicationRepo
	users    basicRepo.BasicUserRepo
	auth     basicRepo.AuthRepo
	sAdm     gameRepo.GameShopAdminRepo
}

func NewShopApplicationService(mgr plaza.Manager, userRepo gameRepo.UserApplicationRepo, users basicRepo.BasicUserRepo, auth basicRepo.AuthRepo, sAdm gameRepo.GameShopAdminRepo) *ShopApplicationService {
	return &ShopApplicationService{mgr: mgr, userRepo: userRepo, users: users, auth: auth, sAdm: sAdm}
}

func (s *ShopApplicationService) RegisterRouter(r *gin.RouterGroup) {
	shops := r.Group("/shops").Use(middleware.JWTAuth())
	apps := r.Group("/applications").Use(middleware.JWTAuth())

	// 路由不带参数，全部从 body 取
	shops.POST("/applications/list", middleware.RequirePerm("shop:apply:view"), s.List)
	shops.POST("/applications/applyAdmin", s.ApplyAdmin)
	shops.POST("/applications/applyJoin", s.ApplyJoin)
	// 仅返回当前用户的申请记录，无需额外权限
	shops.POST("/applications/history", s.History)
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

	// 需要该店铺的在线会话（只读：支持共享会话回退）；若没有在线会话，则回退到平台库
	sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
	if !ok || sess == nil {
		if shared, ok2 := s.mgr.GetAnyByHouse(in.HouseGID); ok2 && shared != nil {
			sess = shared
		} else {
			// 无在线会话：直接返回平台库中的待审申请（支持类型与圈主过滤）
			pending := int32(0)
			var typ *int32
			if in.Type != nil {
				t := int32(*in.Type)
				typ = &t
			}
			list, err := s.userRepo.ListByHouse(c.Request.Context(), int32(in.HouseGID), typ, &pending)
			if err != nil {
				c.JSON(http.StatusConflict, response.Body{Code: ecode.Failed, Msg: "session not online and read platform applications failed"})
				return
			}
			// 过滤圈主（仅对加圈申请生效）：默认当前登录用户
			var adminFilter *int
			if in.AdminUID != nil {
				adminFilter = in.AdminUID
			} else {
				u := int(claims.BaseClaims.UserID)
				adminFilter = &u
			}
			out := make([]*resp.ApplicationItemVO, 0, len(list))
			for _, a := range list {
				if typ != nil && *typ == 2 { // join
					if adminFilter != nil && int32(*adminFilter) != a.AdminUID {
						continue
					}
				}
				out = append(out, &resp.ApplicationItemVO{
					ID:          int(a.Id),
					Status:      int(a.Status),
					ApplierID:   int(a.Applicant),
					ApplierGID:  0,
					ApplierName: "",
					HouseGID:    int(a.HouseGID),
					Type:        int(a.Type),
					AdminUserID: int(a.AdminUID),
					CreatedAt:   a.CreatedAt.UnixMilli(),
				})
			}
			// 批量昵称
			if len(out) > 0 {
				ids := make([]int32, 0, len(out))
				for _, it := range out {
					ids = append(ids, int32(it.ApplierID))
				}
				rows, _ := s.users.SelectByPK(c.Request.Context(), ids)
				nameMap := map[int32]string{}
				for _, u := range rows {
					if u != nil {
						nameMap[u.Id] = u.NickName
					}
				}
				for _, it := range out {
					it.ApplierName = nameMap[int32(it.ApplierID)]
				}
			}
			response.Success(c, &resp.ApplicationsVO{Items: out})
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
	// 若在线会话缓存为空，补充平台侧待审申请（支持类型与圈主过滤）
	if len(out) == 0 {
		pending := int32(0)
		var typ *int32
		if in.Type != nil {
			t := int32(*in.Type)
			typ = &t
		}
		list, err := s.userRepo.ListByHouse(c.Request.Context(), int32(in.HouseGID), typ, &pending)
		if err == nil && len(list) > 0 {
			// 批量查询昵称
			ids := make([]int32, 0, len(list))
			for _, a := range list {
				ids = append(ids, a.Applicant)
			}
			// SelectByPK 支持数组
			rows, _ := s.users.SelectByPK(c.Request.Context(), ids)
			nameMap := map[int32]string{}
			for _, u := range rows {
				if u != nil {
					nameMap[u.Id] = u.NickName
				}
			}
			// 过滤圈主（仅对加圈申请生效）：默认当前登录用户
			var adminFilter *int
			if in.AdminUID != nil {
				adminFilter = in.AdminUID
			} else {
				u := int(claims.BaseClaims.UserID)
				adminFilter = &u
			}
			for _, a := range list {
				if typ != nil && *typ == 2 { // join
					if adminFilter != nil && int32(*adminFilter) != a.AdminUID {
						continue
					}
				}
				out = append(out, &resp.ApplicationItemVO{
					ID:          int(a.Id),
					Status:      int(a.Status),
					ApplierID:   int(a.Applicant),
					ApplierGID:  0,
					ApplierName: nameMap[a.Applicant],
					HouseGID:    int(a.HouseGID),
					Type:        int(a.Type),
					AdminUserID: int(a.AdminUID),
					CreatedAt:   a.CreatedAt.UnixMilli(),
				})
			}
		}
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

// ensureGroupAdminMapping 已废弃，使用新的圈子系统
// 圈子会在店铺管理员首次访问时自动创建
func (s *ShopApplicationService) ensureGroupAdminMapping(c *gin.Context, houseGID, userID int32) {
	// 不再需要手动创建圈子映射
}

// ensureShopAdmin 在管理员审批通过后，写入/更新 game_shop_admin（幂等）
func (s *ShopApplicationService) ensureShopAdmin(c *gin.Context, houseGID, userID int32) {
	if s.sAdm == nil {
		return
	}
	_ = s.sAdm.Assign(c.Request.Context(), &gameModel.GameShopAdmin{HouseGID: houseGID, UserID: userID, Role: "admin"})
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

	// 优先：遍历该用户的所有在线会话（不同店铺），找到这条申请
	sessions := s.mgr.GetByUser(int(claims.BaseClaims.UserID))
	if len(sessions) == 0 {
		// 平台端兜底：直接修改平台申请状态
		if s.platformDecide(c, in.ID, agree) {
			return
		}
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
			if agree && ai.ApplyType == 1 { // 管理员申请
				// 仅保留管理员角色
				_ = s.auth.EnsureUserHasOnlyRoleByCode(c.Request.Context(), int32(ai.AplierId), "shop_admin")
				// 确保圈管理员映射：选取该店铺可见的第一个圈ID；若不可见则置0占位
				s.ensureGroupAdminMapping(c, int32(ai.HouseGid), int32(ai.AplierId))
				// 同步写入店铺管理员表（兼容旧逻辑）
				s.ensureShopAdmin(c, int32(ai.HouseGid), int32(ai.AplierId))
			}
			response.Success(c, &resp.AckVO{OK: true})
			return
		}
	}
	// 在线会话未找到该消息：平台端兜底
	if s.platformDecide(c, in.ID, agree) {
		return
	}
	c.JSON(http.StatusNotFound, response.Body{Code: ecode.Failed, Msg: "application not found (session cache and platform store)"})
}

// platformDecide 平台端审批兜底：将库中 status 从 0 改为 1/2
func (s *ShopApplicationService) platformDecide(c *gin.Context, id int, agree bool) bool {
	m, err := s.userRepo.GetByID(c.Request.Context(), int32(id))
	if err != nil || m == nil {
		return false
	}
	if m.Status != 0 { // 非待审
		response.Success(c, &resp.AckVO{OK: true})
		return true
	}
	newStatus := int32(2)
	if agree {
		newStatus = 1
	}
	if _, err := s.userRepo.UpdateStatusByID(c.Request.Context(), m.Id, newStatus); err != nil {
		return false
	}
	if agree && m.Type == 1 { // 管理员申请通过 -> 仅保留管理员角色
		_ = s.auth.EnsureUserHasOnlyRoleByCode(c.Request.Context(), m.Applicant, "shop_admin")
		s.ensureGroupAdminMapping(c, m.HouseGID, m.Applicant)
		s.ensureShopAdmin(c, m.HouseGID, m.Applicant)
	}
	response.Success(c, &resp.AckVO{OK: true})
	return true
}

// History 平台端申请记录（不依赖会话），支持类型/状态/时间筛选
func (s *ShopApplicationService) History(c *gin.Context) {
	var in struct {
		HouseGID int32   `json:"house_gid" binding:"required"`
		Type     *int32  `json:"type"`
		Status   *int32  `json:"status"`
		StartAt  *string `json:"start_at"`
		EndAt    *string `json:"end_at"`
	}
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid params")
		return
	}
	var sPtr, ePtr *time.Time
	parse := func(p *string) *time.Time {
		if p == nil || *p == "" {
			return nil
		}
		if t, err := time.Parse(time.RFC3339, *p); err == nil {
			return &t
		}
		return nil
	}
	sPtr, ePtr = parse(in.StartAt), parse(in.EndAt)
	uid := utils.GetUserID(c)
	list, err := s.userRepo.ListHistory(c.Request.Context(), in.HouseGID, &uid, in.Type, in.Status, sPtr, ePtr)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	// 昵称
	ids := make([]int32, 0, len(list))
	for _, a := range list {
		ids = append(ids, a.Applicant)
	}
	rows, _ := s.users.SelectByPK(c.Request.Context(), ids)
	nameMap := map[int32]string{}
	for _, u := range rows {
		if u != nil {
			nameMap[u.Id] = u.NickName
		}
	}

	out := make([]*resp.ApplicationItemVO, 0, len(list))
	for _, a := range list {
		out = append(out, &resp.ApplicationItemVO{
			ID:          int(a.Id),
			Status:      int(a.Status),
			ApplierID:   int(a.Applicant),
			ApplierGID:  0,
			ApplierName: nameMap[a.Applicant],
			HouseGID:    int(a.HouseGID),
			Type:        int(a.Type),
			AdminUserID: int(a.AdminUID),
			CreatedAt:   a.CreatedAt.UnixMilli(),
		})
	}
	response.Success(c, &resp.ApplicationsVO{Items: out})
}

// ApplyAdmin 发起管理员申请（平台侧持久化）
func (s *ShopApplicationService) ApplyAdmin(c *gin.Context) {
	var in req.ApplyAdminRequest
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid house_gid")
		return
	}
	uid := utils.GetUserID(c)
	// 禁止重复待审
	if ok, err := s.userRepo.ExistsPending(c.Request.Context(), int32(in.HouseGID), uid, 1, 0); err == nil && ok {
		c.JSON(http.StatusConflict, response.Body{Code: ecode.Failed, Msg: "duplicate pending admin application"})
		return
	}
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
	// 禁止重复待审（同一圈主）
	if ok, err := s.userRepo.ExistsPending(c.Request.Context(), int32(in.HouseGID), uid, 2, int32(in.AdminUserID)); err == nil && ok {
		c.JSON(http.StatusConflict, response.Body{Code: ecode.Failed, Msg: "duplicate pending join application"})
		return
	}
	app := &gameModel.UserApplication{HouseGID: in.HouseGID, Applicant: uid, Type: 2, AdminUID: in.AdminUserID, Note: in.Note, Status: 0}
	if err := s.userRepo.Insert(c.Request.Context(), app); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}
