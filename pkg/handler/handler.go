package handler

import (
	"github.com/evstrART/taskFlow/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.LoadHTMLGlob("web/templates/*html")
	router.Static("/static", "./web/static")
	main := router.Group("/")
	{
		main.GET("/", h.mainGet)
		main.GET("/projects", h.getProjectsPage)
		main.GET("/projects/:id", h.projectGet)
		main.GET("/profile", h.profileGet)
		main.GET("/projects/:id/tasks/:task_id", h.taskGet)
		main.GET("/forgot-password", h.forgotPasswordGet)
		main.GET("/reset-password", h.resetPasswordGet)
		admin := router.Group("/admin")
		{
			admin.GET("/projects", h.getAdminProjectsPage)
			admin.GET("/projects/:id", h.adminProjectGet)
			admin.GET("/profile", h.adminProfileGet)
			admin.GET("/projects/:id/tasks/:task_id", h.adminTaskGet)
			admin.GET("/reports", h.adminReportsGet)
			admin.GET("/users", h.adminUsersGet)
		}
	}
	auth := router.Group("/auth")
	{
		auth.GET("/sign-in", h.signInGet)
		auth.POST("/sign-in", h.signIn)
		auth.GET("/sign-up", h.signUpGet)
		auth.POST("/sign-up", h.signUp)
		auth.POST("/forgot-password", h.forgotPassword)
		auth.POST("/reset-password", h.changeForgotPassword)
		auth.POST("/check-admin", h.checkAdmin)

	}
	api := router.Group("/api", h.userIdentity)
	{
		project := api.Group("/projects")
		{
			project.POST("/", h.createProject)
			project.GET("/", h.getAllProjects)
			project.GET("/user", h.getAllProjectsForUser)
			project.GET("/:id", h.getProjectById)
			project.PUT("/:id", h.updateProject)
			project.DELETE("/:id", h.deleteProject)
			project.POST("/:id/complete", h.completeProject)

			members := project.Group(":id/members")
			members.GET("/", h.getMembers)
			members.POST("/", h.addMembers)
			members.DELETE("/:member_id", h.deleteMembers)

			task := project.Group(":id/tasks")
			task.POST("/", h.createTask)
			task.GET("/", h.getAllTasks)
			task.GET("/:task_id", h.getTaskById)
			task.PUT("/:task_id", h.updateTask)
			task.DELETE("/:task_id", h.deleteTask)
			task.PUT("/:task_id/complete", h.completeTask)

			tag := task.Group(":task_id/tags")
			tag.GET("/", h.getTags)
			tag.POST("/", h.createTag)
			tag.DELETE("/:tag_id", h.deleteTag)
			tag.PUT("/:tag_id", h.updateTag)

			comment := task.Group(":task_id/comments")
			comment.POST("/", h.addComment)
			comment.GET("/", h.getAllComments)
			comment.GET("/:comment_id", h.getCommentById)
			comment.DELETE("/:comment_id", h.deleteComment)
			comment.PUT("/:comment_id", h.updateComment)
		}
		tasks := api.Group("/tasks")
		{
			tasks.GET("/", h.getAllTasksForUser)
		}
		comments := api.Group("/comments")
		{
			comments.GET("/", h.getAllCommentsForUser)
			comments.DELETE("/:comment_id", h.deleteComment)
			comments.PUT("/:comment_id", h.updateComment)
		}
		tags := api.Group("/tags")
		{
			tags.GET("/", h.getAllTags)
		}
		users := api.Group("/users")
		{
			users.GET("/", h.getUsers)
			users.GET("/:user_id", h.getUserById)
			users.GET("/profile", h.getUserProfile)
			users.PUT("/:user_id", h.updateUser)
			users.DELETE("/:user_id", h.deleteUser)
			users.PUT(":user_id/change-password", h.changePassword)
		}
	}
	admin := router.Group("/admin", h.userIdentity, h.userIdentity, h.checkAdmin)
	{
		admin.GET("/", h.getAllProjects)
		admin.GET("/reports/:format", h.getReport)
		admin.GET("/export", h.exportJSON)
		admin.POST("/import", h.importJSON)
	}
	return router
}
