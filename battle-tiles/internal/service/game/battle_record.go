// internal/service/game/battle_record.go
package game

import (
	gameBiz "battle-tiles/internal/biz/game"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"time"

	"github.com/gin-gonic/gin"
)

type BattleRecordService struct {
	uc *gameBiz.BattleRecordUseCase
}

func NewBattleRecordService(uc *gameBiz.BattleRecordUseCase) *BattleRecordService {
	return &BattleRecordService{uc: uc}
}

func (s *BattleRecordService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/battle-query").Use(middleware.JWTAuth())

	// 用户查看自己的战绩
	g.POST("/my/battles", s.ListMyBattles)
	g.POST("/my/balances", s.GetMyBalances)
	g.POST("/my/stats", s.GetMyStats)

	// 管理员查看圈子战绩（需要权限）
	g.POST("/group/battles", middleware.RequirePerm("battles:view"), s.ListGroupBattles)
	g.POST("/group/balances", middleware.RequirePerm("battles:view"), s.ListGroupMemberBalances)
	g.POST("/group/stats", middleware.RequirePerm("battles:view"), s.GetGroupStats)

	// 超级管理员查看店铺统计
	g.POST("/house/stats", middleware.RequirePerm("battles:view"), s.GetHouseStats)
}

// ListMyBattlesRequest 查询我的战绩请求
type ListMyBattlesRequest struct {
	HouseGID  int32  `json:"house_gid"`  // 可选，如果不传则查询所有店铺
	StartTime *int64 `json:"start_time"` // Unix timestamp
	GroupID   *int32 `json:"group_id"`   // 可选，如果不传则查询所有圈子
	EndTime   *int64 `json:"end_time"`   // Unix timestamp
	Page      int32  `json:"page"`
	Size      int32  `json:"size"`
}

// ListMyBattles
// @Summary      查询我的战绩
// @Tags         战绩
// @Accept       json
// @Produce      json
// @Param        in body ListMyBattlesRequest true "查询参数"
// @Success      200 {object} response.Body{data=[]model.GameBattleRecord}
// @Router       /battles/my/list [post]
func (s *BattleRecordService) ListMyBattles(c *gin.Context) {
	var in ListMyBattlesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	// 从 JWT 中获取用户 ID
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	// 转换时间参数
	var start, end *time.Time
	if in.StartTime != nil {
		t := time.Unix(*in.StartTime, 0)
		start = &t
	}
	if in.EndTime != nil {
		t := time.Unix(*in.EndTime, 0)
		end = &t
	}

	// 设置默认分页参数
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}

	// 查询战绩
	records, total, err := s.uc.ListMyBattleRecords(
		c.Request.Context(),
		claims.UserID,
		in.HouseGID,
		in.GroupID,
		start,
		end,
		in.Page,
		in.Size,
	)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, gin.H{
		"list":  records,
		"total": total,
		"page":  in.Page,
		"size":  in.Size,
	})
}

// GetMyStatsRequest 查询我的统计请求
type GetMyStatsRequest struct {
	HouseGID  int32  `json:"house_gid"` // 可选，如果不传则查询所有店铺
	GroupID   *int32 `json:"group_id"`
	StartTime *int64 `json:"start_time"` // Unix timestamp
	EndTime   *int64 `json:"end_time"`   // Unix timestamp
}

// GetMyStats
// @Summary      查询我的战绩统计
// @Tags         战绩
// @Accept       json
// @Produce      json
// @Param        in body GetMyStatsRequest true "查询参数"
// @Success      200 {object} response.Body{data=object}
// @Router       /battles/my/stats [post]
func (s *BattleRecordService) GetMyStats(c *gin.Context) {
	var in GetMyStatsRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	// 转换时间参数
	var start, end *time.Time
	if in.StartTime != nil {
		t := time.Unix(*in.StartTime, 0)
		start = &t
	}
	if in.EndTime != nil {
		t := time.Unix(*in.EndTime, 0)
		end = &t
	}

	// 查询统计
	totalGames, totalScore, totalFee, err := s.uc.GetMyBattleStats(
		c.Request.Context(),
		claims.UserID,
		in.HouseGID,
		in.GroupID,
		start,
		end,
	)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	// 计算平均分
	avgScore := 0.0
	if totalGames > 0 {
		avgScore = float64(totalScore) / float64(totalGames)
	}

	response.Success(c, gin.H{
		"total_games": totalGames,
		"total_score": totalScore,
		"total_fee":   totalFee,
		"avg_score":   avgScore,
	})
}

