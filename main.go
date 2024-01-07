package main

import (
	"todo-app/config"
	"todo-app/app/router"
	"todo-app/app/service"

	"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	config.InitLogger()

	env, err := config.LoadConfig(".")

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"event": "config_load_fail",
			"error": err.Error(),
		}).Fatal("Failed to load config")
	}

	db := config.ConnectToDB(env)

	defer db.Close()

	todoService := service.NewTodoService(db)

	router := router.Init(todoService)

	if err := router.Run("localhost:8080"); err != nil {
		logrus.WithFields(logrus.Fields{
			"event": "server_run_fail",
			"error": err.Error(),
		}).Fatal("Failed to run server")
	}
}
