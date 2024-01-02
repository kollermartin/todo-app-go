package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"todo-app/service"
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

func GetTodos(service *service.TodoService, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventErrorKey := "todo_get_all_fail"
		eventKey := "todo_get_all"
		todos, err := service.GetAllTodos()

		if err != nil {

			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), nil)

			return
		}

		mappedTodos := make([]types.TodoResponse, len(todos))

		for i, td := range todos {
			mappedTodos[i] = *utils.MapTodoResponse(&td)
		}

		log.WithFields(logrus.Fields{
			"event":   eventKey,
			"handler": c.HandlerName(),
		}).Info("Get all todos")

		c.IndentedJSON(http.StatusOK, mappedTodos)
	}
}

func CreateTodo(service *service.TodoService, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input types.TodoInput
		eventKey := "todo_create"
		eventErrorKey := "todo_create_fail"

		if err := c.ShouldBindJSON(&input); err != nil {

			logAndRespondError(log, c, eventErrorKey, http.StatusBadRequest, err.Error(), nil)

			return
		}

		newTodo, err := service.CreateTodo(input.Title)

		if err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), nil)

			return
		}

		log.WithFields(logrus.Fields{
			"event":       eventKey,
			"external_id": newTodo.ExternalID,
			"handler":     c.HandlerName(),
		}).Info("Created new todo")

		c.IndentedJSON(http.StatusCreated, utils.MapTodoResponse(newTodo))
	}
}

func GetTodoByID(service *service.TodoService, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventErrorKey := "todo_get_fail"

		id := c.Param("id")

		if !isValidUUID(id) {
			logAndRespondError(log, c, eventErrorKey, http.StatusBadRequest, "Invalid ID", nil)

			return
		}

		todo, err := service.GetTodoByID(id)

		if err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusNotFound, err.Error(), logrus.Fields{"external_id": id})

			return
		}

		c.IndentedJSON(http.StatusOK, utils.MapTodoResponse(todo))
	}
}

func UpdateTodo(service *service.TodoService, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todoInput types.TodoInput

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

		updatedTodo, err := service.UpdateTodo(id, todoInput.Title)

		if err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"external_id": id})

			return
		}

		log.WithFields(logrus.Fields{
			"event":       eventKey,
			"external_id": id,
			"handler":     c.HandlerName(),
		}).Info("Updated todo")

		c.IndentedJSON(http.StatusOK, utils.MapTodoResponse(updatedTodo))
	}
}

func DeleteTodo(service *service.TodoService, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventKey := "todo_delete"
		eventErrorKey := "todo_delete_fail"

		id := c.Param("id")

		if !isValidUUID(id) {
			logAndRespondError(log, c, eventErrorKey, http.StatusBadRequest, "Invalid ID", nil)

			return
		}

		err := service.DeleteTodo(id)

		if err != nil {
			logAndRespondError(log, c, eventErrorKey, http.StatusInternalServerError, err.Error(), logrus.Fields{"external_id": id})

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
