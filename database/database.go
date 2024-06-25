package database

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
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

func initTestDB(t *testing.T) (Database, func()) {
	db, err := NewDatabase("postgresql://aweayo:aweayo@localhost:5432/blogging_api_test?sslmode=disable")
	require.NoError(t, err)

	closeFn := func() {
		_, err := db.GetDB().Exec("TRUNCATE TABLE articles;")
		require.NoError(t, err)
		db.GetDB().Close()
	}
	return db, closeFn
}
