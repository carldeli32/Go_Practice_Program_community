package controllers

import (
	"net/http"
	"strconv"

	"community/backend/data"
	"community/backend/models"
	"community/backend/service"

	"github.com/gin-gonic/gin"
)

// ─── 用户列表 ───
func AdminListUsers(c *gin.Context) {
	users, err := data.ListUsers(c.Query("q"))
	if err != nil {
		models.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}

	type userItem struct {
		ID       uint     `json:"id"`
		Username string   `json:"username"`
		IsAdmin  bool     `json:"is_admin"`
		IsBanned bool     `json:"is_banned"`
		Motto    string   `json:"motto"`
		Roles    []string `json:"roles"`
	}
	items := make([]userItem, len(users))
	for i, u := range users {
		items[i] = userItem{u.ID, u.Username, u.IsAdminLike(), u.IsBanned, u.Motto, u.RoleNames()}
	}
	models.Success(c, "获取成功", gin.H{"users": items, "total": len(items)})
}

// ─── 创建用户 ───
func AdminCreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=2,max=50"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	user, err := service.AdminCreateUser(req.Username, req.Password)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "用户已创建", gin.H{"id": user.ID, "username": user.Username})
}

// ─── 删除用户 ───
func AdminDeleteUser(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))

	if err := service.AdminDeleteUser(uint(uid)); err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "用户已删除", nil)
}

// ─── 角色分配 ───
func AdminAssignRoles(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Roles []string `json:"roles" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := service.AssignRoles(uint(uid), req.Roles); err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "角色已更新", gin.H{"roles": req.Roles})
}

// ─── 封禁 ───
func BanUser(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))

	username, err := service.BanUser(uint(uid))
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "已封禁 "+username, nil)
}

// ─── 解封 ───
func UnbanUser(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))

	username, err := service.UnbanUser(uint(uid))
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "已解封 "+username, nil)
}

// ─── 公告 ───
func SetAnnouncement(c *gin.Context) {
	var req struct {
		Content    string `json:"content" binding:"required"`
		CategoryID uint   `json:"category_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "内容不能为空")
		return
	}

	if err := service.SetAnnouncement(req.Content, req.CategoryID); err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "公告已发布 📢", nil)
}

func DeleteAnnouncement(c *gin.Context) {
	cid, _ := strconv.Atoi(c.DefaultQuery("category_id", "0"))
	_ = service.DeleteAnnouncement(uint(cid))
	models.Success(c, "公告已删除", nil)
}

func GetAnnouncement(c *gin.Context) {
	categoryID := c.DefaultQuery("category_id", "0")

	announcements, _ := data.ListAnnouncements(categoryID)

	content := ""
	for _, a := range announcements {
		if content != "" {
			content += "\n"
		}
		content += a.Content
	}

	models.Success(c, "获取成功", gin.H{"content": content, "announcements": announcements})
}
