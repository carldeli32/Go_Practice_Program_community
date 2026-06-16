package models

import "time"

// ModeratorCategory 版主管辖分类关联
type ModeratorCategory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"uniqueIndex:idx_mod_cat;not null" json:"user_id"`
	CategoryID uint      `gorm:"uniqueIndex:idx_mod_cat;not null" json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
}
