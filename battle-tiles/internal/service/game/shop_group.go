package game

import (
	"battle-tiles/internal/biz/game"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type ShopGroupService struct {
	groupUC        *game.ShopGroupUseCase
	accountGroupUC *game.GameAccountGroupUseCase
	log            *log.Helper
}

func NewShopGroupService(
	groupUC *game.ShopGroupUseCase,
	accountGroupUC *game.GameAccountGroupUseCase,
	logger log.Logger,
) *ShopGroupService {
	return &ShopGroupService{
		groupUC:        groupUC,
		accountGroupUC: accountGroupUC,
		log:            log.NewHelper(log.With(logger, "module", "service/shop_group")),
	}
}

// RegisterRouter 注册路由
func (s *ShopGroupService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/groups").Use(middleware.JWTAuth())

	g.POST("/create", s.CreateGroup)      // 创建圈子
	g.POST("/my", s.GetMyGroup)           // 获取我的圈子
	g.POST("/list", s.ListGroupsByHouse)  // 获取店铺圈子列表
	g.POST("/options", s.GetGroupOptions) // 获取圈子选项列表（用于下拉框）
	// 已删除 /members/add 接口，统一使用 /shops/members/pull-to-group
	g.POST("/game-accounts/add", s.AddGameAccounts) // 添加游戏账号到圈子（通过游戏账号ID）
	g.POST("/members/remove", s.RemoveMember)       // 从圈子移除成员
	g.POST("/members/list", s.ListMembers)          // 获取圈子成员列表
	g.POST("/my/list", s.ListMyGroups)              // 获取我加入的圈子
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

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	userID := claims.BaseClaims.UserID

	group, err := s.groupUC.CreateGroup(c.Request.Context(), req.HouseGID, userID, req.GroupName, req.Description)
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

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	userID := claims.BaseClaims.UserID

	group, err := s.groupUC.GetMyGroup(c.Request.Context(), req.HouseGID, userID)
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

// 已删除 AddMembersReq 和 AddMembers 方法
// 现在统一使用 /shops/members/pull-to-group 接口进行拉圈操作

// AddGameAccountsReq 添加游戏账号到圈子请求
type AddGameAccountsReq struct {
	GroupID        int32   `json:"group_id" binding:"required"`
	GameAccountIDs []int32 `json:"game_account_ids" binding:"required"`
}

// AddGameAccounts 通过游戏账号ID添加成员到圈子（用于添加游戏内的玩家）
// POST /api/groups/game-accounts/add
func (s *ShopGroupService) AddGameAccounts(c *gin.Context) {
	var req AddGameAccountsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, nil)
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	userID := claims.BaseClaims.UserID

	err = s.groupUC.AddGameAccountsToGroup(c.Request.Context(), req.GroupID, userID, req.GameAccountIDs)
	if err != nil {
		s.log.Errorf("add game accounts failed: %v", err)
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

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	userID := claims.BaseClaims.UserID

	err = s.groupUC.RemoveMemberFromGroup(c.Request.Context(), req.GroupID, userID, req.UserID)
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

// ListMyGroups 获取我加入的所有圈子（通过游戏账号反向查询）
// POST /api/groups/my/list
func (s *ShopGroupService) ListMyGroups(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	userID := claims.BaseClaims.UserID

	// 使用新的游戏账号反向查询逻辑
	accountGroups, err := s.accountGroupUC.ListGroupsByUser(c.Request.Context(), userID)
	if err != nil {
		s.log.Errorf("list my groups failed: %v", err)
		response.Fail(c, ecode.Failed, err.Error())
		return
	}

	// 转换为圈子列表（去重）
	groupMap := make(map[int32]bool)
	var groupIDs []int32
	for _, ag := range accountGroups {
		if !groupMap[ag.GroupID] {
			groupMap[ag.GroupID] = true
			groupIDs = append(groupIDs, ag.GroupID)
		}
	}

	// 查询圈子详情
	var groups []map[string]interface{}
	for _, ag := range accountGroups {
		if groupMap[ag.GroupID] {
			groups = append(groups, map[string]interface{}{
				"id":            ag.GroupID,
				"house_gid":     ag.HouseGID,
				"group_name":    ag.GroupName,
				"admin_user_id": ag.AdminUserID,
				"status":        ag.Status,
				"joined_at":     ag.JoinedAt,
			})
			delete(groupMap, ag.GroupID) // 避免重复
		}
	}

	response.Success(c, groups)
}

// GetGroupOptionsReq 获取圈子选项请求
type GetGroupOptionsReq struct {
	HouseGID int32 `json:"house_gid" binding:"required"`
}

// GetGroupOptions 获取圈子选项列表（用于下拉框）
// POST /api/groups/options
func (s *ShopGroupService) GetGroupOptions(c *gin.Context) {
	var req GetGroupOptionsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, nil)
		return
	}

	groups, err := s.groupUC.ListGroupsByHouse(c.Request.Context(), req.HouseGID)
	if err != nil {
		s.log.Errorf("get group options failed: %v", err)
		response.Fail(c, ecode.Failed, err.Error())
		return
	}

	// 返回只包含 ID 和名称的列表
	options := make([]map[string]interface{}, 0, len(groups))
	for _, group := range groups {
		options = append(options, map[string]interface{}{
			"id":   group.Id,
			"name": group.GroupName,
		})
	}

	response.Success(c, options)
}
