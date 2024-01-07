package router

import (
	"todo-app/app/controller"
	"todo-app/app/middlewares"
	"todo-app/app/service"

	"github.com/gin-gonic/gin"
)

func Init(todoService *service.TodoService) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middlewares.LoggerMiddleware())

	router.GET("/todos", controller.GetTodos(todoService))
	router.POST("/todos", controller.CreateTodo(todoService))
	router.GET("/todos/:id", controller.GetTodoByID(todoService))
	router.PUT("/todos/:id", controller.UpdateTodo(todoService))
	router.DELETE("/todos/:id", controller.DeleteTodo(todoService))

	return router
}
