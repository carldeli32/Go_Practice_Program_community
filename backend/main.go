package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"community/backend/config"
	"community/backend/data"
	"community/backend/routes"
	"community/backend/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	config.InitDB()

	// 初始化文件存储（切云存储只需改这里）
	storage.Store = data.NewLocalStorage(
		config.UploadImgDir(),
		config.UploadFileDir(),
		config.ServeImgPrefix(),
		config.ServeFilePrefix(),
	)

	r := gin.Default()
	r.SetTrustedProxies(nil)
	routes.Setup(r)

	srv := &http.Server{
		Addr:    config.ServerPort(),
		Handler: r,
	}

	// 启动
	go func() {
		fmt.Printf("🚀 服务启动，监听 %s\n", config.ServerPort())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动失败: %v", err)
		}
	}()

	// 等待信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("正在关闭...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("强制关闭: %v", err)
	}

	if sqlDB, err := config.DB.DB(); err == nil {
		sqlDB.Close()
	}

	fmt.Println("已关闭")
}
