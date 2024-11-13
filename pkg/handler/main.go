package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) mainGet(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
