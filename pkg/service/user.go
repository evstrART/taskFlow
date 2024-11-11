package service

import (
	"github.com/evstrART/taskFlow"
	"github.com/evstrART/taskFlow/pkg/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(userId int) (taskFlow.User, error) {
	return s.repo.GetUser(userId)
}

func (s *UserService) GetUsers() ([]taskFlow.User, error) {
	return s.repo.GetUsers()
}
