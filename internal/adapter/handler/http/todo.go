package http

import (
	"todo-app/internal/core/port"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TodoHandler struct {
	svc port.TodoService
}

func NewTodoHandler(svc port.TodoService) *TodoHandler {
	return &TodoHandler{svc}
}

func (th *TodoHandler) GetAllTodos(ctx *gin.Context) {
	todos, err := th.svc.GetAllTodos(ctx)
	if err != nil {
		// TODO handle error
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": todos})
}

func (th *TodoHandler) GetTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	// TODO Better error handling
	if id == "" {
		ctx.JSON(400, gin.H{"error": "id is required"})
		return
	}
	// TODO Better error handling
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	todo, err := th.svc.GetTodo(ctx, parsedUUID)
	if err != nil {
		// TODO handle error
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": todo})
}