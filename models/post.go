package models

import "time"

// Post 帖子模型
type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CategoryID uint      `gorm:"index;not null;default:1" json:"category_id"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Comments   []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
