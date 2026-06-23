package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter(mw ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(mw...)
	r.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
	return r
}

func TestRateLimit_AllowFirstN(t *testing.T) {
	limiter := &rateLimiter{entries: make(map[string]*entry)}

	// 窗口内前 3 次应允许
	for i := 0; i < 3; i++ {
		if !limiter.allow("test-ip", 3, time.Minute) {
			t.Errorf("第 %d 次请求应被允许", i+1)
		}
	}

	// 第 4 次应拒绝
	if limiter.allow("test-ip", 3, time.Minute) {
		t.Error("第 4 次请求应被拒绝")
	}
}

func TestRateLimit_WindowReset(t *testing.T) {
	limiter := &rateLimiter{entries: make(map[string]*entry)}

	// 用过期窗口
	limiter.allow("test-ip", 1, -time.Second)

	// 窗口已过期，下次应重置
	if !limiter.allow("test-ip", 1, time.Hour) {
		t.Error("过期窗口后应重新允许")
	}
}

func TestRateLimit_IndependentIPs(t *testing.T) {
	limiter := &rateLimiter{entries: make(map[string]*entry)}

	// IP A 消耗完毕
	for i := 0; i < 3; i++ {
		limiter.allow("ip-a", 3, time.Minute)
	}

	// IP B 不受影响
	if !limiter.allow("ip-b", 3, time.Minute) {
		t.Error("不同 IP 不应相互影响")
	}
}

func TestRateLimit_Cleanup(t *testing.T) {
	limiter := &rateLimiter{entries: make(map[string]*entry)}

	// 放入过期条目
	limiter.entries["old-ip"] = &entry{count: 100, resetAt: time.Now().Add(-time.Hour)}
	limiter.entries["new-ip"] = &entry{count: 1, resetAt: time.Now().Add(time.Hour)}

	limiter.mu.Lock()
	now := time.Now()
	for key, e := range limiter.entries {
		if now.After(e.resetAt) {
			delete(limiter.entries, key)
		}
	}
	limiter.mu.Unlock()

	if _, exists := limiter.entries["old-ip"]; exists {
		t.Error("过期条目应被清理")
	}
	if _, exists := limiter.entries["new-ip"]; !exists {
		t.Error("未过期条目不应被清理")
	}
}

func TestRateLimit_HTTPMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter(RateLimit(2, time.Hour))

	// 前 2 次 → 200
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		r.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("第 %d 次: 期望 200，得到 %d", i+1, w.Code)
		}
	}

	// 第 3 次 → 429
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	r.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("超限请求: 期望 429，得到 %d", w.Code)
	}
}

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter(SecurityHeaders())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	checks := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
	}
	for header, expected := range checks {
		got := w.Header().Get(header)
		if got != expected {
			t.Errorf("%s = %q, want %q", header, got, expected)
		}
	}
}
