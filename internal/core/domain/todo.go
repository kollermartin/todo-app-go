package domain

import (
	"github.com/google/uuid"
)

type Todo struct {
	UUID uuid.UUID    `json:"uuid"`
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	CreatedAt  string `json:"created_at"`
}