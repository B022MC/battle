package basic

import (
	basicBiz "battle-tiles/internal/biz/basic"
	basicModel "battle-tiles/internal/dal/model/basic"
	"battle-tiles/internal/dal/req"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

// BasicMenuService —— 菜单服务
type BasicMenuService struct {
	log *log.Helper
	uc  *basicBiz.BasicMenuUseCase
}

func NewBasicMenuService(uc *basicBiz.BasicMenuUseCase, logger log.Logger) *BasicMenuService {
	return &BasicMenuService{uc: uc, log: log.NewHelper(log.With(logger, "module", "service/baseMenu"))}
}

// RegisterRouter —— 建议传入挂了 JWT 的私有分组
func (s *BasicMenuService) RegisterRouter(rootRouter *gin.RouterGroup) {
	privateRouter := rootRouter.Group("/basic").Use(middleware.JWTAuth())

	privateRouter.POST("/baseMenu/addOne", middleware.RequirePerm("menu:create"), s.AddOne)
	privateRouter.POST("/baseMenu/updateOne", middleware.RequirePerm("menu:update"), s.UpdateOne)

	privateRouter.GET("/baseMenu/getOne", middleware.RequirePerm("menu:view"), s.GetOne)
	privateRouter.GET("/baseMenu/getPage", middleware.RequirePerm("menu:view"), s.GetPage)
	privateRouter.GET("/baseMenu/getOption", middleware.RequirePerm("menu:view"), s.GetOption)
	privateRouter.GET("/baseMenu/getTree", middleware.RequirePerm("menu:view"), s.GetTree)

	// 用户可见菜单树（按当前登录用户过滤）
	privateRouter.GET("/menu/me/tree", s.MeTree)
	// 用户在某菜单下的按钮权限（返回按钮ID列表；如需 perm_code 列表，可在按钮表上扩展）
	privateRouter.GET("/menu/me/buttons", s.MeButtons)

	privateRouter.POST("/baseMenu/saveTree", middleware.RequirePerm("menu:update"), s.SaveTree)
	privateRouter.GET("/baseMenu/delOne", middleware.RequirePerm("menu:delete"), s.DeleteOne)
	privateRouter.POST("/baseMenu/delMany", middleware.RequirePerm("menu:delete"), s.DeleteMany)
}

func activeScope(db *gorm.DB) *gorm.DB { // 过滤软删
	return db.Where("deleted_at IS NULL")
}

func sortScope(db *gorm.DB) *gorm.DB { // rank（可能为 NULL）升序 + id 升序
	return db.Order("COALESCE(rank,'') ASC").Order("id ASC")
}

// AddOne
// @Summary      BaseMenu 添加单条
// @Description  Add by model
// @Tags         系统管理/菜单管理
// @Accept       json
// @Produce      json
// @Param        inParam  body      req.AddBasicMenuReq  true  "请求参数"
// @Success      200      {object}  response.Body{data=basic.BasicMenuDoc,msg=string}
// @Router       /basic/baseMenu/addOne [post]
func (s *BasicMenuService) AddOne(c *gin.Context) {
	var in req.AddBasicMenuReq
	defaults.SetDefaults(&in)
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	var m basicModel.BasicMenu
	if err := copier.Copy(&m, &in); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	// 默认父级
	if m.ParentId == 0 {
		m.ParentId = -1
	}
	if err := s.uc.Create(c.Request.Context(), &m); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, m)
}

// UpdateOne
// @Summary      BaseMenu 更新单条
// @Description  Update by ID + fields
// @Tags         系统管理/菜单管理
// @Accept       json
// @Produce      json
// @Param        inParam  body      req.UpdateBasicMenuReq  true  "请求参数"
// @Success      200      {object}  response.Body{msg=string}
// @Router       /basic/baseMenu/updateOne [post]
func (s *BasicMenuService) UpdateOne(c *gin.Context) {
	var in req.UpdateBasicMenuReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	if in.Id == 0 {
		response.Fail(c, ecode.ParamsFailed, "id is required")
		return
	}
	fields := make(map[string]interface{})
	_ = mapstructure.Decode(in, &fields)
	if err := s.uc.Update(c.Request.Context(), in.Id, fields); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}

// GetOne
// @Summary      获取单条
// @Description  Get by id
// @Tags         系统管理/菜单管理
// @Accept       json
// @Produce      json
// @Param        inParam  query     req.ReqById  true  "请求参数"
// @Success      200      {object}  response.Body{data=basic.BasicMenuDoc,msg=string}
// @Router       /basic/baseMenu/getOne [get]
func (s *BasicMenuService) GetOne(c *gin.Context) {
	var in req.ReqById
	if err := c.ShouldBindQuery(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	data, err := s.uc.SelectOne(
		c.Request.Context(),
		func(db *gorm.DB) *gorm.DB { return activeScope(db).Where("id = ?", in.ID) },
	)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, data)
}

// GetPage
// @Summary      分页列表
// @Description  Get list by page
// @Tags         系统管理/菜单管理
// @Accept       json
// @Produce      json
// @Param        inParam  query     req.PageInfo  true  "请求参数"
// @Success      200      {object}  response.Body{data=utils.PageResult{list=[]basic.BasicMenuDoc},msg=string}
// @Router       /basic/baseMenu/getPage [get]
func (s *BasicMenuService) GetPage(c *gin.Context) {
	var in req.PageInfo
	if err := c.ShouldBindQuery(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	list, total, err := s.uc.SelectPage(
		c.Request.Context(),
		in.Page, in.PageSize, true,
		activeScope, sortScope,
	)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, utils.PageResult{
		List:     list,
		Total:    total,
		PageNo:   int32(in.Page),
		PageSize: int32(in.PageSize),
		NotPage:  false,
	})
}

