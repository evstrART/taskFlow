package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/evstrART/taskFlow"
	"github.com/evstrART/taskFlow/pkg/repository"
	"github.com/golang-jwt/jwt"
	"time"
)

const (
	salt             = "taskFlow123"
	tokenTTL         = time.Hour * 12
	signingKey       = "dvkdvkdfvsfvhg"
	tokenTTLForEmail = 15 * time.Minute
)

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
func checkPasswordHash(password, hash string) bool {
	generatedHash := generatePasswordHash(password)
	return generatedHash == hash
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.UserID,
	})
	return token.SignedString([]byte(signingKey))
}
func (a *AuthService) GenerateTokenForReset(username, email string) (taskFlow.User, string, error) {
	user, err := a.repo.GetUserByNameAndEmail(username, email)
	if err != nil {
		return taskFlow.User{}, "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTLForEmail).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.UserID,
	})
	tokenString, err := token.SignedString([]byte(signingKey))
	return user, tokenString, err
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type")
	}
	return claims.UserId, nil
}

func (s *AuthService) ChangePassword(userId int, newPasswordInput string) error {
	newPasswordHex := generatePasswordHash(newPasswordInput)
	return s.repo.ChangePassword(userId, newPasswordHex)
}
func (s *AuthService) UserExistsForReset(input taskFlow.ResetPasswordInput) (bool, error) {
	return s.repo.UserExistsForReset(input)
}
