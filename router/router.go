package router

import (
	"database/sql"
	"todo-app/api"
	"todo-app/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Init(db *sql.DB, log *logrus.Logger) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middlewares.LoggerMiddleware(log))

	router.GET("/todos", api.GetTodos(db, log))
	router.POST("/todos", api.CreateTodo(db, log))
	router.GET("/todos/:id", api.GetTodoByID(db, log))
	router.PUT("/todos/:id", api.UpdateTodo(db, log))
	router.DELETE("/todos/:id", api.DeleteTodo(db, log))

	return router
}