package handler

import (
	"github.com/evstrART/taskFlow"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) singIn(ctx *gin.Context) {
	var input taskFlow.User
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.AutorisationService.CreateUser(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"user_id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"` // Имя пользователя.
	Password string `json:"password" binding:"required"` // Пароль (в хэшированном виде).
}

func (h *Handler) singUp(ctx *gin.Context) {
	var input signInInput
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.services.AutorisationService.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
