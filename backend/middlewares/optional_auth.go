package middlewares

import (
	"strings"

	"community/backend/config"
	"community/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// OptionalAuth 可选认证：解析 JWT 但不强制要求，解析成功则注入 context
// 用于公开接口（如评论列表）需要根据登录状态显示不同按钮的场景
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWTKey(), nil
		})
		if err != nil || !token.Valid {
			c.Next()
			return
		}

		var user models.User
		if err := config.DB.Preload("Roles").First(&user, claims.UserID).Error; err != nil {
			c.Next()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("username", claims.Username)
		c.Set("roles", user.RoleNames())
		c.Next()
	}
}
