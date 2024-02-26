package todo

import (
	"net/http"
	"todo-app/internal/domain/entity"
	"todo-app/internal/domain/port"
	"todo-app/internal/ui/http/request"
	"todo-app/internal/ui/http/response"

	// "todo-app/internal/domain/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TodoHandler struct {
	repo port.TodoRepository
}

func NewTodoHandler(repo port.TodoRepository) *TodoHandler {
	return &TodoHandler{repo}
}

func (th *TodoHandler) GetAllTodos(ctx *gin.Context) {
	todoListRsp := []response.TodoResponse{}

	todos, err := th.repo.GetAllTodos(ctx)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	for _, todo := range todos {
		todoListRsp = append(todoListRsp, response.NewTodoResponse(&todo))
	}

	ctx.JSON(http.StatusOK, todoListRsp)
}

func (th *TodoHandler) GetTodo(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		response.HandleValidationError(ctx, "ID is required")
		return
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		response.HandleValidationError(ctx, "Invalid ID")
		return
	}

	todo, err := th.repo.GetTodo(ctx, parsedUUID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewTodoResponse(todo)

	ctx.JSON(http.StatusOK, rsp)
}

func (th *TodoHandler) CreateTodo(ctx *gin.Context) {
	var req request.CreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(ctx, err.Error())
		return
	}

	todo := entity.Todo{
		Title: req.Title,
	}

	createdTodo, err := th.repo.CreateTodo(ctx, &todo)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewTodoResponse(createdTodo)

	ctx.JSON(http.StatusCreated, rsp)
}

func (th *TodoHandler) UpdateTodo(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		response.HandleValidationError(ctx, "ID is required")
		return
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		response.HandleValidationError(ctx, "Invalid ID")
		return
	}

	var req request.UpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(ctx, err.Error())
		return
	}

	todo := entity.Todo{
		UUID:  parsedUUID,
		Title: req.Title,
	}

	updatedTodo, err := th.repo.UpdateTodo(ctx, &todo)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewTodoResponse(updatedTodo)

	ctx.JSON(http.StatusOK, rsp)
}

func (th *TodoHandler) DeleteTodo(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		response.HandleValidationError(ctx, "ID is required")
		return
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		response.HandleValidationError(ctx, "Invalid ID")
		return
	}

	err = th.repo.DeleteTodo(ctx, parsedUUID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
