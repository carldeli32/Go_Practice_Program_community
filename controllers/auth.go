package controllers

import (
	"net/http"
	"strconv"

	"community/config"
	"community/middlewares"
	"community/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ========== 注册 ==========
// POST /api/register
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=2,max=50"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email"`
		Gender   string `json:"gender"`
		Age      int    `json:"age"`
		Job      string `json:"job"`
		Motto    string `json:"motto"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 查重
	var exist models.User
	if err := config.DB.Where("username = ?", req.Username).First(&exist).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已被注册"})
		return
	}

	// bcrypt 加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Gender:   req.Gender,
		Age:      req.Age,
		Job:      req.Job,
		Motto:    req.Motto,
	}
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功 🎉",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"gender":   user.Gender,
			"age":      user.Age,
			"job":      user.Job,
			"motto":    user.Motto,
		},
	})
}

// ========== 登录 ==========
// POST /api/login
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if user.IsBanned {
		c.JSON(http.StatusForbidden, gin.H{"error": "账号已被封禁，请联系管理员"})
		return
	}

	token, err := middlewares.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token 生成失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功 👋",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"gender":   user.Gender,
			"age":      user.Age,
			"job":      user.Job,
			"motto":    user.Motto,
			"is_admin": user.IsAdmin,
		},
	})
}

// ========== 用户主页 ==========
// GET /api/users/:id
func GetUserProfile(c *gin.Context) {
	id := c.Param("id")

	uid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户 ID"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 统计发帖数和评论数
	var postCount, commentCount int64
	config.DB.Model(&models.Post{}).Where("user_id = ?", user.ID).Count(&postCount)
	config.DB.Model(&models.Comment{}).Where("user_id = ?", user.ID).Count(&commentCount)

	// 统计关注数和粉丝数
	var followerCount, followingCount, isFollowing int64
	config.DB.Model(&models.Follow{}).Where("followee_id = ?", user.ID).Count(&followerCount)
	config.DB.Model(&models.Follow{}).Where("follower_id = ?", user.ID).Count(&followingCount)
	// 检查当前登录用户是否关注了此人
	if uid, exists := c.Get("user_id"); exists {
		config.DB.Model(&models.Follow{}).
			Where("follower_id = ? AND followee_id = ?", uid, user.ID).
			Count(&isFollowing)
	}

	// 查该用户的帖子列表
	var posts []models.Post
	config.DB.Where("user_id = ?", user.ID).Order("created_at DESC").Limit(10).Find(&posts)

	level := user.Level(config.DB)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":              user.ID,
			"username":        user.Username,
			"email":           user.Email,
			"gender":          user.Gender,
			"age":             user.Age,
			"job":             user.Job,
			"motto":           user.Motto,
			"created_at":      user.CreatedAt,
			"post_count":      postCount,
			"comment_count":   commentCount,
			"follower_count":  followerCount,
			"following_count": followingCount,
			"is_following":    isFollowing > 0,
			"level":           level,
		},
		"posts": posts,
	})
}
