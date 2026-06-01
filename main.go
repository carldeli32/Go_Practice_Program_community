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
	routes.Setup(r)
	r.Run(":8080")
}
