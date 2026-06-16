package service

import (
	"community/data"
	"community/models"
)

// ─── 关注 ───

// FollowUser 关注用户（含自关注检查 + 目标存在检查 + 重复检查）
func FollowUser(followerID, followeeID uint) error {
	if followerID == followeeID {
		return ErrCannotSelfFollow
	}

	// 对方存在
	if _, err := data.FindUserByID(followeeID); err != nil {
		return ErrNotFound("用户")
	}

	// 已关注
	if ok, _ := data.ExistsFollow(followerID, followeeID); ok {
		return ErrAlreadyFollowing
	}

	follow := &models.Follow{FollowerID: followerID, FolloweeID: followeeID}
	if err := data.CreateFollow(follow); err != nil {
		return ErrDBOpFail
	}
	return nil
}

// UnfollowUser 取消关注
func UnfollowUser(followerID, followeeID uint) error {
	if err := data.DeleteFollow(followerID, followeeID); err != nil {
		return ErrNotFollowing
	}
	return nil
}
