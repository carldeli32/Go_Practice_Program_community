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

// ─── 带分类查询 ───

// CommentWithCategory 评论 + 所属帖子分类 ID（用于权限判断）
type CommentWithCategory struct {
	models.Comment
	CategoryID uint `json:"category_id"`
}

// FindCommentByIDWithCategory 按 ID 查评论，JOIN post 获取 category_id
func FindCommentByIDWithCategory(id uint) (*CommentWithCategory, error) {
	var result CommentWithCategory
	err := config.DB.Table("comments").
		Select("comments.*, posts.category_id").
		Joins("JOIN posts ON posts.id = comments.post_id").
		Where("comments.id = ?", id).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ─── 删除 ───

// DeleteComment 删除单条评论
func DeleteComment(id uint) error {
	return config.DB.Delete(&models.Comment{}, id).Error
}
