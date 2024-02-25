package service

import (
	"context"
	"todo-app/internal/domain/entity"
	"todo-app/internal/domain/errors"
	"todo-app/internal/domain/port"

	"github.com/google/uuid"
)
//TODO refactor to usecases/operations
type TodoService struct {
	repo port.TodoRepository
}

func NewTodoService(repo port.TodoRepository) *TodoService {
	return &TodoService{repo}
}

func (ts *TodoService) GetAllTodos(ctx context.Context) ([]entity.Todo, error) {
	todos, err := ts.repo.GetAllTodos(ctx)
	if err != nil {
		return nil, handleErrorSelection(err)
	}

	return todos, nil
}

func (ts *TodoService) GetTodo(ctx context.Context, uuid uuid.UUID) (*entity.Todo, error) {
	todo, err := ts.repo.GetTodo(ctx, uuid)
	if err != nil {
		return nil, handleErrorSelection(err)
	}

	return todo, nil
}

func (ts *TodoService) CreateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	todo.UUID = uuid.New()

	todo, err := ts.repo.CreateTodo(ctx, todo)
	if err != nil {
		return nil, handleErrorSelection(err)
	}

	return todo, nil
}

func (ts *TodoService) UpdateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	todo, err := ts.repo.UpdateTodo(ctx, todo)
	if err != nil {
		return nil, handleErrorSelection(err)
	}

	return todo, nil
}

func (ts *TodoService) DeleteTodo(ctx context.Context, uuid uuid.UUID) error {
	err := ts.repo.DeleteTodo(ctx, uuid)
	if err != nil {
		return handleErrorSelection(err)
	}

	return nil
}

func handleErrorSelection(err error) error {
	if err == errors.ErrTicketNotFound {
		return errors.ErrTicketNotFound
	}

	return errors.ErrInternal
}
