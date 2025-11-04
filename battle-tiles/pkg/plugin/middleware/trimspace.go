package middleware

import (
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
)

// TrimSpacesMiddleware 去除请求参数中的前后空格
func TrimSpacesMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 遍历请求中的所有参数，并去除其前后空格
		for key, values := range c.Request.URL.Query() {
			for i, value := range values {
				values[i] = strings.TrimSpace(value)
			}
			c.Request.URL.Query().Set(key, strings.Join(values, ","))
		}

		// 遍历请求中的表单参数，并去除其前后空格
		c.Request.ParseForm()
		for key, values := range c.Request.PostForm {
			for i, value := range values {
				values[i] = strings.TrimSpace(value)
			}
			c.Request.PostForm.Set(key, strings.Join(values, ","))
		}
		if c.Request.Method == "POST" && c.Request.Header.Get("Content-Type") == "application/json" {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.Abort()
				return
			}

			var jsonMap map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &jsonMap); err != nil {
				c.Abort()
				response.Fail(c, ecode.ParamsFailed, err)
				return
			}

			trimSpaces(jsonMap)

			trimmedBodyBytes, err := json.Marshal(jsonMap)
			if err != nil {
				c.Abort()
				response.Fail(c, ecode.ParamsFailed, err)
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewBuffer(trimmedBodyBytes))
		}
		c.Next()
	}
}

func trimSpaces(data map[string]interface{}) {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			data[key] = strings.TrimSpace(v)
		case map[string]interface{}:
			trimSpaces(v)
		case []interface{}:
			for i, item := range v {
				if str, ok := item.(string); ok {
					v[i] = strings.TrimSpace(str)
				} else if nestedMap, ok := item.(map[string]interface{}); ok {
					trimSpaces(nestedMap)
				}
			}
		}
	}
}
