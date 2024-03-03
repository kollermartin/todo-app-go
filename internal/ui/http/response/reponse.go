package response

import (
	"time"
	"todo-app/internal/domain/entity"
)

type TodoResponse struct {
	ID        string    `json:"id"`
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
