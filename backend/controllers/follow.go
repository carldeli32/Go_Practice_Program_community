package controllers

import (
	"net/http"

	"community/backend/data"
	"community/backend/middlewares"
	"community/backend/models"
	"community/backend/service"

	"github.com/gin-gonic/gin"
)

// ─── 关注 ───
func FollowUser(c *gin.Context) {
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := service.FollowUser(c.GetUint("user_id"), req.UserID); err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "关注成功 🤝", nil)
}

// ─── 取消关注 ───
func UnfollowUser(c *gin.Context) {
	myID := c.GetUint("user_id")
	targetID := middlewares.StrToUint(c.Param("user_id"))

	if err := service.UnfollowUser(myID, targetID); err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "已取消关注", nil)
}

// ─── 我的关注 ───
func GetMyFollowing(c *gin.Context) {
	follows, _ := data.ListFollowing(c.GetUint("user_id"))

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

// ─── 我的粉丝 ───
func GetMyFollowers(c *gin.Context) {
	follows, _ := data.ListFollowers(c.GetUint("user_id"))

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
