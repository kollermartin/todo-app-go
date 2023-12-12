package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Todo struct {
	ID string `json:"id"`
	Title string `json:"title"`
	CreatedAt string `json:"created_at"`
}

type TodoInput struct {
	Title string `json:"title"`
}

var todos []Todo = []Todo{
	{ID: "155215as4", Title: "Belajar Golang", CreatedAt: "2021-01-01"},
	{ID: "5a585fsg", Title: "Belajar Gin", CreatedAt: "2021-01-02"},
	{ID: "8ag5fsss5", Title: "Belajar Gorm", CreatedAt: "2021-01-03"},
	{ID: "5a5g5f5", Title: "Belajar Golannga", CreatedAt: "2021-01-06"},
	{ID: "5a5g5fdas5", Title: "Frd Golannga", CreatedAt: "2021-05-04"},
}

func getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)
}

func postTodo(c *gin.Context) {
	var input TodoInput

	if err:= c.BindJSON(&input); err != nil {
		return
	}

	newTodo:= Todo{
		ID: uuid.New().String(),
		Title: input.Title,
		CreatedAt: time.Now().Format("2006-01-02"),
	}

	todos = append(todos, newTodo)
	c.IndentedJSON(http.StatusCreated, newTodo)
}

func getTodoByID(c *gin.Context) {
	id:= c.Param("id")

	for _, todo:= range todos {
		if todo.ID == id {
			c.IndentedJSON(http.StatusOK, todo)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
}

func updateTodo(c *gin.Context) {
	var input TodoInput
	id:= c.Param("id")

	if err:= c.BindJSON(&input); err != nil {
		return
	}

	for i, todo:= range todos {
		if todo.ID == id {
			todos[i].Title = input.Title
			c.IndentedJSON(http.StatusOK, todos[i])
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
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

func initializeHandlers(router *gin.Engine) {
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodoByID)
	router.POST("/todos", postTodo)
	router.PUT("/todos/:id", updateTodo)

	router.Run("localhost:8080")
}

func initializeConfig() {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("json") 

	viper.ReadInConfig();

}

func main() {
	logger := logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Level = logrus.InfoLevel

	router:= gin.New()
	router.Use(gin.Recovery())
	router.Use(LoggerMiddleware(logger))

	initializeHandlers(router)
	initializeConfig()
}
