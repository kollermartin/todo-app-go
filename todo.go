package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
    _ "github.com/lib/pq"
)

type Todo struct {
	ID string `json:"id"`
	Title string `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type TodoInput struct {
	Title string `json:"title"`
}

type Config struct {
    Env      string `mapstructure:"ENV"`
	Port     string `mapstructure:"PORT"`
	DBType   string `mapstructure:"DB_TYPE"`
	DBHost   string `mapstructure:"DB_HOST"`
	DBPort   int `mapstructure:"DB_PORT"`
	DBUser   string `mapstructure:"DB_USER"`
	DBPass   string `mapstructure:"DB_PASS"`
	DBName   string `mapstructure:"DB_NAME"`
}

func getTodos(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, title, created_at FROM todos")

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		defer rows.Close()

		var todos []Todo = []Todo{}
		for rows.Next() {
			var todo Todo
			if err := rows.Scan(&todo.ID, &todo.Title, &todo.CreatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}

			todos = append(todos, todo)
		}

		c.IndentedJSON(http.StatusOK, todos)
	}
}

func postTodo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input TodoInput

		if err:= c.BindJSON(&input); err != nil {
			return
		}

		newTodo:= Todo{
			ID: uuid.New().String(),
			Title: input.Title,
			CreatedAt: time.Now(),
		}

		fmt.Println(newTodo.CreatedAt)

	_, err := db.Exec("INSERT INTO todos (id, title, created_at) VALUES ($1, $2, $3)", newTodo.ID, newTodo.Title, newTodo.CreatedAt)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newTodo)

	}
}

func getTodoByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id");

		var todo Todo

		err := db.QueryRow("SELECT * from todos where id = $1", id).Scan(&todo.ID, &todo.Title, &todo.CreatedAt)


		if err != nil {
			if err == sql.ErrNoRows {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}
			return
		}

		c.IndentedJSON(http.StatusOK, todo)
	}
}

func updateTodo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todoInput TodoInput;
		var updatedTodo Todo;

		id := c.Param("id")

		if err:= c.BindJSON(&todoInput); err != nil {
			return
		}

		result, err:= db.Exec("UPDATE todos SET title = $1 WHERE id = $2", todoInput.Title, id)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if rowsAffected == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
			return
		}

		err = db.QueryRow("SELECT * from todos where id = $1", id).Scan(&updatedTodo.ID, &updatedTodo.Title, &updatedTodo.CreatedAt)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, updatedTodo)
	}
}

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

func loadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env") 

	if err = viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	if err = viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, err
}

func initializeLogger() *logrus.Logger {
	logger := logrus.New();
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Level = logrus.InfoLevel

	return logger
}

func initializeDB(config Config) (*sql.DB, error) {
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


	router.GET("/todos", getTodos(db))
	router.POST("/todos", postTodo(db))
	router.GET("/todos/:id", getTodoByID(db))
	router.PUT("/todos/:id", updateTodo(db))

	router.Run("localhost:8080")
	fmt.Println(config)
}
