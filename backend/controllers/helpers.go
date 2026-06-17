package controllers

import (
	"github.com/gin-gonic/gin"
)

// getRoles 从 gin context 获取角色列表
func getRoles(c *gin.Context) []string {
	roles := c.GetStringSlice("roles")
	if roles == nil {
		return []string{}
	}
	return roles
}
