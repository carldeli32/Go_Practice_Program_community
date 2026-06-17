package data

import (
	"community/backend/config"
	"community/backend/models"
)

// ─── 写入 ───

// CreateComment 创建评论
func CreateComment(comment *models.Comment) error {
	return config.DB.Create(comment).Error
}

// ─── 查询 ───

// FindCommentByID 按 ID 查评论（带作者）
func FindCommentByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := config.DB.Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// FindCommentByIDRaw 按 ID 查评论（不带关联）
func FindCommentByIDRaw(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := config.DB.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListCommentsByPost 列出帖子的所有评论（带作者，按时间正序）
func ListCommentsByPost(postID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := config.DB.Where("post_id = ?", postID).
		Preload("User").
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}

// ─── 更新 ───

// UpdateCommentContent 更新评论内容
func UpdateCommentContent(id uint, content string) error {
	return config.DB.Model(&models.Comment{}).Where("id = ?", id).Update("content", content).Error
}

// ─── 删除 ───

// DeleteComment 删除单条评论
func DeleteComment(id uint) error {
	return config.DB.Delete(&models.Comment{}, id).Error
}
