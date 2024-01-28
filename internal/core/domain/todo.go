package domain

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	UUID uuid.UUID    `json:"uuid"`
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}