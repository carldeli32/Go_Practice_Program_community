package controllers

import (
	"net/http"
	"strconv"

	"community/config"
	"community/models"

	"github.com/gin-gonic/gin"
)

// ========== 创建帖子 ==========
// POST /api/posts  (需登录)
func CreatePost(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required,min=1,max=200"`
		Content string `json:"content" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}
	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发帖失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "发帖成功 ✍️",
		"post":    post,
	})
}

// ========== 帖子列表（分页）==========
// GET /api/posts?page=1&page_size=10
func GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var posts []models.Post
	var total int64

	config.DB.Model(&models.Post{}).Count(&total)
	config.DB.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"posts":     posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ========== 帖子详情 ==========
// GET /api/posts/:id
func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	if err := config.DB.Preload("User").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// ========== 更新帖子 ==========
// PUT /api/posts/:id  (需登录，只能改自己的)
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")

	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "只能编辑自己的帖子"})
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有要更新的内容"})
		return
	}

	config.DB.Model(&post).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "post": post})
}

// ========== 删除帖子 ==========
// DELETE /api/posts/:id  (需登录，只能删自己的)
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")

	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "只能删除自己的帖子"})
		return
	}

	// 删除帖子的同时删除关联评论
	config.DB.Where("post_id = ?", post.ID).Delete(&models.Comment{})
	config.DB.Delete(&post)

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
