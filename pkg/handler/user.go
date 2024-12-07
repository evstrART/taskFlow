package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getUserProfile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}
	user, err := h.services.User.GetUser(userId)
	if err != nil {
		if err.Error() == "user not found" {
			newErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) getUsers(c *gin.Context) {
	users, err := h.services.User.GetUsers()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}
