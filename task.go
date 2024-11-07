package taskFlow

import (
	"database/sql"
	"time"
)

type Task struct {
	TaskID      int          `json:"task_id" db:"task_id"`                // Уникальный идентификатор задачи.
	ProjectID   int          `json:"project_id" db:"project_id"`          // Идентификатор проекта, к которому относится задача.
	Title       string       `json:"title" db:"title" binding:"required"` // Заголовок задачи.
	Description string       `json:"description"`                         // Подробное описание задачи.
	AssignedTo  int          `json:"assigned_to" db:"assigned_to"`        // Идентификатор пользователя, которому назначена задача.
	Status      string       `json:"status" binding:"required"`           // Статус задачи.
	Priority    string       `json:"priority" binding:"required"`         // Приоритет задачи.
	DueDate     time.Time    `json:"due_date" db:"due_date"`              // Дата завершения задачи.
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`          // Дата создания задачи.
	UpdatedAt   sql.NullTime `json:"updated_at" db:"updated_at"`          // Дата последнего обновления задачи.
}

type UpdateTaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	AssignedTo  *int    `json:"assigned_to"`
	Status      *string `json:"status"`
	Priority    *string `json:"priority"`
	DueDate     *string `json:"due_date"`
}
