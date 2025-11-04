package server

import (
	"battle-tiles/internal/conf"
	rbacstore "battle-tiles/internal/dal/repo/rbac" // ★ 新增
	"battle-tiles/internal/infra"
	"battle-tiles/internal/router"
	"battle-tiles/pkg/plugin/middleware"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/tx7do/kratos-transport/transport/gin"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, logger log.Logger, router *router.RootRouter, data *infra.Data) *gin.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := gin.NewServer(
		gin.WithAddress(c.Http.Addr),
		gin.WithLogger(logger),
	)

	// ★ 全局中间件
	srv.Engine.Use(middleware.CORS())

	// ★ 先完成 RBAC 绑定（在注册路由/挂鉴权中间件之前）
	initPlugin(srv, logger, data)

	// ★ 再注册你的业务路由（里面会用到 RequirePerm / RequireAnyPerm）
	router.InitRouter(srv.Group(""))

	return srv
}

// 把“插件/全局绑定”都放这里，保持 server 是“组装层”
func initPlugin(srv *gin.Server, logger log.Logger, data *infra.Data) {
	ps := rbacstore.NewStore(data, logger)

	middleware.BindPermissionStore(ps)

}
