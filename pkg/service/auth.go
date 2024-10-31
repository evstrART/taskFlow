package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/evstrART/taskFlow"
	"github.com/evstrART/taskFlow/pkg/repository"
)

const salt = "taskFlow123"

type AuthService struct {
	repo repository.AutorisationService
}

func NewAuthService(repo repository.AutorisationService) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user taskFlow.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
