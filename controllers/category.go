package controllers

import (
	"net/http"

	"community/data"
	"community/models"
	"community/service"

	"github.com/gin-gonic/gin"
)

// GetCategories 获取分类列表
func GetCategories(c *gin.Context) {
	categories, err := data.ListCategories(c.DefaultQuery("q", ""))
	if err != nil {
		models.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}
	models.Success(c, "获取成功", gin.H{"categories": categories})
}

// CreateCategory 创建分类
func CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required,min=1,max=50"`
		Description string `json:"description" binding:"max=200"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	cat, err := service.CreateCategory(req.Name, req.Description)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "分类已创建", gin.H{"category": cat})
}

// DeleteCategory 删除分类
func DeleteCategory(c *gin.Context) {
	id := strToUint(c.Param("id"))

	if err := service.DeleteCategory(id); err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "分类已删除（帖子已迁移至综合讨论）", nil)
}
