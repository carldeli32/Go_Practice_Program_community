package controllers

import (
	"net/http"
	"strconv"

	"community/config"
	"community/models"

	"github.com/gin-gonic/gin"
)

// ========== 创建对话主题 ==========
// POST /api/threads  (需登录)
// Body: { "with_user_id": 2, "title": "旅行计划" }
func CreateThread(c *gin.Context) {
	var req struct {
		WithUserID uint   `json:"with_user_id" binding:"required"`
		Title      string `json:"title" binding:"required,min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	myID := c.GetUint("user_id")
	if myID == req.WithUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能和自己对话"})
		return
	}

	// 标准化 user_a < user_b 避免重复主题判断问题
	a, b := myID, req.WithUserID
	if a > b {
		a, b = b, a
	}

	thread := models.Thread{Title: req.Title, UserAID: a, UserBID: b}
	if err := config.DB.Create(&thread).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "主题已创建 📌", "thread": thread})
}

// ========== 获取与某人的所有主题 ==========
// GET /api/threads?with=2  (需登录)
func GetThreads(c *gin.Context) {
	myID := c.GetUint("user_id")
	withID, _ := strconv.Atoi(c.Query("with"))

	var threads []models.Thread
	config.DB.Where(
		"(user_a_id = ? AND user_b_id = ?) OR (user_a_id = ? AND user_b_id = ?)",
		myID, withID, withID, myID,
	).Order("created_at ASC").Find(&threads)

	// 统计每个主题的消息数
	type threadItem struct {
		models.Thread
		MessageCount int64 `json:"message_count"`
	}
	items := make([]threadItem, len(threads))
	for i, t := range threads {
		var count int64
		config.DB.Model(&models.Message{}).Where("thread_id = ?", t.ID).Count(&count)
		items[i] = threadItem{Thread: t, MessageCount: count}
	}

	c.JSON(http.StatusOK, gin.H{"threads": items})
}

// ========== 删除主题 ==========
// DELETE /api/threads/:id  (需登录)
func DeleteThread(c *gin.Context) {
	myID := c.GetUint("user_id")
	id := c.Param("id")

	var thread models.Thread
	if err := config.DB.First(&thread, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "主题不存在"})
		return
	}
	if thread.UserAID != myID && thread.UserBID != myID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作"})
		return
	}

	config.DB.Where("thread_id = ?", thread.ID).Delete(&models.Message{})
	config.DB.Delete(&thread)
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}

// ========== 发送私信（指定主题）==========
// POST /api/messages  (需登录)
// Body: { "to_user_id": 2, "thread_id": 1, "content": "你好！" }
func SendMessage(c *gin.Context) {
	var req struct {
		ToUserID uint   `json:"to_user_id" binding:"required"`
		ThreadID uint   `json:"thread_id"`
		Content  string `json:"content" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	fromUserID := c.GetUint("user_id")
	if fromUserID == req.ToUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能给自己发私信"})
		return
	}

	var toUser models.User
	if err := config.DB.First(&toUser, req.ToUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	msg := models.Message{
		ThreadID:   req.ThreadID,
		FromUserID: fromUserID,
		ToUserID:   req.ToUserID,
		Content:    req.Content,
	}
	if err := config.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送失败"})
		return
	}
	config.DB.Preload("FromUser").Preload("ToUser").First(&msg, msg.ID)
	c.JSON(http.StatusOK, gin.H{"message": "发送成功 ✉️", "data": msg})
}

// ========== 会话列表（按主题分组）==========
func GetConversations(c *gin.Context) {
	userID := c.GetUint("user_id")

	var messages []models.Message
	config.DB.Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Preload("FromUser").Preload("ToUser").
		Order("created_at DESC").
		Find(&messages)

	type Conversation struct {
		Partner     models.User    `json:"partner"`
		LastMessage models.Message `json:"last_message"`
		UnreadCount int64          `json:"unread_count"`
	}

	seen := make(map[uint]bool)
	var conversations []Conversation

	for _, msg := range messages {
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

		var unread int64
		config.DB.Model(&models.Message{}).
			Where("from_user_id = ? AND to_user_id = ? AND is_read = false", partner.ID, userID).
			Count(&unread)

		conversations = append(conversations, Conversation{
			Partner:     partner,
			LastMessage: msg,
			UnreadCount: unread,
		})
	}

	c.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

// ========== 与某人的某主题对话详情 ==========
// GET /api/messages/:user_id?thread=1  (需登录)
func GetConversation(c *gin.Context) {
	userID := c.GetUint("user_id")
	partnerID, _ := strconv.Atoi(c.Param("user_id"))
	threadID, _ := strconv.Atoi(c.DefaultQuery("thread", "0"))

	var messages []models.Message
	query := config.DB.Where(
		"(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
		userID, partnerID, partnerID, userID,
	)
	if threadID > 0 {
		query = query.Where("thread_id = ?", threadID)
	}
	query.Preload("FromUser").Preload("ToUser").Order("created_at ASC").Find(&messages)

	var partner models.User
	config.DB.First(&partner, partnerID)

	c.JSON(http.StatusOK, gin.H{
		"partner":  gin.H{"id": partner.ID, "username": partner.Username},
		"messages": messages,
	})
}

func GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")
	var count int64
	config.DB.Model(&models.Message{}).Where("to_user_id = ? AND is_read = false", userID).Count(&count)
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func MarkAllRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	config.DB.Model(&models.Message{}).Where("to_user_id = ? AND is_read = false", userID).Update("is_read", true)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func MarkMessagesRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	partnerID := c.Param("user_id")
	config.DB.Model(&models.Message{}).
		Where("from_user_id = ? AND to_user_id = ? AND is_read = false", partnerID, userID).
		Update("is_read", true)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
