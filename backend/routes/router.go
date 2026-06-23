package routes

import (
	"time"

	"community/backend/controllers"
	"community/backend/middlewares"
	"community/backend/models"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.Use(middlewares.SecurityHeaders())
	r.GET("/ping", func(c *gin.Context) { models.Success(c, "pong", nil) })

	// 静态文件服务：前端 + 上传文件
	r.Static("/uploads/images", "../uploads/images")
	r.Static("/uploads/files", "../uploads/files")
	r.StaticFile("/", "../frontend/dist/index.html")
	r.Static("/assets", "../frontend/dist/assets")
	r.NoRoute(func(c *gin.Context) {
		c.File("../frontend/dist/index.html")
	})

	api := r.Group("/api")
	{
		api.POST("/login", middlewares.RateLimit(5, time.Minute), controllers.Login)
		api.GET("/users", controllers.SearchUsers)
		api.GET("/users/:id", controllers.GetUserProfile)
		api.GET("/posts", controllers.GetPosts)
		api.GET("/posts/:id", controllers.GetPost)
		api.GET("/posts/:id/comments", controllers.GetComments)
		api.GET("/announcement", controllers.GetAnnouncement)
		api.GET("/categories", controllers.GetCategories)

		auth := api.Group("")
		auth.Use(middlewares.AuthRequired())
		{
			// 上传
			auth.POST("/upload/image", controllers.UploadImage)
			auth.POST("/upload/file", controllers.UploadFile)

			// 发帖 / 评论（所有登录用户）
			auth.POST("/posts", controllers.CreatePost)
			auth.PUT("/posts/:id", controllers.UpdatePost)
			auth.DELETE("/posts/:id", controllers.DeletePost)

			auth.POST("/posts/:id/comments", controllers.CreateComment)
			auth.PUT("/comments/:id", controllers.UpdateComment)
			auth.DELETE("/comments/:id", controllers.DeleteComment)

			// 私信
			auth.POST("/threads", controllers.CreateThread)
			auth.GET("/threads", controllers.GetThreads)
			auth.DELETE("/threads/:id", controllers.DeleteThread)

			auth.POST("/messages", controllers.SendMessage)
			auth.GET("/messages", controllers.GetConversations)
			auth.GET("/messages/stream", middlewares.TokenFromQuery(), controllers.MessageStream)
			auth.GET("/messages/unread-count", controllers.GetUnreadCount)
			auth.PUT("/messages/read-all", controllers.MarkAllRead)
			auth.GET("/messages/:user_id", controllers.GetConversation)
			auth.PUT("/messages/:user_id/read", controllers.MarkMessagesRead)
			auth.PUT("/messages/:id/recall", controllers.RecallMessage)

			// 关注
			auth.POST("/follow", controllers.FollowUser)
			auth.DELETE("/follow/:user_id", controllers.UnfollowUser)
			auth.GET("/following", controllers.GetMyFollowing)
			auth.GET("/followers", controllers.GetMyFollowers)

			// 管理员面板
			admin := auth.Group("/admin")
			{
				// 管理员级别（admin + super_admin）
				adminLevel := admin.Group("")
				adminLevel.Use(middlewares.RequireAdmin())
				{
					adminLevel.PUT("/users/:id/ban", controllers.BanUser)
					adminLevel.PUT("/users/:id/unban", controllers.UnbanUser)
					adminLevel.GET("/users", controllers.AdminListUsers)
					adminLevel.POST("/announcement", controllers.SetAnnouncement)
					adminLevel.DELETE("/announcement", controllers.DeleteAnnouncement)
				}

				// 仅 super_admin：创建用户
				admin.POST("/users", middlewares.RequirePerm("user.create"), controllers.AdminCreateUser)

				// 仅 super_admin：删除用户
				admin.DELETE("/users/:id", middlewares.RequirePerm("user.delete"), controllers.AdminDeleteUser)

				// 仅 super_admin：角色管理
				admin.PUT("/users/:id/roles", middlewares.RequirePerm("role.assign"), controllers.AdminAssignRoles)

				// 仅 super_admin：分类管理
				admin.POST("/categories", middlewares.RequirePerm("category.manage"), controllers.CreateCategory)
				admin.DELETE("/categories/:id", middlewares.RequirePerm("category.manage"), controllers.DeleteCategory)
			}
		}
	}
}
