package taskFlow

type Tag struct {
	TagID int    `json:"tag_id" db:"tag_id" ` // Уникальный идентификатор метки.
	Name  string `json:"name"`                // Название метки.
	Color string `json:"color"`               // Цвет метки (например, #FF5733).
}

type TagInput struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}
