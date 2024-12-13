package handler

import (
	"github.com/gin-gonic/gin"
	"log"
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
	format := ctx.Param("format")

	if format == "pdf" {
		path, err = h.services.Admin.GetReportPDF()
	} else if format == "excel" {
		path, err = h.services.Admin.GetReportExcel()
	} else {
		newErrorResponse(ctx, http.StatusBadRequest, "Invalid format")
		return
	}

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.File(path)
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
func (h *Handler) importJSON(ctx *gin.Context) {
	file, err := ctx.FormFile("json")
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.Admin.ImportDBInJSON(file)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}
func (h *Handler) GetCompletedTasksByProject(c *gin.Context) {
	stats, err := h.services.Admin.GetCompletedTasksByProject()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Completed Tasks By Project: %+v\n", stats)

	// Возвращаем только проекты
	c.JSON(http.StatusOK, gin.H{"projects": stats})
}
func (h *Handler) GetCreatedTasksByUser(c *gin.Context) {
	stats, err := h.services.Admin.GetCreatedTasksByUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Created Tasks By User: %+v\n", stats)

	// Возвращаем только пользователей
	c.JSON(http.StatusOK, gin.H{"users": stats})
}
