package data

import (
	"community/backend/config"
	"community/backend/models"
)

// ─── 写入 ───

// CreatePost 创建帖子
func CreatePost(post *models.Post) error {
	return config.DB.Create(post).Error
}

// ─── 查询 ───

// FindPostByID 按 ID 查帖子（带作者和分类）
func FindPostByID(id uint) (*models.Post, error) {
	var post models.Post
	err := config.DB.Preload("User").Preload("Category").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// FindPostByIDRaw 按 ID 查帖子（不带关联，用于权限检查）
func FindPostByIDRaw(id uint) (*models.Post, error) {
	var post models.Post
	err := config.DB.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// ─── 列表 ───

// ListPosts 分页列帖子（支持分类筛选和关键词搜索）
func ListPosts(page, pageSize int, categoryID, keyword string) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	base := config.DB.Model(&models.Post{})
	if categoryID != "" {
		base = base.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		base = base.Where("title LIKE ?", "%"+keyword+"%")
	}
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := base.Preload("User").Preload("Category").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

// FindPostsByUser 查找某用户最近 N 条帖子
func FindPostsByUser(userID uint, limit int) ([]models.Post, error) {
	var posts []models.Post
	err := config.DB.Where("user_id = ?", userID).
		Order("created_at DESC").Limit(limit).
		Find(&posts).Error
	return posts, err
}

// ─── 更新 ───

// UpdatePost 更新帖子字段
func UpdatePost(id uint, updates map[string]interface{}) error {
	return config.DB.Model(&models.Post{}).Where("id = ?", id).Updates(updates).Error
}

// ─── 删除 ───

// DeletePost 删除帖子（仅删帖子表）
func DeletePost(id uint) error {
	return config.DB.Delete(&models.Post{}, id).Error
}

// DeleteCommentsByPost 删除帖子下所有评论
func DeleteCommentsByPost(postID uint) error {
	return config.DB.Where("post_id = ?", postID).Delete(&models.Comment{}).Error
}

// ─── 迁移 ───

// MigratePostsToCategory 将某分类下帖子全部迁移到另一个分类
func MigratePostsToCategory(fromCategoryID, toCategoryID uint) error {
	return config.DB.Model(&models.Post{}).
		Where("category_id = ?", fromCategoryID).
		Update("category_id", toCategoryID).Error
}
