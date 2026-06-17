package service

import (
	"community/backend/data"
	"community/backend/models"
	"time"
)

// ─── 发帖 ───

// CreatePost 禁言检查 → 创建帖子
func CreatePost(userID uint, title, content string, categoryID uint) (*models.Post, error) {
	// 禁言检查
	if err := checkMuted(userID, categoryID); err != nil {
		return nil, err
	}

	post := &models.Post{
		Title:      title,
		Content:    content,
		UserID:     userID,
		CategoryID: categoryID,
	}
	if err := data.CreatePost(post); err != nil {
		return nil, ErrDBOpFail
	}
	return post, nil
}

// ─── 编辑 ───

// UpdatePost 查帖子 → 权限检查 → 更新
func UpdatePost(postID, userID uint, title, content string, roles []string) (*models.Post, error) {
	post, err := data.FindPostByIDRaw(postID)
	if err != nil {
		return nil, ErrNotFound("帖子")
	}

	if !CanManagePost(userID, roles, post) {
		return nil, ErrForbidden
	}

	updates := map[string]interface{}{}
	if title != "" {
		updates["title"] = title
	}
	if content != "" {
		updates["content"] = content
	}
	if len(updates) == 0 {
		return nil, ErrNoUpdateContent
	}

	if err := data.UpdatePost(postID, updates); err != nil {
		return nil, ErrDBOpFail
	}
	// 重新查询以返回完整数据
	return data.FindPostByID(postID)
}

// ─── 删除 ───

// DeletePost 查帖子 → 权限检查 → 删评论 → 删帖子
func DeletePost(postID, userID uint, roles []string) error {
	post, err := data.FindPostByIDRaw(postID)
	if err != nil {
		return ErrNotFound("帖子")
	}

	if !CanManagePost(userID, roles, post) {
		return ErrForbidden
	}

	_ = data.DeleteCommentsByPost(postID)
	_ = data.DeletePost(postID)
	return nil
}

// ─── 权限判断 ───

// CanManagePost 作者本人 or 具有 manage_any 权限 or 管辖区版内
func CanManagePost(userID uint, roles []string, post *models.Post) bool {
	if post.UserID == userID {
		return true
	}
	if models.HasPerm(roles, "post.manage_any") {
		return true
	}
	if models.HasPerm(roles, "post.manage_category") {
		ok, _ := data.ExistsModeratorCategory(userID, post.CategoryID)
		return ok
	}
	return false
}

// ─── 禁言检查 ───

func checkMuted(userID, categoryID uint) *AppError {
	now := time.Now().Format("2006-01-02 15:04")

	mute, err := data.FindActiveMute(userID, categoryID)
	if err != nil {
		return nil // 无有效禁言
	}

	until := mute.MutedUntil.Format("2006-01-02 15:04")
	_ = now

	if mute.CategoryID == 0 {
		return ErrMuted(until)
	}
	return ErrMutedCategory(until)
}
