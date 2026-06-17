package models

import "time"

type Message struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ThreadID   uint      `gorm:"index;default:0" json:"thread_id"`
	FromUserID uint      `gorm:"index;not null" json:"from_user_id"`
	FromUser   User      `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
	ToUserID   uint      `gorm:"index;not null" json:"to_user_id"`
	ToUser     User      `gorm:"foreignKey:ToUserID" json:"to_user,omitempty"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	IsRead     bool      `gorm:"default:false" json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}
