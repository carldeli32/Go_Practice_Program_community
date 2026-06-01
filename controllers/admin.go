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
	c.JSON(http.StatusOK, gin.H{"users": items, "total": len(items)})
}

// ========== 管理员创建用户 ==========
// POST /api/admin/users
func AdminCreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=2,max=50"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var exist models.User
	if err := config.DB.Where("username = ?", req.Username).First(&exist).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{Username: req.Username, Password: string(hashed)}
	config.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"message": "用户已创建", "user": gin.H{"id": user.ID, "username": user.Username}})
}
func BanUser(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	if err := config.DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	if user.IsAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能封禁管理员"})
		return
	}
	config.DB.Model(&user).Update("is_banned", true)
	c.JSON(http.StatusOK, gin.H{"message": "已封禁 " + user.Username})
}

func UnbanUser(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	if err := config.DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	config.DB.Model(&user).Update("is_banned", false)
	c.JSON(http.StatusOK, gin.H{"message": "已解封 " + user.Username})
}

// ========== 公告 ==========
func SetAnnouncement(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "内容不能为空"})
		return
	}
	config.DB.Where("1 = 1").Delete(&models.Announcement{})
	config.DB.Create(&models.Announcement{Content: req.Content})
	c.JSON(http.StatusOK, gin.H{"message": "公告已发布 📢"})
}

func DeleteAnnouncement(c *gin.Context) {
	config.DB.Where("1 = 1").Delete(&models.Announcement{})
	c.JSON(http.StatusOK, gin.H{"message": "公告已删除"})
}

func GetAnnouncement(c *gin.Context) {
	var announcement models.Announcement
	if err := config.DB.First(&announcement).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"content": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": announcement.Content})
}
