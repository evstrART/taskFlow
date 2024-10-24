package repository

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

func NewRepository() *Repository {
	return &Repository{}
}
