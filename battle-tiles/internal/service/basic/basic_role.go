package basic

import (
	"battle-tiles/internal/infra"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"time"

	"github.com/gin-gonic/gin"
)

// BasicRoleService 提供角色查询接口
type BasicRoleService struct {
	data *infra.Data
}

func NewBasicRoleService(data *infra.Data) *BasicRoleService {
	return &BasicRoleService{data: data}
}

func (s *BasicRoleService) RegisterRouter(root *gin.RouterGroup) {
	r := root.Group("/basic/role").Use(middleware.JWTAuth())
	
	// 查询接口
	r.GET("/list", middleware.RequireAnyPerm("role:view"), s.List)
	r.GET("/getOne", middleware.RequireAnyPerm("role:view"), s.GetOne)
	r.GET("/all", middleware.RequireAnyPerm("role:view"), s.All)
	
	// 管理接口（需要权限）
	r.POST("/create", middleware.RequirePerm("role:create"), s.Create)
	r.POST("/update", middleware.RequirePerm("role:update"), s.Update)
	r.POST("/delete", middleware.RequirePerm("role:delete"), s.Delete)
	
	// 角色菜单关联
	r.GET("/menus", middleware.RequireAnyPerm("role:view"), s.GetRoleMenus)
	r.POST("/menus/assign", middleware.RequirePerm("role:update"), s.AssignMenus)
}

