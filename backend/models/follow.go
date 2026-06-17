package models

import "time"

// Follow 关注关系
type Follow struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	FollowerID uint      `gorm:"uniqueIndex:idx_follow;not null" json:"follower_id"`
	Follower   User      `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	FolloweeID uint      `gorm:"uniqueIndex:idx_follow;not null" json:"followee_id"`
	Followee   User      `gorm:"foreignKey:FolloweeID" json:"followee,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
