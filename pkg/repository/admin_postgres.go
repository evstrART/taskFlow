package repository

import (
	"github.com/evstrART/taskFlow"
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

func (r *AdminPostgres) SelectActivityLogs() ([]taskFlow.ActivityLog, error) {
	var logs []taskFlow.ActivityLog

	// SQL-запрос для извлечения данных из таблицы ActivityLogs
	query := `
        SELECT log_id, user_id, action, timestamp, related_entity, entity_id
        FROM activitylogs
        ORDER BY timestamp DESC`

	// Выполнение запроса и сканирование результатов в структуру logs
	err := r.db.Select(&logs, query)
	if err != nil {
		return nil, err // Возвращаем ошибку, если запрос не удался
	}

	return logs, nil // Возвращаем срез логов активности и nil как ошибку
}

func (a *AdminPostgres) SelectTable(tableName string) ([]map[string]interface{}, error) {
	// Подготовка SQL-запроса
	query := "SELECT * FROM " + tableName

	// Выполнение запроса
	rows, err := a.db.Queryx(query)
	if err != nil {
		return nil, err // Обработка ошибки, если запрос не удался
	}
	defer rows.Close() // Закрываем rows после завершения работы с ними

	// Преобразование данных в формат JSON
	var result []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		// Сканируем данные в карту
		if err := rows.MapScan(row); err != nil {
			return nil, err // Обработка ошибки при сканировании строки
		}
		result = append(result, row) // Добавляем строку в результат
	}

	if err := rows.Err(); err != nil {
		return nil, err // Проверяем на ошибки после завершения итерации
	}

	return result, nil // Возвращаем результат
}
