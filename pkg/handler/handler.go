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
		main.GET("/", h.mainGet) // Главная страница
		main.GET("/projects", h.getProjectsPage)
		main.GET("/projects/:id", h.projectGet)
		main.GET("/profile", h.profileGet)
		main.GET("/projects/:id/tasks/:task_id", h.taskGet)
		main.GET("/forgot-password", h.forgotPasswordGet)
		main.GET("/reset-password", h.resetPasswordGet)
	}
	auth := router.Group("/auth")
	{
		auth.GET("/sign-in", h.signInGet) // Маршрут для страницы входа
		auth.POST("/sign-in", h.signIn)
		auth.GET("/sign-up", h.signUpGet)
		auth.POST("/sign-up", h.signUp)
		auth.POST("/forgot-password", h.forgotPassword)
		auth.POST("/reset-password", h.changeForgotPassword)

	}
	api := router.Group("/api", h.userIdentity)
	{
		project := api.Group("/projects")
		{
			project.POST("/", h.createProject)
			project.GET("/", h.getAllProjects)
			project.GET("/:id", h.getProjectById)
			project.PUT("/:id", h.updateProject)
			project.DELETE("/:id", h.deleteProject)

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
			users.GET("/:user_id", h.getUserProfile)
			users.PUT("/:user_id", h.updateUser)
			users.DELETE("/:user_id", h.deleteUser)
			users.PUT(":user_id/change-password", h.changePassword)
		}
	}

	return router
}
