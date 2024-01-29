package http

import "github.com/gin-gonic/gin"

type Router struct {
	Engine *gin.Engine
}

func NewRouter(
	todoHandler *TodoHandler,
) (*Router, error) {
	router := gin.New()

	router.Use(gin.Recovery())
	//TODO add logger middleware
	// router.Use(LoggerMiddleware())

	router.GET("/todos", todoHandler.GetAllTodos)
	router.POST("/todos", todoHandler.CreateTodo)
	router.GET("/todos/:id", todoHandler.GetTodo)
	router.PUT("/todos/:id", todoHandler.UpdateTodo)
	router.DELETE("/todos/:id", todoHandler.DeleteTodo)

	return &Router{router}, nil
}