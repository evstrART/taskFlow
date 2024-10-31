package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	UserTable          = "users"
	TaskTable          = "tasks"
	TagTable           = "tags"
	ProjectTable       = "projects"
	CommentTable       = "comments"
	ActivityLogsTable  = "activitylogs"
	TaskTagTable       = "task_tags"
	ProjectMemberTable = "project_members"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
