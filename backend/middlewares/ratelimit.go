package middlewares

import (
	"net/http"
	"sync"
	"time"

	"community/backend/models"

	"github.com/gin-gonic/gin"
)

type entry struct {
	count   int
	resetAt time.Time
}

type rateLimiter struct {
	mu      sync.Mutex
	entries map[string]*entry
}

var limiter = &rateLimiter{entries: make(map[string]*entry)}

func init() {
	go limiter.cleanup(10 * time.Minute)
}

func (rl *rateLimiter) cleanup(every time.Duration) {
	for range time.Tick(every) {
		rl.mu.Lock()
		now := time.Now()
		for key, e := range rl.entries {
			if now.After(e.resetAt) {
				delete(rl.entries, key)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) allow(key string, limit int, window time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	e, exists := rl.entries[key]
	if !exists || now.After(e.resetAt) {
		rl.entries[key] = &entry{count: 1, resetAt: now.Add(window)}
		return true
	}

	if e.count >= limit {
		return false
	}

	e.count++
	return true
}

// RateLimit 返回限流中间件
// limit: 窗口内允许次数, window: 时间窗口
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.allow(ip, limit, window) {
			models.Error(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}