type roleDoc struct {
	ID          int32   `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	ParentID    int32   `json:"parent_id"`
	Remark      *string `json:"remark"`
	CreatedAt   string  `json:"created_at"`
	CreatedUser *int32  `json:"created_user"`
	UpdatedAt   *string `json:"updated_at"`
	UpdatedUser *int32  `json:"updated_user"`
	FirstLetter string  `json:"first_letter"`
	PinyinCode  string  `json:"pinyin_code"`
	Enable      bool    `json:"enable"`
	IsDeleted   bool    `json:"is_deleted"`
}

// List 分页查询角色
// @Summary      角色列表
// @Tags         基础管理/角色
// @Security     BearerAuth
// @Produce      json
// @Router       /basic/role/list [get]
func (s *BasicRoleService) List(c *gin.Context) {
	type query struct {
		Keyword  string `form:"keyword"`
		Enable   *bool  `form:"enable"`
		PageNo   int    `form:"page_no,default=1"`
		PageSize int    `form:"page_size,default=20"`
	}
	var q query
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, 400, err)
		return
	}

	db := s.data.GetDBWithContext(c.Request.Context())
	if db.Error != nil {
		response.Fail(c, 500, db.Error)
		return
	}

	queryDB := db.Table("basic_role").Where("is_deleted = false")
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		queryDB = queryDB.Where("code ILIKE ? OR name ILIKE ?", like, like)
	}
	if q.Enable != nil {
		queryDB = queryDB.Where("enable = ?", *q.Enable)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		response.Fail(c, 500, err)
		return
	}

	if q.PageNo <= 0 {
		q.PageNo = 1
	}
	if q.PageSize <= 0 || q.PageSize > 1000 {
		q.PageSize = 20
	}

	var list []roleDoc
	if err := queryDB.Order("id DESC").Offset((q.PageNo - 1) * q.PageSize).Limit(q.PageSize).Find(&list).Error; err != nil {
		response.Fail(c, 500, err)
		return
	}

	response.Success(c, gin.H{
		"list":      list,
		"page_no":   q.PageNo,
		"page_size": q.PageSize,
		"total":     total,
	})
}

// GetOne 查询单条角色
// @Summary      角色详情
// @Tags         基础管理/角色
// @Security     BearerAuth
// @Produce      json
// @Router       /basic/role/getOne [get]
func (s *BasicRoleService) GetOne(c *gin.Context) {
	type idq struct {
		ID int32 `form:"id" binding:"required"`
	}
	var in idq
	if err := c.ShouldBindQuery(&in); err != nil {
		response.Fail(c, 400, err)
		return
	}
	db := s.data.GetDBWithContext(c.Request.Context())
	if db.Error != nil {
		response.Fail(c, 500, db.Error)
		return
	}
	var one roleDoc
	if err := db.Table("basic_role").Where("id = ? AND is_deleted = false", in.ID).First(&one).Error; err != nil {
		response.Fail(c, 404, err)
		return
	}
	response.Success(c, one)
}

// All 返回全部启用且未删除的角色（不分页）
// @Summary      角色全量（启用）
// @Tags         基础管理/角色
// @Security     BearerAuth
// @Produce      json
// @Router       /basic/role/all [get]
func (s *BasicRoleService) All(c *gin.Context) {
	db := s.data.GetDBWithContext(c.Request.Context())
	if db.Error != nil {
		response.Fail(c, 500, db.Error)
		return
	}
	var list []roleDoc
	if err := db.Table("basic_role").Where("is_deleted = false AND enable = true").Order("id").Find(&list).Error; err != nil {
		response.Fail(c, 500, err)
		return
	}
	response.Success(c, gin.H{"list": list})
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Code   string  `json:"code" binding:"required"`
	Name   string  `json:"name" binding:"required"`
	Remark *string `json:"remark"`
}

// Create 创建角色
// @Summary      创建角色
// @Tags         基础管理/角色
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        in body CreateRoleRequest true "角色信息"
// @Success      200 {object} response.Body
// @Router       /basic/role/create [post]
func (s *BasicRoleService) Create(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	db := s.data.GetDBWithContext(c.Request.Context())
	if db.Error != nil {
		response.Fail(c, ecode.Failed, db.Error)
		return
	}

	// 检查编码是否已存在
	var count int64
	if err := db.Table("basic_role").Where("code = ? AND is_deleted = false", req.Code).Count(&count).Error; err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	if count > 0 {
		response.Fail(c, ecode.Failed, "角色编码已存在")
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	var createdUser *int32
	if exists {
		if uid, ok := userID.(int32); ok {
			createdUser = &uid
		}
	}

	now := time.Now()
	role := map[string]interface{}{
		"code":         req.Code,
		"name":         req.Name,
		"parent_id":    -1,
		"remark":       req.Remark,
		"created_at":   now,
		"created_user": createdUser,
		"updated_at":   now,
		"updated_user": createdUser,
		"enable":       true,
		"is_deleted":   false,
	}

	if err := db.Table("basic_role").Create(&role).Error; err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, role)
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	ID     int32   `json:"id" binding:"required"`
	Name   string  `json:"name"`
	Remark *string `json:"remark"`
	Enable *bool   `json:"enable"`
}

// Update 更新角色
// @Summary      更新角色
// @Tags         基础管理/角色
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        in body UpdateRoleRequest true "角色信息"
// @Success      200 {object} response.Body
// @Router       /basic/role/update [post]
func (s *BasicRoleService) Update(c *gin.Context) {
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	db := s.data.GetDBWithContext(c.Request.Context())
	if db.Error != nil {
		response.Fail(c, ecode.Failed, db.Error)
		return
	}

	// 检查角色是否存在
	var exists int64
	if err := db.Table("basic_role").Where("id = ? AND is_deleted = false", req.ID).Count(&exists).Error; err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	if exists == 0 {
		response.Fail(c, ecode.Failed, "角色不存在")
		return
	}

	// 获取当前用户ID
	userID, existsUser := c.Get("user_id")
	var updatedUser *int32
	if existsUser {
		if uid, ok := userID.(int32); ok {
			updatedUser = &uid
		}
	}

	updates := map[string]interface{}{
		"updated_at":   time.Now(),
		"updated_user": updatedUser,
	}

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Remark != nil {
		updates["remark"] = req.Remark
	}
	if req.Enable != nil {
		updates["enable"] = *req.Enable
	}

	if err := db.Table("basic_role").Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, nil)
}

// DeleteRoleRequest 删除角色请求
type DeleteRoleRequest struct {
	ID int32 `json:"id" binding:"required"`
}

// Delete 删除角色（软删除）
// @Summary      删除角色
// @Tags         基础管理/角色
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        in body DeleteRoleRequest true "角色ID"
// @Success      200 {object} response.Body
// @Router       /basic/role/delete [post]
func (s *BasicRoleService) Delete(c *gin.Context) {
	var req DeleteRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	// 不允许删除系统预定义角色
	if req.ID == 1 || req.ID == 2 || req.ID == 3 {
		response.Fail(c, ecode.Failed, "不能删除系统预定义角色")
		return
	}

	db := s.data.GetDBWithContext(c.Request.Context())
	if db.Error != nil {
		response.Fail(c, ecode.Failed, db.Error)
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	var updatedUser *int32
	if exists {
		if uid, ok := userID.(int32); ok {
			updatedUser = &uid
		}
	}

	updates := map[string]interface{}{
		"is_deleted":   true,
		"updated_at":   time.Now(),
		"updated_user": updatedUser,
	}

	if err := db.Table("basic_role").Where("id = ? AND is_deleted = false", req.ID).Updates(updates).Error; err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, nil)
}

// GetRoleMenusRequest 查询角色菜单请求
type GetRoleMenusRequest struct {
	RoleID int32 `form:"role_id" binding:"required"`
}

// GetRoleMenus 查询角色的菜单权限
// @Summary      查询角色菜单
// @Tags         基础管理/角色
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        role_id query int32 true "角色ID"
// @Success      200 {object} response.Body
// @Router       /basic/role/menus [get]
func (s *BasicRoleService) GetRoleMenus(c *gin.Context) {
	var req GetRoleMenusRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	db := s.data.GetDBWithContext(c.Request.Context())
	if db.Error != nil {
		response.Fail(c, ecode.Failed, db.Error)
		return
	}

	var menuIDs []int32
	if err := db.Table("basic_role_menu_rel").
		Where("role_id = ?", req.RoleID).
		Pluck("menu_id", &menuIDs).Error; err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, gin.H{"menu_ids": menuIDs})
}

// AssignMenusRequest 分配菜单请求
type AssignMenusRequest struct {
	RoleID  int32   `json:"role_id" binding:"required"`
	MenuIDs []int32 `json:"menu_ids" binding:"required"`
}

// AssignMenus 为角色分配菜单
// @Summary      为角色分配菜单
// @Tags         基础管理/角色
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        in body AssignMenusRequest true "分配信息"
// @Success      200 {object} response.Body
// @Router       /basic/role/menus/assign [post]
func (s *BasicRoleService) AssignMenus(c *gin.Context) {
	var req AssignMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	db := s.data.GetDBWithContext(c.Request.Context())
	if db.Error != nil {
		response.Fail(c, ecode.Failed, db.Error)
		return
	}

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除现有的菜单关联
	if err := tx.Table("basic_role_menu_rel").Where("role_id = ?", req.RoleID).Delete(nil).Error; err != nil {
		tx.Rollback()
		response.Fail(c, ecode.Failed, err)
		return
	}

	// 添加新的菜单关联
	if len(req.MenuIDs) > 0 {
		var relations []map[string]interface{}
		for _, menuID := range req.MenuIDs {
			relations = append(relations, map[string]interface{}{
				"role_id": req.RoleID,
				"menu_id": menuID,
			})
		}
		if err := tx.Table("basic_role_menu_rel").Create(&relations).Error; err != nil {
			tx.Rollback()
			response.Fail(c, ecode.Failed, err)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	response.Success(c, nil)
}
