package route

import (
	"net/http"
	"todo-app/internal/application/todo"
	"todo-app/internal/ui/http/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const DeleteTodoPath = "/todos/:id"

func NewDeleteTodoRoute(todoHandler *todo.TodoHandler) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		if id == "" {
			response.HandleValidationError(ctx, "ID is required")
			return
		}

		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			response.HandleValidationError(ctx, "Invalid ID")
			return
		}

		err = todoHandler.DeleteTodo(ctx, parsedUUID)
		if err != nil {
			response.HandleError(ctx, err)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}