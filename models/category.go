package models

import "time"

// Category 帖子分类
type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string    `gorm:"size:200" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
