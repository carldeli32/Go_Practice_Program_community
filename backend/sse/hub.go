package sse

import (
	"encoding/json"
	"sync"
)

// Event SSE 推送事件
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Hub SSE 连接中心
type Hub struct {
	mu      sync.RWMutex
	clients map[uint][]chan []byte
}

// DefaultHub 全局实例（main.go 中初始化）
var DefaultHub *Hub

// NewHub 创建 Hub
func NewHub() *Hub {
	return &Hub{clients: make(map[uint][]chan []byte)}
}

// Subscribe 为用户建立 SSE 连接，返回消息 channel
func (h *Hub) Subscribe(userID uint) chan []byte {
	h.mu.Lock()
	defer h.mu.Unlock()
	ch := make(chan []byte, 32)
	h.clients[userID] = append(h.clients[userID], ch)
	return ch
}

// Unsubscribe 断开用户的某条 SSE 连接
func (h *Hub) Unsubscribe(userID uint, ch chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	chs := h.clients[userID]
	for i, c := range chs {
		if c == ch {
			h.clients[userID] = append(chs[:i], chs[i+1:]...)
			close(c)
			return
		}
	}
}

// Publish 推送事件给指定用户（非阻塞，32 缓冲满了就丢弃）
func (h *Hub) Publish(userID uint, eventType string, data interface{}) {
	payload, err := json.Marshal(Event{Type: eventType, Data: data})
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, ch := range h.clients[userID] {
		select {
		case ch <- payload:
		default:
			// 缓冲区满，丢弃
		}
	}
}
