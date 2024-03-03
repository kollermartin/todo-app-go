package request

type CreateTodoRequest struct {
	Title string `json:"title" binding:"required"`
}

type UpdateTodoRequest struct {
	Title string `json:"title" binding:"required"`
}
