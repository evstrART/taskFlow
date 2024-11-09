package repository

import (
	"fmt"
	"github.com/evstrART/taskFlow"
	"github.com/jmoiron/sqlx"
	"strings"
)

type TagPostgres struct {
	db *sqlx.DB
}

func NewTagPostgres(db *sqlx.DB) *TagPostgres {
	return &TagPostgres{db: db}
}

func (r *TagPostgres) CreateTag(taskId int, tag taskFlow.Tag) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err // Возвращаем ошибку, если не удалось начать транзакцию
	}

	// Сначала добавляем тег в таблицу tags
	query := "INSERT INTO tags (name, color) VALUES ($1, $2) RETURNING tag_id"
	var tagId int
	err = tx.QueryRow(query, tag.Name, tag.Color).Scan(&tagId)
	if err != nil {
		tx.Rollback() // Откатываем транзакцию в случае ошибки
		return 0, err // Возвращаем ошибку, если не удалось вставить тег
	}

	// Затем добавляем запись в связную таблицу task_tags
	linkQuery := "INSERT INTO tasktags (task_id, tag_id) VALUES ($1, $2)"
	_, err = tx.Exec(linkQuery, taskId, tagId)
	if err != nil {
		tx.Rollback() // Откатываем транзакцию в случае ошибки
		return 0, err // Возвращаем ошибку, если не удалось связать тег с задачей
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		return 0, err // Возвращаем ошибку, если не удалось подтвердить транзакцию
	}

	return tagId, nil // Возвращаем ID нового тега
}

func (r *TagPostgres) GetAllTags() ([]taskFlow.Tag, error) {
	var tags []taskFlow.Tag
	query := fmt.Sprintf("SELECT * FROM %s", TagTable)
	err := r.db.Select(&tags, query)
	if err != nil {
		return nil, err
	}
	return tags, err
}
func (r *TagPostgres) GetTags(taskId int) ([]taskFlow.Tag, error) {
	var tags []taskFlow.Tag
	// Исправленный запрос с запятой между t.name и t.color
	query := "SELECT t.tag_id, t.name, t.color FROM tags t INNER JOIN tasktags tt ON t.tag_id = tt.tag_id WHERE tt.task_id = $1"
	err := r.db.Select(&tags, query, taskId)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagPostgres) AddTag(taskId int, tag taskFlow.Tag) (int, error) {
	/*tagId, err := r.CreateTag(tag)
	if err != nil {
		return 0, err
	}

	// Затем свяжем тег с задачей
	query := "INSERT INTO task_tags (task_id, tag_id) VALUES ($1, $2)"
	_, err = r.db.Exec(query, taskId, tagId)
	if err != nil {
		return 0, err
	}

	return tagId, nil*/
	return 0, nil
}
func (r *TagPostgres) DeleteTag(tagId int) error {
	query := "DELETE FROM tags WHERE tag_id = $1"
	_, err := r.db.Exec(query, tagId)
	return err
}
func (r *TagPostgres) ChangeTag(taskId, newTag int) error {
	query := fmt.Sprintf("UPDATE %s SET tag_id = $1 WHERE task_id = $2", TaskTagTable)
	_, err := r.db.Exec(query, newTag, taskId)
	return err
}
func (r *TagPostgres) UpdateTag(tagId int, input taskFlow.TagInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	// Проверяем поля и добавляем их в запрос
	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Color != nil {
		setValues = append(setValues, fmt.Sprintf("color=$%d", argId))
		args = append(args, *input.Color)
		argId++
	}

	// Формируем строку обновления
	setQuery := strings.Join(setValues, ", ")
	if setQuery == "" {
		return nil // Если нет полей для обновления, ничего не делаем
	}

	// Формируем полный запрос
	query := fmt.Sprintf("UPDATE tags SET %s WHERE tag_id = $%d", setQuery, argId)
	args = append(args, tagId)

	// Выполняем запрос
	_, err := r.db.Exec(query, args...)
	return err // Возвращаем ошибку, если она возникла
}
