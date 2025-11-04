package utils

import (
	"battle-tiles/pkg/utils/request"
	"net"
	"net/http"
	"strings"

	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	// CookieLoginToken 登录验证 Cookie中传递的参数
	CookieLoginToken = "x-token"

	// HeaderSignToken 签名验证 Authorization，Header 中传递的参数
	HeaderSignToken = "Authorization"
	// WSHeaderSignToken 签名验证 WebSocket Sec-WebSocket-Protocol，Header 中传递的参数
	WSHeaderSignToken = "Sec-WebSocket-Protocol"
	HeaderTokenPrefix = "Bearer"
)

func ClearToken(c *gin.Context) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie(CookieLoginToken, "", -1, "/", "", false, false)
	} else {
		c.SetCookie(CookieLoginToken, "", -1, "/", host, false, false)
	}
}

func GetToken(c *gin.Context) string {
	token := c.Request.Header.Get(HeaderSignToken)
	if token == "" {
		token, _ = c.Cookie(CookieLoginToken)
	}
	return strings.ReplaceAll(token, fmt.Sprintf("%s ", HeaderTokenPrefix), "")
}

func GetHttpToken(req *http.Request) string {
	token := req.URL.Query().Get(WSHeaderSignToken)
	if req.Header.Get(WSHeaderSignToken) != "" {
		token = req.Header.Get(WSHeaderSignToken)
	}
	return strings.ReplaceAll(token, "Bearer ", "")

}

func GetWSToken(c *gin.Context) string {
	token := c.Request.Header.Get(WSHeaderSignToken)
	if token == "" {
		token, _ = c.Cookie(CookieLoginToken)
	}
	return strings.ReplaceAll(token, fmt.Sprintf("%s ", HeaderTokenPrefix), "")
}

func GetClaims(c *gin.Context) (*request.CustomClaims, error) {
	token := GetToken(c)
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {

	}
	return claims, err
}

// GetOrgID 从Gin的Context中获取从jwt解析出来的用户ID
//func GetOrgID(c *gin.Context) int32 {
//	if claims, exists := c.Get("claims"); !exists {
//		if cl, err := GetClaims(c); err != nil {
//			return 0
//		} else {
//			return cl.BaseClaims.OrgID
//		}
//	} else {
//		waitUse := claims.(*request.CustomClaims)
//		return waitUse.BaseClaims.OrgID
//	}
//}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) int32 {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.BaseClaims.UserID
		}
	} else {
		waitUse := claims.(*request.CustomClaims)
		return waitUse.BaseClaims.UserID
	}
}

// GetUserAuthorityId 从Gin的Context中获取从jwt解析出来的用户角色id
//func GetUserAuthorityId(c *gin.Context) int32 {
//	if claims, exists := c.Get("claims"); !exists {
//		if cl, err := GetClaims(c); err != nil {
//			return 0
//		} else {
//			return cl.RoleId
//		}
//	} else {
//		waitUse := claims.(*request.CustomClaims)
//		return waitUse.RoleId
//	}
//}

// GetUserInfo 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserInfo(c *gin.Context) *request.CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*request.CustomClaims)
		return waitUse
	}
}

// GetUserName 从Gin的Context中获取从jwt解析出来的用户名
func GetUserName(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.Username
		}
	} else {
		waitUse := claims.(*request.CustomClaims)
		return waitUse.Username
	}
}
