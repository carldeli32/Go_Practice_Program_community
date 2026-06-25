package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"community/backend/data"
	"community/backend/middlewares"
	"community/backend/models"
	"community/backend/service"
	"community/backend/sse"

	"github.com/gin-gonic/gin"
)

// ─── 创建对话主题 ───
func CreateThread(c *gin.Context) {
	var req struct {
		WithUserID uint   `json:"with_user_id" binding:"required"`
		Title      string `json:"title" binding:"required,min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	thread, err := service.CreateThread(c.GetUint("user_id"), req.WithUserID, req.Title)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "主题已创建 📌", gin.H{"thread": thread})
}

// ─── 获取与某人的所有主题 ───
func GetThreads(c *gin.Context) {
	myID := c.GetUint("user_id")
	withID, _ := strconv.Atoi(c.Query("with"))

	threads, err := service.GetThreads(myID, uint(withID))
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "获取成功", gin.H{"threads": threads})
}

// ─── 删除主题 ───
func DeleteThread(c *gin.Context) {
	myID := c.GetUint("user_id")

	if err := service.DeleteThread(middlewares.StrToUint(c.Param("id")), myID); err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "已删除", nil)
}

// ─── 发送私信 ───
func SendMessage(c *gin.Context) {
	var req struct {
		ToUserID uint   `json:"to_user_id" binding:"required"`
		ThreadID uint   `json:"thread_id"`
		Content  string `json:"content" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	msg, err := service.SendMessage(c.GetUint("user_id"), req.ToUserID, req.ThreadID, req.Content)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "发送成功 ✉️", gin.H{"message": msg})
}

// ─── 会话列表 ───
func GetConversations(c *gin.Context) {
	userID := c.GetUint("user_id")

	conversations, err := service.GetConversations(userID)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "获取成功", gin.H{"conversations": conversations})
}

// ─── 对话详情 ───
func GetConversation(c *gin.Context) {
	userID := c.GetUint("user_id")
	partnerID, _ := strconv.Atoi(c.Param("user_id"))
	threadID, _ := strconv.Atoi(c.DefaultQuery("thread", "0"))

	result, err := service.GetConversation(userID, uint(partnerID), uint(threadID))
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "获取成功", gin.H{
		"partner":  result.Partner,
		"messages": result.Messages,
	})
}

// ─── 未读统计 ───
func GetUnreadCount(c *gin.Context) {
	count, _ := data.CountUnread(c.GetUint("user_id"))
	models.Success(c, "获取成功", gin.H{"count": count})
}

// ─── 全部标记已读 ───
func MarkAllRead(c *gin.Context) {
	data.MarkAllRead(c.GetUint("user_id"))
	models.Success(c, "已标记", nil)
}

// ─── 标记某人的消息已读 ───
func MarkMessagesRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	partnerID := c.Param("user_id")
	uid, _ := strconv.Atoi(partnerID)
	data.MarkReadFrom(uint(uid), userID)
	models.Success(c, "已标记", nil)
}

// ─── 撤回消息 ───
func RecallMessage(c *gin.Context) {
	msgID := middlewares.StrToUint(c.Param("id"))
	if err := service.RecallMessage(msgID, c.GetUint("user_id")); err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}
	models.Success(c, "已撤回", nil)
}

// ─── SSE 消息流 ───
func MessageStream(c *gin.Context) {
	userID := c.GetUint("user_id")

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // 禁用 nginx 缓冲
	c.Writer.Flush()

	ch := sse.DefaultHub.Subscribe(userID)
	defer sse.DefaultHub.Unsubscribe(userID, ch)

	for {
		select {
		case data := <-ch:
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}
