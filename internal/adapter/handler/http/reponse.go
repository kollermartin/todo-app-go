package http

import (
	"net/http"
	"time"
	"todo-app/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var errorStatusMap = map[error]int{
	domain.ErrInternal: http.StatusInternalServerError,
	domain.ErrTicketNotFound: http.StatusNotFound,
}

type TodoResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewTodoResponse(todo *domain.Todo) TodoResponse {
	return TodoResponse{
		ID:        todo.UUID,
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
		case domain.ErrTicketNotFound:
			ctx.AbortWithStatusJSON(statusCode, gin.H{"message": "Not Found"})
			return
		case domain.ErrInternal:
			ctx.AbortWithStatusJSON(statusCode, gin.H{"message": "Internal Server Error"})
			return
		default:
			ctx.AbortWithStatusJSON(statusCode, gin.H{"message": "Internal Server Error"})

	}
}
