package routes

import (
	"community/controllers"
	"community/middlewares"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })

	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.GET("/users/:id", controllers.GetUserProfile)
		api.GET("/posts", controllers.GetPosts)
		api.GET("/posts/:id", controllers.GetPost)
		api.GET("/posts/:id/comments", controllers.GetComments)
		api.GET("/announcement", controllers.GetAnnouncement)

		auth := api.Group("")
		auth.Use(middlewares.AuthRequired())
		{
			auth.POST("/posts", controllers.CreatePost)
			auth.PUT("/posts/:id", controllers.UpdatePost)
			auth.DELETE("/posts/:id", controllers.DeletePost)

			auth.POST("/posts/:id/comments", controllers.CreateComment)
			auth.PUT("/comments/:id", controllers.UpdateComment)
			auth.DELETE("/comments/:id", controllers.DeleteComment)

			auth.POST("/threads", controllers.CreateThread)
			auth.GET("/threads", controllers.GetThreads)
			auth.DELETE("/threads/:id", controllers.DeleteThread)

			auth.POST("/messages", controllers.SendMessage)
			auth.GET("/messages", controllers.GetConversations)
			auth.GET("/messages/unread-count", controllers.GetUnreadCount)
			auth.PUT("/messages/read-all", controllers.MarkAllRead)
			auth.GET("/messages/:user_id", controllers.GetConversation)
			auth.PUT("/messages/:user_id/read", controllers.MarkMessagesRead)

			auth.POST("/follow", controllers.FollowUser)
			auth.DELETE("/follow/:user_id", controllers.UnfollowUser)
			auth.GET("/following", controllers.GetMyFollowing)
			auth.GET("/followers", controllers.GetMyFollowers)

			admin := auth.Group("/admin")
			admin.Use(middlewares.AdminRequired())
			{
				admin.GET("/users", controllers.AdminListUsers)
				admin.POST("/users", controllers.AdminCreateUser)
				admin.PUT("/users/:id/ban", controllers.BanUser)
				admin.PUT("/users/:id/unban", controllers.UnbanUser)
				admin.POST("/announcement", controllers.SetAnnouncement)
				admin.DELETE("/announcement", controllers.DeleteAnnouncement)
			}
		}
	}
}
