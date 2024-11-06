package repository

import (
	"github.com/evstrART/taskFlow"
	"github.com/jmoiron/sqlx"
)

type AutorisationService interface {
	CreateUser(user taskFlow.User) (int, error)
	GetUser(username, password string) (taskFlow.User, error)
}

type Task interface {
}

type User interface {
}

type Project interface {
	Create(ownerID int, project taskFlow.Project) (int, error)
	GetAllProjects() ([]taskFlow.Project, error)
	GetProjectById(id int) (taskFlow.Project, error)
	DeleteProject(id int) error
}

type Repository struct {
	AutorisationService
	Project
	Task
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AutorisationService: NewAuthPostgres(db),
		Project:             NewProjectPostgres(db),
	}
}