// ListHouseBattlesRequest 查询店铺战绩请求
type ListHouseBattlesRequest struct {
	HouseGID  int32  `json:"house_gid" binding:"required"`
	GroupID   *int32 `json:"group_id"`
	GameID    *int32 `json:"game_id"`
	StartTime *int64 `json:"start_time"` // Unix timestamp
	EndTime   *int64 `json:"end_time"`   // Unix timestamp
	Page      int32  `json:"page"`
	Size      int32  `json:"size"`
}

// ListGroupBattles
// @Summary      查询圈子战绩（管理员）
// @Tags         战绩
// @Accept       json
// @Produce      json
// @Param        in body ListHouseBattlesRequest true "查询参数"
// @Success      200 {object} response.Body{data=[]model.GameBattleRecord}
// @Router       /battle-query/group/battles [post]
func (s *BattleRecordService) ListGroupBattles(c *gin.Context) {
	var in ListHouseBattlesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	// 转换时间参数
	var start, end *time.Time
	if in.StartTime != nil {
		t := time.Unix(*in.StartTime, 0)
		start = &t
	}
	if in.EndTime != nil {
		t := time.Unix(*in.EndTime, 0)
		end = &t
	}

	// 设置默认分页参数
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}

	// 查询战绩
	records, total, err := s.uc.ListHouseBattleRecords(
		c.Request.Context(),
		in.HouseGID,
		in.GroupID,
		in.GameID,
		start,
		end,
		in.Page,
		in.Size,
	)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, gin.H{
		"list":  records,
		"total": total,
		"page":  in.Page,
		"size":  in.Size,
	})
}

// ListHouseBattles 已弃用，使用 ListGroupBattles 替代
func (s *BattleRecordService) ListHouseBattles(c *gin.Context) {
	s.ListGroupBattles(c)
}

// GetPlayerStatsRequest 查询玩家统计请求
type GetPlayerStatsRequest struct {
	HouseGID     int32  `json:"house_gid" binding:"required"`
	PlayerGameID int32  `json:"player_game_id" binding:"required"`
	StartTime    *int64 `json:"start_time"` // Unix timestamp
	EndTime      *int64 `json:"end_time"`   // Unix timestamp
}

// GetPlayerStats
// @Summary      查询玩家战绩统计（管理员）
// @Tags         战绩
// @Accept       json
// @Produce      json
// @Param        in body GetPlayerStatsRequest true "查询参数"
// @Success      200 {object} response.Body{data=object}
// @Router       /battles/house/player/stats [post]
func (s *BattleRecordService) GetPlayerStats(c *gin.Context) {
	var in GetPlayerStatsRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	// 转换时间参数
	var start, end *time.Time
	if in.StartTime != nil {
		t := time.Unix(*in.StartTime, 0)
		start = &t
	}
	if in.EndTime != nil {
		t := time.Unix(*in.EndTime, 0)
		end = &t
	}

	// 查询统计
	totalGames, totalScore, totalFee, err := s.uc.GetPlayerBattleStats(
		c.Request.Context(),
		in.HouseGID,
		in.PlayerGameID,
		start,
		end,
	)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, gin.H{
		"total_games": totalGames,
		"total_score": totalScore,
		"total_fee":   totalFee,
	})
}

// GetGroupStatsRequest 查询圈子统计请求（重用 GetGroupStatsRequest）
type GetGroupStatsRequestBody struct {
	HouseGID  int32  `json:"house_gid" binding:"required"`
	GroupID   int32  `json:"group_id" binding:"required"`
	StartTime *int64 `json:"start_time"`
	EndTime   *int64 `json:"end_time"`
}

// GetMyBalancesRequest 查询我的余额请求
type GetMyBalancesRequest struct {
	HouseGID int32  `json:"house_gid"`
	GroupID  *int32 `json:"group_id"`
}

