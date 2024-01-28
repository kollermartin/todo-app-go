package port

import (
	"context"
	"todo-app/internal/core/domain"

	"github.com/google/uuid"
)

type TodoRepository interface {
	GetAllTodos(ctx context.Context)  ([]domain.Todo, error)
	GetTodo(ctx context.Context, uuid uuid.UUID ) (*domain.Todo, error)
	CreateTodo (ctx context.Context, todo *domain.Todo) (*domain.Todo, error)
	UpdateTodo (ctx context.Context, todo *domain.Todo) (*domain.Todo, error)
	DeleteTodo (ctx context.Context, uuid uuid.UUID ) error
}

type TodoService interface {
	GetAllTodos(ctx context.Context)  ([]domain.Todo, error)
	GetTodo(ctx context.Context, uuid uuid.UUID ) (*domain.Todo, error)
	CreateTodo (ctx context.Context, todo *domain.Todo) (*domain.Todo, error)
	UpdateTodo (ctx context.Context, todo *domain.Todo) (*domain.Todo, error)
	DeleteTodo (ctx context.Context, uuid uuid.UUID ) error
}