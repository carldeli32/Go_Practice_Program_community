package controllers

import (
	"community/config"
	"community/models"
	"time"

	"github.com/gin-gonic/gin"
)

// canManagePost 检查当前用户是否有权管理帖子（作者 or manage_any or manage_category）
func canManagePost(c *gin.Context, post models.Post) bool {
	userID := c.GetUint("user_id")

	// 作者本人
	if post.UserID == userID {
		return true
	}

	roles := c.GetStringSlice("roles")

	// 全站管理
	if models.HasPerm(roles, "post.manage_any") {
		return true
	}

	// 版内管理：检查是否管辖此分类
	if models.HasPerm(roles, "post.manage_category") {
		var count int64
		config.DB.Model(&models.ModeratorCategory{}).
			Where("user_id = ? AND category_id = ?", userID, post.CategoryID).
			Count(&count)
		return count > 0
	}

	return false
}

// canManageComment 检查当前用户是否有权管理评论
func canManageComment(c *gin.Context, comment models.Comment) bool {
	userID := c.GetUint("user_id")

	if comment.UserID == userID {
		return true
	}

	roles := c.GetStringSlice("roles")
	if models.HasPerm(roles, "comment.manage_any") {
		return true
	}

	if models.HasPerm(roles, "comment.manage_category") {
		var post models.Post
		if err := config.DB.First(&post, comment.PostID).Error; err != nil {
			return false
		}
		var count int64
		config.DB.Model(&models.ModeratorCategory{}).
			Where("user_id = ? AND category_id = ?", userID, post.CategoryID).
			Count(&count)
		return count > 0
	}

	return false
}

// checkMuted 检查用户是否被禁言（全站 or 指定分类）
func checkMuted(c *gin.Context, post models.Post) (bool, string) {
	userID := c.GetUint("user_id")
	now := time.Now()

	// 全站禁言
	var globalMute models.Mute
	if err := config.DB.Where("user_id = ? AND category_id = 0 AND muted_until > ?", userID, now).First(&globalMute).Error; err == nil {
		return true, "你已被全站禁言，解禁时间: " + globalMute.MutedUntil.Format("2006-01-02 15:04")
	}

	// 该分类禁言
	var catMute models.Mute
	if err := config.DB.Where("user_id = ? AND category_id = ? AND muted_until > ?", userID, post.CategoryID, now).First(&catMute).Error; err == nil {
		return true, "你已被该版块禁言，解禁时间: " + catMute.MutedUntil.Format("2006-01-02 15:04")
	}

	return false, ""
}

// getRoles 从 context 中获取角色列表
func getRoles(c *gin.Context) []string {
	roles := c.GetStringSlice("roles")
	if roles == nil {
		return []string{}
	}
	return roles
}
