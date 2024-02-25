package http

import (
	"net/http"
	"todo-app/internal/domain/entity"
	"todo-app/internal/domain/port"

	// "todo-app/internal/domain/vo"

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
	todoListRsp := []TodoResponse{}

	todos, err := th.svc.GetAllTodos(ctx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	for _, todo := range todos {
		todoListRsp = append(todoListRsp, NewTodoResponse(&todo))
	}

	ctx.JSON(http.StatusOK, todoListRsp)
}

func (th *TodoHandler) GetTodo(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		HandleValidationError(ctx, "ID is required")
		return
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		HandleValidationError(ctx, "Invalid ID")
		return
	}

	todo, err := th.svc.GetTodo(ctx, parsedUUID)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	rsp := NewTodoResponse(todo)

	ctx.JSON(http.StatusOK, rsp)
}

func (th *TodoHandler) CreateTodo(ctx *gin.Context) {
	var req CreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleValidationError(ctx, err.Error())
		return
	}

	todo := entity.Todo{
		Title: req.Title,
	}

	createdTodo, err := th.svc.CreateTodo(ctx, &todo)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	rsp := NewTodoResponse(createdTodo)

	ctx.JSON(http.StatusCreated, rsp)
}

func (th *TodoHandler) UpdateTodo(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		HandleValidationError(ctx, "ID is required")
		return
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		HandleValidationError(ctx, "Invalid ID")
		return
	}

	var req UpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleValidationError(ctx, err.Error())
		return
	}

	todo := entity.Todo{
		UUID:  parsedUUID,
		Title: req.Title,
	}

	updatedTodo, err := th.svc.UpdateTodo(ctx, &todo)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	rsp := NewTodoResponse(updatedTodo)

	ctx.JSON(http.StatusOK, rsp)
}

func (th *TodoHandler) DeleteTodo(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		HandleValidationError(ctx, "ID is required")
		return
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		HandleValidationError(ctx, "Invalid ID")
		return
	}

	err = th.svc.DeleteTodo(ctx, parsedUUID)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
