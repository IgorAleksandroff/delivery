package ddd

import (
	"github.com/google/uuid"
)

type DomainEvent interface {
	ID() uuid.UUID
	Name() string
}

type AggregateRoot interface {
	GetDomainEvents() []DomainEvent
	ClearDomainEvents()
}
