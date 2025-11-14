package game

import (
	"battle-tiles/internal/biz/game"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type ShopGroupService struct {
	groupUC *game.ShopGroupUseCase
	log     *log.Helper
}

func NewShopGroupService(groupUC *game.ShopGroupUseCase, logger log.Logger) *ShopGroupService {
	return &ShopGroupService{
		groupUC: groupUC,
		log:     log.NewHelper(log.With(logger, "module", "service/shop_group")),
	}
}

// RegisterRouter 注册路由
func (s *ShopGroupService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/groups").Use(middleware.JWTAuth())

	g.POST("/create", s.CreateGroup)          // 创建圈子
	g.POST("/my", s.GetMyGroup)               // 获取我的圈子
	g.POST("/list", s.ListGroupsByHouse)      // 获取店铺圈子列表
	g.POST("/members/add", s.AddMembers)      // 添加成员到圈子
	g.POST("/members/remove", s.RemoveMember) // 从圈子移除成员
	g.POST("/members/list", s.ListMembers)    // 获取圈子成员列表
	g.POST("/my/list", s.ListMyGroups)        // 获取我加入的圈子
}

// CreateGroupReq 创建圈子请求
type CreateGroupReq struct {
	HouseGID    int32  `json:"house_gid" binding:"required"`
	GroupName   string `json:"group_name" binding:"required"`
	Description string `json:"description"`
}

// CreateGroup 创建圈子
// POST /api/groups/create
func (s *ShopGroupService) CreateGroup(c *gin.Context) {
	var req CreateGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, nil)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, ecode.TokenFailed, nil)
		return
	}

	group, err := s.groupUC.CreateGroup(c.Request.Context(), req.HouseGID, userID.(int32), req.GroupName, req.Description)
	if err != nil {
		s.log.Errorf("create group failed: %v", err)
		response.Fail(c, ecode.Failed, err.Error())
		return
	}

	response.Success(c, group)
}

// GetMyGroupReq 获取我的圈子请求
type GetMyGroupReq struct {
	HouseGID int32 `json:"house_gid" binding:"required"`
}

// GetMyGroup 获取我的圈子
// POST /api/groups/my
func (s *ShopGroupService) GetMyGroup(c *gin.Context) {
	var req GetMyGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, nil)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, ecode.TokenFailed, nil)
		return
	}

	group, err := s.groupUC.GetMyGroup(c.Request.Context(), req.HouseGID, userID.(int32))
	if err != nil {
		s.log.Errorf("get my group failed: %v", err)
		response.Fail(c, ecode.Failed, err.Error())
		return
	}

	// 获取成员数量
	_, count, err := s.groupUC.GetGroupWithMemberCount(c.Request.Context(), group.Id)
	if err != nil {
		s.log.Errorf("get member count failed: %v", err)
	}

	result := map[string]interface{}{
		"id":            group.Id,
		"house_gid":     group.HouseGID,
		"group_name":    group.GroupName,
		"admin_user_id": group.AdminUserID,
		"description":   group.Description,
		"member_count":  count,
		"created_at":    group.CreatedAt,
	}

	response.Success(c, result)
}

// ListGroupsByHouseReq 获取店铺圈子列表请求
type ListGroupsByHouseReq struct {
	HouseGID int32 `json:"house_gid" binding:"required"`
}

// ListGroupsByHouse 获取店铺的所有圈子
// POST /api/groups/list
func (s *ShopGroupService) ListGroupsByHouse(c *gin.Context) {
	var req ListGroupsByHouseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, nil)
		return
	}

	groups, err := s.groupUC.ListGroupsByHouse(c.Request.Context(), req.HouseGID)
	if err != nil {
		s.log.Errorf("list groups failed: %v", err)
		response.Fail(c, ecode.Failed, err.Error())
		return
	}

	// 为每个圈子添加成员数量
	result := make([]map[string]interface{}, 0, len(groups))
	for _, group := range groups {
		_, count, _ := s.groupUC.GetGroupWithMemberCount(c.Request.Context(), group.Id)
		result = append(result, map[string]interface{}{
			"id":            group.Id,
			"house_gid":     group.HouseGID,
			"group_name":    group.GroupName,
			"admin_user_id": group.AdminUserID,
			"description":   group.Description,
			"member_count":  count,
			"created_at":    group.CreatedAt,
		})
	}

	response.Success(c, result)
}

// AddMembersReq 添加成员请求
type AddMembersReq struct {
	GroupID int32   `json:"group_id" binding:"required"`
	UserIDs []int32 `json:"user_ids" binding:"required"`
}

// AddMembers 添加成员到圈子
// POST /api/groups/members/add
func (s *ShopGroupService) AddMembers(c *gin.Context) {
	var req AddMembersReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, nil)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, ecode.TokenFailed, nil)
		return
	}

	err := s.groupUC.AddMembersToGroup(c.Request.Context(), req.GroupID, userID.(int32), req.UserIDs)
	if err != nil {
		s.log.Errorf("add members failed: %v", err)
		response.Fail(c, ecode.Failed, err.Error())
		return
	}

	response.Success(c, "添加成功")
}

// RemoveMemberReq 移除成员请求
type RemoveMemberReq struct {
	GroupID int32 `json:"group_id" binding:"required"`
	UserID  int32 `json:"user_id" binding:"required"`
}

// RemoveMember 从圈子移除成员
// POST /api/groups/members/remove
func (s *ShopGroupService) RemoveMember(c *gin.Context) {
	var req RemoveMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, nil)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, ecode.TokenFailed, nil)
		return
	}

	err := s.groupUC.RemoveMemberFromGroup(c.Request.Context(), req.GroupID, userID.(int32), req.UserID)
	if err != nil {
		s.log.Errorf("remove member failed: %v", err)
		response.Fail(c, ecode.Failed, err.Error())
		return
	}

	response.Success(c, "移除成功")
}

// ListMembersReq 获取圈子成员列表请求
type ListMembersReq struct {
	GroupID int32 `json:"group_id" binding:"required"`
	Page    int32 `json:"page"`
	Size    int32 `json:"size"`
}

// ListMembers 获取圈子成员列表
// POST /api/groups/members/list
func (s *ShopGroupService) ListMembers(c *gin.Context) {
	var req ListMembersReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, nil)
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	users, total, err := s.groupUC.ListGroupMembers(c.Request.Context(), req.GroupID, req.Page, req.Size)
	if err != nil {
		s.log.Errorf("list members failed: %v", err)
		response.Fail(c, ecode.Failed, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"items": users,
		"total": total,
		"page":  req.Page,
		"size":  req.Size,
	})
}

// ListMyGroups 获取我加入的所有圈子
// POST /api/groups/my/list
func (s *ShopGroupService) ListMyGroups(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, ecode.TokenFailed, nil)
		return
	}

	groups, err := s.groupUC.ListMyGroups(c.Request.Context(), userID.(int32))
	if err != nil {
		s.log.Errorf("list my groups failed: %v", err)
		response.Fail(c, ecode.Failed, err.Error())
		return
	}

	response.Success(c, groups)
}
