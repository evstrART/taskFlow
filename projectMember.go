package taskFlow

import "time"

type ProjectMember struct {
	MemberID  int       `json:"member_id" db:"member_id"'` // Уникальный идентификатор записи.
	ProjectID int       `json:"project_id"`                // Идентификатор проекта.
	UserID    int       `json:"user_id"`                   // Идентификатор пользователя.
	Role      string    `json:"role"`                      // Роль пользователя в проекте.
	CreatedAt time.Time `json:"created_at"`                // Дата добавления пользователя в проект.
}
