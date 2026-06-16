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
	query := config.DB.Model(&models.User{}).Preload("Roles")
	if q != "" {
		query = query.Where("username LIKE ?", "%"+q+"%")
	}
	query.Order("id ASC").Find(&users)

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

// ========== 管理员创建用户 ==========
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

	if err := config.DB.Create(&user).Error; err != nil {
		models.Error(c, http.StatusInternalServerError, "创建失败")
		return
	}

	// 赋予默认 user 角色
	var userRole models.Role
	config.DB.Where("name = ?", "user").First(&userRole)
	config.DB.Model(&user).Association("Roles").Append(&userRole)

	models.Success(c, "用户已创建", gin.H{"id": user.ID, "username": user.Username})
}

// ========== 删除用户 ==========
func AdminDeleteUser(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	if err := config.DB.First(&user, uid).Error; err != nil {
		models.Error(c, http.StatusNotFound, "用户不存在")
		return
	}
	if user.IsAdminLike() {
		models.Error(c, http.StatusBadRequest, "不能删除管理员")
		return
	}

	// 清理关联数据
	config.DB.Where("user_id = ?", uid).Delete(&models.Post{})
	config.DB.Where("user_id = ?", uid).Delete(&models.Comment{})
	config.DB.Where("from_user_id = ? OR to_user_id = ?", uid, uid).Delete(&models.Message{})
	config.DB.Where("follower_id = ? OR followee_id = ?", uid, uid).Delete(&models.Follow{})
	config.DB.Delete(&user)

	models.Success(c, "用户已删除", nil)
}

// ========== 赋予角色 ==========
func AdminAssignRoles(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Roles []string `json:"roles" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var user models.User
	if err := config.DB.Preload("Roles").First(&user, uid).Error; err != nil {
		models.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	// 清空旧角色，设置新角色
	config.DB.Model(&user).Association("Roles").Clear()

	var roles []models.Role
	for _, name := range req.Roles {
		var role models.Role
		if err := config.DB.Where("name = ?", name).First(&role).Error; err == nil {
			roles = append(roles, role)
		}
	}
	config.DB.Model(&user).Association("Roles").Append(roles)

	models.Success(c, "角色已更新", gin.H{"roles": req.Roles})
}

// ========== 封禁 / 解封 ==========
func BanUser(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	if err := config.DB.Preload("Roles").First(&user, uid).Error; err != nil {
		models.Error(c, http.StatusNotFound, "用户不存在")
		return
	}
	if user.IsAdminLike() {
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
		Content    string `json:"content" binding:"required"`
		CategoryID uint   `json:"category_id"` // 0=全站，非 0=版内
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "内容不能为空")
		return
	}

	// 覆盖同名分类的公告（每分类仅保留一条）
	config.DB.Where("category_id = ?", req.CategoryID).Delete(&models.Announcement{})
	config.DB.Create(&models.Announcement{Content: req.Content, CategoryID: req.CategoryID})
	models.Success(c, "公告已发布 📢", nil)
}

func DeleteAnnouncement(c *gin.Context) {
	categoryID := c.DefaultQuery("category_id", "0")
	config.DB.Where("category_id = ?", categoryID).Delete(&models.Announcement{})
	models.Success(c, "公告已删除", nil)
}

func GetAnnouncement(c *gin.Context) {
	categoryID := c.DefaultQuery("category_id", "0")

	var announcements []models.Announcement

	// 全站公告 + 指定分类公告
	if categoryID == "0" {
		config.DB.Where("category_id = 0").Order("created_at DESC").Find(&announcements)
	} else {
		config.DB.Where("category_id = 0 OR category_id = ?", categoryID).
			Order("created_at DESC").Find(&announcements)
	}

	// 合并内容返回
	content := ""
	for _, a := range announcements {
		if content != "" {
			content += "\n"
		}
		content += a.Content
	}

	models.Success(c, "获取成功", gin.H{"content": content, "announcements": announcements})
}
