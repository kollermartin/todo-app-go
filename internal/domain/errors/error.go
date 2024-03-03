package errors

import "errors"

var (
	ErrInternal       = errors.New("internal server error occured when processing ticket data")
	ErrTicketNotFound = errors.New("data for the ticket not found")
)
