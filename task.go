package taskFlow

import "time"

type Task struct {
	TaskID      int       `json:"task_id"`     // Уникальный идентификатор задачи.
	ProjectID   int       `json:"project_id"`  // Идентификатор проекта, к которому относится задача.
	Title       string    `json:"title"`       // Заголовок задачи.
	Description string    `json:"description"` // Подробное описание задачи.
	AssignedTo  int       `json:"assigned_to"` // Идентификатор пользователя, которому назначена задача.
	Status      string    `json:"status"`      // Статус задачи.
	Priority    string    `json:"priority"`    // Приоритет задачи.
	DueDate     time.Time `json:"due_date"`    // Дата завершения задачи.
	CreatedAt   time.Time `json:"created_at"`  // Дата создания задачи.
	UpdatedAt   time.Time `json:"updated_at"`  // Дата последнего обновления задачи.
}
