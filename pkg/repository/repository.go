package repository

import (
	"github.com/evstrART/taskFlow"
	"github.com/jmoiron/sqlx"
)

type AutorisationService interface {
	CreateUser(user taskFlow.User) (int, error)
}

type TaskRepository interface {
}

type UserRepository interface {
}

type ProjectRepository interface {
}

type Repository struct {
	AutorisationService
	TaskRepository
	UserRepository
	ProjectRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AutorisationService: NewAuthPostgres(db),
	}
}
