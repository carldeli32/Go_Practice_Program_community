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

// UpdatePost 更新帖子内容（权限由中间件保证）
func UpdatePost(post *models.Post, title, content string, categoryID uint) (*models.Post, error) {
	updates := map[string]interface{}{}
	if title != "" {
		updates["title"] = title
	}
	if content != "" {
		updates["content"] = content
	}
	if categoryID > 0 {
		updates["category_id"] = categoryID
	}
	if len(updates) == 0 {
		return nil, ErrNoUpdateContent
	}

	if err := data.UpdatePost(post.ID, updates); err != nil {
		return nil, ErrDBOpFail
	}
	return data.FindPostByID(post.ID)
}

// ─── 删除 ───

// DeletePost 删评论 → 删帖子 → 清理上传文件（权限由中间件保证）
func DeletePost(post *models.Post) error {
	// 收集需要清理的上传文件
	var urlsToDelete []string
	urlsToDelete = append(urlsToDelete, extractUploadURLs(post.Content)...)

	comments, _ := data.ListCommentsByPost(post.ID)
	for _, c := range comments {
		urlsToDelete = append(urlsToDelete, extractUploadURLs(c.Content)...)
	}

	// 删数据库
	_ = data.DeleteCommentsByPost(post.ID)
	_ = data.DeletePost(post.ID)

	// 清磁盘
	for _, url := range urlsToDelete {
		storage.Store.DeleteImage(url)
	}
	storage.Store.DeletePostDir(post.ID)

	return nil
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
