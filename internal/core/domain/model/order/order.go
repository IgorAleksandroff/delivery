package order

import (
	"errors"

	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
)

type Status string

const (
	StatusCreated   Status = "created"
	StatusAssigned  Status = "assigned"
	StatusCompleted Status = "completed"
)

type Order struct {
	ID        uuid.UUID
	Location  kernel.Location
	Status    Status
	CourierID *uuid.UUID
}

var (
	ErrOrderAlreadyAssigned = errors.New("order is already assigned to courier")
	ErrOrderNotAssigned     = errors.New("order is not assigned to courier")
	ErrOrderCompleted       = errors.New("order is already completed")
)

func NewOrder(id uuid.UUID, location kernel.Location) *Order {
	return &Order{
		ID:        id,
		Location:  location,
		Status:    StatusCreated,
		CourierID: nil,
	}
}

func (o *Order) AssignToCourier(courierId uuid.UUID) error {
	if o.Status == StatusCompleted {
		return ErrOrderCompleted
	}

	if o.Status == StatusAssigned && o.CourierID != nil && *o.CourierID != courierId {
		return ErrOrderAlreadyAssigned
	}

	o.Status = StatusAssigned
	o.CourierID = &courierId

	return nil
}

func (o *Order) Complete() error {
	if o.Status != StatusAssigned || o.CourierID == nil {
		return ErrOrderNotAssigned
	}

	o.Status = StatusCompleted

	return nil
}

func (o *Order) IsAssigned() bool {
	return o.Status == StatusAssigned && o.CourierID != nil
}

func (o *Order) IsCompleted() bool {
	return o.Status == StatusCompleted
}
