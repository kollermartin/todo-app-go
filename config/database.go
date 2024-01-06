package config

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/sirupsen/logrus"

	"todo-app/app/types"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func ConnectToDB(config types.Config, log *logrus.Logger) *sql.DB {
	cnt := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)

	db, err := sql.Open(config.DBType, cnt)
	if err != nil {
		log.WithFields(logrus.Fields{
			"event": "db_init_fail",
			"error": err.Error(),
		}).Fatal("Failed to initialize database")
	}

	if err = db.Ping(); err != nil {
		log.WithFields(logrus.Fields{
			"event": "db_init_fail",
			"error": err.Error(),
		}).Fatal("Failed to ping database")
	}

	if err := runMigrations(db, config.MigrationsPath); err != nil {
		log.WithFields(logrus.Fields{
			"event": "db_migrations_fail",
			"error": err.Error(),
		}).Fatal("Failed to run migrations")
	}

	return db
}

func runMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
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
