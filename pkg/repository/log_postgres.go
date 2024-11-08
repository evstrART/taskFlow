package repository

import (
	"github.com/jmoiron/sqlx"
)

type LogPostgres struct {
	db *sqlx.DB
}

func NewLogPostgres(db *sqlx.DB) *LogPostgres {
	return &LogPostgres{db: db}
}
