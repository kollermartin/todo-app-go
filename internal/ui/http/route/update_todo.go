package route

import (
	"net/http"
	"todo-app/internal/application/todo"
	"todo-app/internal/ui/http/request"
	"todo-app/internal/ui/http/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const UpdateTodoPath = "/todos/:id"

func NewUpdateTodoRoute(th *todo.TodoHandler) func(ctx *gin.Context) {
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
	
		var req request.UpdateRequest
	
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.HandleValidationError(ctx, err.Error())
			return
		}
	
		updatedTodo, err := th.UpdateTodo(ctx, parsedUUID, &req)
		if err != nil {
			response.HandleError(ctx, err)
			return
		}
	
		rsp := response.NewTodoResponse(updatedTodo)
	
		ctx.JSON(http.StatusOK, rsp)
	}
}