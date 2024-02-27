package todo

import (
	"context"
	"todo-app/internal/domain/entity"
	"todo-app/internal/domain/repository"
	"todo-app/internal/ui/http/request"

	// "todo-app/internal/domain/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TodoHandler struct {
	repo repository.TodoRepository
}
//TODo zbavit se ctx z application vrstvy
func NewTodoHandler(repo repository.TodoRepository) *TodoHandler {
	return &TodoHandler{repo}
}

func (th *TodoHandler) GetAllTodos(ctx context.Context) ([]entity.Todo, error) {
	todos, err := th.repo.GetAllTodos(ctx)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (th *TodoHandler) GetTodo(ctx *gin.Context, id uuid.UUID) (*entity.Todo, error) {
	todo, err := th.repo.GetTodo(ctx, id)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (th *TodoHandler) CreateTodo(ctx *gin.Context, todoReq *request.CreateRequest) (*entity.Todo, error) {
	todo := entity.Todo{
		Title: todoReq.Title,
	}

	createdTodo, err := th.repo.CreateTodo(ctx, &todo)
	if err != nil {
		return nil, err
	}

	return createdTodo, nil
}

func (th *TodoHandler) UpdateTodo(ctx *gin.Context, id uuid.UUID , todoReq *request.UpdateRequest) (*entity.Todo, error) {
	todo := entity.Todo{
		Title: todoReq.Title,
		UUID: id,
	}

	updatedTodo, err := th.repo.UpdateTodo(ctx, &todo)
	if err != nil {
		return nil, err
	}

	return updatedTodo, nil
	// id := ctx.Param("id")

	// if id == "" {
	// 	response.HandleValidationError(ctx, "ID is required")
	// 	return
	// }

	// parsedUUID, err := uuid.Parse(id)
	// if err != nil {
	// 	response.HandleValidationError(ctx, "Invalid ID")
	// 	return
	// }

	// var req request.UpdateRequest

	// if err := ctx.ShouldBindJSON(&req); err != nil {
	// 	response.HandleValidationError(ctx, err.Error())
	// 	return
	// }

	// todo := entity.Todo{
	// 	UUID:  parsedUUID,
	// 	Title: req.Title,
	// }

	// updatedTodo, err := th.repo.UpdateTodo(ctx, &todo)
	// if err != nil {
	// 	response.HandleError(ctx, err)
	// 	return
	// }

	// rsp := response.NewTodoResponse(updatedTodo)

	// ctx.JSON(http.StatusOK, rsp)
}

func (th *TodoHandler) DeleteTodo(ctx *gin.Context, id uuid.UUID) error {
	err := th.repo.DeleteTodo(ctx, id)
	if err != nil {
		return err
	}
	
	return nil
}
