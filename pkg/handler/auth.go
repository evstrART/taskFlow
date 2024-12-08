package handler

import (
	"github.com/evstrART/taskFlow"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) signUp(ctx *gin.Context) {
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
	token, err := h.services.AutorisationService.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"user_id": id,
		"token":   token,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"` // Имя пользователя.
	Password string `json:"password" binding:"required"` // Пароль (в хэшированном виде).
}

func (h *Handler) signIn(ctx *gin.Context) {
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

func (h *Handler) signInGet(c *gin.Context) {
	c.HTML(http.StatusOK, "sign-in.html", nil) // Отображаем страницу входа
}
func (h *Handler) signUpGet(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "sign-up.html", nil) // Отображаем страницу регистрации
}

func (h *Handler) changePassword(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var input taskFlow.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// Проверяем старый пароль
	isValid, err := h.services.User.CheckOldPassword(userId, input.OldPassword)
	if err != nil || !isValid {
		newErrorResponse(c, http.StatusUnauthorized, "Old password is incorrect")
		return
	}

	err = h.services.AutorisationService.ChangePassword(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
