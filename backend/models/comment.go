package models

import "time"

// Comment 评论模型
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PostID    uint      `gorm:"index;not null" json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}
