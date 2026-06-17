package controllers

import (
	"net/http"
	"strconv"

	"community/data"
	"community/models"
	"community/service"

	"github.com/gin-gonic/gin"
)

// ─── 注册 ───
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
		models.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	user, err := service.Register(req.Username, req.Password, req.Email, req.Gender, req.Motto, req.Job, req.Age)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "注册成功 🎉", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"gender":   user.Gender,
		"age":      user.Age,
		"job":      user.Job,
		"motto":    user.Motto,
	})
}

// ─── 登录 ───
// POST /api/login
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	token, user, err := service.Login(req.Username, req.Password)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "登录成功 👋", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"gender":   user.Gender,
			"age":      user.Age,
			"job":      user.Job,
			"motto":    user.Motto,
			"is_admin": user.IsAdminLike(),
			"roles":    user.RoleNames(),
		},
	})
}

// ─── 用户主页 ───
// GET /api/users/:id
func GetUserProfile(c *gin.Context) {
	id := c.Param("id")
	uid, err := strconv.Atoi(id)
	if err != nil {
		models.Error(c, http.StatusBadRequest, "无效的用户 ID")
		return
	}

	// 获取当前登录用户 ID（可选）
	var currentUID uint
	if v, exists := c.Get("user_id"); exists {
		currentUID = v.(uint)
	}

	result, err := service.GetUserProfile(uint(uid), currentUID)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "获取成功", gin.H{
		"user": gin.H{
			"id":              result.User.ID,
			"username":        result.User.Username,
			"email":           result.User.Email,
			"gender":          result.User.Gender,
			"age":             result.User.Age,
			"job":             result.User.Job,
			"motto":           result.User.Motto,
			"created_at":      result.User.CreatedAt,
			"post_count":      result.PostCount,
			"comment_count":   result.CommentCount,
			"follower_count":  result.FollowerCount,
			"following_count": result.FollowingCount,
			"is_following":    result.IsFollowing,
			"level":           result.Level,
		},
		"posts": result.Posts,
	})
}

// ========== 搜索用户 ==========
// GET /api/users?q=xxx
func SearchUsers(c *gin.Context) {
	q := c.Query("q")
	users, err := data.ListUsers(q)
	if err != nil {
		models.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}

	type userItem struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Motto    string `json:"motto"`
	}
	items := make([]userItem, len(users))
	for i, u := range users {
		items[i] = userItem{u.ID, u.Username, u.Motto}
	}

	models.Success(c, "获取成功", gin.H{"users": items})
}
