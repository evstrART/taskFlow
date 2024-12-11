package repository

import (
	"github.com/jmoiron/sqlx"
)

type AdminPostgres struct {
	db *sqlx.DB
}

func NewAdminPostgres(db *sqlx.DB) *AdminPostgres {
	return &AdminPostgres{db: db}
}

func (r *AdminPostgres) SelectAdminId(userId int) ([]int, error) {
	// Подготовим список для хранения идентификаторов администраторов
	var adminIds []int

	// SQL-запрос для получения идентификаторов администраторов
	query := "SELECT user_id FROM Users WHERE role = $1"

	// Выполняем запрос
	rows, err := r.db.Query(query, "admin")
	if err != nil {
		return nil, err // Возвращаем ошибку, если запрос не удался
	}
	defer rows.Close() // Закрываем rows после завершения работы с ними

	// Проходим по результатам запроса
	for rows.Next() {
		var adminId int
		if err := rows.Scan(&adminId); err != nil {
			return nil, err // Возвращаем ошибку, если сканирование не удалось
		}
		adminIds = append(adminIds, adminId) // Добавляем идентификатор в список
	}

	// Проверяем наличие ошибок после завершения итерации
	if err := rows.Err(); err != nil {
		return nil, err // Возвращаем ошибку, если она возникла во время итерации
	}

	return adminIds, nil // Возвращаем список идентификаторов администраторов
}
