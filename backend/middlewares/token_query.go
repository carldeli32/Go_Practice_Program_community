package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// TokenFromQuery 从 URL query 参数中提取 token，注入 Authorization header
// 用于 SSE / EventSource 等不支持自定义 header 的场景
func TokenFromQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token := c.Query("token"); token != "" {
			c.Request.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))
		}
		c.Next()
	}
}
