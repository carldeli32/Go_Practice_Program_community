package controllers

import (
	"net/http"

	"community/config"
	"community/models"

	"github.com/gin-gonic/gin"
)

// ========== 关注某人 ==========
// POST /api/follow  (需登录)
// Body: { "user_id": 2 }
func FollowUser(c *gin.Context) {
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	myID := c.GetUint("user_id")

	if myID == req.UserID {
		models.Error(c, http.StatusBadRequest, "不能关注自己")
		return
	}

	// 检查对方是否存在
	var target models.User
	if err := config.DB.First(&target, req.UserID).Error; err != nil {
		models.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	// 检查是否已关注
	var exist models.Follow
	if err := config.DB.Where("follower_id = ? AND followee_id = ?", myID, req.UserID).First(&exist).Error; err == nil {
		models.Error(c, http.StatusConflict, "已关注该用户")
		return
	}

	follow := models.Follow{FollowerID: myID, FolloweeID: req.UserID}
	if err := config.DB.Create(&follow).Error; err != nil {
		models.Error(c, http.StatusInternalServerError, "关注失败")
		return
	}

	models.Success(c, "关注成功 🤝", nil)
}

// ========== 取消关注 ==========
// DELETE /api/follow/:user_id  (需登录)
func UnfollowUser(c *gin.Context) {
	myID := c.GetUint("user_id")
	targetID := c.Param("user_id")

	result := config.DB.Where("follower_id = ? AND followee_id = ?", myID, targetID).Delete(&models.Follow{})
	if result.RowsAffected == 0 {
		models.Error(c, http.StatusNotFound, "未关注该用户")
		return
	}

	models.Success(c, "已取消关注", nil)
}

// ========== 我关注的人 ==========
// GET /api/following  (需登录)
func GetMyFollowing(c *gin.Context) {
	myID := c.GetUint("user_id")

	var follows []models.Follow
	config.DB.Where("follower_id = ?", myID).
		Preload("Followee").
		Order("created_at DESC").
		Find(&follows)

	// 提取用户列表
	users := make([]gin.H, len(follows))
	for i, f := range follows {
		users[i] = gin.H{
			"id":          f.FolloweeID,
			"username":    f.Followee.Username,
			"motto":       f.Followee.Motto,
			"followed_at": f.CreatedAt,
		}
	}

	models.Success(c, "获取成功", gin.H{"users": users, "total": len(users)})
}

// ========== 关注我的人（粉丝）==========
// GET /api/followers  (需登录)
func GetMyFollowers(c *gin.Context) {
	myID := c.GetUint("user_id")

	var follows []models.Follow
	config.DB.Where("followee_id = ?", myID).
		Preload("Follower").
		Order("created_at DESC").
		Find(&follows)

	users := make([]gin.H, len(follows))
	for i, f := range follows {
		users[i] = gin.H{
			"id":          f.FollowerID,
			"username":    f.Follower.Username,
			"motto":       f.Follower.Motto,
			"followed_at": f.CreatedAt,
		}
	}

	models.Success(c, "获取成功", gin.H{"users": users, "total": len(users)})
}
