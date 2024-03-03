package errors

import (
	"fmt"

	"github.com/google/uuid"
)

type Code string

const (
	ErrCodeInternal       Code = "INTERNAL"
	ErrCodeTicketNotFound Code = "TICKET_NOT_FOUND"
)

type TodoError struct {
	Code   Code
	Err error
}

func (e TodoError) Error() string {
	return e.Err.Error()
}

func NewTodoNotFoundError(id uuid.UUID) *TodoError {
	return &TodoError{
		Code:   ErrCodeTicketNotFound,
		Err: fmt.Errorf("todo with id %s not found", id),
	}
}