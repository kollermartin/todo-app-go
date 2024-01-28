package http

import (
	"time"
	"todo-app/internal/core/domain"

	"github.com/google/uuid"
)

type TodoResponse struct {
	ID uuid.UUID `json:"id"`
	Title string `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func NewTodoResponse(todo *domain.Todo) TodoResponse {
	return TodoResponse{
		ID: todo.UUID,
		Title: todo.Title,
		CreatedAt: todo.CreatedAt,
	}
}