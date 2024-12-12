package handler

import (
	"fmt"
	"github.com/evstrART/taskFlow"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"strconv"
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
	userId, err := strconv.Atoi(c.Param("user_id"))
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

	err = h.services.AutorisationService.ChangePassword(userId, input.NewPassword)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
func sendEmail(to string, resetLink string) error {
	m := gomail.NewMessage()

	from := "art.evstr@gmail.com"     // Укажите свой адрес электронной почты
	password := "ckgy yjxy vuki qokz" // Используйте пароль приложения

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Сброс пароля")
	body := fmt.Sprintf(`
        <p>Здравствуйте!</p>
        <p>Чтобы сбросить ваш пароль, перейдите по следующей <a href="%s">ссылке</a>.</p>
        <p>Если вы не инициировали этот запрос, просто проигнорируйте это письмо.</p>
        <p>С уважением,<br>Команда Task Flow.</p>
    `, resetLink)

	m.SetBody("text/html", body) // Устанавливаем тело сообщения в формате HTML

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	// Отправка почты
	err := d.DialAndSend(m)
	if err != nil {
		log.Printf("Ошибка при отправке почты: %v", err) // Логируем ошибку
		return err
	}

	return nil
}

func (h *Handler) forgotPassword(c *gin.Context) {
	var input taskFlow.ResetPasswordInput

	// Привязываем входящие данные к структуре
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Проверяем, существует ли пользователь
	exist, err := h.services.AutorisationService.UserExistsForReset(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке существования пользователя")
		return
	}

	if !exist {
		newErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
		return
	}

	// Генерируем токен для сброса пароля
	user, token, err := h.services.AutorisationService.GenerateTokenForReset(input.Username, input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Ошибка при генерации токена")
		return
	}

	// Создаем ссылку для сброса пароля с токеном
	resetLink := fmt.Sprintf("http://localhost:8080/reset-password?token=%s&username=%s", token, user.Username)

	emailErr := sendEmail(input.Email, resetLink)
	if emailErr != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Не удалось отправить электронное письмо")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ссылка для сброса пароля отправлена на вашу электронную почту."})
}

func (h *Handler) changeForgotPassword(c *gin.Context) {
	var input taskFlow.ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userId, err := h.services.AutorisationService.ParseToken(input.Token)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	err = h.services.AutorisationService.ChangePassword(userId, input.NewPassword)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})

}
