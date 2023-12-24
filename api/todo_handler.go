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

func GetTodos(db *sql.DB, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, title, created_at FROM todos")

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		defer rows.Close()

		todos := []types.Todo{}
		for rows.Next() {
			var todo types.Todo
			if err := rows.Scan(&todo.ID, &todo.Title, &todo.CreatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}

			todos = append(todos, todo)
		}

		c.IndentedJSON(http.StatusOK, todos)
	}
}

func PostTodo(db *sql.DB, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input types.TodoInput

		if err := c.BindJSON(&input); err != nil {
			return
		}

		newTodo := types.Todo{
			ID:        uuid.New().String(),
			Title:     input.Title,
			CreatedAt: time.Now(),
		}

		_, err := db.Exec("INSERT INTO todos (id, title, created_at) VALUES ($1, $2, $3)", newTodo.ID, newTodo.Title, newTodo.CreatedAt)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusCreated, newTodo)

	}
}

func GetTodoByID(db *sql.DB, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var todo types.Todo

		err := db.QueryRow("SELECT * from todos where id = $1", id).Scan(&todo.ID, &todo.Title, &todo.CreatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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

		id := c.Param("id")

		if err := c.BindJSON(&todoInput); err != nil {
			return
		}

		result, err := db.Exec("UPDATE todos SET title = $1 WHERE id = $2", todoInput.Title, id)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if rowsAffected == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
			return
		}

		err = db.QueryRow("SELECT * from todos where id = $1", id).Scan(&updatedTodo.ID, &updatedTodo.Title, &updatedTodo.CreatedAt)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, updatedTodo)
	}
}
