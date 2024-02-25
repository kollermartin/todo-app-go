package domain

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	UUID      uuid.UUID
	ID        int
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
