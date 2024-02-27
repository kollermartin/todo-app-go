package route

import (
	"net/http"
	"todo-app/internal/application/todo"
	"todo-app/internal/ui/http/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const GetTodoPath = "/todos/:id"

func NewGetTodoRoute(todoHandler *todo.TodoHandler) func(ctx *gin.Context) {
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

		todo, err := todoHandler.GetTodo(ctx, parsedUUID)
		if err != nil {
			response.HandleError(ctx, err)
			return
		}

		rsp := response.NewTodoResponse(todo)

		ctx.JSON(http.StatusOK, rsp)
	}
}