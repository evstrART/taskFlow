package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) mainGet(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"breadcrumb": "Home",
	})
}

func (h *Handler) projectGet(c *gin.Context) {
	c.HTML(http.StatusOK, "project.html", gin.H{})
}
func (h *Handler) getProjectsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "projects.html", gin.H{
		"breadcrumb": "Projects",
	})
}
func (h *Handler) profileGet(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", gin.H{})
}
func (h *Handler) taskGet(c *gin.Context) {
	c.HTML(http.StatusOK, "task.html", gin.H{})
}
func (h *Handler) forgotPasswordGet(c *gin.Context) {
	c.HTML(http.StatusOK, "forgotPass.html", gin.H{})
}
func (h *Handler) resetPasswordGet(c *gin.Context) {
	c.HTML(http.StatusOK, "resetPass.html", gin.H{})
}

func (h *Handler) getAdminProjectsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-projects.html", gin.H{
		"breadcrumb": "Projects",
	})
}
func (h *Handler) adminProjectGet(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-project.html", gin.H{})
}
func (h *Handler) adminTaskGet(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-task.html", gin.H{})
}
func (h *Handler) adminProfileGet(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-profile.html", gin.H{})
}
