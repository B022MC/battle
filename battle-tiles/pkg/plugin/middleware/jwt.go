package middleware

import (
	pdb "battle-tiles/pkg/plugin/dbx"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := utils.GetToken(c)
		if token == "" {
			response.Fail(c, ecode.TokenValidateFailed, ecode.Text(ecode.TokenValidateFailed))
			c.Abort()
			return
		}
		j := utils.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				response.Fail(c, ecode.TokenValidateFailed, err)
				utils.ClearToken(c)
				c.Abort()
				return
			}
			response.Fail(c, ecode.TokenValidateFailed, err)
			utils.ClearToken(c)
			c.Abort()
			return
		}
		dbName := c.GetHeader(pdb.CtxDBKey)

		if dbName == "" && claims.Platform == "" {
			c.Abort()
			response.Fail(c, ecode.Failed, "Database name not provided")
			return
		}

		if dbName == "" {
			dbName = claims.Platform
		}
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, pdb.CtxDBKey, dbName)
		c.Request = c.Request.WithContext(ctx)
		c.Set("claims", claims)
		c.Next()

	}
}

func WSJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := utils.GetWSToken(c)
		if token == "" {
			response.Fail(c, ecode.TokenValidateFailed, ecode.Text(ecode.TokenValidateFailed))
			c.Abort()
			return
		}
		j := utils.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				response.Fail(c, ecode.TokenValidateFailed, err)
				utils.ClearToken(c)
				c.Abort()
				return
			}
			response.Fail(c, ecode.TokenValidateFailed, err)
			utils.ClearToken(c)
			c.Abort()
			return
		}
		dbName := c.GetHeader(pdb.CtxDBKey)
		if dbName == "" && claims.Platform == "" {
			c.Abort()
			response.Fail(c, ecode.Failed, "Database name not provided")
			return
		}
		if claims.Platform != "" {
			dbName = claims.Platform
		}
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, pdb.CtxDBKey, dbName)
		c.Request = c.Request.WithContext(ctx)
		c.Set("claims", claims)
		c.Next()

	}
}
