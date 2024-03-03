package repo

import (
	"context"
	"database/sql"
	"todo-app/internal/domain/entity"
	"todo-app/internal/domain/errors"
	"todo-app/internal/infrastructure/postgre"

	// "todo-app/internal/domain/vo"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type TodoRepository struct {
	db *postgre.DB
}

func NewTodoRepository(db *postgre.DB) *TodoRepository {
	return &TodoRepository{db}
}

func (tr *TodoRepository) GetAllTodos(ctx context.Context) ([]entity.Todo, error) {
	var todo entity.Todo
	var todos []entity.Todo

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

func (tr *TodoRepository) GetTodo(ctx context.Context, uuid uuid.UUID) (*entity.Todo, error) {
	var todo entity.Todo

	err := tr.db.SqlDB.QueryRow("SELECT * from todos where uuid = $1", uuid).Scan(&todo.ID, &todo.UUID, &todo.Title, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrTicketNotFound
		}

		return nil, err
	}

	return &todo, nil
}

func (tr *TodoRepository) CreateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	newTodo := entity.Todo{
		Title: todo.Title,
	}

	err := tr.db.SqlDB.QueryRow("INSERT INTO todos (title) VALUES ($1) RETURNING id, uuid, title, created_at, updated_at", newTodo.Title).Scan(&newTodo.ID, &newTodo.UUID, &newTodo.Title, &newTodo.CreatedAt, &newTodo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &newTodo, nil
}

func (tr *TodoRepository) UpdateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	var updatedTodo entity.Todo

	err := tr.db.SqlDB.QueryRow("UPDATE todos SET title = $1, updated_at = now() WHERE uuid = $2 RETURNING id, uuid, title, created_at, updated_at", todo.Title, todo.UUID).Scan(&updatedTodo.ID, &updatedTodo.UUID, &updatedTodo.Title, &updatedTodo.CreatedAt, &updatedTodo.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.ErrTicketNotFound
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
		return errors.ErrTicketNotFound
	}

	return nil
}
