package service

import (
	"community/backend/data"
	"community/backend/middlewares"
	"community/backend/models"
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

	return data.FindCommentByID(comment.ID)
}

// ─── 编辑 ───

// UpdateComment 更新评论内容（权限由中间件保证）
func UpdateComment(comment *models.Comment, content string) (*models.Comment, error) {
	if err := data.UpdateCommentContent(comment.ID, content); err != nil {
		return nil, ErrDBOpFail
	}
	return data.FindCommentByID(comment.ID)
}

// ─── 删除 ───

// DeleteComment 删除评论（权限由中间件保证）
func DeleteComment(comment *models.Comment) error {
	return data.DeleteComment(comment.ID)
}

// ─── 权限计算 ───

// BuildCommentPermissions 给评论列表计算 can_edit / can_delete
// 使用统一的 CheckRelation 引擎，与中间件保持逻辑一致
func BuildCommentPermissions(comments []models.Comment, currentUserID uint, roles []string) []CommentWithPerm {
	items := make([]CommentWithPerm, len(comments))
	catCache := make(map[uint]uint) // postID → categoryID，避免 N+1 查询
	for i, c := range comments {
		catID, ok := catCache[c.PostID]
		if !ok {
			catID = postCategoryID(c.PostID)
			catCache[c.PostID] = catID
		}
		items[i] = CommentWithPerm{
			Comment:   c,
			CanEdit:   middlewares.CheckRelation(currentUserID, c.UserID, catID, roles, "comment.manage_any", "comment.manage_category"),
			CanDelete: middlewares.CheckRelation(currentUserID, c.UserID, catID, roles, "comment.manage_any", "comment.manage_category"),
		}
	}
	return items
}

// ─── 辅助 ───

// postCategoryID 获取帖子所属分类 ID
func postCategoryID(postID uint) uint {
	post, err := data.FindPostByIDRaw(postID)
	if err != nil {
		return 0
	}
	return post.CategoryID
}
