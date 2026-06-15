package controllers

import (
	"net/http"
	"strconv"

	"community/config"
	"community/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ========== 用户检索 / 列表 ==========
func AdminListUsers(c *gin.Context) {
	q := c.Query("q")
	var users []models.User
	query := config.DB.Model(&models.User{})
	if q != "" {
		query = query.Where("username LIKE ?", "%"+q+"%")
	}
	query.Order("id ASC").Find(&users)

	type userItem struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		IsAdmin  bool   `json:"is_admin"`
		IsBanned bool   `json:"is_banned"`
		Motto    string `json:"motto"`
	}
	items := make([]userItem, len(users))
	for i, u := range users {
		items[i] = userItem{u.ID, u.Username, u.IsAdmin, u.IsBanned, u.Motto}
	}
	models.Success(c, "获取成功", gin.H{"users": items, "total": len(items)})
}

// ========== 管理员创建用户 ==========
// POST /api/admin/users
func AdminCreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=2,max=50"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var exist models.User
	if err := config.DB.Where("username = ?", req.Username).First(&exist).Error; err == nil {
		models.Error(c, http.StatusConflict, "用户名已存在")
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{Username: req.Username, Password: string(hashed)}
	config.DB.Create(&user)

	models.Success(c, "用户已创建", gin.H{"id": user.ID, "username": user.Username})
}

func BanUser(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	if err := config.DB.First(&user, uid).Error; err != nil {
		models.Error(c, http.StatusNotFound, "用户不存在")
		return
	}
	if user.IsAdmin {
		models.Error(c, http.StatusBadRequest, "不能封禁管理员")
		return
	}
	config.DB.Model(&user).Update("is_banned", true)
	models.Success(c, "已封禁 "+user.Username, nil)
}

func UnbanUser(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	if err := config.DB.First(&user, uid).Error; err != nil {
		models.Error(c, http.StatusNotFound, "用户不存在")
		return
	}
	config.DB.Model(&user).Update("is_banned", false)
	models.Success(c, "已解封 "+user.Username, nil)
}

// ========== 公告 ==========
func SetAnnouncement(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "内容不能为空")
		return
	}
	config.DB.Where("1 = 1").Delete(&models.Announcement{})
	config.DB.Create(&models.Announcement{Content: req.Content})
	models.Success(c, "公告已发布 📢", nil)
}

func DeleteAnnouncement(c *gin.Context) {
	config.DB.Where("1 = 1").Delete(&models.Announcement{})
	models.Success(c, "公告已删除", nil)
}

func GetAnnouncement(c *gin.Context) {
	var announcement models.Announcement
	if err := config.DB.First(&announcement).Error; err != nil {
		models.Success(c, "获取成功", gin.H{"content": ""})
		return
	}
	models.Success(c, "获取成功", gin.H{"content": announcement.Content})
}
