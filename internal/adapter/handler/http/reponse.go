package http

import (
	"net/http"
	"time"
	"todo-app/internal/domain/entity"
	"todo-app/internal/domain/errors"

	"github.com/gin-gonic/gin"
)

var errorStatusMap = map[error]int{
	errors.ErrInternal: http.StatusInternalServerError,
	errors.ErrTicketNotFound: http.StatusNotFound,
}

type TodoResponse struct {
	ID        string `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewTodoResponse(todo *entity.Todo) TodoResponse {
	return TodoResponse{
		ID:        todo.UUID.String(),
		Title:     todo.Title,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}
}

func HandleValidationError(ctx *gin.Context, errMsg string) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errMsg})
}

func HandleError(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]

	if !ok {
		statusCode = http.StatusInternalServerError
	}

	//TODO create UI errors

	switch err {
		case errors.ErrTicketNotFound:
			ctx.AbortWithStatusJSON(statusCode, gin.H{"message": "Not Found"})
			return
		case errors.ErrInternal:
			ctx.AbortWithStatusJSON(statusCode, gin.H{"message": "Internal Server Error"})
			return
		default:
			ctx.AbortWithStatusJSON(statusCode, gin.H{"message": "Internal Server Error"})

	}
}
