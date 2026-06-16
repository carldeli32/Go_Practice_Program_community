package data

import (
	"community/config"
	"community/models"
)

// ─── 写入 ───

// CreateFollow 创建关注关系
func CreateFollow(follow *models.Follow) error {
	return config.DB.Create(follow).Error
}

// ─── 删除 ───

// DeleteFollow 取消关注（返回受影响行数）
func DeleteFollow(followerID, followeeID uint) error {
	return config.DB.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&models.Follow{}).Error
}

// ─── 判断 ───

// ExistsFollow 检查关注关系是否存在
func ExistsFollow(followerID, followeeID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Follow{}).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&count).Error
	return count > 0, err
}

// ─── 列表 ───

// ListFollowing 我关注的人（带 Followee 用户数据）
func ListFollowing(userID uint) ([]models.Follow, error) {
	var follows []models.Follow
	err := config.DB.Where("follower_id = ?", userID).
		Preload("Followee").
		Order("created_at DESC").
		Find(&follows).Error
	return follows, err
}

// ListFollowers 关注我的人（粉丝，带 Follower 用户数据）
func ListFollowers(userID uint) ([]models.Follow, error) {
	var follows []models.Follow
	err := config.DB.Where("followee_id = ?", userID).
		Preload("Follower").
		Order("created_at DESC").
		Find(&follows).Error
	return follows, err
}
