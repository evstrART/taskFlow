package taskFlow

import "time"

type Project struct {
	ProjectID   int       `json:"project_id"`  // Уникальный идентификатор проекта.
	Name        string    `json:"name"`        // Название проекта.
	Description string    `json:"description"` // Описание проекта.
	OwnerID     int       `json:"owner_id"`    // Идентификатор пользователя, создавшего проект.
	StartDate   time.Time `json:"start_date"`  // Дата начала проекта.
	EndDate     time.Time `json:"end_date"`    // Дата завершения проекта.
	Status      string    `json:"status"`      // Статус проекта.
	CreatedAt   time.Time `json:"created_at"`  // Дата создания проекта.
	UpdatedAt   time.Time `json:"updated_at"`  // Дата последнего обновления проекта.
}
