package middlewares

import (
	"net/http"

	"community/config"
	"community/models"

	"github.com/gin-gonic/gin"
)

// AdminRequired 管理员鉴权中间件（需要在 AuthRequired 之后使用）
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil || !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}
