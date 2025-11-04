package middleware

import (
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 创建一个限流中间件
func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.Abort()
			response.Fail(c, ecode.Failed, "request out of limit")
			return
		}
		c.Next()
	}
}