// GetMyBalances
// @Summary      查询我的余额
// @Tags         战绩
// @Accept       json
// @Produce      json
// @Param        in body GetMyBalancesRequest true "查询参数"
// @Success      200 {object} response.Body{data=object}
// @Router       /battle-query/my/balances [post]
func (s *BattleRecordService) GetMyBalances(c *gin.Context) {
	var in GetMyBalancesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	// 从 JWT 中获取用户 ID
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}

	// 查询我的余额
	balances, err := s.uc.GetMyBalances(
		c.Request.Context(),
		claims.UserID,
		in.HouseGID,
		in.GroupID,
	)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, gin.H{
		"balances": balances,
	})
}

// ListGroupMemberBalancesRequest 查询圈子成员余额请求
type ListGroupMemberBalancesRequest struct {
	HouseGID int32  `json:"house_gid" binding:"required"`
	GroupID  int32  `json:"group_id" binding:"required"`
	MinYuan  *int32 `json:"min_yuan"`
	MaxYuan  *int32 `json:"max_yuan"`
	Page     int32  `json:"page"`
	Size     int32  `json:"size"`
}

// ListGroupMemberBalances
// @Summary      查询圈子成员余额
// @Tags         战绩
// @Accept       json
// @Produce      json
// @Param        in body ListGroupMemberBalancesRequest true "查询参数"
// @Success      200 {object} response.Body{data=object}
// @Router       /battle-query/group/balances [post]
func (s *BattleRecordService) ListGroupMemberBalances(c *gin.Context) {
	var in ListGroupMemberBalancesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	// 设置默认分页参数
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}

	// 查询圈子成员余额
	balances, total, err := s.uc.ListGroupMemberBalances(
		c.Request.Context(),
		in.HouseGID,
		in.GroupID,
		in.MinYuan,
		in.MaxYuan,
		in.Page,
		in.Size,
	)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, gin.H{
		"list":  balances,
		"total": total,
		"page":  in.Page,
		"size":  in.Size,
	})
}

// GetGroupStatsRequest 查询圈子统计请求
type GetGroupStatsRequest struct {
	HouseGID  int32  `json:"house_gid" binding:"required"`
	GroupID   int32  `json:"group_id" binding:"required"`
	StartTime *int64 `json:"start_time"`
	EndTime   *int64 `json:"end_time"`
}

// GetGroupStats
// @Summary      查询圈子统计
// @Tags         战绩
// @Accept       json
// @Produce      json
// @Param        in body GetGroupStatsRequest true "查询参数"
// @Success      200 {object} response.Body{data=object}
// @Router       /battle-query/group/stats [post]
func (s *BattleRecordService) GetGroupStats(c *gin.Context) {
	var in GetGroupStatsRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	// 转换时间参数
	var start, end *time.Time
	if in.StartTime != nil {
		t := time.Unix(*in.StartTime, 0)
		start = &t
	}
	if in.EndTime != nil {
		t := time.Unix(*in.EndTime, 0)
		end = &t
	}

	// 查询圈子统计
	stats, err := s.uc.GetGroupStats(
		c.Request.Context(),
		in.HouseGID,
		in.GroupID,
		start,
		end,
	)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, stats)
}

// GetHouseStatsRequest 查询店铺统计请求
type GetHouseStatsRequest struct {
	HouseGID  int32  `json:"house_gid" binding:"required"`
	StartTime *int64 `json:"start_time"`
	EndTime   *int64 `json:"end_time"`
}

// GetHouseStats
// @Summary      查询店铺统计
// @Tags         战绩
// @Accept       json
// @Produce      json
// @Param        in body GetHouseStatsRequest true "查询参数"
// @Success      200 {object} response.Body{data=object}
// @Router       /battle-query/house/stats [post]
func (s *BattleRecordService) GetHouseStats(c *gin.Context) {
	var in GetHouseStatsRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	// 转换时间参数
	var start, end *time.Time
	if in.StartTime != nil {
		t := time.Unix(*in.StartTime, 0)
		start = &t
	}
	if in.EndTime != nil {
		t := time.Unix(*in.EndTime, 0)
		end = &t
	}

	// 查询店铺统计
	stats, err := s.uc.GetHouseStats(
		c.Request.Context(),
		in.HouseGID,
		start,
		end,
	)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, stats)
}
