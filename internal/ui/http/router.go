package http

import (
	"todo-app/internal/application/todo"
	"todo-app/internal/ui/http/middleware"
	"todo-app/internal/ui/http/route"

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

	getAllTodosRoute := route.NewGetTodosRoute(todoHandler)
	getTodoRoute := route.NewGetTodoRoute(todoHandler)
	createTodoRoute := route.NewCreateTodoRoute(todoHandler)
	updateTodoRoute := route.NewUpdateTodoRoute(todoHandler)
	deleteTodoRoute := route.NewDeleteTodoRoute(todoHandler)

	router.GET(route.GetTodosPath, getAllTodosRoute)
	router.POST(route.CreateTodoPath, createTodoRoute)
	router.GET(route.GetTodoPath, getTodoRoute)
	router.PUT(route.UpdateTodoPath, updateTodoRoute)
	router.DELETE(route.DeleteTodoPath, deleteTodoRoute)

	return &Router{router}, nil
}

func (r *Router) Run(address string) error {
	return r.Engine.Run(address)
}
