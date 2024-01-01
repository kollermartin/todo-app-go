package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"todo-app/types"
	"todo-app/utils"
)

func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)

	return err == nil
}

func logAndRespondError(log *logrus.Logger, c *gin.Context, eventKey string, httpStatus int, errMsg string, extraFields logrus.Fields) {
	logFields := logrus.Fields{
		"event":   eventKey,
		"handler": c.HandlerName(),
	}

	for k, v := range extraFields {
		logFields[k] = v
	}

	log.WithFields(logFields).Error(errMsg)

	c.AbortWithStatusJSON(httpStatus, gin.H{"message": errMsg})
}

func GetTodos(db *sql.DB, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventErrorKey := "todo_get_all_fail"
		eventKey := "todo_get_all"

		rows, err := db.Query("SELECT * FROM todos")

		if err != nil {

			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), nil)

			return
		}

		defer rows.Close()

		todos := []types.TodoResponse{}
		for rows.Next() {
			var todo types.Todo
			if err := rows.Scan(&todo.ID, &todo.ExternalID, &todo.Title, &todo.CreatedAt); err != nil {
				logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), nil)

				return
			}

			todos = append(todos, *utils.MapTodoResponse(&todo))
		}

		log.WithFields(logrus.Fields{
			"event":   eventKey,
			"handler": c.HandlerName(),
		}).Info("Get all todos")

		c.IndentedJSON(http.StatusOK, todos)
	}
}

func CreateTodo(db *sql.DB, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input types.TodoInput
		eventKey := "todo_create"
		eventErrorKey := "todo_create_fail"

		if err := c.ShouldBindJSON(&input); err != nil {

			logAndRespondError(log, c, eventErrorKey, http.StatusBadRequest, err.Error(), nil)

			return
		}

		newTodo := types.Todo{
			ExternalID: uuid.New().String(),
			Title:      input.Title,
			CreatedAt:  time.Now(),
		}

		err := db.QueryRow("INSERT INTO todos (external_id, title, created_at) VALUES ($1, $2, $3) RETURNING id", newTodo.ExternalID, newTodo.Title, newTodo.CreatedAt).Scan(&newTodo.ID)

		if err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"external_id": newTodo.ExternalID})

			return
		}

		log.WithFields(logrus.Fields{
			"event":       eventKey,
			"external_id": newTodo.ExternalID,
			"handler":     c.HandlerName(),
		}).Info("Created new todo")

		c.IndentedJSON(http.StatusCreated, utils.MapTodoResponse(&newTodo))

	}
}

func GetTodoByID(db *sql.DB, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todo types.Todo

		eventErrorKey := "todo_get_fail"

		id := c.Param("id")

		if !isValidUUID(id) {
			logAndRespondError(log, c, eventErrorKey, http.StatusBadRequest, "Invalid ID", nil)

			return
		}

		err := db.QueryRow("SELECT * from todos where external_id = $1", id).Scan(&todo.ID, &todo.ExternalID, &todo.Title, &todo.CreatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				logAndRespondError(log, c, eventErrorKey, http.StatusNotFound, "Todo not found", logrus.Fields{"external_id": id})
			} else {
				logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"external_id": id})
			}
			return
		}

		c.IndentedJSON(http.StatusOK, utils.MapTodoResponse(&todo))
	}
}

func UpdateTodo(db *sql.DB, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todoInput types.TodoInput
		var updatedTodo types.Todo

		eventKey := "todo_update"
		eventErrorKey := "todo_update_fail"

		id := c.Param("id")

		if !isValidUUID(id) {
			logAndRespondError(log, c, eventErrorKey, http.StatusBadRequest, "Invalid ID", nil)

			return
		}

		if err := c.ShouldBindJSON(&todoInput); err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusBadRequest, err.Error(), logrus.Fields{"external_id": id})

			return
		}

		result, err := db.Exec("UPDATE todos SET title = $1 WHERE external_id = $2", todoInput.Title, id)

		if err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"external_id": id})

			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"external_id": id})

			return
		}

		if rowsAffected == 0 {
			logAndRespondError(log, c, eventErrorKey, http.StatusNotFound, "Todo not found", logrus.Fields{"external_id": id})

			return
		}

		err = db.QueryRow("SELECT * from todos where external_id = $1", id).Scan(&updatedTodo.ID, &updatedTodo.ExternalID, &updatedTodo.Title, &updatedTodo.CreatedAt)
		if err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"external_id": id})

			return
		}

		log.WithFields(logrus.Fields{
			"event":       eventKey,
			"external_id": id,
			"handler":     c.HandlerName(),
		}).Info("Updated todo")

		c.IndentedJSON(http.StatusOK, utils.MapTodoResponse(&updatedTodo))
	}
}

func DeleteTodo(db *sql.DB, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventKey := "todo_delete"
		eventErrorKey := "todo_delete_fail"

		id := c.Param("id")

		if !isValidUUID(id) {
			logAndRespondError(log, c, eventErrorKey, http.StatusBadRequest, "Invalid ID", nil)

			return
		}

		result, err := db.Exec("DELETE FROM todos WHERE external_id = $1", id)

		if err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"external_id": id})

			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"external_id": id})

			return
		}

		if rowsAffected == 0 {
			logAndRespondError(log, c, eventErrorKey, http.StatusNotFound, "Todo not found", logrus.Fields{"external_id": id})

			return
		}

		log.WithFields(logrus.Fields{
			"event":       eventKey,
			"external_id": id,
			"handler":     c.HandlerName(),
		}).Info("Deleted todo")

		c.Status(http.StatusNoContent)
	}
}
