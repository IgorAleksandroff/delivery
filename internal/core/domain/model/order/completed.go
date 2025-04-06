package order

import (
	"reflect"

	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain"
)

var _ domain.Event = &CompletedDomainEvent{}

type CompletedDomainEvent struct {
	// base
	ID   uuid.UUID
	Name string

	// payload
	OrderID     uuid.UUID
	OrderStatus string

	isSet bool
}

func (e CompletedDomainEvent) GetID() uuid.UUID { return e.ID }

func (e CompletedDomainEvent) GetName() string {
	return e.Name
}

func (e CompletedDomainEvent) GetOrderID() uuid.UUID {
	return e.OrderID
}

func (e CompletedDomainEvent) GetOrderStatus() string {
	return e.OrderStatus
}

func NewCompletedDomainEvent(aggregate *Order) *CompletedDomainEvent {
	event := CompletedDomainEvent{
		ID: uuid.New(),

		OrderID:     aggregate.ID(),
		OrderStatus: aggregate.Status().String(),

		isSet: true,
	}
	event.Name = reflect.TypeOf(event).Name()
	return &event
}

func (e CompletedDomainEvent) IsEmpty() bool {
	return !e.isSet
}
