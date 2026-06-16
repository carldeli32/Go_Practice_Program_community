package service

import (
	"community/data"
	"community/models"
)

// ─── 对话主题 ───

// ThreadWithCount 主题 + 消息数
type ThreadWithCount struct {
	models.Thread
	MessageCount int64 `json:"message_count"`
}

// CreateThread 创建对话主题（标准化 a < b）
func CreateThread(userA, userB uint, title string) (*models.Thread, error) {
	if userA == userB {
		return nil, ErrCannotSelfThread
	}

	// 标准化
	a, b := userA, userB
	if a > b {
		a, b = b, a
	}

	thread := &models.Thread{Title: title, UserAID: a, UserBID: b}
	if err := data.CreateThread(thread); err != nil {
		return nil, ErrDBOpFail
	}
	return thread, nil
}

// GetThreads 获取两人之间的所有主题（带消息数）
func GetThreads(userA, userB uint) ([]ThreadWithCount, error) {
	threads, err := data.FindThreadsBetween(userA, userB)
	if err != nil {
		return nil, ErrDBOpFail
	}

	items := make([]ThreadWithCount, len(threads))
	for i, t := range threads {
		count, _ := data.CountMessagesInThread(t.ID)
		items[i] = ThreadWithCount{Thread: t, MessageCount: count}
	}
	return items, nil
}

// DeleteThread 删除主题（需检查是否参与者）
func DeleteThread(threadID, userID uint) error {
	thread, err := data.FindThreadByID(threadID)
	if err != nil {
		return ErrNotFound("主题")
	}

	if thread.UserAID != userID && thread.UserBID != userID {
		return ErrForbidden
	}

	_ = data.DeleteMessagesInThread(threadID)
	_ = data.DeleteThread(threadID)
	return nil
}

// ─── 发送私信 ───

// SendMessage 发送私信
func SendMessage(fromID, toID, threadID uint, content string) (*models.Message, error) {
	if fromID == toID {
		return nil, ErrCannotSelfMessage
	}

	// 验证对方存在
	if _, err := data.FindUserByID(toID); err != nil {
		return nil, ErrNotFound("用户")
	}

	msg := &models.Message{
		ThreadID:   threadID,
		FromUserID: fromID,
		ToUserID:   toID,
		Content:    content,
	}
	if err := data.CreateMessage(msg); err != nil {
		return nil, ErrDBOpFail
	}

	return data.FindMessageByID(msg.ID)
}

// ─── 会话列表 ───

// Conversation 会话摘要
type Conversation struct {
	Partner     models.User    `json:"partner"`
	LastMessage models.Message `json:"last_message"`
	UnreadCount int64          `json:"unread_count"`
}

// GetConversations 获取用户的会话列表（按最后消息去重）
func GetConversations(userID uint) ([]Conversation, error) {
	msgs, err := data.ListAllRelatedMessages(userID)
	if err != nil {
		return nil, ErrDBOpFail
	}

	seen := make(map[uint]bool)
	var conversations []Conversation

	for _, msg := range msgs {
		var partner models.User
		if msg.FromUserID == userID {
			partner = msg.ToUser
		} else {
			partner = msg.FromUser
		}
		if seen[partner.ID] {
			continue
		}
		seen[partner.ID] = true

		unread, _ := data.CountUnreadFrom(partner.ID, userID)

		conversations = append(conversations, Conversation{
			Partner:     partner,
			LastMessage: msg,
			UnreadCount: unread,
		})
	}
	return conversations, nil
}

// ─── 对话详情 ───

// ConversationDetail 与某人的某主题对话
type ConversationDetail struct {
	Partner  models.User       `json:"partner"`
	Messages []models.Message  `json:"messages"`
}

// GetConversation 获取与某人的对话详情
func GetConversation(userID, partnerID, threadID uint) (*ConversationDetail, error) {
	messages, err := data.ListMessagesBetween(userID, partnerID, threadID)
	if err != nil {
		return nil, ErrDBOpFail
	}

	partner, err := data.FindUserByID(partnerID)
	if err != nil {
		return nil, ErrNotFound("用户")
	}

	// 只返回必要的用户字段
	return &ConversationDetail{
		Partner:  models.User{ID: partner.ID, Username: partner.Username},
		Messages: messages,
	}, nil
}
