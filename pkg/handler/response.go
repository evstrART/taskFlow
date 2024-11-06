package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(ctx *gin.Context, statuscode int, message string) {
	logrus.Error(message)
	ctx.AbortWithStatusJSON(statuscode, errorResponse{message})
}

type statusResponse struct {
	Status string `json:"status"`
}
