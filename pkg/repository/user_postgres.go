package repository

import (
	"fmt"
	"github.com/evstrART/taskFlow"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetUser(userId int) (taskFlow.User, error) {
	var user taskFlow.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", UserTable)
	err := r.db.Get(&user, query, userId)
	return user, err
}

func (r *UserPostgres) GetUsers() ([]taskFlow.User, error) {
	var users []taskFlow.User
	query := fmt.Sprintf("SELECT * FROM %s", UserTable)
	err := r.db.Select(&users, query)
	return users, err
}
func (r *UserPostgres) UpdateUser(userId int, input taskFlow.UpdateUserInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Username != nil {
		setValues = append(setValues, fmt.Sprintf("username=$%d", argId))
		args = append(args, *input.Username)
		argId++
	}

	if input.Email != nil {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argId))
		args = append(args, *input.Email)
		argId++
	}

	// Если ничего не обновляется, возвращаем nil
	if len(setValues) == 0 {
		return nil
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`
        UPDATE Users
        SET %s, updated_at = CURRENT_TIMESTAMP
        WHERE user_id = $%d
    `, setQuery, argId)

	args = append(args, userId)

	// Начинаем транзакцию
	tx, err := r.db.Beginx() // Используем Beginx()
	if err != nil {
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
func (r *UserPostgres) DeleteUser(userId int) error {
	// Проверяем наличие незавершенных проектов
	var unfinishedProjectsCount int
	err := r.db.QueryRow("SELECT COUNT(*) FROM Projects WHERE owner_id = $1 AND status != 'completed'", userId).Scan(&unfinishedProjectsCount)
	if err != nil {
		return err
	}

	// Проверяем наличие незавершенных задач
	var unfinishedTasksCount int
	err = r.db.QueryRow("SELECT COUNT(*) FROM Tasks WHERE assigned_to = $1 AND status != 'completed'", userId).Scan(&unfinishedTasksCount)
	if err != nil {
		return err
	}

	// Если есть незавершенные проекты или задачи, возвращаем ошибку
	if unfinishedProjectsCount > 0 || unfinishedTasksCount > 0 {
		return fmt.Errorf("у вас остались незавершенные проекты или задачи. Выполните или удалите их перед удалением профиля.")
	}

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

	// Удаляем записи из ActivityLogs
	_, err = tx.Exec("DELETE FROM ActivityLogs WHERE user_id = $1", userId)
	if err != nil {
		return err
	}

	// Удаляем пользователя
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", UserTable)
	_, err = tx.Exec(query, userId)
	if err != nil {
		return err
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *UserPostgres) GetUserByNameAndEmail(username, email string) (taskFlow.User, error) {
	var user taskFlow.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1 AND email = $2", UserTable)
	err := r.db.Get(&user, query, username, email)
	return user, err
}
