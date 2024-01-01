package utils

import "todo-app/types"

func MapTodoResponse(todo *types.Todo) *types.TodoResponse {
	return &types.TodoResponse{
		ID:        todo.ExternalID,
		Title:     todo.Title,
		CreatedAt: todo.CreatedAt,
	}
}
