package models

import "time"

// Role 角色（super_admin / admin / moderator / user）
type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:30;uniqueIndex;not null" json:"name"`
	Description string    `gorm:"size:200" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
