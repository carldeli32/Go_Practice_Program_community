package controllers

import (
	"net/http"
	"strconv"

	"community/config"
	"community/models"

	"github.com/gin-gonic/gin"
)

// ========== 创建帖子 ==========
func CreatePost(c *gin.Context) {
	var req struct {
		Title      string `json:"title" binding:"required,min=1,max=200"`
		Content    string `json:"content" binding:"required,min=1"`
		CategoryID *uint  `json:"category_id"` // nil 时默认综合讨论(1)
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	categoryID := uint(1)
	if req.CategoryID != nil {
		categoryID = *req.CategoryID
	}

	userID := c.GetUint("user_id")
	post := models.Post{Title: req.Title, Content: req.Content, UserID: userID, CategoryID: categoryID}

	// 禁言检查
	if muted, msg := checkMuted(c, post); muted {
		models.Error(c, http.StatusForbidden, msg)
		return
	}

	if err := config.DB.Create(&post).Error; err != nil {
		models.Error(c, http.StatusInternalServerError, "发帖失败")
		return
	}

	models.Success(c, "发帖成功 ✍️", gin.H{"post": post})
}

// ========== 帖子列表 ==========
func GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	categoryIDStr := c.DefaultQuery("category_id", "")
	q := c.DefaultQuery("q", "")
	offset := (page - 1) * pageSize

	var posts []models.Post
	var total int64

	baseQuery := config.DB.Model(&models.Post{})
	if categoryIDStr != "" {
		baseQuery = baseQuery.Where("category_id = ?", categoryIDStr)
	}
	if q != "" {
		baseQuery = baseQuery.Where("title LIKE ?", "%"+q+"%")
	}
	baseQuery.Count(&total)
	baseQuery.Preload("User").Preload("Category").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&posts)

	models.Success(c, "获取成功", gin.H{"posts": posts, "total": total, "page": page, "page_size": pageSize})
}

// ========== 帖子详情 ==========
func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := config.DB.Preload("User").Preload("Category").First(&post, id).Error; err != nil {
		models.Error(c, http.StatusNotFound, "帖子不存在")
		return
	}
	models.Success(c, "获取成功", gin.H{"post": post})
}

// ========== 更新帖子 ==========
func UpdatePost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		models.Error(c, http.StatusNotFound, "帖子不存在")
		return
	}
	if !canManagePost(c, post) {
		models.Error(c, http.StatusForbidden, "无权编辑此帖子")
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if len(updates) == 0 {
		models.Error(c, http.StatusBadRequest, "没有要更新的内容")
		return
	}

	config.DB.Model(&post).Updates(updates)
	models.Success(c, "更新成功", gin.H{"post": post})
}

// ========== 删除帖子 ==========
func DeletePost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		models.Error(c, http.StatusNotFound, "帖子不存在")
		return
	}
	if !canManagePost(c, post) {
		models.Error(c, http.StatusForbidden, "无权删除此帖子")
		return
	}

	config.DB.Where("post_id = ?", post.ID).Delete(&models.Comment{})
	config.DB.Delete(&post)
	models.Success(c, "删除成功", nil)
}
