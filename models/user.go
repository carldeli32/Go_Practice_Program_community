package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Email     string    `gorm:"size:100" json:"email"`
	Gender    string    `gorm:"size:10;default:''" json:"gender"`
	Age       int       `gorm:"default:0" json:"age"`
	Job       string    `gorm:"size:100;default:''" json:"job"`
	Motto     string    `gorm:"size:200;default:''" json:"motto"`
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
	IsBanned  bool      `gorm:"default:false" json:"is_banned"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserLevel struct {
	Name  string `json:"name"`
	Badge string `json:"badge"`
	Level int    `json:"level"`
}

// Level 计算用户等级（需传入 DB，由调用方提供）
func (u *User) Level(db *gorm.DB) UserLevel {
	var postCount, commentCount int64
	db.Model(&Post{}).Where("user_id = ?", u.ID).Count(&postCount)
	db.Model(&Comment{}).Where("user_id = ?", u.ID).Count(&commentCount)
	total := postCount + commentCount
	switch {
	case total >= 51:
		return UserLevel{Name: "社区长老", Badge: "👑", Level: 4}
	case total >= 21:
		return UserLevel{Name: "资深会员", Badge: "💎", Level: 3}
	case total >= 6:
		return UserLevel{Name: "活跃用户", Badge: "🔥", Level: 2}
	default:
		return UserLevel{Name: "初来乍到", Badge: "🌱", Level: 1}
	}
}
