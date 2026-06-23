package data

import (
	"community/backend/config"
	"community/backend/models"
)

// ─── Thread ───

// CreateThread 创建对话主题
func CreateThread(thread *models.Thread) error {
	return config.DB.Create(thread).Error
}

// FindThreadByID 按 ID 查主题
func FindThreadByID(id uint) (*models.Thread, error) {
	var thread models.Thread
	err := config.DB.First(&thread, id).Error
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

// FindThreadsBetween 查找两人之间的所有主题
func FindThreadsBetween(userA, userB uint) ([]models.Thread, error) {
	var threads []models.Thread
	err := config.DB.Where(
		"(user_a_id = ? AND user_b_id = ?) OR (user_a_id = ? AND user_b_id = ?)",
		userA, userB, userB, userA,
	).Order("created_at ASC").Find(&threads).Error
	return threads, err
}

// DeleteThread 删除主题
func DeleteThread(id uint) error {
	return config.DB.Delete(&models.Thread{}, id).Error
}

// CountMessagesInThread 统计主题内的消息数
func CountMessagesInThread(threadID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Message{}).Where("thread_id = ?", threadID).Count(&count).Error
	return count, err
}

// DeleteMessagesInThread 删除主题下的所有消息
func DeleteMessagesInThread(threadID uint) error {
	return config.DB.Where("thread_id = ?", threadID).Delete(&models.Message{}).Error
}

// ─── Message ───

// CreateMessage 发送私信
func CreateMessage(msg *models.Message) error {
	return config.DB.Create(msg).Error
}

// FindMessageByID 按 ID 查私信（带收发双方）
func FindMessageByID(id uint) (*models.Message, error) {
	var msg models.Message
	err := config.DB.Preload("FromUser").Preload("ToUser").First(&msg, id).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// ListMessagesBetween 列出两人之间（可选某主题）的私信
func ListMessagesBetween(u1, u2, threadID uint) ([]models.Message, error) {
	var messages []models.Message
	query := config.DB.Where(
		"(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
		u1, u2, u2, u1,
	)
	if threadID > 0 {
		query = query.Where("thread_id = ?", threadID)
	}
	err := query.Preload("FromUser").Preload("ToUser").
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

// ListAllRelatedMessages 列出与用户相关的所有私信（收+发，用于会话聚合）
func ListAllRelatedMessages(userID uint) ([]models.Message, error) {
	var messages []models.Message
	err := config.DB.Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Preload("FromUser").Preload("ToUser").
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}

// ─── 未读 ───

// CountUnread 统计用户未读私信数
func CountUnread(userID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Message{}).
		Where("to_user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}

// CountUnreadFrom 统计来自特定用户的未读私信数
func CountUnreadFrom(fromUserID, toUserID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Message{}).
		Where("from_user_id = ? AND to_user_id = ? AND is_read = false", fromUserID, toUserID).
		Count(&count).Error
	return count, err
}

// MarkAllRead 标记某用户所有未读为已读
func MarkAllRead(userID uint) error {
	return config.DB.Model(&models.Message{}).
		Where("to_user_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}

// MarkReadFrom 标记来自特定用户的消息为已读
func MarkReadFrom(fromUserID, toUserID uint) error {
	return config.DB.Model(&models.Message{}).
		Where("from_user_id = ? AND to_user_id = ? AND is_read = false", fromUserID, toUserID).
		Update("is_read", true).Error
}

// RecallMessage 撤回消息（仅发送者，标记 is_recalled）
func RecallMessage(msgID, userID uint) error {
	return config.DB.Model(&models.Message{}).
		Where("id = ? AND from_user_id = ?", msgID, userID).
		Update("is_recalled", true).Error
}
