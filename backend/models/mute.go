package models

import "time"

// Mute 禁言记录（CategoryID=0 表示全站禁言）
type Mute struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"uniqueIndex:idx_mute_uc;not null" json:"user_id"`
	CategoryID uint      `gorm:"uniqueIndex:idx_mute_uc;not null" json:"category_id"`
	MutedUntil time.Time `gorm:"not null" json:"muted_until"`
	CreatedAt  time.Time `json:"created_at"`
}
