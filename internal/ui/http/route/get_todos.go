package route

import (
	"net/http"
	"todo-app/internal/application/todo"
	"todo-app/internal/ui/http/errors"
	"todo-app/internal/ui/http/response"

	"github.com/gin-gonic/gin"
)

const GetTodosPath = "/todos"

func NewGetTodosRoute(todoHandler *todo.TodoHandler) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		todos, err := todoHandler.GetAllTodos(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(errors.GetStatusAndHttpError(err))
			return
		}

		todoListRsp := make([]response.TodoResponse, 0, len(todos))

		for _, todo := range todos {
			todoListRsp = append(todoListRsp, response.NewTodoResponse(&todo))
		}

		ctx.JSON(http.StatusOK, todoListRsp)
	}
}