// GetOption
// @Summary      下拉字典
// @Description  Option list (label/value)
// @Tags         系统管理/菜单管理
// @Accept       json
// @Produce      json
// @Param        inParam  query     req.PageInfo  true  "请求参数"
// @Success      200      {object}  response.Body{data=utils.PageResult{list=[]common.Option},msg=string}
// @Router       /basic/baseMenu/getOption [get]
func (s *BasicMenuService) GetOption(c *gin.Context) {
	var in req.PageInfo
	if err := c.ShouldBindQuery(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	list, total, err := s.uc.SelectOption(
		c.Request.Context(),
		in.Page, in.PageSize, true,
		activeScope, sortScope,
	)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, utils.PageResult{
		List:     list,
		Total:    total,
		PageNo:   int32(in.Page),
		PageSize: int32(in.PageSize),
		NotPage:  false,
	})
}

// GetTree
// @Summary      菜单树
// @Description  支持 keyword 模糊筛选后再组树
// @Tags         系统管理/菜单管理
// @Accept       json
// @Produce      json
// @Param        inParam  query     req.ReqInfo  true  "请求参数"
// @Success      200      {object}  response.Body{data=resp.MenuTree,msg=string}
// @Router       /basic/baseMenu/getTree [get]
func (s *BasicMenuService) GetTree(c *gin.Context) {
	var in req.ReqInfo
	if err := c.ShouldBindQuery(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	tree, err := s.uc.SelectAllToTree(c.Request.Context(), in.Keyword, activeScope, sortScope)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, tree)
}

// MeTree 返回当前用户可见的菜单树（基于其角色聚合可见菜单，并按树结构返回）
// @Summary      我的菜单树
// @Tags         系统管理/菜单管理
// @Security     BearerAuth
// @Produce      json
// @Success      200      {object}  response.Body
// @Router       /basic/menu/me/tree [get]
func (s *BasicMenuService) MeTree(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	tree, err := s.uc.SelectMeTreeFiltered(c.Request.Context(), claims.BaseClaims.UserID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, tree)
}

// MeButtons 返回当前用户在某菜单下的按钮集合（按钮ID列表）。如需返回按钮 perm_code，请扩展按钮维表再 join。
// @Summary      我的菜单按钮
// @Tags         系统管理/菜单管理
// @Security     BearerAuth
// @Produce      json
// @Param        menu_id query int true "菜单ID"
// @Success      200      {object}  response.Body
// @Router       /basic/menu/me/buttons [get]
func (s *BasicMenuService) MeButtons(c *gin.Context) {
	type q struct {
		MenuID int32 `form:"menu_id" binding:"required"`
	}
	var in q
	if err := c.ShouldBindQuery(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	btnIDs, err := s.uc.ListMyButtonIDs(c.Request.Context(), claims.BaseClaims.UserID, in.MenuID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, gin.H{"btn_ids": btnIDs})
}

// SaveTree
// @Summary      保存整棵菜单树（覆盖式）
// @Description  成功后返回最新菜单树
// @Tags         系统管理/菜单管理
// @Accept       json
// @Produce      json
// @Param        inParam  body      req.SaveMenuTree  true  "请求参数"
// @Success      200      {object}  response.Body{data=resp.MenuTree,msg=string}
// @Router       /basic/baseMenu/saveTree [post]
func (s *BasicMenuService) SaveTree(c *gin.Context) {
	var in req.SaveMenuTree
	defaults.SetDefaults(&in)
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	tree, err := s.uc.SaveAllTree(c.Request.Context(), in.MenuTree, activeScope, sortScope)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, tree)
}

// DeleteOne
// @Summary      删除单条
// @Description  Delete by id
// @Tags         系统管理/菜单管理
// @Accept       json
// @Produce      json
// @Param        inParam  query     req.ReqById  true  "请求参数"
// @Success      200      {object}  response.Body{msg=string}
// @Router       /basic/baseMenu/delOne [get]
func (s *BasicMenuService) DeleteOne(c *gin.Context) {
	var in req.ReqById
	if err := c.ShouldBindQuery(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	if in.ID == 0 {
		response.Fail(c, ecode.ParamsFailed, "id is required")
		return
	}
	if err := s.uc.DeleteByIDs(c.Request.Context(), in.ID); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}

// DeleteMany
// @Summary      批量删除
// @Description  Delete by ids
// @Tags         系统管理/菜单管理
// @Accept       json
// @Produce      json
// @Param        inParam  body      req.ReqByIds  true  "请求参数"
// @Success      200      {object}  response.Body{msg=string}
// @Router       /basic/baseMenu/delMany [post]
func (s *BasicMenuService) DeleteMany(c *gin.Context) {
	var in req.ReqByIds
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	if len(in.IDs) == 0 {
		response.Fail(c, ecode.ParamsFailed, "ids is empty")
		return
	}
	if err := s.uc.DeleteByIDs(c.Request.Context(), in.IDs...); err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(c)
}
