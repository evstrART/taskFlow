package taskFlow

type TaskTag struct {
	TaskID int `json:"task_id"` // Идентификатор задачи.
	TagID  int `json:"tag_id"`  // Идентификатор метки.
}
