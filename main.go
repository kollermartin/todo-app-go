package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/lib/pq"


	"todo-app/api"
	"todo-app/middlewares"
	"todo-app/types"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func initRouter(logger *logrus.Logger) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middlewares.LoggerMiddleware(logger))

	return router
}

func loadConfig(path string) (config types.Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err = viper.ReadInConfig(); err != nil {
		return types.Config{}, err
	}

	if err = viper.Unmarshal(&config); err != nil {
		return types.Config{}, err
	}

	return config, err
}

func initLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Level = logrus.InfoLevel

	return logger
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

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}

func initDB(config types.Config, cString string) (*sql.DB, error) {
	db, err := sql.Open(config.DBType, cString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if err := runMigrations(db, config.MigrationsPath); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	logger := initLogger()
	router := initRouter(logger)

	config, err := loadConfig(".")

	if err != nil {
		logger.WithFields(logrus.Fields{
			"event": "config_load_fail",
			"error": err.Error(),
		}).Fatal("Failed to load config")
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)

	db, err := initDB(config, connStr)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"event": "db_init_fail",
			"error": err.Error(),
		}).Fatal("Failed to initialize database")
	}

	defer db.Close()

	router.GET("/todos", api.GetTodos(db, logger))
	router.POST("/todos", api.PostTodo(db, logger))
	router.GET("/todos/:id", api.GetTodoByID(db, logger))
	router.PUT("/todos/:id", api.UpdateTodo(db, logger))

	if err := router.Run("localhost:8080"); err != nil {
		logger.WithFields(logrus.Fields{
			"event": "server_run_fail",
			"error": err.Error(),
		}).Fatal("Failed to run server")
	}
}
