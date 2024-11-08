package repository

import (
	"fmt"
	"github.com/evstrART/taskFlow"
	"github.com/jmoiron/sqlx"
)

type CommentPostgres struct {
	db *sqlx.DB
}

func NewCommentPostgres(db *sqlx.DB) *CommentPostgres {
	return &CommentPostgres{db: db}
}

func (s *CommentPostgres) AddComment(taskId, userId int, input taskFlow.CommentInput) (int, error) {
	var commentId int
	query := fmt.Sprintf("INSERT INTO %s (task_id, user_id, content) VALUES ($1, $2, $3) RETURNING comment_id", CommentTable)
	err := s.db.QueryRow(query, taskId, userId, input.Comment).Scan(&commentId)
	if err != nil {
		return 0, err
	}
	return commentId, err
}

func (s *CommentPostgres) GetComments(taskId int) ([]taskFlow.Comment, error) {
	var comments []taskFlow.Comment
	query := fmt.Sprintf("SELECT * FROM %s WHERE task_id = $1", CommentTable)
	err := s.db.Select(&comments, query, taskId)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *CommentPostgres) GetCommentById(taskId, commentId int) (taskFlow.Comment, error) {
	var comment taskFlow.Comment
	query := fmt.Sprintf("SELECT * FROM %s WHERE task_id = $1 AND comment_id = $2", CommentTable)
	err := s.db.Get(&comment, query, taskId, commentId)
	if err != nil {
		return taskFlow.Comment{}, err
	}

	return comment, nil
}

func (s *CommentPostgres) GetAllCommentsForUser(userId int) ([]taskFlow.Comment, error) {
	var comments []taskFlow.Comment
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", CommentTable)
	err := s.db.Select(&comments, query, userId)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *CommentPostgres) DeleteComment(commentId, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE comment_id = $1 AND user_id = $2", CommentTable)
	_, err := s.db.Exec(query, commentId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *CommentPostgres) UpdateComment(commentId, userId int, input taskFlow.CommentInput) error {
	query := fmt.Sprintf("UPDATE %s SET content = $1 WHERE comment_id = $2 AND user_id = $3", CommentTable)
	_, err := s.db.Exec(query, input.Comment, commentId, userId)
	if err != nil {
		return err
	}
	return nil
}
