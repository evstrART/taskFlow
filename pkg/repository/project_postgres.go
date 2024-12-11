package repository

import (
	"fmt"
	"github.com/evstrART/taskFlow"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
)

type ProjectPostgres struct {
	db *sqlx.DB
}

func NewProjectPostgres(db *sqlx.DB) *ProjectPostgres {
	return &ProjectPostgres{db: db}
}

func (r *ProjectPostgres) Create(ownerID int, project taskFlow.Project) (int, error) {
	var projectID int

	// Начинаем транзакцию
	tx, err := r.db.Beginx()
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

	_, err = tx.Exec(fmt.Sprintf("SET LOCAL myapp.user_id = %d", ownerID))
	if err != nil {
		return 0, err
	}
	// Запрос для вставки проекта
	query := fmt.Sprintf(`
        INSERT INTO %s (name, description, owner_id, start_date, end_date, status)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING project_id`, ProjectTable)

	// Подготовка параметров
	params := []interface{}{
		project.Name,
		project.Description,
		ownerID,
		project.StartDate,
		project.EndDate,
		project.Status,
	}

	err = tx.QueryRow(query, params...).Scan(&projectID)
	if err != nil {
		return 0, err
	}

	// Запрос для вставки участника проекта
	memberQuery := fmt.Sprintf(`
	           INSERT INTO %s (project_id, user_id, role)
	           VALUES ($1, $2, $3)`, ProjectMemberTable)

	// Выполняем вставку с текущей датой
	_, err = tx.Exec(memberQuery, projectID, ownerID, "owner")
	if err != nil {
		return 0, err
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return projectID, nil
}

func (r *ProjectPostgres) GetAllProjects() ([]taskFlow.Project, error) {
	var projects []taskFlow.Project
	query := fmt.Sprintf("SELECT * FROM %s", ProjectTable)
	err := r.db.Select(&projects, query)
	return projects, err
}

func (r *ProjectPostgres) GetProjectById(id int) (taskFlow.Project, error) {
	var project taskFlow.Project
	query := fmt.Sprintf("SELECT * FROM %s WHERE project_id = $1", ProjectTable)
	err := r.db.Get(&project, query, id)
	return project, err
}

func (r *ProjectPostgres) DeleteProject(userID, id int) error {
	tx, err := r.db.Beginx()
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
	_, err = tx.Exec(fmt.Sprintf("SET LOCAL myapp.user_id = %d", userID))
	if err != nil {
		return err
	}
	// Запрос для удаления проекта
	query := fmt.Sprintf(`DELETE FROM %s WHERE project_id = $1`, ProjectTable)
	_, err = tx.Exec(query, id)
	if err != nil {
		return err // Если ошибка, возвращаем её
	}

	// Запрос для удаления участников проекта
	memberQuery := fmt.Sprintf(`DELETE FROM %s WHERE project_id = $1`, ProjectMemberTable)
	_, err = tx.Exec(memberQuery, id)
	if err != nil {
		return err // Если ошибка, возвращаем её
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProjectPostgres) UpdateProject(userID, id int, input taskFlow.UpdateProjectInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	// Проверяем поля и добавляем их в запрос
	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.StartDate != nil {
		setValues = append(setValues, fmt.Sprintf("start_date=$%d", argId))
		args = append(args, *input.StartDate)
		argId++
	}

	if input.EndDate != nil {
		setValues = append(setValues, fmt.Sprintf("end_date=$%d", argId))
		args = append(args, *input.EndDate)
		argId++
	}

	if input.Status != nil {
		setValues = append(setValues, fmt.Sprintf("status=$%d", argId))
		args = append(args, *input.Status)
		argId++
	}

	// Формируем строку обновления
	setQuery := strings.Join(setValues, ", ")
	if setQuery == "" {
		return nil // Если нет полей для обновления, ничего не делаем
	}

	// Формируем полный запрос
	query := fmt.Sprintf("UPDATE %s SET %s, updated_at = NOW() WHERE project_id = $%d", ProjectTable,
		setQuery, argId)
	args = append(args, id)

	// Начинаем транзакцию
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	// Устанавливаем user_id из JWT токена
	_, err = tx.Exec(fmt.Sprintf("SET LOCAL myapp.user_id = %d", userID))
	if err != nil {
		return err
	}

	// Логирование
	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %v", args)

	// Выполняем запрос
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback() // Откат транзакции в случае ошибки
		return err
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProjectPostgres) AddMembers(projectId int, input taskFlow.AddMemberRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (project_id, user_id, role) VALUES ($1, $2, $3)", ProjectMemberTable)
	_, err := r.db.Exec(query, projectId, input.UserId, input.Role)
	if err != nil {
		return fmt.Errorf("failed to add member to project: %w", err)
	}
	return nil

}

func (r *ProjectPostgres) GetMembers(projectId int) ([]taskFlow.User, error) {
	var users []taskFlow.User
	query := `
		SELECT u.user_id, u.username, u.email, pm.role 
		FROM ProjectMembers pm 
	 JOIN Users u ON pm.user_id = u.user_id 
		WHERE pm.project_id = $1`
	err := r.db.Select(&users, query, projectId)
	return users, err
}

// add delete member
func (r *ProjectPostgres) DeleteMember(projectId, memberId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE project_id = $1 AND user_id = $2", ProjectMemberTable)
	result, err := r.db.Exec(query, projectId, memberId)
	if err != nil {
		return fmt.Errorf("failed to delete member from project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no member found with user_id %d in project %d", memberId, projectId)
	}

	return nil
}

func (r *ProjectPostgres) CompleteProject(projectId, userID int) error {
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

	_, err = tx.Exec(fmt.Sprintf("SET LOCAL myapp.user_id = %d", userID))
	if err != nil {
		return err
	}
	// Обновляем статус проекта на "Завершен"
	updateQuery := fmt.Sprintf("UPDATE %s SET status = 'Completed' WHERE project_id = $1", ProjectTable)
	_, err = tx.Exec(updateQuery, projectId)
	if err != nil {
		return err
	}

	// Записываем действие в ActivityLogs
	logQuery := `INSERT INTO ActivityLogs (user_id, action, related_entity, entity_id)
                  VALUES ($1, 'COMPLETE', 'projects', $2)`
	_, err = tx.Exec(logQuery, userID, projectId)
	if err != nil {
		return err
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
