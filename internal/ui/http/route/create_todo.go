package route

import (
	"net/http"
	"todo-app/internal/application/todo"
	"todo-app/internal/ui/http/request"
	"todo-app/internal/ui/http/response"

	"github.com/gin-gonic/gin"
)

const CreateTodoPath = "/todos"

func NewCreateTodoRoute(todoHandler *todo.TodoHandler) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req request.CreateRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.HandleValidationError(ctx, err.Error())
			return
		}



		createdTodo, err := todoHandler.CreateTodo(ctx, &req)
		if err != nil {
			response.HandleError(ctx, err)
			return
		}

		rsp := response.NewTodoResponse(createdTodo)

		ctx.JSON(http.StatusCreated, rsp)
	}
}