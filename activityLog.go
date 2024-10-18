package taskFlow

import "time"

type ActivityLog struct {
	LogID         int       `json:"log_id"`         // Уникальный идентификатор записи.
	UserID        int       `json:"user_id"`        // Идентификатор пользователя, выполнившего действие.
	Action        string    `json:"action"`         // Описание действия.
	Timestamp     time.Time `json:"timestamp"`      // Время выполнения действия.
	RelatedEntity string    `json:"related_entity"` // Тип сущности (например, “task”, “project”).
	EntityID      int       `json:"entity_id"`      // Идентификатор сущности, связанной с действием.
}
