package service

import (
	cloudBiz "battle-tiles/internal/biz/cloud"
	"battle-tiles/pkg/utils/response"
	"context"

	"github.com/gin-gonic/gin"
)

// PlatformService 提供平台相关公共接口
type PlatformService struct{ uc *cloudBiz.PlatformUsecase }

func NewPlatformService(uc *cloudBiz.PlatformUsecase) *PlatformService {
	return &PlatformService{uc: uc}
}

func (s *PlatformService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/platforms")
	g.GET("/list", s.List)
}

// List 列出全部平台
// @Summary      平台列表
// @Description  返回全部平台
// @Tags         Public
// @Produce      json
// @Success      200 {object} response.Body
// @Router       /platforms/list [get]
func (s *PlatformService) List(c *gin.Context) {
	items, err := s.uc.ListAll(context.Background())
	if err != nil {
		response.Fail(c, 500, err)
		return
	}
	response.Success(c, items)
}
