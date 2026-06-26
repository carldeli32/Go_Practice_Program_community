package middlewares

import (
	"net/http"

	"community/backend/models"

	"github.com/gin-gonic/gin"
)

// RequirePerm 检查当前用户是否拥有指定权限
func RequirePerm(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, _ := c.Get("roles")
		roleNames, ok := roles.([]string)
		if !ok {
			models.Error(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}
		if !models.HasPerm(roleNames, perm) {
			models.Error(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}


