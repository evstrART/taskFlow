package taskFlow

import "time"

type Comment struct {
	CommentID int       `json:"comment_id"` // Уникальный идентификатор комментария.
	TaskID    int       `json:"task_id"`    // Идентификатор задачи, к которой относится комментарий.
	UserID    int       `json:"user_id"`    // Идентификатор пользователя, оставившего комментарий.
	Content   string    `json:"content"`    // Текст комментария.
	CreatedAt time.Time `json:"created_at"` // Дата создания комментария.
}
