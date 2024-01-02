package service

import (
	"database/sql"
	"time"
	"todo-app/types"

	"github.com/google/uuid"
)

type TodoService struct {
	DB *sql.DB
}

func NewTodoService(db *sql.DB) *TodoService {
	return &TodoService{
		DB: db,
	}
}

func (service *TodoService) GetAllTodos() ([]types.Todo, error) {
	var todos []types.Todo

	rows, err := service.DB.Query("SELECT * FROM todos")
	if err != nil {

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var todo types.Todo
		if err := rows.Scan(&todo.ID, &todo.ExternalID, &todo.Title, &todo.CreatedAt); err != nil {

			return nil, err
		}

		todos = append(todos, *&todo)
	}

	return todos, nil
}

func (service *TodoService) GetTodoByID(id string) (*types.Todo, error) {
	var todo types.Todo

	err := service.DB.QueryRow("SELECT * from todos where external_id = $1", id).Scan(&todo.ID, &todo.ExternalID, &todo.Title, &todo.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (service *TodoService) CreateTodo(title string) (*types.Todo, error) {
	newTodo := types.Todo{
		ExternalID: uuid.New().String(),
		Title:      title,
		CreatedAt:  time.Now(),
	}

	err := service.DB.QueryRow("INSERT INTO todos (external_id, title, created_at) VALUES ($1, $2, $3) RETURNING id", newTodo.ExternalID, newTodo.Title, newTodo.CreatedAt).Scan(&newTodo.ID)

	if err != nil {
		return nil, err
	}

	return &newTodo, nil
}

func (service *TodoService) UpdateTodo(id string, title string) (*types.Todo, error) {
	var updatedTodo types.Todo

	result, err := service.DB.Exec("UPDATE todos SET title = $1 WHERE external_id = $2", title, id)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, err
	}

	err = service.DB.QueryRow("SELECT * from todos where external_id = $1", id).Scan(&updatedTodo.ID, &updatedTodo.ExternalID, &updatedTodo.Title, &updatedTodo.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &updatedTodo, nil
}

func (service *TodoService) DeleteTodo(id string) error {
	result, err := service.DB.Exec("DELETE FROM todos WHERE external_id = $1", id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return err
	}

	return nil
}
