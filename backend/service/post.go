package service

import (
	"community/backend/data"
	"community/backend/models"
	"community/backend/storage"
	"regexp"
	"time"
)

// mdImageRegex 匹配 Markdown 图片语法 ![](url)
var mdImageRegex = regexp.MustCompile(`!\[.*?\]\((/uploads/images/[^)]+)\)`)

// extractUploadURLs 从 Markdown 内容中提取所有上传文件 URL
func extractUploadURLs(content string) []string {
	matches := mdImageRegex.FindAllStringSubmatch(content, -1)
	urls := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) > 1 {
			urls = append(urls, m[1])
		}
	}
	return urls
}

// ─── 发帖 ───

// CreatePost 禁言检查 → 创建帖子
func CreatePost(userID uint, title, content string, categoryID uint, status string) (*models.Post, error) {
	if err := checkMuted(userID, categoryID); err != nil {
		return nil, err
	}

	post := &models.Post{
		Title:      title,
		Content:    content,
		Status:     status,
		UserID:     userID,
		CategoryID: categoryID,
	}
	if err := data.CreatePost(post); err != nil {
		return nil, ErrDBOpFail
	}
	return post, nil
}

// ─── 编辑 ───

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
	return data.FindPostByID(postID)
}

// ─── 删除 ───

// DeletePost 查帖子 → 权限检查 → 收集上传文件 → 删评论 → 删帖子 → 清理文件
func DeletePost(postID, userID uint, roles []string) error {
	post, err := data.FindPostByIDRaw(postID)
	if err != nil {
		return ErrNotFound("帖子")
	}

	if !CanManagePost(userID, roles, post) {
		return ErrForbidden
	}

	// 收集需要清理的上传文件
	var urlsToDelete []string
	urlsToDelete = append(urlsToDelete, extractUploadURLs(post.Content)...)

	comments, _ := data.ListCommentsByPost(postID)
	for _, c := range comments {
		urlsToDelete = append(urlsToDelete, extractUploadURLs(c.Content)...)
	}

	// 删数据库
	_ = data.DeleteCommentsByPost(postID)
	_ = data.DeletePost(postID)

	// 清磁盘
	for _, url := range urlsToDelete {
		storage.Store.DeleteImage(url)
	}
	storage.Store.DeletePostDir(postID)

	return nil
}

// ─── 权限判断 ───

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
	mute, err := data.FindActiveMute(userID, categoryID)
	if err != nil {
		return nil
	}

	until := mute.MutedUntil.Format("2006-01-02 15:04")
	_ = now()

	if mute.CategoryID == 0 {
		return ErrMuted(until)
	}
	return ErrMutedCategory(until)
}

func now() string {
	return time.Now().Format("2006-01-02 15:04")
}
