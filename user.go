package taskFlow

import (
	"database/sql"
	"time"
)

type User struct {
	UserID    int          `json:"user_id" db:"user_id"`        // Уникальный идентификатор пользователя.
	Username  string       `json:"username" binding:"required"` // Имя пользователя.
	Password  string       `json:"password" binding:"required"` // Пароль (в хэшированном виде).
	Email     string       `json:"email" binding:"required"`    // Электронная почта пользователя.
	Role      string       `json:"role"`                        // Роль пользователя (например, “admin”, “user”).
	CreatedAt time.Time    `json:"created_at" db:"created_at"`  // Дата и время создания пользователя.
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at"`  // Дата и время последнего обновления данных пользователя.
}
type UpdateUserInput struct {
	Username *string `db:"username"` // Добавляем теги для связывания с колонками
	Email    *string `db:"email"`
}
type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
type ResetPasswordInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}
