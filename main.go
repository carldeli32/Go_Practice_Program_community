package main

import (
	"community/config"
	"community/models"
	"community/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// ========== 加载配置 ==========
	config.Init()

	// ========== 初始化数据库 ==========
	config.InitDB()

	// ========== 自动迁移（根据模型自动建表）==========
	// 注意：生产环境应该用独立的 migration 脚本，这里为了学习方便直接用 AutoMigrate
	config.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	// ========== 创建 Gin 引擎 ==========
	r := gin.Default()

	// ========== 注册路由 ==========
	routes.Setup(r)

	// ========== 启动服务器 ==========
	r.Run(":8080")
}
