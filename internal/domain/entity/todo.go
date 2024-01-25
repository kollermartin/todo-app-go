package entity

import (
	"todo-app/internal/domain/value"

	"github.com/google/uuid"
)

type Todo struct {
	ID         int       `json:"id"`
	ExternalID uuid.UUID    `json:"external_id"`
	Title      value.CreatedAt    `json:"title"`
	CreatedAt  value.Title `json:"created_at"`
}