package repository

import (
	"fmt"
	"github.com/evstrART/taskFlow"
	"github.com/jmoiron/sqlx"
	"log"
)

type CommentPostgres struct {
	db *sqlx.DB
}

func NewCommentPostgres(db *sqlx.DB) *CommentPostgres {
	return &CommentPostgres{db: db}
}

func (s *CommentPostgres) AddComment(taskId, userId int, input taskFlow.CommentInput) (int, error) {
	var commentId int

	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	if _, err = tx.Exec(fmt.Sprintf("SET LOCAL myapp.user_id = %d", userId)); err != nil {
		return 0, err
	}

	// Запрос для вставки комментария
	query := fmt.Sprintf("INSERT INTO %s (task_id, user_id, content) VALUES ($1, $2, $3) RETURNING comment_id", CommentTable)
	err = tx.QueryRow(query, taskId, userId, input.Comment).Scan(&commentId)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return commentId, nil
}

func (s *CommentPostgres) GetComments(taskId int) ([]taskFlow.Comment, error) {
	var comments []taskFlow.Comment
	query := fmt.Sprintf(`
		SELECT c.comment_id, c.task_id, c.user_id, c.content, c.created_at, u.username 
		FROM %s c 
		JOIN %s u ON c.user_id = u.user_id 
		WHERE c.task_id = $1`, CommentTable, UserTable)
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
	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	if _, err = tx.Exec(fmt.Sprintf("SET LOCAL myapp.user_id = %d", userId)); err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE comment_id = $1 AND user_id = $2", CommentTable)
	_, err = tx.Exec(query, commentId, userId)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *CommentPostgres) UpdateComment(commentId, userId int, input taskFlow.CommentInput) error {
	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Обеспечиваем откат в случае ошибки
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	// Устанавливаем локальный user_id
	if _, err = tx.Exec(fmt.Sprintf("SET LOCAL myapp.user_id = %d", userId)); err != nil {
		return err
	}

	// Запрос для обновления комментария
	query := fmt.Sprintf("UPDATE %s SET content = $1 WHERE comment_id = $2 AND user_id = $3", CommentTable)
	_, err = tx.Exec(query, input.Comment, commentId, userId)
	if err != nil {
		return err
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
