// internal/service/game/shop_table_service.go
package game

import (
	"battle-tiles/internal/dal/req"
	resp "battle-tiles/internal/dal/resp"
	"battle-tiles/internal/infra/plaza"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShopTableService struct {
	mgr plaza.Manager
}

func NewShopTableService(mgr plaza.Manager) *ShopTableService { return &ShopTableService{mgr: mgr} }

func (s *ShopTableService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/shops").Use(middleware.JWTAuth())

	g.POST("/tables/list", middleware.RequirePerm("shop:table:view"), s.List)
	g.POST("/tables/dismiss", middleware.RequirePerm("shop:table:dismiss"), s.Dismiss)
	g.POST("/tables/check", middleware.RequirePerm("shop:table:view"), s.Check)
	g.POST("/tables/detail", middleware.RequirePerm("shop:table:view"), s.Detail)
	g.POST("/tables/pull", middleware.RequirePerm("shop:table:view"), s.PullTables)
}

// List
// @Summary      会话中的房间列表
// @Description  需要管理员在该店铺下已建立在线会话；返回最近一次 SUB_GA_TABLE_LIST 的本地快照。
// @Tags         店铺/房间
// @Accept       json
// @Produce      json
// @Param        in body req.ListTablesRequest true "入参在body（house_gid）"
// @Success      200 {object} response.Body{data=resp.ShopTableListResponse} "data: { items: []TableInfoVO }"
// @Failure      400 {object} response.Body "参数错误"
// @Failure      401 {object} response.Body "未授权/Token 失效"
// @Failure      409 {object} response.Body "会话不存在或未在线"
// @Router       /shops/tables/list [post]
func (s *ShopTableService) List(c *gin.Context) {
	var in req.ListTablesRequest
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid house_gid")
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	// 取会话：先按当前用户，其次回退到“任意用户在该店铺下的共享会话”（只读）
	sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
	if !ok || sess == nil {
		if shared, ok2 := s.mgr.GetAnyByHouse(in.HouseGID); ok2 && shared != nil {
			sess = shared
		} else {
			c.JSON(http.StatusConflict, response.Body{
				Code: (ecode.Failed),
				Msg:  "session not found or not online for this house",
			})
			return
		}
	}

	// 从会话的快照读取（由 SUB_GA_TABLE_LIST 刷新）
	tables := sess.ListTables()
	out := make([]resp.TableInfoVO, 0, len(tables))
	for _, t := range tables {
		out = append(out, resp.TableInfoVO{
			TableID:   t.TableID,
			MappedNum: t.MappedNum,
			GroupID:   t.GroupID,
			KindID:    t.KindID,
			BaseScore: t.BaseScore,
		})
	}
	response.Success(c, resp.ShopTableListResponse{Items: out})
}

// Dismiss
// @Summary      解散房间
// @Description  `kind_id` 可选；若不传则从当前会话缓存中按 `mapped_num` 推断，不存在则返回 422。
// @Tags         店铺/房间
// @Accept       json
// @Produce      json
// @Param        in body req.DismissTableRequest true "入参在body（house_gid, mapped_num, kind_id 可选）"
// @Success      200 {object} response.Body
// @Failure      400 {object} response.Body "参数错误"
// @Failure      401 {object} response.Body "未授权/Token 失效"
// @Failure      409 {object} response.Body "会话不存在或未在线"
// @Failure      422 {object} response.Body "需要提供 kind_id（缓存中未找到且未提交）"
// @Router       /shops/tables/dismiss [post]
func (s *ShopTableService) Dismiss(c *gin.Context) {
	var in req.DismissTableRequest
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 || in.MappedNum <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid house_gid or mapped_num")
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
	if !ok || sess == nil {
		c.JSON(http.StatusConflict, response.Body{
			Code: (ecode.Failed),
			Msg:  "session not found or not online for this house",
		})
		return
	}

	kindID := in.KindID
	if kindID == 0 {
		if tables := sess.ListTables(); len(tables) > 0 {
			for _, t := range tables {
				if t.MappedNum == in.MappedNum {
					kindID = t.KindID
					break
				}
			}
		}
	}
	if kindID == 0 {
		c.JSON(http.StatusUnprocessableEntity, response.Body{
			Code: (ecode.ParamsFailed),
			Msg:  "kind_id required (not found in session cache and not provided in body)",
		})
		return
	}

	if err := s.mgr.DismissTable(int(claims.BaseClaims.UserID), in.HouseGID, kindID, in.MappedNum); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}

// Check
// @Summary      查询单桌玩家（直透 Plaza）
// @Description  调用会话 QueryTable；仅触发下发，不保证立即返回玩家列表，返回触发状态和最近一次房间快照中该桌是否存在。
// @Tags         店铺/房间
// @Accept       json
// @Produce      json
// @Param        in body req.QueryTableRequest true "house_gid, mapped_num"
// @Success      200 {object} response.Body{data=resp.ShopTableCheckResponse} "data: { triggered: bool, exists_in_cache: bool }"
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Failure      409 {object} response.Body "session not online"
// @Router       /shops/tables/check [post]
func (s *ShopTableService) Check(c *gin.Context) {
	var in req.QueryTableRequest
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 || in.MappedNum <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid house_gid or mapped_num")
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
	if !ok || sess == nil {
		if shared, ok2 := s.mgr.GetAnyByHouse(in.HouseGID); ok2 && shared != nil {
			sess = shared
		} else {
			c.JSON(http.StatusConflict, response.Body{Code: ecode.Failed, Msg: "session not found or not online for this house"})
			return
		}
	}

	// 触发查询
	sess.QueryTable(in.MappedNum)

	// 简要返回：是否在本地快照中存在该桌 + 尝试返回快照中的桌信息
	exists := false
	var tableSnap *resp.TableInfoVO
	if tables := sess.ListTables(); len(tables) > 0 {
		for _, t := range tables {
			if t.MappedNum == in.MappedNum {
				exists = true
				tableSnap = &resp.TableInfoVO{TableID: t.TableID, MappedNum: t.MappedNum, GroupID: t.GroupID, KindID: t.KindID, BaseScore: t.BaseScore}
				break
			}
		}
	}
	response.Success(c, resp.ShopTableCheckResponse{Triggered: true, ExistsInCache: exists, Table: tableSnap})
}

// Detail
// @Summary      查桌详情（快照+触发刷新）
// @Description  从会话缓存读取该桌信息，并触发一次 QueryTable(mapped_num)；若缓存中不存在则返回 404。
// @Tags         店铺/房间
// @Accept       json
// @Produce      json
// @Param        in body req.QueryTableRequest true "house_gid, mapped_num"
// @Success      200 {object} response.Body{data=resp.ShopTableDetailResponse} "data: { table: TableInfoVO, triggered: bool }"
// @Failure      400 {object} response.Body
// @Failure      401 {object} response.Body
// @Failure      404 {object} response.Body "table not found in cache"
// @Failure      409 {object} response.Body "session not online"
// @Router       /shops/tables/detail [post]
func (s *ShopTableService) Detail(c *gin.Context) {
	var in req.QueryTableRequest
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 || in.MappedNum <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid house_gid or mapped_num")
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
	if !ok || sess == nil {
		if shared, ok2 := s.mgr.GetAnyByHouse(in.HouseGID); ok2 && shared != nil {
			sess = shared
		} else {
			c.JSON(http.StatusConflict, response.Body{Code: ecode.Failed, Msg: "session not found or not online for this house"})
			return
		}
	}

	// 找到当前快照中的桌信息
	var tableSnap *resp.TableInfoVO
	if tables := sess.ListTables(); len(tables) > 0 {
		for _, t := range tables {
			if t.MappedNum == in.MappedNum {
				tableSnap = &resp.TableInfoVO{TableID: t.TableID, MappedNum: t.MappedNum, GroupID: t.GroupID, KindID: t.KindID, BaseScore: t.BaseScore}
				break
			}
		}
	}
	if tableSnap == nil {
		// 不再触发底层查询，避免协议异常；仅返回 404
		c.JSON(http.StatusNotFound, response.Body{Code: ecode.Failed, Msg: "table not found in cache"})
		return
	}

	// 可选：如需轻量刷新，可在未来加白名单映射后触发
	response.Success(c, resp.ShopTableDetailResponse{Table: tableSnap, Triggered: false})
}

// PullTables
// @Summary      手动刷新房间列表（触发拉取）
// @Tags         店铺/房间
// @Accept       json
// @Produce      json
// @Param        in body req.ListTablesRequest true "house_gid"
// @Success      200 {object} response.Body
// @Router       /shops/tables/pull [post]
func (s *ShopTableService) PullTables(c *gin.Context) {
	var in req.ListTablesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID)
	if !ok || sess == nil {
		if shared, ok2 := s.mgr.GetAnyByHouse(in.HouseGID); ok2 && shared != nil {
			sess = shared
		} else {
			response.Fail(c, ecode.Failed, "session not found or not online for this house")
			return
		}
	}

	// 改为触发进入圈，促使服务端下发成员与房间列表
	sess.GetGroupMembers()
	response.SuccessWithOK(c)
}
