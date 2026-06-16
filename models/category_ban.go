package models

import "time"

// CategoryBan 版内封禁（禁止访问指定分类内容）
type CategoryBan struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"uniqueIndex:idx_catban_uc;not null" json:"user_id"`
	CategoryID uint      `gorm:"uniqueIndex:idx_catban_uc;not null" json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
}
