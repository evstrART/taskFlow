package service

import (
	"github.com/evstrART/taskFlow"
	"github.com/evstrART/taskFlow/pkg/repository"
)

type AutorisationService interface {
	CreateUser(user taskFlow.User) (int, error)
}

type TaskService interface {
}

type UserService interface {
}

type ProjectService interface {
}

type Service struct {
	AutorisationService
	TaskService
	UserService
	ProjectService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		AutorisationService: NewAuthService(repos.AutorisationService),
	}
}
