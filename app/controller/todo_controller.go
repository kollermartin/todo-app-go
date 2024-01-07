package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"todo-app/app/service"
	"todo-app/app/types"
	"todo-app/app/utils"
)

func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)

	return err == nil
}

func respondError(c *gin.Context, httpStatus int, errMsg string) {
	c.AbortWithStatusJSON(httpStatus, gin.H{"message": errMsg})
}

func GetTodos(todoService *service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		todos, err := todoService.GetAllTodos()

		if err != nil {
			if todoErr, ok := err.(service.TodoError); ok {
				if todoErr.Reason == service.ReasonUnknown {
					respondError(c, http.StatusInternalServerError, todoErr.Message)

					return
				}
			}

			respondError(c, http.StatusInternalServerError, err.Error())

			return
		}

		mappedTodos := make([]types.TodoResponse, len(todos))

		for i, td := range todos {
			mappedTodos[i] = *utils.MapTodoResponse(&td)
		}

		c.IndentedJSON(http.StatusOK, mappedTodos)
	}
}

func CreateTodo(todoService *service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input types.TodoInput

		if err := c.ShouldBindJSON(&input); err != nil {

			respondError(c, http.StatusBadRequest, err.Error())

			return
		}

		newTodo, err := todoService.CreateTodo(input.Title)

		if err != nil {
			if todoErr, ok := err.(service.TodoError); ok {
				if todoErr.Reason == service.ReasonUnknown {
					respondError(c, http.StatusInternalServerError, todoErr.Message)

					return
				}
			}

			respondError(c, http.StatusInternalServerError, err.Error())

			return
		}

		c.IndentedJSON(http.StatusCreated, utils.MapTodoResponse(newTodo))
	}
}

func GetTodoByID(todoService *service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		if !isValidUUID(id) {
			respondError(c, http.StatusBadRequest, "Invalid ID")

			return
		}

		todo, err := todoService.GetTodoByID(id)

		if err != nil {
			if todoErr, ok := err.(service.TodoError); ok {
				switch todoErr.Reason {
				case service.ReasonNotFound:
					respondError(c, http.StatusNotFound, todoErr.Message)
				case service.ReasonUnknown:
					respondError(c, http.StatusInternalServerError, todoErr.Message)
				}

				return
			}

			respondError(c, http.StatusInternalServerError, err.Error())

			return
		}

		c.IndentedJSON(http.StatusOK, utils.MapTodoResponse(todo))
	}
}

func UpdateTodo(todoService *service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todoInput types.TodoInput

		id := c.Param("id")

		if !isValidUUID(id) {
			respondError(c, http.StatusBadRequest, "Invalid ID")

			return
		}

		if err := c.ShouldBindJSON(&todoInput); err != nil {
			respondError(c, http.StatusBadRequest, err.Error())

			return
		}

		updatedTodo, err := todoService.UpdateTodo(id, todoInput.Title)

		if err != nil {
			if todoErr, ok := err.(service.TodoError); ok {
				switch todoErr.Reason {
				case service.ReasonNotFound:
					respondError(c, http.StatusNotFound, todoErr.Message)
				case service.ReasonUnknown:
					respondError(c, http.StatusInternalServerError, todoErr.Message)
				}

				return
			}

			respondError(c, http.StatusInternalServerError, err.Error())

			return
		}

		c.IndentedJSON(http.StatusOK, utils.MapTodoResponse(updatedTodo))
	}
}

func DeleteTodo(todoService *service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		if !isValidUUID(id) {
			respondError(c, http.StatusBadRequest, "Invalid ID")

			return
		}

		err := todoService.DeleteTodo(id)

		if err != nil {
			if todoErr, ok := err.(service.TodoError); ok {
				switch todoErr.Reason {
				case service.ReasonNotFound:
					respondError(c, http.StatusNotFound, todoErr.Message)
				case service.ReasonUnknown:
					respondError(c, http.StatusInternalServerError, todoErr.Message)
				}

				return
			}

			respondError(c, http.StatusInternalServerError, err.Error())

			return
		}

		c.Status(http.StatusNoContent)
	}
}
