package router

import (
	"battle-tiles/internal/service/basic"

	"github.com/gin-gonic/gin"
)

type BasicRouter struct {
	userService  *basic.BasicUserService
	loginService *basic.BasicLoginService
	menuService  *basic.BasicMenuService
	roleService  *basic.BasicRoleService
}

func (r *BasicRouter) InitRouter(root *gin.RouterGroup) {
	r.userService.RegisterRouter(root)
	r.loginService.RegisterRouter(root)
	r.menuService.RegisterRouter(root)
	r.roleService.RegisterRouter(root)
}

func NewBasicRouter(
	userService *basic.BasicUserService,
	loginService *basic.BasicLoginService,
	menuService *basic.BasicMenuService,
	roleService *basic.BasicRoleService,
) *BasicRouter {
	return &BasicRouter{
		userService:  userService,
		loginService: loginService,
		menuService:  menuService,
		roleService:  roleService,
	}
}
