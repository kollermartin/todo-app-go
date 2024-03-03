package route

import (
	"errors"
	"net/http"
	"todo-app/internal/application/todo"
	uiErrors "todo-app/internal/ui/http/errors"
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
			ctx.AbortWithStatusJSON(http.StatusBadRequest, uiErrors.NewBadRequest(errors.New("ID is required")))
			return
		}

		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, uiErrors.NewBadRequest(errors.New("invalid ID")))
			return
		}

		var req request.UpdateTodoRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, uiErrors.NewBadRequest(err))
			return
		}

		updatedTodo, err := th.UpdateTodo(ctx, parsedUUID, &req)
		if err != nil {
			ctx.AbortWithStatusJSON(uiErrors.GetStatusAndHttpError(err))
			return
		}

		rsp := response.NewTodoResponse(updatedTodo)

		ctx.JSON(http.StatusOK, rsp)
	}
}
