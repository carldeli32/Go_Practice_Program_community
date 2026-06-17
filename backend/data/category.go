package data

import (
	"community/backend/config"
	"community/backend/models"
)

// ─── 查询 ───

// FindCategoryByID 按 ID 查分类
func FindCategoryByID(id uint) (*models.Category, error) {
	var cat models.Category
	err := config.DB.First(&cat, id).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// FindCategoryByName 按名称查分类
func FindCategoryByName(name string) (*models.Category, error) {
	var cat models.Category
	err := config.DB.Where("name = ?", name).First(&cat).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// ListCategories 列出所有分类（支持搜索）
func ListCategories(search string) ([]models.Category, error) {
	var categories []models.Category
	query := config.DB.Order("id ASC")
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	err := query.Find(&categories).Error
	return categories, err
}

// ─── 写入 ───

// CreateCategory 创建分类
func CreateCategory(cat *models.Category) error {
	return config.DB.Create(cat).Error
}

// ─── 删除 ───

// DeleteCategory 删除分类
func DeleteCategory(id uint) error {
	return config.DB.Delete(&models.Category{}, id).Error
}
