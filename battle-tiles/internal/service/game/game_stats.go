package game

import (
	gameBiz "battle-tiles/internal/biz/game"
	"battle-tiles/internal/dal/req"
	"battle-tiles/internal/infra/plaza"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type GameStatsService struct {
	uc          *gameBiz.GameStatsUseCase
	shopAdminUC *gameBiz.ShopAdminUseCase
	mgr         plaza.Manager
}

func NewGameStatsService(uc *gameBiz.GameStatsUseCase, shopAdminUC *gameBiz.ShopAdminUseCase, mgr plaza.Manager) *GameStatsService {
	return &GameStatsService{uc: uc, shopAdminUC: shopAdminUC, mgr: mgr}
}
func (s *GameStatsService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/stats").Use(middleware.JWTAuth())
	g.POST("/today", middleware.RequirePerm("stats:view"), s.Today)
	g.POST("/week", middleware.RequirePerm("stats:view"), s.Week)
	g.POST("/yesterday", middleware.RequirePerm("stats:view"), s.Yesterday)
	g.POST("/lastweek", middleware.RequirePerm("stats:view"), s.LastWeek)
	// 用户维度
	g.POST("/member/today", middleware.RequirePerm("stats:view"), s.MemberToday)
	g.POST("/member/yesterday", middleware.RequirePerm("stats:view"), s.MemberYesterday)
	g.POST("/member/thisweek", middleware.RequirePerm("stats:view"), s.MemberThisWeek)
	g.POST("/member/lastweek", middleware.RequirePerm("stats:view"), s.MemberLastWeek)
	// 按店铺会话活跃
	g.GET("/sessions/activeByHouse", s.ActiveByHouse)
}

// Today
// @Summary      今日统计
// @Tags         统计
// @Accept       json
// @Produce      json
// @Param        in body req.StatsBaseRequest true "house_gid"
// @Success      200 {object} response.Body{data=resp.StatsVO}
// @Router       /stats/today [post]
func (s *GameStatsService) Today(c *gin.Context) {
	var in req.StatsBaseRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	out, err := s.handleStats(c, in)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, out)
}

// Week
// @Summary      本周统计（近7天）
// @Tags         统计
// @Accept       json
// @Produce      json
// @Param        in body req.StatsBaseRequest true "house_gid"
// @Success      200 {object} response.Body{data=resp.StatsVO}
// @Router       /stats/week [post]
func (s *GameStatsService) Week(c *gin.Context) {
	var in req.StatsBaseRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	out, err := s.handleStats(c, in)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, out)
}

// Yesterday
// @Summary      昨日统计
// @Tags         统计
// @Accept       json
// @Produce      json
// @Param        in body req.StatsBaseRequest true "house_gid"
// @Success      200 {object} response.Body{data=resp.StatsVO}
// @Router       /stats/yesterday [post]
func (s *GameStatsService) Yesterday(c *gin.Context) {
	var in req.StatsBaseRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	out, err := s.handleStats(c, in)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, out)
}

// LastWeek
// @Summary      上周统计（上周一至本周一）
// @Tags         统计
// @Accept       json
// @Produce      json
// @Param        in body req.StatsBaseRequest true "house_gid"
// @Success      200 {object} response.Body{data=resp.StatsVO}
// @Router       /stats/lastweek [post]
func (s *GameStatsService) LastWeek(c *gin.Context) {
	var in req.StatsBaseRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	out, err := s.handleStats(c, in)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, out)
}

