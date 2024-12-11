package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) checkAdmin(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized: unable to get user ID")
		return
	}
	isAdmin, err := h.services.Admin.CheckAdmin(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if isAdmin {
		c.JSON(http.StatusOK, gin.H{"isAdmin": true})
	} else {
		newErrorResponse(c, http.StatusUnauthorized, "you are not admin")
		return
	}
}

func (h *Handler) getReport(ctx *gin.Context) {
	var path string
	var err error
	format := ctx.Param("format") // Получаем формат из параметров URL

	if format == "pdf" {
		path, err = h.services.Admin.GetReportPDF() // Генерация PDF-отчета
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	} else if format == "excel" {
		path, err = h.services.Admin.GetReportExcel() // Генерация Excel-отчета
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		newErrorResponse(ctx, http.StatusBadRequest, "Invalid format")
		return
	}

	ctx.File(path) // Отправляем файл клиенту
	ctx.Status(http.StatusOK)
}

func (h *Handler) exportJSON(ctx *gin.Context) {
	filePath, err := h.services.Admin.GetDBInJSON()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.File(filePath)
	ctx.Status(http.StatusOK)
}
