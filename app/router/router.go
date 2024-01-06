package router

import (
	"todo-app/app/api"
	"todo-app/app/middlewares"
	"todo-app/app/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Init(todoService *service.TodoService, log *logrus.Logger) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middlewares.LoggerMiddleware(log))

	router.GET("/todos", api.GetTodos(todoService, log))
	router.POST("/todos", api.CreateTodo(todoService, log))
	router.GET("/todos/:id", api.GetTodoByID(todoService, log))
	router.PUT("/todos/:id", api.UpdateTodo(todoService, log))
	router.DELETE("/todos/:id", api.DeleteTodo(todoService, log))

	return router
}
