package controllers

import (
	"net/http"
	"strings"

	"community/backend/data"
	"community/backend/middlewares"
	"community/backend/models"
	"community/backend/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ─── 发表评论 ───
func CreateComment(c *gin.Context) {
	postID := strToUint(c.Param("id"))

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

// ─── 评论列表（公开，但解析 token 确定权限）───
func GetComments(c *gin.Context) {
	postID := strToUint(c.Param("id"))

	comments, err := data.ListCommentsByPost(postID)
	if err != nil {
		models.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}

	// 尝试从 token 获取当前用户
	currentUserID := uint(0)
	roleNames := []string{}
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims := &middlewares.Claims{}
			token, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (interface{}, error) {
				return middlewares.JWTKey(), nil
			})
			if err == nil && token.Valid {
				currentUserID = claims.UserID
				if u, err := data.FindUserByIDWithRoles(currentUserID); err == nil {
					roleNames = u.RoleNames()
				}
			}
		}
	}

	items := service.BuildCommentPermissions(comments, currentUserID, roleNames)

	models.Success(c, "获取成功", gin.H{"comments": items, "total": len(items)})
}

// ─── 编辑评论 ───
func UpdateComment(c *gin.Context) {
	id := strToUint(c.Param("id"))
	userID := c.GetUint("user_id")
	roles := c.GetStringSlice("roles")

	var req struct {
		Content string `json:"content" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	comment, err := service.UpdateComment(id, userID, req.Content, roles)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "已更新", gin.H{"comment": comment})
}

// ─── 删除评论 ───
func DeleteComment(c *gin.Context) {
	id := strToUint(c.Param("id"))
	userID := c.GetUint("user_id")
	roles := c.GetStringSlice("roles")

	if err := service.DeleteComment(id, userID, roles); err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "已删除", nil)
}
