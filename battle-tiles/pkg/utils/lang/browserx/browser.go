package browserx

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"strings"
)

// BrowserFingerprint 获取浏览器指纹
func BrowserFingerprint(c *gin.Context) string {
	// 从请求头中提取可以标识浏览器的信息
	userAgent := c.GetHeader("User-Agent")
	accept := c.GetHeader("Accept")
	acceptEncoding := c.GetHeader("Accept-Encoding")
	acceptLanguage := c.GetHeader("Accept-Language")

	// 可以添加更多浏览器特征，如屏幕分辨率等
	// 注意：在实际应用中可能需要前端配合收集更多信息

	// 生成指纹哈希
	h := sha256.New()
	h.Write([]byte(userAgent))
	h.Write([]byte(accept))
	h.Write([]byte(acceptEncoding))
	h.Write([]byte(acceptLanguage))

	return hex.EncodeToString(h.Sum(nil))[:16] // 取前16位作为指纹
}

// GenerateBrowserKey 生成防抖Key
func GenerateBrowserKey(ctx *gin.Context) string {

	var keyParts []string
	keyParts = append(keyParts, BrowserFingerprint(ctx))

	// 添加请求方法+路径作为区分不同接口的依据
	keyParts = append(keyParts, ctx.Request.Method+"_"+ctx.Request.URL.Path)

	return strings.Join(keyParts, "|")
}
