package main

import (
	"community/config"
	"community/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	config.InitDB()

	r := gin.Default()
	r.SetTrustedProxies(nil) // 本地开发，禁用代理信任检查
	routes.Setup(r)
	r.Run(":8080")
}
