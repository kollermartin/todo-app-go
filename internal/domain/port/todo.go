package port

import (
	"context"
	"todo-app/internal/domain/entity"

	"github.com/google/uuid"
	// "todo-app/internal/domain/vo"
)

//TODO Refactor this to usecases/operations
type TodoRepository interface {
	GetAllTodos(ctx context.Context) ([]entity.Todo, error)
	GetTodo(ctx context.Context, uuid uuid.UUID) (*entity.Todo, error)
	CreateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
	UpdateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
	DeleteTodo(ctx context.Context, uuid uuid.UUID) error
}

type TodoService interface {
	GetAllTodos(ctx context.Context) ([]entity.Todo, error)
	GetTodo(ctx context.Context, uuid uuid.UUID) (*entity.Todo, error)
	CreateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
	UpdateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
	DeleteTodo(ctx context.Context, uuid uuid.UUID) error
}
