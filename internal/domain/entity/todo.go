package entity

import (
	"time"
	// "todo-app/internal/domain/vo"

	"github.com/google/uuid"
)

type Todo struct {
	UUID      uuid.UUID
	ID        int
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
