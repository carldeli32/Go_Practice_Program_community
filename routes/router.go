package routes

import (
	"community/controllers"
	"community/middlewares"
	"community/models"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) { models.Success(c, "pong", nil) })

	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.GET("/users/:id", controllers.GetUserProfile)
		api.GET("/posts", controllers.GetPosts)
		api.GET("/posts/:id", controllers.GetPost)
		api.GET("/posts/:id/comments", controllers.GetComments)
		api.GET("/announcement", controllers.GetAnnouncement)
		api.GET("/categories", controllers.GetCategories)

		auth := api.Group("")
		auth.Use(middlewares.AuthRequired())
		{
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
			auth.GET("/messages/unread-count", controllers.GetUnreadCount)
			auth.PUT("/messages/read-all", controllers.MarkAllRead)
			auth.GET("/messages/:user_id", controllers.GetConversation)
			auth.PUT("/messages/:user_id/read", controllers.MarkMessagesRead)

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
