package http

import (
	"net/http"
	"todo-app/internal/core/domain"
	"todo-app/internal/core/port"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TodoHandler struct {
	svc port.TodoService
}

type createRequest struct {
	Title string `json:"title" binding:"required"`
}

type updateRequest struct {
	Title string `json:"title" binding:"required"`
}

func NewTodoHandler(svc port.TodoService) *TodoHandler {
	return &TodoHandler{svc}
}

func (th *TodoHandler) GetAllTodos(ctx *gin.Context) {
	var todoListRsp []TodoResponse

	todos, err := th.svc.GetAllTodos(ctx)
	if err != nil {
		// TODO handle error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, todo := range todos {
		todoListRsp = append(todoListRsp, NewTodoResponse(&todo))
	}

	ctx.JSON(http.StatusOK, todoListRsp)
}

func (th *TodoHandler) GetTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	// TODO Better error handling
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	// TODO Better error handling
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	todo, err := th.svc.GetTodo(ctx, parsedUUID)
	if err != nil {
		// TODO handle error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rsp := NewTodoResponse(todo)

	ctx.JSON(http.StatusOK, rsp)
}

func (th *TodoHandler) CreateTodo(ctx *gin.Context) {
	var req createRequest

	// TODO Better error handling
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo := domain.Todo{
		Title: req.Title,
	}

	createdTodo, err := th.svc.CreateTodo(ctx, &todo)
	if err != nil {
		// TODO handle error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rsp := NewTodoResponse(createdTodo)

	ctx.JSON(http.StatusCreated, rsp)
}

func (th *TodoHandler) updateTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	// TODO Better error handling
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	// TODO Better error handling
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateRequest

	// TODO Better error handling
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo := domain.Todo{
		UUID:  parsedUUID,
		Title: req.Title,
	}

	updatedTodo, err := th.svc.UpdateTodo(ctx, &todo)
	if err != nil {
		// TODO handle error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rsp := NewTodoResponse(updatedTodo)

	ctx.JSON(http.StatusOK, rsp)
}

func (th *TodoHandler) deleteTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	// TODO Better error handling
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	// TODO Better error handling
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = th.svc.DeleteTodo (ctx, parsedUUID)
	if err != nil {
		// TODO handle error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}