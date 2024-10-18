package taskFlow

import "time"

type User struct {
	UserID    int       `json:"user_id"`    // Уникальный идентификатор пользователя.
	Username  string    `json:"username"`   // Имя пользователя.
	Password  string    `json:"password"`   // Пароль (в хэшированном виде).
	Email     string    `json:"email"`      // Электронная почта пользователя.
	Role      string    `json:"role"`       // Роль пользователя (например, “admin”, “user”).
	CreatedAt time.Time `json:"created_at"` // Дата и время создания пользователя.
	UpdatedAt time.Time `json:"updated_at"` // Дата и время последнего обновления данных пользователя.
}
