package service

import (
	"github.com/evstrART/taskFlow"
	"github.com/evstrART/taskFlow/pkg/repository"
)

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) AddComment(taskId, userId int, input taskFlow.CommentInput) (int, error) {
	return s.repo.AddComment(taskId, userId, input)
}

func (s *CommentService) GetComments(taskId int) ([]taskFlow.Comment, error) {
	return s.repo.GetComments(taskId)
}

func (s *CommentService) GetCommentById(taskID, id int) (taskFlow.Comment, error) {
	return s.repo.GetCommentById(taskID, id)
}
func (s *CommentService) GetAllCommentsForUser(userId int) ([]taskFlow.Comment, error) {
	return s.repo.GetAllCommentsForUser(userId)
}

func (s *CommentService) DeleteComment(commentId, userId int) error {
	return s.repo.DeleteComment(commentId, userId)
}

func (s *CommentService) UpdateComment(commentId, userId int, input taskFlow.CommentInput) error {
	return s.repo.UpdateComment(commentId, userId, input)
}
