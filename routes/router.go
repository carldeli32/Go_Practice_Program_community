package routes

import (
	"community/controllers"
	"community/middlewares"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	//	检查系统状态是否正常
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	//	路由组
	api := r.Group("/api")
	{
		// 公开接口（无需登录）
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		// 帖子（列表和详情公开，方便游客浏览）
		api.GET("/posts", controllers.GetPosts)
		api.GET("/posts/:id", controllers.GetPost)

		// 评论列表公开
		api.GET("/posts/:id/comments", controllers.GetComments)

		// 需要登录的接口
		auth := api.Group("")
		auth.Use(middlewares.AuthRequired())
		{
			// 帖子
			auth.POST("/posts", controllers.CreatePost)
			auth.PUT("/posts/:id", controllers.UpdatePost)
			auth.DELETE("/posts/:id", controllers.DeletePost)

			// 评论
			auth.POST("/posts/:id/comments", controllers.CreateComment)
		}
	}
}
