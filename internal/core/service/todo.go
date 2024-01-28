package service

import (
	"context"
	"todo-app/internal/core/domain"
	"todo-app/internal/core/port"

	"github.com/google/uuid"
)

type TodoService struct {
	repo port.TodoRepository
}

func NewTodoService(repo port.TodoRepository) *TodoService {
	return &TodoService{repo}
}

func (ts *TodoService) GetAllTodos(ctx context.Context) ([]domain.Todo, error) {
	todos, err := ts.repo.GetAllTodos(ctx)
	if (err != nil) {
		return nil, domain.ErrInternal
	}
	
	return todos, nil
}

func (ts *TodoService) GetTodo(ctx context.Context, uuid uuid.UUID) (*domain.Todo, error) {
	todo, err := ts.repo.GetTodo(ctx, uuid)
	if (err != nil) {
		return nil, domain.ErrInternal
	}
	
	return todo, nil
}

func (ts *TodoService) CreateTodo (ctx context.Context, todo *domain.Todo) (*domain.Todo, error) {
	todo.UUID = uuid.New()
	
	// TODO CreateTodo db should have set default for created_at and updated_at
	todo, err := ts.repo.CreateTodo(ctx, todo)
	if (err != nil) {
		return nil, domain.ErrInternal
	}
	
	return todo, nil
}

func (ts *TodoService) UpdateTodo (ctx context.Context, todo *domain.Todo) (*domain.Todo, error) {
	todo, err := ts.repo.UpdateTodo (ctx, todo)
	if (err != nil) {
		return nil, domain.ErrInternal
	}
	
	return todo, nil
}

func (ts *TodoService) DeleteTodo (ctx context.Context, uuid uuid.UUID) error {
	err := ts.repo.DeleteTodo (ctx, uuid)
	if (err != nil) {
		return domain.ErrInternal
	}
	
	return nil
}