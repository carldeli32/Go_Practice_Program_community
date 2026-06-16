package service

import (
	"community/data"
	"community/models"
)

// ─── 创建 ───

// CreateCategory 创建分类（含重名检查）
func CreateCategory(name, description string) (*models.Category, error) {
	if _, err := data.FindCategoryByName(name); err == nil {
		return nil, ErrCategoryExists
	}

	cat := &models.Category{Name: name, Description: description}
	if err := data.CreateCategory(cat); err != nil {
		return nil, ErrDBOpFail
	}
	return cat, nil
}

// ─── 删除 ───

// DeleteCategory 删除分类（综合讨论不可删，帖子迁移）
func DeleteCategory(categoryID uint) error {
	if categoryID == 1 {
		return ErrCannotDeleteRoot
	}

	if _, err := data.FindCategoryByID(categoryID); err != nil {
		return ErrNotFound("分类")
	}

	// 迁移帖子到综合讨论
	_ = data.MigratePostsToCategory(categoryID, 1)
	return data.DeleteCategory(categoryID)
}
