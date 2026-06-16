package service

import (
	"community/data"
	"community/models"
)

// ─── 评论完整视图（含权限标记） ───

// CommentWithPerm 带权限标记的评论
type CommentWithPerm struct {
	models.Comment
	CanEdit   bool `json:"can_edit"`
	CanDelete bool `json:"can_delete"`
}

// ─── 创建 ───

// CreateComment 禁言检查 → 创建评论
func CreateComment(postID, userID uint, content string) (*models.Comment, error) {
	// 检查帖子存在
	if _, err := data.FindPostByIDRaw(postID); err != nil {
		return nil, ErrNotFound("帖子")
	}

	// 禁言检查
	if err := checkMuted(userID, postCategoryID(postID)); err != nil {
		return nil, err
	}

	comment := &models.Comment{
		Content: content,
		UserID:  userID,
		PostID:  postID,
	}
	if err := data.CreateComment(comment); err != nil {
		return nil, ErrDBOpFail
	}

	// 返回带作者的完整评论
	return data.FindCommentByID(comment.ID)
}

// ─── 编辑 ───

// UpdateComment 查评论 → 权限检查 → 更新
func UpdateComment(commentID, userID uint, content string, roles []string) (*models.Comment, error) {
	comment, err := data.FindCommentByIDRaw(commentID)
	if err != nil {
		return nil, ErrNotFound("评论")
	}

	if !CanManageComment(userID, roles, comment) {
		return nil, ErrForbidden
	}

	if err := data.UpdateCommentContent(commentID, content); err != nil {
		return nil, ErrDBOpFail
	}
	return data.FindCommentByID(commentID)
}

// ─── 删除 ───

// DeleteComment 查评论 → 权限检查 → 删除
func DeleteComment(commentID, userID uint, roles []string) error {
	comment, err := data.FindCommentByIDRaw(commentID)
	if err != nil {
		return ErrNotFound("评论")
	}

	if !CanManageComment(userID, roles, comment) {
		return ErrForbidden
	}

	return data.DeleteComment(commentID)
}

// ─── 权限计算 ───

// CanManageComment 作者本人 or manage_any or manage_category（管辖区）
func CanManageComment(userID uint, roles []string, comment *models.Comment) bool {
	if comment.UserID == userID {
		return true
	}
	if models.HasPerm(roles, "comment.manage_any") {
		return true
	}
	if models.HasPerm(roles, "comment.manage_category") {
		catID := postCategoryID(comment.PostID)
		ok, _ := data.ExistsModeratorCategory(userID, catID)
		return ok
	}
	return false
}

// BuildCommentPermissions 给评论列表计算 can_edit / can_delete
func BuildCommentPermissions(comments []models.Comment, currentUserID uint, roles []string) []CommentWithPerm {
	items := make([]CommentWithPerm, len(comments))
	for i, c := range comments {
		canEdit := c.UserID == currentUserID
		canDelete := c.UserID == currentUserID

		if !canEdit && models.HasPerm(roles, "comment.manage_any") {
			canEdit, canDelete = true, true
		}
		if !canEdit && models.HasPerm(roles, "comment.manage_category") {
			catID := postCategoryID(c.PostID)
			ok, _ := data.ExistsModeratorCategory(currentUserID, catID)
			if ok {
				canEdit, canDelete = true, true
			}
		}

		items[i] = CommentWithPerm{
			Comment:   c,
			CanEdit:   canEdit,
			CanDelete: canDelete,
		}
	}
	return items
}

// ─── 辅助 ───

// postCategoryID 获取帖子所属分类 ID（用于禁言检查）
func postCategoryID(postID uint) uint {
	post, err := data.FindPostByIDRaw(postID)
	if err != nil {
		return 0
	}
	return post.CategoryID
}
