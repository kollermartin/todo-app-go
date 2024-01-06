package main

import (
	"todo-app/config"
	"todo-app/app/router"
	"todo-app/app/service"
	"todo-app/app/types"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

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

func main() {
	logger := initLogger()

	env, err := loadConfig(".")

	if err != nil {
		logger.WithFields(logrus.Fields{
			"event": "config_load_fail",
			"error": err.Error(),
		}).Fatal("Failed to load config")
	}

	db := config.ConnectToDB(env, logger)
	todoService := service.NewTodoService(db)

	defer db.Close()

	// TODO Zbavit se zavislosti db, logger
	router := router.Init(todoService, logger)

	if err := router.Run("localhost:8080"); err != nil {
		logger.WithFields(logrus.Fields{
			"event": "server_run_fail",
			"error": err.Error(),
		}).Fatal("Failed to run server")
	}
}
