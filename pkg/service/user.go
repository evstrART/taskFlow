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
func (s *UserService) UpdateUser(userId int, input taskFlow.UpdateUserInput) error {
	return s.repo.UpdateUser(userId, input)
}
func (s *UserService) CheckOldPassword(userId int, oldPassword string) (bool, error) {
	user, err := s.repo.GetUser(userId) // Получаем пользователя из базы данных
	if err != nil {
		return false, err
	}
	return checkPasswordHash(oldPassword, user.Password), nil // Используем новую функцию
}

func (s *UserService) DeleteUser(userId int) error {
	return s.repo.DeleteUser(userId)
}
