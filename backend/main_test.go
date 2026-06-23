package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"
)

// TestGracefulShutdown 验证优雅关闭流程：启动 → 可达 → 关闭 → 拒绝
func TestGracefulShutdown(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()

	srv := &http.Server{Handler: mux}
	go srv.Serve(l)
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://" + addr + "/ping")
	if err != nil {
		t.Fatalf("服务未就绪: %v", err)
	}
	resp.Body.Close()
	t.Logf("✅ /ping 返回 %d", resp.StatusCode)

	shutdownDone := make(chan error, 1)
	go func() {
		shutdownDone <- srv.Shutdown(nil)
	}()

	select {
	case err := <-shutdownDone:
		if err != nil {
			t.Errorf("Shutdown 失败: %v", err)
		} else {
			t.Log("✅ 优雅关闭完成")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Shutdown 超时")
	}

	_, err = http.Get("http://" + addr + "/ping")
	if err == nil {
		t.Error("关闭后服务仍可达")
	} else {
		t.Log("✅ 关闭后连接被拒绝")
	}
}

// TestSignalNotify 验证 SIGINT/SIGTERM 信号监听注册
func TestSignalNotify(t *testing.T) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	p.Signal(syscall.SIGTERM)

	select {
	case sig := <-ch:
		t.Logf("✅ 收到信号: %v", sig)
	case <-time.After(2 * time.Second):
		t.Error("未收到 SIGTERM 信号")
	}

	signal.Stop(ch)
	close(ch)
}

// TestShutdownTimeout 验证：即使 handler 死循环，超时也能强制退出
func TestShutdownTimeout(t *testing.T) {
	mux := http.NewServeMux()

	// 模拟死循环 handler
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		t.Log("慢请求开始（将阻塞 10 秒）...")
		time.Sleep(10 * time.Second)
		w.Write([]byte("done"))
	})

	l, _ := net.Listen("tcp", "127.0.0.1:0")
	addr := l.Addr().String()

	srv := &http.Server{Handler: mux}
	go srv.Serve(l)
	time.Sleep(50 * time.Millisecond)

	// 发起慢请求
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.Get("http://" + addr + "/slow")
	}()
	time.Sleep(50 * time.Millisecond)

	// 2 秒超时关闭
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	elapsed := time.Since(start)

	if err == nil {
		t.Error("期望超时，但 Shutdown 成功了（不可能，handler 还在 sleep 10s）")
	} else {
		t.Logf("Shutdown 返回: %v", err)
	}

	if elapsed > 5*time.Second {
		t.Errorf("关闭耗时 %v，超过 5 秒，超时失效", elapsed)
	} else {
		t.Logf("✅ 超时生效！关闭耗时 %v（限制 2s，handler 本应阻塞 10s）", elapsed)
	}

	wg.Wait()
}
