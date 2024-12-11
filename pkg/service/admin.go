package service

import "github.com/evstrART/taskFlow/pkg/repository"

type AdminService struct {
	repo repository.Admin
}

func NewAdminService(repo repository.Admin) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) CheckAdmin(userId int) (bool, error) {
	// Получаем список идентификаторов администраторов
	adminIds, err := s.repo.SelectAdminId(userId)
	if err != nil {
		return false, err // Возвращаем ошибку, если запрос не удался
	}

	// Проверяем, есть ли userId в списке администраторов
	for _, id := range adminIds {
		if id == userId {
			return true, nil // Пользователь является администратором
		}
	}

	return false, nil // Пользователь не является администратором
}
