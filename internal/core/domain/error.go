package domain

import "errors"
// TODO improve errors
var (
	ErrInternal = errors.New("internal server error occured when processing ticket data")
	// TODO Improve, pass in ticket ID
	ErrTicketNotFound = errors.New("data for the ticket not found")
)
