package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"todo-app/types"
	"todo-app/api"

    _ "github.com/lib/pq"
)

func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime:= time.Now()
		
		c.Next()

		logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path": c.Request.RequestURI,
			"status": c.Writer.Status(),
			"latency": time.Since(startTime),
			"ip": c.ClientIP(),
		}).Info("Handled request")
	}
}

func initializeRouter(logger *logrus.Logger) *gin.Engine {
	router:= gin.New()
	
	router.Use(gin.Recovery())
	router.Use(LoggerMiddleware(logger))

	return router;
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

func initializeLogger() *logrus.Logger {
	logger := logrus.New();
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Level = logrus.InfoLevel

	return logger
}

func initializeDB(config types.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)
	db, err := sql.Open(config.DBType, connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Example table creation
	createTableSQL := `CREATE TABLE IF NOT EXISTS todos (
		id UUID PRIMARY KEY,
		title TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	logger := initializeLogger()
	router := initializeRouter(logger)

	config, err := loadConfig(".")

	if err != nil {
		logger.Fatal("cannot load config:", err)
	}

	db, err := initializeDB(config)

	if err != nil {
		logger.Fatal("cannot initialize db:", err)
	}

    defer db.Close()


	router.GET("/todos", api.GetTodos(db))
	router.POST("/todos", api.PostTodo(db))
	router.GET("/todos/:id", api.GetTodoByID(db))
	router.PUT("/todos/:id", api.UpdateTodo(db))

	router.Run("localhost:8080")
	fmt.Println(config)
}
