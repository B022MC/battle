package basic

import (
	basicModel "battle-tiles/internal/dal/model/basic"
	basicRepo "battle-tiles/internal/dal/repo/basic"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

// BasicPermissionService 权限服务
type BasicPermissionService struct {
	permRepo basicRepo.PermissionRepo
}

// NewBasicPermissionService 创建权限服务
func NewBasicPermissionService(permRepo basicRepo.PermissionRepo) *BasicPermissionService {
	return &BasicPermissionService{permRepo: permRepo}
}

// RegisterRouter 注册路由
func (s *BasicPermissionService) RegisterRouter(rootRouter *gin.RouterGroup) {
	privateRouter := rootRouter.Group("/basic/permission").Use(middleware.JWTAuth())

	// 权限管理（仅超级管理员）
	privateRouter.GET("/list", middleware.RequireAnyPerm("permission:view"), s.ListPermissions)
	privateRouter.GET("/listAll", middleware.RequireAnyPerm("permission:view"), s.ListAllPermissions)
	privateRouter.POST("/create", middleware.RequirePerm("permission:create"), s.CreatePermission)
	privateRouter.POST("/update", middleware.RequirePerm("permission:update"), s.UpdatePermission)
	privateRouter.POST("/delete", middleware.RequirePerm("permission:delete"), s.DeletePermission)

	// 角色权限管理
	privateRouter.GET("/role/permissions", middleware.RequireAnyPerm("permission:view", "role:view"), s.GetRolePermissions)
	privateRouter.POST("/role/assign", middleware.RequirePerm("permission:assign"), s.AssignPermissionsToRole)
	privateRouter.POST("/role/remove", middleware.RequirePerm("permission:assign"), s.RemovePermissionsFromRole)
}

// ListPermissionsRequest 查询权限列表请求
type ListPermissionsRequest struct {
	Category string `form:"category" json:"category"` // 可选，按分类过滤
}

// ListPermissions 查询权限列表
// @Summary      查询权限列表
// @Tags         基础管理/权限
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        category query string false "权限分类"
// @Success      200 {object} response.Body{data=[]basicModel.BasicPermission}
// @Router       /basic/permission/list [get]
func (s *BasicPermissionService) ListPermissions(c *gin.Context) {
	var req ListPermissionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	permissions, err := s.permRepo.List(c.Request.Context(), req.Category)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, permissions)
}

// ListAllPermissions 查询所有权限
// @Summary      查询所有权限
// @Tags         基础管理/权限
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} response.Body{data=[]basicModel.BasicPermission}
// @Router       /basic/permission/listAll [get]
func (s *BasicPermissionService) ListAllPermissions(c *gin.Context) {
	permissions, err := s.permRepo.ListAll(c.Request.Context())
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, permissions)
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Description string `json:"description"`
}

// CreatePermission 创建权限
// @Summary      创建权限
// @Tags         基础管理/权限
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        in body CreatePermissionRequest true "权限信息"
// @Success      200 {object} response.Body
// @Router       /basic/permission/create [post]
func (s *BasicPermissionService) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	permission := &basicModel.BasicPermission{
		Code:        req.Code,
		Name:        req.Name,
		Category:    req.Category,
		Description: req.Description,
	}

	if err := s.permRepo.Create(c.Request.Context(), permission); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, permission)
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	ID          int32  `json:"id" binding:"required"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// UpdatePermission 更新权限
// @Summary      更新权限
// @Tags         基础管理/权限
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        in body UpdatePermissionRequest true "权限信息"
// @Success      200 {object} response.Body
// @Router       /basic/permission/update [post]
func (s *BasicPermissionService) UpdatePermission(c *gin.Context) {
	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	// 获取现有权限记录
	existing, err := s.permRepo.GetByID(c.Request.Context(), req.ID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	if existing == nil {
		response.Fail(c, ecode.Failed, "permission not found")
		return
	}

	// 更新字段
	existing.Name = req.Name
	existing.Category = req.Category
	existing.Description = req.Description

	if err := s.permRepo.Update(c.Request.Context(), existing); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, nil)
}

// DeletePermissionRequest 删除权限请求
type DeletePermissionRequest struct {
	ID int32 `json:"id" binding:"required"`
}

// DeletePermission 删除权限
// @Summary      删除权限
// @Tags         基础管理/权限
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        in body DeletePermissionRequest true "权限ID"
// @Success      200 {object} response.Body
// @Router       /basic/permission/delete [post]
func (s *BasicPermissionService) DeletePermission(c *gin.Context) {
	var req DeletePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	if err := s.permRepo.Delete(c.Request.Context(), req.ID); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, nil)
}

// GetRolePermissionsRequest 查询角色权限请求
type GetRolePermissionsRequest struct {
	RoleID int32 `form:"role_id" json:"role_id" binding:"required"`
}

// GetRolePermissions 查询角色权限
// @Summary      查询角色权限
// @Tags         基础管理/权限
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        role_id query int32 true "角色ID"
// @Success      200 {object} response.Body{data=[]basicModel.BasicPermission}
// @Router       /basic/permission/role/permissions [get]
func (s *BasicPermissionService) GetRolePermissions(c *gin.Context) {
	var req GetRolePermissionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	permissions, err := s.permRepo.GetRolePermissions(c.Request.Context(), req.RoleID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, permissions)
}

// AssignPermissionsRequest 分配权限请求
type AssignPermissionsRequest struct {
	RoleID        int32   `json:"role_id" binding:"required"`
	PermissionIDs []int32 `json:"permission_ids" binding:"required"`
}

// AssignPermissionsToRole 为角色分配权限
// @Summary      为角色分配权限
// @Tags         基础管理/权限
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        in body AssignPermissionsRequest true "分配信息"
// @Success      200 {object} response.Body
// @Router       /basic/permission/role/assign [post]
func (s *BasicPermissionService) AssignPermissionsToRole(c *gin.Context) {
	var req AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	if err := s.permRepo.AssignPermissionsToRole(c.Request.Context(), req.RoleID, req.PermissionIDs); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, nil)
}

// RemovePermissionsRequest 移除权限请求
type RemovePermissionsRequest struct {
	RoleID        int32   `json:"role_id" binding:"required"`
	PermissionIDs []int32 `json:"permission_ids" binding:"required"`
}

// RemovePermissionsFromRole 从角色移除权限
// @Summary      从角色移除权限
// @Tags         基础管理/权限
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        in body RemovePermissionsRequest true "移除信息"
// @Success      200 {object} response.Body
// @Router       /basic/permission/role/remove [post]
func (s *BasicPermissionService) RemovePermissionsFromRole(c *gin.Context) {
	var req RemovePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	if err := s.permRepo.RemovePermissionsFromRole(c.Request.Context(), req.RoleID, req.PermissionIDs); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, nil)
}
