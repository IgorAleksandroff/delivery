package order

import (
	"errors"

	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
)

type Status string

func (s Status) String() string {
	return string(s)
}

const (
	StatusCreated   Status = "created"
	StatusAssigned  Status = "assigned"
	StatusCompleted Status = "completed"
)

type Order struct {
	id        uuid.UUID
	location  kernel.Location
	status    Status
	courierID *uuid.UUID
}

var (
	ErrOrderAlreadyAssigned = errors.New("order is already assigned to courier")
	ErrOrderNotAssigned     = errors.New("order is not assigned to courier")
	ErrOrderCompleted       = errors.New("order is already completed")
	ErrInvalidLocation      = errors.New("invalid Location")
	ErrInvalidOrderId       = errors.New("invalid order id")
)

func NewOrder(id uuid.UUID, location kernel.Location) (*Order, error) {
	if id == uuid.Nil {
		return nil, ErrInvalidOrderId
	}

	if location.IsEmpty() {
		return nil, ErrInvalidLocation
	}

	return &Order{
		id:        id,
		location:  location,
		status:    StatusCreated,
		courierID: nil,
	}, nil
}

func MustNewOrder(id uuid.UUID, location kernel.Location) *Order {
	t, err := NewOrder(id, location)
	if err != nil {
		panic(err)
	}
	return t
}

func (o *Order) AssignToCourier(courierId uuid.UUID) error {
	if o.IsCompleted() {
		return ErrOrderCompleted
	}

	if o.IsAssigned() && *o.courierID != courierId {
		return ErrOrderAlreadyAssigned
	}

	o.status = StatusAssigned
	o.courierID = &courierId

	return nil
}

func (o *Order) Complete() error {
	if !o.IsAssigned() {
		return ErrOrderNotAssigned
	}

	o.status = StatusCompleted

	return nil
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) AssignedCourier() *uuid.UUID {
	return o.courierID
}

func (o *Order) IsAssigned() bool {
	return o.status == StatusAssigned && o.courierID != nil
}

func (o *Order) IsCompleted() bool {
	return o.status == StatusCompleted
}
