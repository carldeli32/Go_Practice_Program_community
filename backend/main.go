package main

import (
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
		"../uploads/images",
		"../uploads/files",
		"/uploads/images",
		"/uploads/files",
	)

	r := gin.Default()
	r.SetTrustedProxies(nil)
	routes.Setup(r)
	r.Run(":8080")
}
