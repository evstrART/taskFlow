package repository

import (
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
