package http

import (
	"todo-app/internal/application/todo"
	"todo-app/internal/ui/http/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
}

func NewRouter(
	todoHandler *todo.TodoHandler,
) (*Router, error) {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())

	router.GET("/todos", todoHandler.GetAllTodos)
	router.POST("/todos", todoHandler.CreateTodo)
	router.GET("/todos/:id", todoHandler.GetTodo)
	router.PUT("/todos/:id", todoHandler.UpdateTodo)
	router.DELETE("/todos/:id", todoHandler.DeleteTodo)

	return &Router{router}, nil
}

func (r *Router) Run(address string) error {
	return r.Engine.Run(address)
}
