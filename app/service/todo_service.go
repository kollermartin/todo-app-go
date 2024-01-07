package service

import (
	"database/sql"
	"fmt"
	"time"

	"todo-app/app/types"

	"github.com/google/uuid"
)

type TodoService struct {
	DB *sql.DB
}

type TodoError struct {
	Message string
	Reason  TodoErrorReason
}

type TodoErrorReason int

const (
	ReasonNotFound TodoErrorReason = iota
	ReasonUnknown
)

const ErrMsgInternalServer = "Internal server error"

func (e TodoError) Error() string {
	return e.Message
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

		return nil, TodoError{Message: ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	defer rows.Close()

	for rows.Next() {
		var todo types.Todo
		if err := rows.Scan(&todo.ID, &todo.ExternalID, &todo.Title, &todo.CreatedAt); err != nil {

			return nil, TodoError{Message: ErrMsgInternalServer, Reason: ReasonUnknown}
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (service *TodoService) GetTodoByID(id string) (*types.Todo, error) {
	var todo types.Todo

	err := service.DB.QueryRow("SELECT * from todos where external_id = $1", id).Scan(&todo.ID, &todo.ExternalID, &todo.Title, &todo.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
            return nil, TodoError{Message: fmt.Sprintf("Todo with id %s not found", id), Reason: ReasonNotFound}
        }
        return nil, TodoError{Message: ErrMsgInternalServer, Reason: ReasonUnknown}
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
		return nil, TodoError{Message: ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	return &newTodo, nil
}

func (service *TodoService) UpdateTodo(id string, title string) (*types.Todo, error) {
	var updatedTodo types.Todo

	result, err := service.DB.Exec("UPDATE todos SET title = $1 WHERE external_id = $2", title, id)

	if err != nil {
		return nil, TodoError{Message: ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, TodoError{Message: ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	if rowsAffected == 0 {
		return nil, TodoError{Message: fmt.Sprintf("Todo with id %s not found", id), Reason: ReasonNotFound}
	}

	err = service.DB.QueryRow("SELECT * from todos where external_id = $1", id).Scan(&updatedTodo.ID, &updatedTodo.ExternalID, &updatedTodo.Title, &updatedTodo.CreatedAt)
	if err != nil {
		return nil, TodoError{Message: ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	return &updatedTodo, nil
}

func (service *TodoService) DeleteTodo(id string) error {
	result, err := service.DB.Exec("DELETE FROM todos WHERE external_id = $1", id)

	if err != nil {
		return TodoError{Message: ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return TodoError{Message: ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	if rowsAffected == 0 {
		return TodoError{Message: fmt.Sprintf("Todo with id %s not found", id), Reason: ReasonNotFound}
	}

	return nil
}
