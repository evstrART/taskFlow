package service

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/evstrART/taskFlow/pkg/repository"
	"github.com/jung-kurt/gofpdf"
	"os"
)

const (
	UserTable          = "users"
	TaskTable          = "tasks"
	TagTable           = "tags"
	ProjectTable       = "projects"
	CommentTable       = "comments"
	ActivityLogsTable  = "activitylogs"
	TaskTagTable       = "tasktags"
	ProjectMemberTable = "projectmembers"
)

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

func (a *AdminService) GetReportPDF() (string, error) {
	logs, err := a.repo.SelectActivityLogs() // Получаем логи активности
	if err != nil {
		return "", err // Обработка ошибки при получении логов
	}

	pdf := gofpdf.New("P", "mm", "A4", "") // Создаем новый PDF документ
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Activity Logs Report") // Заголовок отчета
	pdf.Ln(20)                              // Перенос строки

	// Устанавливаем шрифт для содержимого
	pdf.SetFont("Arial", "", 12)

	// Заголовки таблицы с учетом ширины ячеек
	headers := []string{"Log ID", "User ID", "Action", "Timestamp", "Related Entity", "Entity ID"}

	// Устанавливаем ширину ячеек для заголовков
	headerWidths := []float64{30, 30, 30, 50, 30, 30} // Ширина для каждого заголовка

	for i, header := range headers {
		pdf.Cell(headerWidths[i], 10, header) // Используем ширину из массива
	}
	pdf.Ln(-1) // Перенос строки

	// Добавляем данные логов
	for _, log := range logs {
		pdf.Cell(headerWidths[0], 10, fmt.Sprintf("%d", log.LogID))  // Log ID
		pdf.Cell(headerWidths[1], 10, fmt.Sprintf("%d", log.UserID)) // User ID
		pdf.Cell(headerWidths[2], 10, log.Action)                    // Action

		// Увеличиваем расстояние между Timestamp и Related Entity
		pdf.Cell(headerWidths[3], 10, log.Timestamp.Format("2006-01-02 15:04:05")) // Timestamp
		pdf.Cell(headerWidths[4], 10, log.RelatedEntity)                           // Related Entity

		pdf.Cell(headerWidths[5], 10, fmt.Sprintf("%d", log.EntityID)) // Entity ID
		pdf.Ln(-1)                                                     // Перенос строки после каждой записи
	}

	// Указываем путь для сохранения PDF-файла
	pathPDF := "/Users/Учеба/3 курс/CУБД/reports/report.pdf"

	if err := pdf.OutputFileAndClose(pathPDF); err != nil {
		return "", err // Обработка ошибки при сохранении PDF-файла
	}

	return pathPDF, nil // Возврат пути к созданному PDF-файлу
}

func (a *AdminService) GetReportExcel() (string, error) {
	file := excelize.NewFile() // Создаем новый файл Excel
	sheetName := "Activity Logs"

	index := file.NewSheet(sheetName) // Создаем новый лист
	file.SetActiveSheet(index)

	// Заголовки таблицы
	headers := []string{"Log ID", "User ID", "Action", "Timestamp", "Related Entity", "Entity ID"}
	for colNum, header := range headers {
		cell := excelize.ToAlphaString(colNum) + "1"
		file.SetCellValue(sheetName, cell, header)
	}

	logs, err := a.repo.SelectActivityLogs() // Получаем логи активности
	if err != nil {
		return "", err // Обработка ошибки при получении логов
	}

	for rowNum, log := range logs {
		colNum := 0
		file.SetCellValue(sheetName, excelize.ToAlphaString(colNum)+fmt.Sprintf("%d", rowNum+2), log.LogID)
		colNum++
		file.SetCellValue(sheetName, excelize.ToAlphaString(colNum)+fmt.Sprintf("%d", rowNum+2), log.UserID)
		colNum++
		file.SetCellValue(sheetName, excelize.ToAlphaString(colNum)+fmt.Sprintf("%d", rowNum+2), log.Action)
		colNum++
		file.SetCellValue(sheetName, excelize.ToAlphaString(colNum)+fmt.Sprintf("%d", rowNum+2), log.Timestamp.Format("2006-01-02 15:04:05"))
		colNum++
		file.SetCellValue(sheetName, excelize.ToAlphaString(colNum)+fmt.Sprintf("%d", rowNum+2), log.RelatedEntity)
		colNum++
		file.SetCellValue(sheetName, excelize.ToAlphaString(colNum)+fmt.Sprintf("%d", rowNum+2), log.EntityID)
	}

	pathExcel := "/Users/Учеба/3 курс/CУБД/reports/report3.xlsx" // Указываем путь для сохранения Excel-файла

	if err := file.SaveAs(pathExcel); err != nil {
		return "", err // Обработка ошибки при сохранении файла
	}

	return pathExcel, nil // Возврат пути к созданному Excel-файлу
}

func (a *AdminService) GetDBInJSON() (string, error) {
	// Получение данных из различных таблиц
	resultUsersTable, err := a.repo.SelectTable(UserTable)
	if err != nil {
		return "", err
	}

	resultProjectsTable, err := a.repo.SelectTable(ProjectTable)
	if err != nil {
		return "", err
	}

	resultTasksTable, err := a.repo.SelectTable(TaskTable)
	if err != nil {
		return "", err
	}

	resultCommentsTable, err := a.repo.SelectTable(CommentTable)
	if err != nil {
		return "", err
	}

	resultProjectMembersTable, err := a.repo.SelectTable(ProjectMemberTable)
	if err != nil {
		return "", err
	}

	resultActivityLogsTable, err := a.repo.SelectTable(ActivityLogsTable)
	if err != nil {
		return "", err
	}

	resultTagsTable, err := a.repo.SelectTable(TagTable)
	if err != nil {
		return "", err
	}

	resultTaskTagsTable, err := a.repo.SelectTable(TaskTagTable)
	if err != nil {
		return "", err
	}

	// Объединение результатов в одну структуру
	result := map[string]interface{}{
		UserTable:          resultUsersTable,
		ProjectTable:       resultProjectsTable,
		TaskTable:          resultTasksTable,
		CommentTable:       resultCommentsTable,
		ProjectMemberTable: resultProjectMembersTable,
		ActivityLogsTable:  resultActivityLogsTable,
		TagTable:           resultTagsTable,
		TaskTagTable:       resultTaskTagsTable,
	}

	// Запись данных в JSON файл
	file, err := os.Create("/Users/Учеба/3 курс/CУБД/backups/output.json") // Укажите путь для сохранения JSON-файла
	if err != nil {
		return "", err // Обработка ошибки при создании файла
	}

	defer file.Close() // Закрываем файл после завершения работы

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "   ") // Устанавливаем отступы для читаемости JSON

	if err := encoder.Encode(result); err != nil {
		return "", err // Обработка ошибки при записи JSON в файл
	}

	return "/Users/Учеба/3 курс/CУБД/backups/output.json", nil // Возврат пути к созданному JSON-файлу и nil как ошибку
}
