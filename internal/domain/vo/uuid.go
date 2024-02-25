package vo

import (

	"github.com/google/uuid"
)

// Currently not used anywhere because of errors with psql scan

type UUID uuid.UUID

func NewUUID() UUID {
	return UUID(uuid.New())
}

func ParseUUID(s string) (UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, err
	}

	return UUID(u), nil
}

func (u UUID) String() string {
	return uuid.UUID(u).String()
}