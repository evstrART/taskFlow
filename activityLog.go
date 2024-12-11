package taskFlow

import "time"

type ActivityLog struct {
	LogID         int       `db:"log_id" json:"log_id"`                 // Уникальный идентификатор записи.
	UserID        int       `db:"user_id" json:"user_id"`               // Идентификатор пользователя, выполнившего действие.
	Action        string    `db:"action" json:"action"`                 // Описание действия.
	Timestamp     time.Time `db:"timestamp" json:"timestamp"`           // Время выполнения действия.
	RelatedEntity string    `db:"related_entity" json:"related_entity"` // Тип сущности.
	EntityID      int       `db:"entity_id" json:"entity_id"`           // Идентификатор сущности, связанной с действием.
}
