package router

import (
	"todo-app/app/controller"
	"todo-app/app/middlewares"
	"todo-app/app/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Init(todoService *service.TodoService, log *logrus.Logger) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middlewares.LoggerMiddleware(log))

	router.GET("/todos", controller.GetTodos(todoService, log))
	router.POST("/todos", controller.CreateTodo(todoService, log))
	router.GET("/todos/:id", controller.GetTodoByID(todoService, log))
	router.PUT("/todos/:id", controller.UpdateTodo(todoService, log))
	router.DELETE("/todos/:id", controller.DeleteTodo(todoService, log))

	return router
}
