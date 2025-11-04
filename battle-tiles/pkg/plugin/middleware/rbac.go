package middleware

import (
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// AND：全部命中
func RequirePerm(perms ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := utils.GetClaims(c)
		if err != nil {
			response.Fail(c, ecode.TokenValidateFailed, err)
			c.Abort()
			return
		}
		// 超级管理员直接放行（角色ID=1）
		if claims.BaseClaims.IsSuperAdmin() {
			c.Next()
			return
		}
		if store() == nil {
			response.Fail(c, ecode.Failed, "rbac store not initialized")
			c.Abort()
			return
		}
		set, err := store().GetUserPermCodes(c.Request.Context(), claims.BaseClaims.UserID)
		if err != nil {
			response.Fail(c, ecode.Failed, err)
			c.Abort()
			return
		}
		for _, p := range perms {
			p = strings.ToLower(strings.TrimSpace(p))
			if _, ok := set[p]; !ok {
				response.Fail(c, ecode.Failed, "permission denied: "+p)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// OR：任意命中
func RequireAnyPerm(perms ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := utils.GetClaims(c)
		if err != nil {
			response.Fail(c, ecode.TokenValidateFailed, err)
			c.Abort()
			return
		}
		// 超级管理员直接放行
		if claims.BaseClaims.IsSuperAdmin() {
			c.Next()
			return
		}
		if store() == nil {
			response.Fail(c, ecode.Failed, "rbac store not initialized")
			c.Abort()
			return
		}
		set, err := store().GetUserPermCodes(c.Request.Context(), claims.BaseClaims.UserID)
		if err != nil {
			response.Fail(c, ecode.Failed, err)
			c.Abort()
			return
		}
		for _, p := range perms {
			p = strings.ToLower(strings.TrimSpace(p))
			if _, ok := set[p]; ok {
				c.Next()
				return
			}
		}
		response.Fail(c, ecode.Failed, "permission denied")
		c.Abort()
	}
}
