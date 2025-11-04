package router

import (
	"battle-tiles/internal/service"

	"github.com/gin-gonic/gin"
)

type OpsRouter struct {
	opsService *service.OpsService
}

func (r *OpsRouter) InitRouter(root *gin.RouterGroup) {
	r.opsService.RegisterRouter(root)
}

func NewOpsRouter(opsService *service.OpsService) *OpsRouter {
	return &OpsRouter{opsService: opsService}
}
