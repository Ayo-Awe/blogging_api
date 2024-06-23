package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database interface {
	GetDB() *sqlx.DB
}

type database struct {
	db *sqlx.DB
}

func (d *database) GetDB() *sqlx.DB {
	return d.db
}

func NewDatabase(dsn string) (Database, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &database{db}, nil
}
