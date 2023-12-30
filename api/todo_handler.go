package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"todo-app/types"
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

		rows, err := db.Query("SELECT id, title, created_at FROM todos")

		if err != nil {

			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), nil)

			return
		}

		defer rows.Close()

		todos := []types.Todo{}
		for rows.Next() {
			var todo types.Todo
			if err := rows.Scan(&todo.ID, &todo.Title, &todo.CreatedAt); err != nil {
				logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), nil)

				return
			}

			todos = append(todos, todo)
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
			ID:        uuid.New().String(),
			Title:     input.Title,
			CreatedAt: time.Now(),
		}

		_, err := db.Exec("INSERT INTO todos (id, title, created_at) VALUES ($1, $2, $3)", newTodo.ID, newTodo.Title, newTodo.CreatedAt)

		if err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"id": newTodo.ID})

			return
		}

		log.WithFields(logrus.Fields{
			"event":   eventKey,
			"id":      newTodo.ID,
			"handler": c.HandlerName(),
		}).Info("Created new todo")

		c.IndentedJSON(http.StatusCreated, newTodo)

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

		err := db.QueryRow("SELECT * from todos where id = $1", id).Scan(&todo.ID, &todo.Title, &todo.CreatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				logAndRespondError(log, c, eventErrorKey, http.StatusNotFound, "Todo not found", logrus.Fields{"id": id})
			} else {
				logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"id": id})
			}
			return
		}

		c.IndentedJSON(http.StatusOK, todo)
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
			logAndRespondError(log, c, eventErrorKey, http.StatusBadRequest, err.Error(), logrus.Fields{"id": id})

			return
		}

		result, err := db.Exec("UPDATE todos SET title = $1 WHERE id = $2", todoInput.Title, id)

		if err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"id": id})

			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"id": id})

			return
		}

		if rowsAffected == 0 {
			logAndRespondError(log, c, eventErrorKey, http.StatusNotFound, "Todo not found", logrus.Fields{"id": id})

			return
		}

		err = db.QueryRow("SELECT * from todos where id = $1", id).Scan(&updatedTodo.ID, &updatedTodo.Title, &updatedTodo.CreatedAt)
		if err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"id": id})

			return
		}

		log.WithFields(logrus.Fields{
			"event":   eventKey,
			"id":      id,
			"handler": c.HandlerName(),
		}).Info("Updated todo")

		c.IndentedJSON(http.StatusOK, updatedTodo)
	}
}
