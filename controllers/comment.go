package controllers

import (
	"net/http"

	"community/config"
	"community/models"

	"github.com/gin-gonic/gin"
)

//	发表评论
//
// POST /api/posts/:id/comments  (需登录)
func CreateComment(c *gin.Context) {
	postID := c.Param("id")

	// 检查帖子是否存在
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID := c.GetUint("user_id")

	comment := models.Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  post.ID,
	}
	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
		return
	}

	// 重新查询带用户信息
	config.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "评论成功 💬",
		"comment": comment,
	})
}

//	获取评论列表
//
// GET /api/posts/:id/comments
func GetComments(c *gin.Context) {
	postID := c.Param("id")

	var comments []models.Comment
	config.DB.Where("post_id = ?", postID).
		Preload("User").
		Order("created_at ASC").
		Find(&comments)

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"total":    len(comments),
	})
}
