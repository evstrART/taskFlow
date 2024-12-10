package handler

import (
	"github.com/evstrART/taskFlow"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) createTask(c *gin.Context) {
	userID, err := getUserId(c)
	projectId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid project id")
		return
	}
	var input taskFlow.Task
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Task.CreateTask(userID, projectId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) updateTask(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user id")
	}
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var input taskFlow.UpdateTaskInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.Task.UpdateTask(userID, projectID, taskID, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
func (h *Handler) deleteTask(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.services.Task.DeleteTask(userID, projectID, taskID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "task not found")
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})

}
func (h *Handler) getAllTasks(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Param("id"))
	tasks, err := h.services.Task.GetAllTasks(projectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, tasks)
}
func (h *Handler) getTaskById(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	task, err := h.services.Task.GetTask(projectID, taskID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "task not found")
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) getAllTasksForUser(c *gin.Context) {
	userID, err := getUserId(c) // Используем вашу функцию для получения userID
	if err != nil {
		return
	}

	tasks, err := h.services.Task.GetAllTasksForUser(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if len(tasks) == 0 {
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "No tasks found",
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) completeTask(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.services.CompleteTask(taskID, userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
