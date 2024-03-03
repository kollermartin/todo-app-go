package errors

import (
	"net/http"
	"todo-app/internal/domain/errors"
)

type HTTPError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewHttpError(code string, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

func NewBadRequest(err error) *HTTPError {
	return NewHttpError("BAD_REQUEST", err.Error())
}

func NewInternalServerError() *HTTPError {
	return NewHttpError("INTERNAL_SERVER_ERROR", "Internal Server Error")
}

func NewNotFound(err error) *HTTPError {
	return NewHttpError("NOT_FOUND", err.Error())
}

func GetStatusAndHttpError(err error) (int, *HTTPError) {
	if err, ok := err.(*errors.TodoError); ok {
		if err.Code == errors.ErrCodeTicketNotFound {
			return http.StatusNotFound, NewNotFound(err)
		}
	}

	return http.StatusInternalServerError, NewInternalServerError()
}
