package handler

import (
	"github.com/evstrART/taskFlow"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) getUserById(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
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

func (h *Handler) updateUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var input taskFlow.UpdateUserInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.User.UpdateUser(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

func (h *Handler) deleteUser(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}
	err = h.services.User.DeleteUser(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
