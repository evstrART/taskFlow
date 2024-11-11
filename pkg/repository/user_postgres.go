package repository

import (
	"fmt"
	"github.com/evstrART/taskFlow"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetUser(userId int) (taskFlow.User, error) {
	var user taskFlow.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", UserTable)
	err := r.db.Get(&user, query, userId)
	return user, err
}

func (r *UserPostgres) GetUsers() ([]taskFlow.User, error) {
	var users []taskFlow.User
	query := fmt.Sprintf("SELECT * FROM %s", UserTable)
	err := r.db.Select(&users, query)
	return users, err
}
