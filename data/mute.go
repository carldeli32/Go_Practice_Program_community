package data

import (
	"community/config"
	"community/models"
	"time"
)

// ─── 禁言 ───

// FindActiveMute 查找某用户的有效禁言记录（全站 category=0 或指定分类）
func FindActiveMute(userID, categoryID uint) (*models.Mute, error) {
	now := time.Now()

	// 先查全站禁言
	var global models.Mute
	err := config.DB.Where("user_id = ? AND category_id = 0 AND muted_until > ?", userID, now).
		First(&global).Error
	if err == nil {
		return &global, nil
	}

	// 再查指定分类禁言
	var catMute models.Mute
	err = config.DB.Where("user_id = ? AND category_id = ? AND muted_until > ?", userID, categoryID, now).
		First(&catMute).Error
	if err == nil {
		return &catMute, nil
	}

	return nil, err
}

// ─── 版主管辖 ───

// ExistsModeratorCategory 检查用户是否是指定分类的版主
func ExistsModeratorCategory(userID, categoryID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.ModeratorCategory{}).
		Where("user_id = ? AND category_id = ?", userID, categoryID).
		Count(&count).Error
	return count > 0, err
}
