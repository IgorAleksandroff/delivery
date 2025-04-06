package order

import (
	"github.com/google/uuid"
)

const CompletedDomainEventName = "order.completed.event"

type CompletedDomainEvent struct {
	// base
	id   uuid.UUID
	name string

	// payload
	orderID     uuid.UUID
	orderStatus string

	isSet bool
}

func (e CompletedDomainEvent) ID() uuid.UUID { return e.id }

func (e CompletedDomainEvent) Name() string {
	return e.name
}

func (e CompletedDomainEvent) OrderID() uuid.UUID {
	return e.orderID
}

func (e CompletedDomainEvent) OrderStatus() string {
	return e.orderStatus
}

func NewCompletedDomainEvent(aggregate *Order) CompletedDomainEvent {
	return CompletedDomainEvent{
		id:   uuid.New(),
		name: CompletedDomainEventName,

		orderID:     aggregate.ID(),
		orderStatus: aggregate.Status().String(),

		isSet: true,
	}
}

func (e CompletedDomainEvent) IsEmpty() bool {
	return !e.isSet
}
