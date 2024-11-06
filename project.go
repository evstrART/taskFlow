package taskFlow

import (
	"database/sql"
	"time"
)

type Project struct {
	ProjectID   int          `json:"project_id" db:"project_id"` // Уникальный идентификатор проекта.
	Name        string       `json:"name" binding:"required"`    // Название проекта.
	Description string       `json:"description"`                // Описание проекта.
	OwnerID     int          `json:"owner_id" db:"owner_id"`     // Идентификатор пользователя, создавшего проект.
	StartDate   time.Time    `json:"start_date" db:"start_date"` // Дата начала проекта.
	EndDate     time.Time    `json:"end_date" db:"end_date"`     // Дата завершения проекта.
	Status      string       `json:"status" binding:"required"`  // Статус проекта.
	CreatedAt   time.Time    `json:"created_at" db:"created_at"` // Дата создания проекта.
	UpdatedAt   sql.NullTime `json:"updated_at" db:"updated_at"` // Дата последнего обновления проекта.
}
