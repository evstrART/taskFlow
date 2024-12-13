package taskFlow

type ProjectStats struct {
	ProjectName    string `db:"project_name" json:"project_name"`
	CompletedTasks int    `db:"completed_tasks" json:"completed_tasks"`
}

type UserStats struct {
	Username     string `db:"username" json:"username"`
	CreatedTasks int    `db:"created_tasks" json:"created_tasks"`
}
