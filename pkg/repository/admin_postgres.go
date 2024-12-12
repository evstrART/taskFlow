package repository

import (
	"github.com/evstrART/taskFlow"
	"github.com/jmoiron/sqlx"
)

type AdminPostgres struct {
	db *sqlx.DB
}

func NewAdminPostgres(db *sqlx.DB) *AdminPostgres {
	return &AdminPostgres{db: db}
}

func (r *AdminPostgres) SelectAdminId(userId int) ([]int, error) {
	var adminIds []int

	query := "SELECT user_id FROM Users WHERE role = $1"

	rows, err := r.db.Query(query, "admin")
	if err != nil {
		return nil, err // Возвращаем ошибку, если запрос не удался
	}
	defer rows.Close() // Закрываем rows после завершения работы с ними

	// Проходим по результатам запроса
	for rows.Next() {
		var adminId int
		if err := rows.Scan(&adminId); err != nil {
			return nil, err // Возвращаем ошибку, если сканирование не удалось
		}
		adminIds = append(adminIds, adminId) // Добавляем идентификатор в список
	}

	// Проверяем наличие ошибок после завершения итерации
	if err := rows.Err(); err != nil {
		return nil, err // Возвращаем ошибку, если она возникла во время итерации
	}

	return adminIds, nil // Возвращаем список идентификаторов администраторов
}

func (r *AdminPostgres) SelectActivityLogs() ([]taskFlow.ActivityLog, error) {
	var logs []taskFlow.ActivityLog

	// SQL-запрос для извлечения данных из таблицы ActivityLogs
	query := `
        SELECT log_id, user_id, action, timestamp, related_entity, entity_id
        FROM activitylogs
        ORDER BY timestamp DESC`

	// Выполнение запроса и сканирование результатов в структуру logs
	err := r.db.Select(&logs, query)
	if err != nil {
		return nil, err // Возвращаем ошибку, если запрос не удался
	}

	return logs, nil // Возвращаем срез логов активности и nil как ошибку
}

func (a *AdminPostgres) SelectTable(tableName string) ([]map[string]interface{}, error) {
	// Подготовка SQL-запроса
	query := "SELECT * FROM " + tableName

	// Выполнение запроса
	rows, err := a.db.Queryx(query)
	if err != nil {
		return nil, err // Обработка ошибки, если запрос не удался
	}
	defer rows.Close() // Закрываем rows после завершения работы с ними

	// Преобразование данных в формат JSON
	var result []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		// Сканируем данные в карту
		if err := rows.MapScan(row); err != nil {
			return nil, err // Обработка ошибки при сканировании строки
		}
		result = append(result, row) // Добавляем строку в результат
	}

	if err := rows.Err(); err != nil {
		return nil, err // Проверяем на ошибки после завершения итерации
	}

	return result, nil // Возвращаем результат
}

func (a *AdminPostgres) InsertUser(user taskFlow.User) error {
	query := `INSERT INTO users (username, password, email, role, created_at, updated_at) VALUES ($1, $2, $3, $4, DEFAULT, $5)`
	_, err := a.db.Exec(query, user.Username, user.Password, user.Email, user.Role, user.UpdatedAt)
	return err
}
func (a *AdminPostgres) InsertProject(project taskFlow.Project) error {
	query := `INSERT INTO projects (name, description, owner_id, start_date, end_date, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, DEFAULT, $7)`
	_, err := a.db.Exec(query, project.Name, project.Description, project.OwnerID, project.StartDate, project.EndDate, project.Status, project.UpdatedAt)
	return err
}
func (a *AdminPostgres) InsertTask(task taskFlow.Task) error {
	query := `INSERT INTO tasks (project_id, title, description, assigned_to, status, priority, due_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, DEFAULT, $8)`
	_, err := a.db.Exec(query, task.ProjectID, task.Title, task.Description, task.AssignedTo, task.Status, task.Priority, task.DueDate, task.UpdatedAt)
	return err
}
func (a *AdminPostgres) InsertComment(comment taskFlow.Comment) error {
	query := `INSERT INTO comments (task_id, user_id, content, created_at) VALUES ($1, $2, $3, DEFAULT)`
	_, err := a.db.Exec(query, comment.TaskID, comment.UserID, comment.Content)
	return err
}
func (a *AdminPostgres) InsertProjectMember(member taskFlow.ProjectMember) error {
	query := `INSERT INTO projectmembers (project_id, user_id, role, created_at) VALUES ($1 ,$2 ,$3 ,DEFAULT)`
	_, err := a.db.Exec(query, member.ProjectID, member.UserID, member.Role)
	return err
}
func (a *AdminPostgres) InsertActivityLog(log taskFlow.ActivityLog) error {
	query := `INSERT INTO activitylogs (user_id , action , timestamp , related_entity , entity_id ) VALUES ($1 ,$2 ,DEFAULT ,$3 ,$4)`
	_, err := a.db.Exec(query, log.UserID, log.Action, log.RelatedEntity, log.EntityID)
	return err
}
func (a *AdminPostgres) InsertTag(tag taskFlow.Tag) error {
	query := `INSERT INTO tags (name,color ) VALUES ($1 ,$2 )`
	_, err := a.db.Exec(query, tag.Name, tag.Color)
	return err
}
func (a *AdminPostgres) InsertTaskTag(taskTag taskFlow.TaskTag) error {
	query := `INSERT INTO tasktags (task_id ,tag_id ) VALUES ($1 ,$2 )`
	_, err := a.db.Exec(query, taskTag.TaskID, taskTag.TagID)
	return err
}
