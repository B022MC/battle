package basic

import (
	"battle-tiles/internal/infra"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils/response"

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
	r.GET("/list", s.List)
	r.GET("/getOne", s.GetOne)
	r.GET("/all", s.All)
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
