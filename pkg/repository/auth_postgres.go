package repository

import (
	"database/sql"
	"fmt"
	"github.com/evstrART/taskFlow"
	"github.com/jmoiron/sqlx"
	"log"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user taskFlow.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password, email, role) values ($1, $2, $3, $4) RETURNING user_id", UserTable)
	row := r.db.QueryRow(query, user.Username, user.Password, user.Email, user.Role)
	if err := row.Scan(&id); err != nil {
		log.Printf("Error inserting user: %v", err)
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (taskFlow.User, error) {
	var user taskFlow.User
	query := fmt.Sprintf("SELECT user_id FROM %s WHERE username = $1 AND password = $2", UserTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}

func (r *AuthPostgres) ChangePassword(userId int, newPasswordHex string) error {
	query := fmt.Sprintf("UPDATE %s SET password = $1, updated_at = CURRENT_TIMESTAMP WHERE user_id = $2", UserTable)
	_, err := r.db.Exec(query, newPasswordHex, userId)
	return err
}
func (r *AuthPostgres) UserExistsForReset(input taskFlow.ResetPasswordInput) (bool, error) {
	var user taskFlow.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1 OR email = $2", UserTable)

	err := r.db.Get(&user, query, input.Username, input.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Пользователь не найден
		}
		return false, err // Возвращаем ошибку, если произошла другая ошибка
	}

	return true, nil // Пользователь найден
}

func (r *AuthPostgres) GetUserByNameAndEmail(username, email string) (taskFlow.User, error) {
	var user taskFlow.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1 AND email = $2", UserTable)
	err := r.db.Get(&user, query, username, email)
	return user, err
}
