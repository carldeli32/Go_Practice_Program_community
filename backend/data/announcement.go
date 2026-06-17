package data

import (
	"community/backend/config"
	"community/backend/models"
)

// ─── 写入 ───

// CreateAnnouncement 创建公告
func CreateAnnouncement(a *models.Announcement) error {
	return config.DB.Create(a).Error
}

// ─── 删除 ───

// DeleteAnnouncementsByCategory 删除某分类的所有公告（categoryID=0 删全站公告）
func DeleteAnnouncementsByCategory(categoryID uint) error {
	return config.DB.Where("category_id = ?", categoryID).Delete(&models.Announcement{}).Error
}

// ─── 查询 ───

// ListAnnouncements 列出公告（categoryID="0" 只查全站，否则全站+指定分类）
func ListAnnouncements(categoryID string) ([]models.Announcement, error) {
	var announcements []models.Announcement
	if categoryID == "0" {
		err := config.DB.Where("category_id = 0").
			Order("created_at DESC").
			Find(&announcements).Error
		return announcements, err
	}
	err := config.DB.Where("category_id = 0 OR category_id = ?", categoryID).
		Order("created_at DESC").
		Find(&announcements).Error
	return announcements, err
}
