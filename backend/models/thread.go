package models

import "time"

// Thread 对话主题（多主题私信）
type Thread struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:100;not null" json:"title"`
	UserAID   uint      `gorm:"index;not null" json:"user_a_id"`
	UserBID   uint      `gorm:"index;not null" json:"user_b_id"`
	CreatedAt time.Time `json:"created_at"`
}
