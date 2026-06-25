package controllers

import (
	"net/http"

	"community/backend/data"
	"community/backend/middlewares"
	"community/backend/models"
	"community/backend/service"

	"github.com/gin-gonic/gin"
)

// ─── 发表评论 ───
func CreateComment(c *gin.Context) {
	postID := middlewares.StrToUint(c.Param("id"))

	var req struct {
		Content string `json:"content" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID := c.GetUint("user_id")
	comment, err := service.CreateComment(postID, userID, req.Content)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "评论成功 💬", gin.H{"comment": comment})
}

// ─── 评论列表（公开，可选认证以计算 can_edit / can_delete）───
func GetComments(c *gin.Context) {
	postID := middlewares.StrToUint(c.Param("id"))

	comments, err := data.ListCommentsByPost(postID)
	if err != nil {
		models.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}

	// 尝试获取当前用户信息（可选）
	currentUserID := uint(0)
	roleNames := []string{}
	if userID, exists := c.Get("user_id"); exists {
		currentUserID = userID.(uint)
	}
	if roles, exists := c.Get("roles"); exists {
		roleNames = roles.([]string)
	}

	items := service.BuildCommentPermissions(comments, currentUserID, roleNames)

	models.Success(c, "获取成功", gin.H{"comments": items, "total": len(items)})
}

// ─── 编辑评论 ───
// 权限由 RequireResource 中间件保证，controller 零权限代码
func UpdateComment(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	comment := c.MustGet("resource").(*models.Comment) // 中间件注入

	result, err := service.UpdateComment(comment, req.Content)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "已更新", gin.H{"comment": result})
}

// ─── 删除评论 ───
func DeleteComment(c *gin.Context) {
	comment := c.MustGet("resource").(*models.Comment) // 中间件注入

	if err := service.DeleteComment(comment); err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "已删除", nil)
}
