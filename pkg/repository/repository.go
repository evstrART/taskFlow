package repository

import (
	"github.com/jmoiron/sqlx"
)

type AutorisationService interface {
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
	return &Repository{}
}
