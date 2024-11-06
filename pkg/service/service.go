package service

import (
	"github.com/evstrART/taskFlow"
	"github.com/evstrART/taskFlow/pkg/repository"
)

type AutorisationService interface {
	CreateUser(user taskFlow.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
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
	UpdateProject(id int, input taskFlow.UpdateProjectInput) error
}

type Service struct {
	AutorisationService
	Project
	Task
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		AutorisationService: NewAuthService(repos.AutorisationService),
		Project:             NewProjectService(repos.Project),
	}
}
