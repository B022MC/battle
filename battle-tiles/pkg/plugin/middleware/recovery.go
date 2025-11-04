package middleware

import (
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
)

func GinRecovery(logger log.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					_ = logger.Log(log.LevelError,
						"path", c.Request.URL.Path,
						"error", err,
						"request", string(httpRequest),
					)
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error))
					c.Abort()
					response.Fail(c, ecode.Failed, err)
					return
				}

				if stack {
					_ = logger.Log(log.LevelError,
						"[Recovery from panic]",
						"error", err,
						"request", string(httpRequest),
						"stack", string(debug.Stack()),
					)
				} else {
					_ = logger.Log(log.LevelError,
						"[Recovery from panic]",
						"error", err,
						"request", string(httpRequest),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
				response.Fail(c, ecode.Failed, err)
			}
		}()
		c.Next()
	}
}
