package service

import "github.com/evstrART/taskFlow/pkg/repository"

type AutorisationService interface {
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
	return &Service{}
}