// 统一处理带可选 group_id 的统计
func (s *GameStatsService) handleStats(c *gin.Context, in req.StatsBaseRequest) (interface{}, error) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		return nil, err
	}
	isSuper := false
	for _, r := range claims.BaseClaims.Roles {
		if r == 1 {
			isSuper = true
			break
		}
	}

	// 若带 group_id，做圈级权限校验并收集成员ID；否则按店铺级处理
	if in.GroupID != nil && *in.GroupID > 0 {
		// 非超管需要 (house, group) 管理权限；若无圈级权限但有店铺管理员也放行（可调节策略）
		if !isSuper {
			okHouse := false
			if s.shopAdminUC != nil {
				if v, err := s.shopAdminUC.IsAdmin(c.Request.Context(), int32(in.HouseGID), claims.BaseClaims.UserID); err == nil && v {
					okHouse = true
				}
			}
			// 圈子权限检查已移除，使用新的圈子系统
			// 暂时只检查店铺管理员权限
			if !okHouse {
				return nil, fmt.Errorf("permission denied for group")
			}
		}
		// 从在线会话快照提取成员ID集合（当前快照无 groupId 维度，暂以全量成员近似）
		memberIDs := []int{}
		if sess, ok := s.mgr.Get(int(claims.BaseClaims.UserID), in.HouseGID); ok && sess != nil {
			for _, m := range sess.ListMembers() {
				memberIDs = append(memberIDs, int(m.MemberID))
			}
		}
		// 根据路由调用对应窗口期
		switch c.FullPath() {
		case "/stats/today":
			now := time.Now()
			y, m, d := now.Date()
			loc := now.Location()
			start := time.Date(y, m, d, 0, 0, 0, 0, loc)
			end := start.AddDate(0, 0, 1)
			return s.uc.AggregateByMembers(c.Request.Context(), in.HouseGID, memberIDs, start, end)
		case "/stats/yesterday":
			now := time.Now()
			y, m, d := now.Date()
			loc := now.Location()
			todayStart := time.Date(y, m, d, 0, 0, 0, 0, loc)
			start := todayStart.AddDate(0, 0, -1)
			end := todayStart
			return s.uc.AggregateByMembers(c.Request.Context(), in.HouseGID, memberIDs, start, end)
		case "/stats/week":
			now := time.Now()
			y, m, d := now.Date()
			loc := now.Location()
			todayStart := time.Date(y, m, d, 0, 0, 0, 0, loc)
			start := todayStart.AddDate(0, 0, -6)
			end := todayStart.AddDate(0, 0, 1)
			return s.uc.AggregateByMembers(c.Request.Context(), in.HouseGID, memberIDs, start, end)
		case "/stats/lastweek":
			// 上周一至本周一
			now := time.Now()
			y, m, d := now.Date()
			loc := now.Location()
			todayStart := time.Date(y, m, d, 0, 0, 0, 0, loc)
			weekday := int(todayStart.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			thisMonday := todayStart.AddDate(0, 0, 1-weekday)
			lastMonday := thisMonday.AddDate(0, 0, -7)
			return s.uc.AggregateByMembers(c.Request.Context(), in.HouseGID, memberIDs, lastMonday, thisMonday)
		}
	} else {
		// 店铺级：非超管需是店铺管理员
		if !isSuper && s.shopAdminUC != nil {
			if v, err := s.shopAdminUC.IsAdmin(c.Request.Context(), int32(in.HouseGID), claims.BaseClaims.UserID); err != nil || !v {
				return nil, fmt.Errorf("permission denied for house")
			}
		}
		switch c.FullPath() {
		case "/stats/today":
			return s.uc.Today(c.Request.Context(), in.HouseGID)
		case "/stats/yesterday":
			return s.uc.Yesterday(c.Request.Context(), in.HouseGID)
		case "/stats/week":
			return s.uc.Week(c.Request.Context(), in.HouseGID)
		case "/stats/lastweek":
			return s.uc.LastWeek(c.Request.Context(), in.HouseGID)
		}
	}
	return nil, fmt.Errorf("unsupported path")
}

// MemberToday
// @Summary      成员-今日统计
// @Tags         统计
// @Accept       json
// @Produce      json
// @Param        in body req.MemberStatsRequest true "house_gid, member_id"
// @Success      200 {object} response.Body{data=resp.StatsVO}
// @Router       /stats/member/today [post]
func (s *GameStatsService) MemberToday(c *gin.Context) {
	var in req.MemberStatsRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	out, err := s.uc.MemberToday(c.Request.Context(), in.HouseGID, in.MemberID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, out)
}

// MemberYesterday
func (s *GameStatsService) MemberYesterday(c *gin.Context) {
	var in req.MemberStatsRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	out, err := s.uc.MemberYesterday(c.Request.Context(), in.HouseGID, in.MemberID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, out)
}

// MemberThisWeek
func (s *GameStatsService) MemberThisWeek(c *gin.Context) {
	var in req.MemberStatsRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	out, err := s.uc.MemberThisWeek(c.Request.Context(), in.HouseGID, in.MemberID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, out)
}

// MemberLastWeek
func (s *GameStatsService) MemberLastWeek(c *gin.Context) {
	var in req.MemberStatsRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	out, err := s.uc.MemberLastWeek(c.Request.Context(), in.HouseGID, in.MemberID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, out)
}

// ActiveByHouse 列出所有店铺的在线会话数
// @Summary      按店铺列出在线会话数
// @Tags         统计
// @Produce      json
// @Success      200 {object} response.Body{data=[]resp.ActiveByHouseVO}
// @Router       /stats/sessions/activeByHouse [get]
func (s *GameStatsService) ActiveByHouse(c *gin.Context) {
	out, err := s.uc.ListActiveByHouse(c.Request.Context())
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, out)
}
