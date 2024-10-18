package taskFlow

type Tag struct {
	TagID int    `json:"tag_id"` // Уникальный идентификатор метки.
	Name  string `json:"name"`   // Название метки.
	Color string `json:"color"`  // Цвет метки (например, #FF5733).
}
