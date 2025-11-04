package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 允许跨域访问（开发环境友好，生产可按需收紧）
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}

		// 允许来源（带凭证时需指定具体 origin，这里按请求回显）
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Vary", "Origin")

		// 允许的方法与头
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Platform, X-Requested-With, Accept")
		c.Header("Access-Control-Expose-Headers", "Content-Disposition, Content-Length, Content-Type")

		// 是否允许携带凭证（如需 Cookie 可置为 true）
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
