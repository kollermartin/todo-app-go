package route

import (
	"errors"
	"net/http"
	"todo-app/internal/application/todo"
	uiErrors "todo-app/internal/ui/http/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const DeleteTodoPath = "/todos/:id"

func NewDeleteTodoRoute(todoHandler *todo.TodoHandler) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		if id == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, uiErrors.NewBadRequest(errors.New("ID is required")))
			return
		}

		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, uiErrors.NewBadRequest(errors.New("invalid ID")))
			return
		}

		err = todoHandler.DeleteTodo(ctx, parsedUUID)
		if err != nil {
			ctx.AbortWithStatusJSON(uiErrors.GetStatusAndHttpError(err))
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}