package service

import (
	"community/backend/data"
	"community/backend/middlewares"
	"community/backend/models"

	"golang.org/x/crypto/bcrypt"
)

// ─── 注册 ───

// Register 查重 → 加密 → 创建
func Register(username, password, email, gender, motto, job string, age int) (*models.User, error) {
	// 查重
	if _, err := data.FindUserByUsername(username); err == nil {
		return nil, ErrUserExists
	}

	// 加密
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrPasswordHashFail
	}

	user := &models.User{
		Username: username,
		Password: string(hashed),
		Email:    email,
		Gender:   gender,
		Age:      age,
		Job:      job,
		Motto:    motto,
	}
	if err := data.CreateUser(user); err != nil {
		return nil, ErrRegisterFail("注册失败: " + err.Error())
	}
	return user, nil
}

// ─── 登录 ───

// Login 查用户 → 验密码 → 封禁检查 → 生成 token
func Login(username, password string) (string, *models.User, error) {
	user, err := data.FindUserByUsernameWithRoles(username)
	if err != nil {
		return "", nil, ErrBadCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, ErrBadCredentials
	}

	if user.IsBanned {
		return "", nil, ErrBanned
	}

	token, err := middlewares.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, ErrTokenGenFail
	}

	return token, user, nil
}

// ─── 用户主页 ───

// UserProfileResult 用户主页聚合结果
type UserProfileResult struct {
	User           models.User        `json:"user"`
	Posts          []models.Post      `json:"posts"`
	PostCount      int64              `json:"post_count"`
	CommentCount   int64              `json:"comment_count"`
	FollowerCount  int64              `json:"follower_count"`
	FollowingCount int64              `json:"following_count"`
	IsFollowing    bool               `json:"is_following"`
	Level          models.UserLevel   `json:"level"`
}

// GetUserProfile 聚合用户主页数据
func GetUserProfile(userID, currentUserID uint) (*UserProfileResult, error) {
	user, err := data.FindUserByID(userID)
	if err != nil {
		return nil, ErrNotFound("用户")
	}

	postCount, _ := data.CountPostsByUser(userID)
	commentCount, _ := data.CountCommentsByUser(userID)
	followerCount, _ := data.CountFollowers(userID)
	followingCount, _ := data.CountFollowing(userID)

	var isFollowing bool
	if currentUserID > 0 {
		isFollowing, _ = data.IsFollowing(currentUserID, userID)
	}

	posts, _ := data.FindPostsByUser(userID, 10)

	return &UserProfileResult{
		User:           *user,
		Posts:          posts,
		PostCount:      postCount,
		CommentCount:   commentCount,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
		IsFollowing:    isFollowing,
		Level:          ComputeUserLevel(postCount, commentCount),
	}, nil
}

// ComputeUserLevel 根据发帖+评论总数计算等级
func ComputeUserLevel(postCount, commentCount int64) models.UserLevel {
	total := postCount + commentCount
	switch {
	case total >= 51:
		return models.UserLevel{Name: "社区长老", Badge: "👑", Level: 4}
	case total >= 21:
		return models.UserLevel{Name: "资深会员", Badge: "💎", Level: 3}
	case total >= 6:
		return models.UserLevel{Name: "活跃用户", Badge: "🔥", Level: 2}
	default:
		return models.UserLevel{Name: "初来乍到", Badge: "🌱", Level: 1}
	}
}
