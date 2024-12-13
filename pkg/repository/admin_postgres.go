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
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var adminId int
		if err := rows.Scan(&adminId); err != nil {
			return nil, err
		}
		adminIds = append(adminIds, adminId)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return adminIds, nil
}

func (r *AdminPostgres) SelectActivityLogs() ([]taskFlow.ActivityLog, error) {
	var logs []taskFlow.ActivityLog

	query := `
        SELECT log_id, user_id, action, timestamp, related_entity, entity_id
        FROM activitylogs
        ORDER BY timestamp DESC`

	err := r.db.Select(&logs, query)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (a *AdminPostgres) SelectTable(tableName string) ([]map[string]interface{}, error) {
	query := "SELECT * FROM " + tableName

	rows, err := a.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		if err := rows.MapScan(row); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
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

func (r *AdminPostgres) GetCompletedTasksByProject() ([]taskFlow.ProjectStats, error) {
	query := `
	SELECT 
		p.name AS project_name,
		COUNT(al.log_id) AS completed_tasks
	FROM 
		ActivityLogs al
	JOIN 
		Tasks t ON al.entity_id = t.task_id
	JOIN 
		Projects p ON t.project_id = p.project_id
	WHERE 
		al.action = 'COMPLETE'
		AND al.related_entity = 'tasks'
	GROUP BY 
		p.name
	ORDER BY 
		completed_tasks DESC
	LIMIT 10;
	`

	var stats []taskFlow.ProjectStats
	err := r.db.Select(&stats, query)
	return stats, err
}

func (r *AdminPostgres) GetCreatedTasksByUser() ([]taskFlow.UserStats, error) {
	query := `
	SELECT 
		u.username,
		COUNT(al.log_id) AS created_tasks
	FROM 
		ActivityLogs al
	JOIN 
		Users u ON al.user_id = u.user_id
	WHERE 
		al.action = 'CREATE'
		AND al.related_entity = 'tasks'
	GROUP BY 
		u.username
	ORDER BY 
		created_tasks DESC
	LIMIT 10;
	`

	var stats []taskFlow.UserStats
	err := r.db.Select(&stats, query)
	return stats, err
}
