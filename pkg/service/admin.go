package service

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/evstrART/taskFlow/pkg/repository"
	"github.com/signintech/gopdf"
	"os"
)

const (
	UserTable          = "users"
	TaskTable          = "tasks"
	TagTable           = "tags"
	ProjectTable       = "projects"
	CommentTable       = "comments"
	ActivityLogsTable  = "activitylogs"
	TaskTagTable       = "task_tags"
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
	var content string
	// Получаем логи активности из базы данных
	logs, err := a.repo.SelectActivityLogs() // Используем метод для получения логов активности
	if err != nil {
		return "", err // Обработка ошибки, если запрос не удался
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	// Установка стандартного шрифта без кастомизации
	err = pdf.SetFont("Helvetica", "", 12) // Используем стандартный шрифт Helvetica
	if err != nil {
		return "", err // Обработка ошибки при установке шрифта
	}

	for _, log := range logs {
		content = fmt.Sprintf("Log ID: %d\nUser ID: %d\nAction: %s\nTimestamp: %s\nRelated Entity: %s\nEntity ID: %d\n",
			log.LogID, log.UserID, log.Action, log.Timestamp.Format("2006-01-02 15:04:05"), log.RelatedEntity, log.EntityID)

		if pdf.GetY() > 800 { // Проверка на переполнение страницы
			pdf.AddPage() // Добавляем новую страницу при необходимости
		}

		err = pdf.Cell(nil, content) // Запись полного содержимого в PDF
		if err != nil {
			return "", err // Обработка ошибки при записи в PDF
		}
		pdf.Br(20) // Перенос строки после каждой записи
	}

	// Указываем путь для сохранения PDF-файла
	pathPDF := "/Users/Учеба/3 курс/CУБД/reports/report.pdf"

	err = pdf.WritePdf(pathPDF) // Запись PDF-файла на диск
	if err != nil {
		return "", err // Обработка ошибки при записи PDF-файла
	}

	return pathPDF, nil // Возврат пути к созданному PDF-файлу и nil как ошибку
}

func (a *AdminService) GetReportExcel() (string, error) {
	file := excelize.NewFile()
	sheetName := "Report"
	index := file.NewSheet(sheetName)

	// Устанавливаем заголовки
	headers := []string{"Log ID", "User ID", "Action", "Timestamp", "Related Entity", "Entity ID"}
	for colNum, header := range headers {
		cell := excelize.ToAlphaString(colNum+1) + "1"
		file.SetCellValue(sheetName, cell, header)
	}

	// Добавляем данные из ActivityLogs
	logs, err := a.repo.SelectActivityLogs() // Получаем логи активности
	if err != nil {
		return "", err // Обработка ошибки
	}

	for rowNum, log := range logs {
		colNum := 1
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

	// Устанавливаем активный лист
	file.SetActiveSheet(index)

	// Указываем путь для сохранения Excel-файла
	pathExcel := "/Users/Учеба/3 курс/CУБД/reports/report.xlsx"

	// Сохраняем файл
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
