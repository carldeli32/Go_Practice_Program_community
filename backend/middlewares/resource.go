package middlewares

import (
	"net/http"

	"community/backend/models"

	"github.com/gin-gonic/gin"
)

// ResourceLoader 从请求中加载资源，返回 ownerID + categoryID + 资源对象
// 资源对象会被注入 context，供 controller 直接使用，避免重复查询
type ResourceLoader func(c *gin.Context) (ownerID, categoryID uint, resource interface{}, err error)

// RequireResource 声明式资源权限中间件
//
// 用法：
//
//	auth.PUT("/posts/:id", RequireResource(postLoader, "post.manage_any", "post.manage_category"), UpdatePost)
//
// 路由声明即权限文档，controller 和 service 层零权限代码。
func RequireResource(loader ResourceLoader, permAny, permCategory string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		roles := c.GetStringSlice("roles")

		ownerID, categoryID, resource, err := loader(c)
		if err != nil {
			models.Error(c, http.StatusNotFound, "资源不存在")
			c.Abort()
			return
		}

		if !CheckRelation(userID, ownerID, categoryID, roles, permAny, permCategory) {
			models.Error(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}

		// 注入 context，controller 免重查
		c.Set("resource", resource)
		c.Next()
	}
}
