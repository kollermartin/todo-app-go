package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"todo-app/internal/adapter/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type DB struct {
	SqlDB *sql.DB
}

func New(ctx context.Context, config *config.Db) (*DB, error) {
	cnt := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Name)

	db, err := sql.Open(config.Type, cnt)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{
		db,
	}, nil
}

func (db *DB) Migrate() error {
	driver, err := postgres.WithInstance(db.SqlDB, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+ "./migrations",
		"postgres", driver,
	)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (db *DB) Close() {
	db.SqlDB.Close()
}