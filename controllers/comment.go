package controllers

import (
	"net/http"
	"strings"

	"community/config"
	"community/middlewares"
	"community/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ========== 发表评论 ==========
func CreateComment(c *gin.Context) {
	postID := c.Param("id")
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
	comment := models.Comment{Content: req.Content, UserID: userID, PostID: post.ID}
	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
		return
	}
	config.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusOK, gin.H{"message": "评论成功 💬", "comment": comment})
}

// ========== 获取评论列表（公开，但解析 token 以确定权限）==========
func GetComments(c *gin.Context) {
	postID := c.Param("id")
	var comments []models.Comment
	config.DB.Where("post_id = ?", postID).Preload("User").Order("created_at ASC").Find(&comments)

	// 尝试从 token 获取当前用户身份
	currentUserID := uint(0)
	isAdmin := false
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims := &middlewares.Claims{}
			token, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (interface{}, error) {
				return middlewares.JWTSecret, nil
			})
			if err == nil && token.Valid {
				currentUserID = claims.UserID
				var u models.User
				if err := config.DB.First(&u, currentUserID).Error; err == nil {
					isAdmin = u.IsAdmin
				}
			}
		}
	}

	type commentItem struct {
		models.Comment
		CanEdit   bool `json:"can_edit"`
		CanDelete bool `json:"can_delete"`
	}

	var items []commentItem
	for _, c := range comments {
		items = append(items, commentItem{
			Comment:   c,
			CanEdit:   c.UserID == currentUserID || isAdmin,
			CanDelete: c.UserID == currentUserID || isAdmin,
		})
	}

	c.JSON(http.StatusOK, gin.H{"comments": items, "total": len(items)})
}

// ========== 编辑评论（作者本人或管理员）==========
func UpdateComment(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")

	var comment models.Comment
	if err := config.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 权限检查
	var user models.User
	config.DB.First(&user, userID)
	if comment.UserID != userID && !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	config.DB.Model(&comment).Update("content", req.Content)
	config.DB.Preload("User").First(&comment, id)
	c.JSON(http.StatusOK, gin.H{"message": "已更新", "comment": comment})
}

// ========== 删除评论（作者本人或管理员）==========
func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")

	var comment models.Comment
	if err := config.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	var user models.User
	config.DB.First(&user, userID)
	if comment.UserID != userID && !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作"})
		return
	}

	config.DB.Delete(&comment)
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}
