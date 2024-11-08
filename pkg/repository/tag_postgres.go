package repository

import (
	"github.com/jmoiron/sqlx"
)

type TagPostgres struct {
	db *sqlx.DB
}

func NewTafPostgres(db *sqlx.DB) *TagPostgres {
	return &TagPostgres{db: db}
}
