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
	r.StaticFile("/manifest.json", "../frontend/dist/manifest.json")
	r.StaticFile("/sw.js", "../frontend/dist/sw.js")
	r.Static("/assets", "../frontend/dist/assets")
	r.NoRoute(func(c *gin.Context) {
		c.File("../frontend/dist/index.html")
	})

	api := r.Group("/api")
	{
		// 公开接口
		api.POST("/login", middlewares.RateLimit(5, time.Minute), controllers.Login)
		api.GET("/users", controllers.SearchUsers)
		api.GET("/users/:id", controllers.GetUserProfile)
		api.GET("/posts", controllers.GetPosts)
		api.GET("/posts/:id", controllers.GetPost)
		api.GET("/posts/:id/comments", middlewares.OptionalAuth(), controllers.GetComments)
		api.GET("/announcement", controllers.GetAnnouncement)
		api.GET("/categories", controllers.GetCategories)

		auth := api.Group("")
		auth.Use(middlewares.AuthRequired())
		{
			// ── 上传 ──
			auth.POST("/upload/image", controllers.UploadImage)
			auth.POST("/upload/file", controllers.UploadFile)

			// ── 帖子：创建 = 角色权限，编辑/删除 = 资源权限（group 只挂一次）──
			auth.POST("/posts", middlewares.RequirePerm("post.create"), controllers.CreatePost)
			postRes := auth.Group("/posts/:id")
			postRes.Use(middlewares.RequireResource(middlewares.PostLoader, "post.manage_any", "post.manage_category"))
			{
				postRes.PUT("", controllers.UpdatePost)
				postRes.DELETE("", controllers.DeletePost)
			}

			// ── 评论：创建 = 角色权限，编辑/删除 = 资源权限（group 只挂一次）──
			auth.POST("/posts/:id/comments", middlewares.RequirePerm("comment.create"), controllers.CreateComment)
			commentRes := auth.Group("/comments/:id")
			commentRes.Use(middlewares.RequireResource(middlewares.CommentLoader, "comment.manage_any", "comment.manage_category"))
			{
				commentRes.PUT("", controllers.UpdateComment)
				commentRes.DELETE("", controllers.DeleteComment)
			}

			// ── 私信 ──
			auth.POST("/threads", controllers.CreateThread)
			auth.GET("/threads", controllers.GetThreads)
			auth.DELETE("/threads/:id", controllers.DeleteThread)

			auth.POST("/messages", controllers.SendMessage)
			auth.GET("/messages", controllers.GetConversations)
			auth.GET("/messages/unread-count", controllers.GetUnreadCount)
			auth.PUT("/messages/read-all", controllers.MarkAllRead)
			auth.GET("/messages/:user_id", controllers.GetConversation)
			auth.PUT("/messages/:user_id/read", controllers.MarkMessagesRead)
			auth.PUT("/messages/recall/:id", controllers.RecallMessage)

			// SSE 流（TokenFromQuery → AuthRequired 顺序不可逆）
			api.GET("/messages/stream", middlewares.TokenFromQuery(), middlewares.AuthRequired(), controllers.MessageStream)

			// ── 关注 ──
			auth.POST("/follow", controllers.FollowUser)
			auth.DELETE("/follow/:user_id", controllers.UnfollowUser)
			auth.GET("/following", controllers.GetMyFollowing)
			auth.GET("/followers", controllers.GetMyFollowers)

			// ── 管理员面板 ──
			admin := auth.Group("/admin")
			{
				// admin / super_admin 均可
				adminLevel := admin.Group("")
				adminLevel.Use(middlewares.RequirePerm("ban.any"))
				{
					adminLevel.PUT("/users/:id/ban", controllers.BanUser)
					adminLevel.PUT("/users/:id/unban", controllers.UnbanUser)
					adminLevel.GET("/users", controllers.AdminListUsers)
					adminLevel.POST("/announcement", controllers.SetAnnouncement)
					adminLevel.DELETE("/announcement", controllers.DeleteAnnouncement)
				}

				// 仅 super_admin
				admin.POST("/users", middlewares.RequirePerm("user.create"), controllers.AdminCreateUser)
				admin.DELETE("/users/:id", middlewares.RequirePerm("user.delete"), controllers.AdminDeleteUser)
				admin.PUT("/users/:id/roles", middlewares.RequirePerm("role.assign"), controllers.AdminAssignRoles)
				admin.POST("/categories", middlewares.RequirePerm("category.manage"), controllers.CreateCategory)
				admin.DELETE("/categories/:id", middlewares.RequirePerm("category.manage"), controllers.DeleteCategory)
			}
		}
	}
}
