package repository

import (
	"fmt"
	"github.com/evstrART/taskFlow"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

type TaskPostgres struct {
	db *sqlx.DB
}

func NewTaskPostgres(db *sqlx.DB) *TaskPostgres {
	return &TaskPostgres{db: db}
}

func (r *TaskPostgres) CreateTask(userID, projectId int, input taskFlow.Task) (int, error) {
	var taskID int

	query := `
        INSERT INTO Tasks (project_id, title, description, assigned_to, status, priority, due_date, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
        RETURNING task_id`

	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// Обеспечиваем откат в случае ошибки
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	// Устанавливаем локальный user_id (без параметров)
	if _, err := tx.Exec(fmt.Sprintf("SET LOCAL myapp.user_id = %d", userID)); err != nil {
		return 0, err
	}

	// Выполняем запрос и возвращаем идентификатор созданной задачи
	err = tx.QueryRow(query, projectId, input.Title, input.Description, input.AssignedTo, input.Status, input.Priority, input.DueDate).Scan(&taskID)
	if err != nil {
		return 0, err
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return taskID, nil
}

func (r *TaskPostgres) GetAllTasks(projectId int) ([]taskFlow.Task, error) {
	var tasks []taskFlow.Task
	query := fmt.Sprintf("SELECT * FROM %s WHERE project_id = $1", TaskTable)
	err := r.db.Select(&tasks, query, projectId)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskPostgres) GetTask(projectId int, taskId int) (taskFlow.Task, error) {
	var task taskFlow.Task
	query := fmt.Sprintf("SELECT * FROM %s WHERE project_id = $1 AND task_id = $2", TaskTable)
	err := r.db.Get(&task, query, projectId, taskId)
	if err != nil {
		return taskFlow.Task{}, err
	}

	return task, nil
}

func (r *TaskPostgres) DeleteTask(userID, projectId, taskId int) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Обеспечиваем откат в случае ошибки
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	// Устанавливаем локальный user_id
	if _, err := tx.Exec(fmt.Sprintf("SET LOCAL myapp.user_id = %d", userID)); err != nil {
		return err
	}

	// Запрос для удаления задачи
	query := fmt.Sprintf("DELETE FROM %s WHERE project_id = $1 AND task_id = $2", TaskTable)
	_, err = tx.Exec(query, projectId, taskId)
	if err != nil {
		return err
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *TaskPostgres) UpdateTask(userID, projectId, taskId int, input taskFlow.UpdateTaskInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.AssignedTo != nil {
		setValues = append(setValues, fmt.Sprintf("assigned_to=$%d", argId))
		args = append(args, *input.AssignedTo)
		argId++
	}

	if input.Status != nil {
		setValues = append(setValues, fmt.Sprintf("status=$%d", argId))
		args = append(args, *input.Status)
		argId++
	}

	if input.Priority != nil {
		setValues = append(setValues, fmt.Sprintf("priority=$%d", argId))
		args = append(args, *input.Priority)
		argId++
	}

	if input.DueDate != nil {
		setValues = append(setValues, fmt.Sprintf("due_date=$%d", argId))
		args = append(args, *input.DueDate)
		argId++
	}

	// Если ничего не обновляется, возвращаем nil
	if len(setValues) == 0 {
		return nil
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`
        UPDATE Tasks
        SET %s
        WHERE project_id = $%d AND task_id = $%d
    `, setQuery, argId, argId+1)

	args = append(args, projectId, taskId)

	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Устанавливаем локальный user_id
	if _, err := tx.Exec(fmt.Sprintf("SET LOCAL myapp.user_id = %d", userID)); err != nil {
		tx.Rollback() // Откат в случае ошибки
		return err
	}

	// Выполняем запрос
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback() // Откат в случае ошибки
		return err
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *TaskPostgres) GetAllTasksForUser(userID int) ([]taskFlow.Task, error) {
	var tasks []taskFlow.Task
	query := fmt.Sprintf("SELECT * FROM %s WHERE assigned_to = $1", TaskTable)
	err := r.db.Select(&tasks, query, userID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskPostgres) CompleteTask(taskId, userId int) error {
	// Начинаем транзакцию
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	// Отменяем транзакцию в случае ошибки
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if _, err := tx.Exec(fmt.Sprintf("SET LOCAL myapp.user_id = %d", userId)); err != nil {
		return err
	}

	// Обновляем статус задачи на "Завершено"
	_, err = tx.Exec(`UPDATE tasks SET status = $1 WHERE task_id = $2`, "Завершено", taskId)
	if err != nil {
		return err
	}

	// Добавляем запись в activitylogs
	_, err = tx.Exec(`INSERT INTO activitylogs (action, related_entity, user_id, entity_id) VALUES ($1, $2, $3, $4)`,
		"COMPLETE", "tasks", userId, taskId)
	if err != nil {
		return err
	}

	// Подтверждаем транзакцию
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
