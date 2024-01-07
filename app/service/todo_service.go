package service

import (
	"database/sql"
	"fmt"
	"time"

	"todo-app/app/constant"
	"todo-app/app/types"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
		logrus.WithFields(logrus.Fields{
			"event": constant.GetTodosLogEventErrorKey,
		}).Error(constant.DbQueryFailMsg)

		return nil, TodoError{Message: constant.ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	defer rows.Close()

	for rows.Next() {
		var todo types.Todo
		if err := rows.Scan(&todo.ID, &todo.ExternalID, &todo.Title, &todo.CreatedAt); err != nil {
			logrus.WithFields(logrus.Fields{
				"event": constant.GetTodosLogEventErrorKey,
			}).Error(constant.DbScanFailMsg)

			return nil, TodoError{Message: constant.ErrMsgInternalServer, Reason: ReasonUnknown}
		}

		todos = append(todos, todo)
	}

	logrus.WithFields(logrus.Fields{
		"event": constant.GetTodosLogEventKey,
	}).Debug("Todos fetched successfully")

	return todos, nil
}

func (service *TodoService) GetTodoByID(id string) (*types.Todo, error) {
	var todo types.Todo

	err := service.DB.QueryRow("SELECT * from todos where external_id = $1", id).Scan(&todo.ID, &todo.ExternalID, &todo.Title, &todo.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithFields(logrus.Fields{
				"event": constant.GetTodoLogEventErrorKey,
				"external_id": id,
			}).Error(constant.DbIdNotFoundMsg)

            return nil, TodoError{Message: fmt.Sprintf("Todo with id %s not found", id), Reason: ReasonNotFound}
        }

		logrus.WithFields(logrus.Fields{
			"event": constant.GetTodoLogEventErrorKey,
			"external_id": id,
		}).Error(constant.DbQueryFailMsg)

        return nil, TodoError{Message: constant.ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	logrus.WithFields(logrus.Fields{
		"event": constant.GetTodoLogEventKey,
		"external_id": id,
	}).Debug("Todo fetched successfully")

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
		logrus.WithFields(logrus.Fields{
			"event": constant.CreateTodoLogEventErrorKey,
		}).Error("Failed to create todo")

		return nil, TodoError{Message: constant.ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	logrus.WithFields(logrus.Fields{
		"event": constant.CreateTodoLogEventKey,
		"external_id": newTodo.ExternalID,
	}).Info("Todo created successfully")

	return &newTodo, nil
}

func (service *TodoService) UpdateTodo(id string, title string) (*types.Todo, error) {
	var updatedTodo types.Todo

	result, err := service.DB.Exec("UPDATE todos SET title = $1 WHERE external_id = $2", title, id)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"event": constant.UpdateTodoLogEventErrorKey,
			"external_id": id,
		}).Error(constant.DbExecFailMsg)

		return nil, TodoError{Message: constant.ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"event": constant.UpdateTodoLogEventErrorKey,
			"external_id": id,
		}).Error(constant.DbRowsAffectedFailMsg)

		return nil, TodoError{Message: constant.ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	if rowsAffected == 0 {
		logrus.WithFields(logrus.Fields{
			"event": constant.UpdateTodoLogEventErrorKey,
			"external_id": id,
		}).Error(constant.DbIdNotFoundMsg)

		return nil, TodoError{Message: fmt.Sprintf("Todo with id %s not found", id), Reason: ReasonNotFound}
	}

	err = service.DB.QueryRow("SELECT * from todos where external_id = $1", id).Scan(&updatedTodo.ID, &updatedTodo.ExternalID, &updatedTodo.Title, &updatedTodo.CreatedAt)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"event": constant.UpdateTodoLogEventErrorKey,
			"external_id": id,
		}).Error(constant.DbQueryFailMsg)

		return nil, TodoError{Message: constant.ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	logrus.WithFields(logrus.Fields{
		"event": constant.UpdateTodoLogEventKey,
		"external_id": id,
	}).Info("Todo updated successfully")

	return &updatedTodo, nil
}

func (service *TodoService) DeleteTodo(id string) error {
	result, err := service.DB.Exec("DELETE FROM todos WHERE external_id = $1", id)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"event": constant.DeleteTodoLogEventErrorKey,
			"external_id": id,
		}).Error(constant.DbExecFailMsg)

		return TodoError{Message: constant.ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"event": constant.DeleteTodoLogEventErrorKey,
			"external_id": id,
		}).Error(constant.DbRowsAffectedFailMsg)

		return TodoError{Message: constant.ErrMsgInternalServer, Reason: ReasonUnknown}
	}

	if rowsAffected == 0 {
		logrus.WithFields(logrus.Fields{
			"event": constant.DeleteTodoLogEventErrorKey,
			"external_id": id,
		}).Error(constant.DbIdNotFoundMsg)

		return TodoError{Message: fmt.Sprintf("Todo with id %s not found", id), Reason: ReasonNotFound}
	}

	logrus.WithFields(logrus.Fields{
		"event": constant.DeleteTodoLogEventKey,
		"external_id": id,
	}).Info("Todo deleted successfully")

	return nil
}
