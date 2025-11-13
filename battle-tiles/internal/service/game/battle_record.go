// internal/service/game/battle_record.go
package game

import (
	gameBiz "battle-tiles/internal/biz/game"
	"battle-tiles/pkg/plugin/middleware"
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
	g := r.Group("/battles").Use(middleware.JWTAuth())
	
	// 用户查看自己的战绩
	g.POST("/my/list", s.ListMyBattles)
	g.POST("/my/stats", s.GetMyStats)
	
	// 管理员查看店铺战绩（需要权限）
	g.POST("/house/list", middleware.RequirePerm("battles:view"), s.ListHouseBattles)
	g.POST("/house/player/stats", middleware.RequirePerm("battles:view"), s.GetPlayerStats)
}

// ListMyBattlesRequest 查询我的战绩请求
type ListMyBattlesRequest struct {
	HouseGID  int32  `json:"house_gid" binding:"required"`
	StartTime *int64 `json:"start_time"` // Unix timestamp
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
	userID, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, ecode.TokenFailed, nil)
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
		userID.(int32),
		in.HouseGID,
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
	HouseGID  int32  `json:"house_gid" binding:"required"`
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

	// 从 JWT 中获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, ecode.TokenFailed, nil)
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
		userID.(int32),
		in.HouseGID,
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

// ListHouseBattles
// @Summary      查询店铺战绩（管理员）
// @Tags         战绩
// @Accept       json
// @Produce      json
// @Param        in body ListHouseBattlesRequest true "查询参数"
// @Success      200 {object} response.Body{data=[]model.GameBattleRecord}
// @Router       /battles/house/list [post]
func (s *BattleRecordService) ListHouseBattles(c *gin.Context) {
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

