package domain

import (
	"github.com/google/uuid"
)

type Event interface {
	ID() uuid.UUID
	Name() string
}
