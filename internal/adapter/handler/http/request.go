package http

type CreateRequest struct {
	Title string `json:"title" binding:"required"`
}

type UpdateRequest struct {
	Title string `json:"title" binding:"required"`
}