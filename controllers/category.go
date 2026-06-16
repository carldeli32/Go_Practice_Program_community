package controllers

import (
	"net/http"

	"community/config"
	"community/models"

	"github.com/gin-gonic/gin"
)

// GetCategories 获取分类列表（支持按名称搜索）
func GetCategories(c *gin.Context) {
	q := c.DefaultQuery("q", "")

	var categories []models.Category
	query := config.DB.Order("id ASC")
	if q != "" {
		query = query.Where("name LIKE ?", "%"+q+"%")
	}
	query.Find(&categories)

	models.Success(c, "获取成功", gin.H{"categories": categories})
}

// CreateCategory 创建分类（super_admin）
func CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required,min=1,max=50"`
		Description string `json:"description" binding:"max=200"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	var exist models.Category
	if err := config.DB.Where("name = ?", req.Name).First(&exist).Error; err == nil {
		models.Error(c, http.StatusConflict, "分类名已存在")
		return
	}

	cat := models.Category{Name: req.Name, Description: req.Description}
	config.DB.Create(&cat)
	models.Success(c, "分类已创建", gin.H{"category": cat})
}

// DeleteCategory 删除分类（super_admin）
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var cat models.Category
	if err := config.DB.First(&cat, id).Error; err != nil {
		models.Error(c, http.StatusNotFound, "分类不存在")
		return
	}

	// 综合讨论不能删
	if cat.ID == 1 {
		models.Error(c, http.StatusBadRequest, "默认分类「综合讨论」不可删除")
		return
	}

	// 把该分类下的帖子迁移到综合讨论
	config.DB.Model(&models.Post{}).Where("category_id = ?", cat.ID).Update("category_id", 1)
	config.DB.Delete(&cat)
	models.Success(c, "分类已删除（帖子已迁移至综合讨论）", nil)
}
