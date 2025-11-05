package game

import (
	gameBiz "battle-tiles/internal/biz/game"
	gameModel "battle-tiles/internal/dal/model/game"
	gameRepo "battle-tiles/internal/dal/repo/game"
	"battle-tiles/internal/infra/plaza"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

// ShopGroupAdminService 管理“我的圈子”与圈子列表
type ShopGroupAdminService struct {
	uc      *gameBiz.ShopGroupAdminUseCase
	mgr     plaza.Manager
	sAdm    gameRepo.GameShopAdminRepo
	grpRepo gameRepo.GameShopGroupAdminRepo
}

func NewShopGroupAdminService(uc *gameBiz.ShopGroupAdminUseCase, mgr plaza.Manager, sAdm gameRepo.GameShopAdminRepo, grpRepo gameRepo.GameShopGroupAdminRepo) *ShopGroupAdminService {
	return &ShopGroupAdminService{uc: uc, mgr: mgr, sAdm: sAdm, grpRepo: grpRepo}
}

func (s *ShopGroupAdminService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/shops").Use(middleware.JWTAuth())
	g.GET("/my/groups", s.MyGroups)
	g.POST("/groups/list", s.ListGroupsByHouse)
	g.POST("/admins/backfill_groups", middleware.RequirePerm("shop:admin:backfill"), s.BackfillGroups)
}

// MyGroups 返回当前用户作为圈管理员的 (house_gid, group_id) 列表
func (s *ShopGroupAdminService) MyGroups(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	arr, err := s.uc.ListByUser(c.Request.Context(), claims.BaseClaims.UserID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	type item struct{ HouseGID, GroupID int32 }
	// 仅返回平台侧映射结果
	out := make([]item, 0, len(arr))
	for _, it := range arr {
		out = append(out, item{HouseGID: it[0], GroupID: it[1]})
	}
	response.Success(c, out)
}

// ListGroupsByHouse 基于在线会话快照，返回指定店铺下当前可见的 group_id 列表
func (s *ShopGroupAdminService) ListGroupsByHouse(c *gin.Context) {
	var in struct {
		HouseGID int `json:"house_gid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil || in.HouseGID <= 0 {
		response.Fail(c, ecode.ParamsFailed, "invalid house_gid")
		return
	}
	// 平台侧：按库中映射返回该店铺下已定义的圈ID（去重）
	rows, err := s.grpRepo.ListByHouse(c.Request.Context(), int32(in.HouseGID))
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	seen := map[int32]struct{}{}
	out := make([]int32, 0, len(rows))
	for _, r := range rows {
		if r == nil || r.GroupID <= 0 {
			continue
		}
		if _, ok := seen[r.GroupID]; ok {
			continue
		}
		seen[r.GroupID] = struct{}{}
		out = append(out, r.GroupID)
	}
	response.Success(c, out)
}

// BackfillGroups 为所有店铺管理员补齐圈管理员映射（若缺失），优先使用共享会话可见圈，否则 group_id=0。
func (s *ShopGroupAdminService) BackfillGroups(c *gin.Context) {
	// 平台策略：每个店铺管理员就是一个圈（唯一），group_id=admin user_id。
	list, err := s.sAdm.ListAll(c.Request.Context())
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	done := 0
	for _, a := range list {
		if a == nil {
			continue
		}
		if s.grpRepo != nil {
			_ = s.grpRepo.Upsert(c.Request.Context(), &gameModel.GameShopGroupAdmin{HouseGID: a.HouseGID, GroupID: a.UserID, UserID: a.UserID, Role: "admin"})
		}
		done++
	}
	response.Success(c, gin.H{"ok": true, "count": done})
}
