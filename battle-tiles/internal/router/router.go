package router

import (
	"battle-tiles/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewRootRouter,
	NewBasicRouter,
	NewGameRouter,
	NewOpsRouter,
)

type RootRouter struct {
	basicRouter *BasicRouter
	gameRouter  *GameRouter
	opsRouter   *OpsRouter
	platformSvc *service.PlatformService
}

func (g *RootRouter) InitRouter(group *gin.RouterGroup) {

	g.basicRouter.InitRouter(group)
	g.gameRouter.InitRouter(group)
	g.opsRouter.InitRouter(group)
	g.platformSvc.RegisterRouter(group)
}

func NewRootRouter(
	basicRouter *BasicRouter,
	gameRouter *GameRouter,
	opsRouter *OpsRouter,
	platformSvc *service.PlatformService,
) *RootRouter {
	return &RootRouter{
		basicRouter: basicRouter,
		gameRouter:  gameRouter,
		opsRouter:   opsRouter,
		platformSvc: platformSvc,
	}
}
