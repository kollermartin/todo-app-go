package main

import (
	"context"
	"fmt"
	"todo-app/config"
	"todo-app/internal/application/todo"
	"todo-app/internal/infrastructure/postgre"
	"todo-app/internal/infrastructure/postgre/repo"
	"todo-app/internal/ui/http"
	"todo-app/pkg/logger"

	"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	config, err := config.New()
	if err != nil {
		logrus.Fatal("Error loading config: ", err)
	}

	logger.Set()

	logrus.Info("Starting the application", " app: ", config.App.Name, " env: ", config.App.Env)

	ctx := context.Background()

	db, err := postgre.New(ctx, config.Db)
	if err != nil {
		logrus.Fatal("Error initializing database", err)
	}

	defer db.Close()

	logrus.Info("Successfully connected to database: ", config.Db.Type)

	err = db.Migrate(config.App)
	if err != nil {
		logrus.Fatal("Error running migrations", err)
	}

	logrus.Info("Successfully ran migrations")

	todoRepo := repo.NewTodoRepository(db)
	todoHandler := todo.NewTodoHandler(todoRepo)

	router, err := http.NewRouter(todoHandler)

	if err != nil {
		logrus.Fatal("Error initializing router", err)
	}

	listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)
	logrus.Info("Starting the HTTP server", listenAddr)

	err = router.Run(listenAddr)
	if err != nil {
		logrus.Fatal("Error starting the HTTP server", err)
	}
}
