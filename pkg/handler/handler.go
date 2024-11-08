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

	auth := router.Group("/auth")
	{
		auth.POST("/sing-in", h.singIn)
		auth.POST("/sing-up", h.singUp)

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
			project.POST("/:id", h.addMembers)

			task := project.Group(":id/tasks")
			task.POST("/", h.createTask)
			task.GET("/", h.getAllTasks)
			task.GET("/:task_id", h.getTaskById)
			task.PUT("/:task_id", h.updateTask)
			task.DELETE("/:task_id", h.deleteTask)

			comment := task.Group(":task_id/comments")
			comment.POST("/", h.addCooment)
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
		}
	}

	return router
}
