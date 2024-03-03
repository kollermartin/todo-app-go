package errors

import (
	"fmt"
	"net/http"
	"todo-app/internal/domain/errors"
)

var ErrorStatusMap = map[error]int{
	errors.ErrInternal:       http.StatusInternalServerError,
	errors.ErrTicketNotFound: http.StatusNotFound,
}

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
	return NewHttpError("BAD_REQUEST", fmt.Sprintf("Bad request %s", err.Error()))
}

func NewInternalServerError() *HTTPError {
	return NewHttpError("INTERNAL_SERVER_ERROR", "Internal Server Error")
}

func NewNotFound() *HTTPError {
	return NewHttpError("NOT_FOUND", "Not Found")
}

func GetStatusAndHttpError(err error) (int, *HTTPError) {
	if code, ok := ErrorStatusMap[err]; ok {
		switch code {
		case http.StatusNotFound:
			return http.StatusNotFound, NewNotFound()
		case http.StatusInternalServerError:
			return http.StatusInternalServerError, NewInternalServerError()
		default:
			return http.StatusInternalServerError, NewInternalServerError()
		}
	}

	return http.StatusInternalServerError, NewInternalServerError()
}
