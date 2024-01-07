package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"todo-app/app/service"
	"todo-app/app/types"
	"todo-app/app/utils"
)

func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)

	return err == nil
}

func logAndRespondError(c *gin.Context, eventKey string, httpStatus int, errMsg string, extraFields logrus.Fields) {
	logFields := logrus.Fields{
		"event":   eventKey,
		"handler": c.HandlerName(),
	}

	for k, v := range extraFields {
		logFields[k] = v
	}

	logrus.WithFields(logFields).Error(errMsg)

	c.AbortWithStatusJSON(httpStatus, gin.H{"message": errMsg})
}

func GetTodos(todoService *service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventErrorKey := "todo_get_all_fail"
		eventKey := "todo_get_all"
		todos, err := todoService.GetAllTodos()

		if err != nil {
			if todoErr, ok := err.(service.TodoError); ok {
				if todoErr.Reason == service.ReasonUnknown {
					logAndRespondError(c, eventErrorKey, http.StatusInternalServerError, todoErr.Message, nil)

					return
				}
			}

			logAndRespondError(c, eventErrorKey, http.StatusInternalServerError, err.Error(), nil)

			return
		}

		mappedTodos := make([]types.TodoResponse, len(todos))

		for i, td := range todos {
			mappedTodos[i] = *utils.MapTodoResponse(&td)
		}

		logrus.WithFields(logrus.Fields{
			"event":   eventKey,
			"handler": c.HandlerName(),
		}).Info("Get all todos")

		c.IndentedJSON(http.StatusOK, mappedTodos)
	}
}

func CreateTodo(service *service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input types.TodoInput
		eventKey := "todo_create"
		eventErrorKey := "todo_create_fail"

		if err := c.ShouldBindJSON(&input); err != nil {

			logAndRespondError(c, eventErrorKey, http.StatusBadRequest, err.Error(), nil)

			return
		}

		newTodo, err := service.CreateTodo(input.Title)

		if err != nil {
			logAndRespondError(c, eventErrorKey, http.StatusInternalServerError, err.Error(), nil)

			return
		}

		logrus.WithFields(logrus.Fields{
			"event":       eventKey,
			"external_id": newTodo.ExternalID,
			"handler":     c.HandlerName(),
		}).Info("Created new todo")

		c.IndentedJSON(http.StatusCreated, utils.MapTodoResponse(newTodo))
	}
}

func GetTodoByID(todoService *service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventErrorKey := "todo_get_fail"

		id := c.Param("id")

		if !isValidUUID(id) {
			logAndRespondError(c, eventErrorKey, http.StatusBadRequest, "Invalid ID", nil)

			return
		}

		todo, err := todoService.GetTodoByID(id)

		if err != nil {
			if todoErr, ok := err.(service.TodoError); ok {
				switch todoErr.Reason {
				case service.ReasonNotFound:
					logAndRespondError(c, eventErrorKey, http.StatusNotFound, todoErr.Message, logrus.Fields{"external_id": id})
				case service.ReasonUnknown:
					logAndRespondError(c, eventErrorKey, http.StatusInternalServerError, todoErr.Message, logrus.Fields{"external_id": id})
				}

				return
			}

			logAndRespondError(c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"external_id": id})

			return
		}

		c.IndentedJSON(http.StatusOK, utils.MapTodoResponse(todo))
	}
}

func UpdateTodo(todoService *service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todoInput types.TodoInput

		eventKey := "todo_update"
		eventErrorKey := "todo_update_fail"

		id := c.Param("id")

		if !isValidUUID(id) {
			logAndRespondError(c, eventErrorKey, http.StatusBadRequest, "Invalid ID", nil)

			return
		}

		if err := c.ShouldBindJSON(&todoInput); err != nil {
			logAndRespondError(c, eventErrorKey, http.StatusBadRequest, err.Error(), logrus.Fields{"external_id": id})

			return
		}

		updatedTodo, err := todoService.UpdateTodo(id, todoInput.Title)

		if err != nil {
			if todoErr, ok := err.(service.TodoError); ok {
				switch todoErr.Reason {
				case service.ReasonNotFound:
					logAndRespondError(c, eventErrorKey, http.StatusNotFound, todoErr.Message, logrus.Fields{"external_id": id})
				case service.ReasonUnknown:
					logAndRespondError(c, eventErrorKey, http.StatusInternalServerError, todoErr.Message, logrus.Fields{"external_id": id})
				}

				return
			}

			logAndRespondError(c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"external_id": id})

			return
		}

		logrus.WithFields(logrus.Fields{
			"event":       eventKey,
			"external_id": id,
			"handler":     c.HandlerName(),
		}).Info("Updated todo")

		c.IndentedJSON(http.StatusOK, utils.MapTodoResponse(updatedTodo))
	}
}

func DeleteTodo(todoService *service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventKey := "todo_delete"
		eventErrorKey := "todo_delete_fail"

		id := c.Param("id")

		if !isValidUUID(id) {
			logAndRespondError(c, eventErrorKey, http.StatusBadRequest, "Invalid ID", nil)

			return
		}

		err := todoService.DeleteTodo(id)

		if err != nil {
			if todoErr, ok := err.(service.TodoError); ok {
				switch todoErr.Reason {
				case service.ReasonNotFound:
					logAndRespondError(c, eventErrorKey, http.StatusNotFound, todoErr.Message, logrus.Fields{"external_id": id})
				case service.ReasonUnknown:
					logAndRespondError(c, eventErrorKey, http.StatusInternalServerError, todoErr.Message, logrus.Fields{"external_id": id})
				}

				return
			}

			logAndRespondError(c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"external_id": id})

			return
		}

		logrus.WithFields(logrus.Fields{
			"event":       eventKey,
			"external_id": id,
			"handler":     c.HandlerName(),
		}).Info("Deleted todo")

		c.Status(http.StatusNoContent)
	}
}
