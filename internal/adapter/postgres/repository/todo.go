package repository

import (
	"context"
	"database/sql"
	"todo-app/internal/adapter/postgres"
	"todo-app/internal/core/domain"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type TodoRepository struct {
	db *postgres.DB
}

func NewTodoRepository(db *postgres.DB) *TodoRepository {
	return &TodoRepository{db}
}

func (tr *TodoRepository) GetAllTodos(ctx context.Context) ([]domain.Todo, error) {
	var todo domain.Todo
	var todos []domain.Todo

	rows, err := tr.db.SqlDB.Query("SELECT * FROM todos")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&todo.ID, &todo.UUID, &todo.Title, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (tr *TodoRepository) GetTodo(ctx context.Context, uuid uuid.UUID) (*domain.Todo, error) {
	var todo domain.Todo

	err := tr.db.SqlDB.QueryRow("SELECT * from todos where uuid = $1", uuid).Scan(&todo.ID, &todo.UUID, &todo.Title, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &todo, nil
}

func (tr *TodoRepository) CreateTodo(ctx context.Context, todo *domain.Todo) (*domain.Todo, error) {
	newTodo := domain.Todo{
		Title: todo.Title,
	}

	err := tr.db.SqlDB.QueryRow("INSERT INTO todos (title) VALUES ($1) RETURNING id, uuid, title, created_at, updated_at", newTodo.Title).Scan(&newTodo.ID, &newTodo.UUID, &newTodo.Title, &newTodo.CreatedAt, &newTodo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &newTodo, nil
}

func (tr *TodoRepository) UpdateTodo(ctx context.Context, todo *domain.Todo) (*domain.Todo, error) {
	var updatedTodo domain.Todo

	err := tr.db.SqlDB.QueryRow("UPDATE todos SET title = $1, updated_at = now() WHERE uuid = $2 RETURNING id, uuid, title, created_at, updated_at", todo.Title, todo.UUID).Scan(&updatedTodo.ID, &updatedTodo.UUID, &updatedTodo.Title, &updatedTodo.CreatedAt, &updatedTodo.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &updatedTodo, nil
}

func (tr *TodoRepository) DeleteTodo(ctx context.Context, uuid uuid.UUID) error {
	result, err := tr.db.SqlDB.Exec("DELETE FROM todos WHERE uuid = $1", uuid)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
