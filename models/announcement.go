package models

import "time"

// Announcement 站内公告（CategoryID=0 全站，非 0 版内）
type Announcement struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CategoryID uint      `gorm:"index;default:0" json:"category_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
