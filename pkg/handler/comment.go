package handler

import (
	"github.com/evstrART/taskFlow"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) addCooment(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var input taskFlow.CommentInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Comment.AddComment(taskID, userID, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id":     id,
		"status": "ok",
	})
}

func (h *Handler) getAllComments(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Param("task_id"))
	comments, err := h.services.Comment.GetComments(taskId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (h *Handler) getCommentById(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	commentId, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	comment, err := h.services.Comment.GetCommentById(taskId, commentId)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "comment not found")
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *Handler) getAllCommentsForUser(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}
	comments, err := h.services.Comment.GetAllCommentsForUser(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, comments)

}

func (h *Handler) deleteComment(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}
	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.services.Comment.DeleteComment(commentID, userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"Status": "ok",
	})
}

func (h *Handler) updateComment(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}
	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var input taskFlow.CommentInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.UpdateComment(commentID, userID, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})

}
