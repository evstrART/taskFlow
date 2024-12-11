package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) checkAdmin(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized: unable to get user ID")
		return
	}
	isAdmin, err := h.services.Admin.CheckAdmin(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if isAdmin {
		c.JSON(http.StatusOK, gin.H{"isAdmin": true})
	} else {
		newErrorResponse(c, http.StatusUnauthorized, "you are not admin")
		return
	}
}
