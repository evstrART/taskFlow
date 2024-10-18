package handler

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sing-in", h.singIn)
		auth.POST("/sing-up", h.singUp)

	}
	api := router.Group("/api")
	{
		project := api.Group("/project")
		{
			project.POST("/", h.CreateProject)
		}
	}

	return router
}
