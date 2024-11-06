package repository

import (
	"fmt"
	"github.com/evstrART/taskFlow"
	"github.com/jmoiron/sqlx"
	"log"
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

func (r *ProjectPostgres) DeleteProject(id int) error {
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
