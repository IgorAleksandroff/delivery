package courier

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

type Status string

const (
	StatusFree Status = "free"
	StatusBusy Status = "busy"
)

type Courier struct {
	ID        uuid.UUID
	Name      string
	Transport *Transport
	Location  kernel.Location
	Status    Status
}

var (
	ErrCourierAlreadyBusy = errors.New("courier is already busy")
	ErrCourierAlreadyFree = errors.New("courier is already free")
	ErrInvalidCourierName = errors.New("invalid courier name")
	ErrInvalidLocation    = errors.New("invalid Location")
)

func NewCourier(name string, transportName string, transportSpeed int, location kernel.Location) (*Courier, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidCourierName
	}

	transport, err := NewTransport(transportName, transportSpeed)
	if err != nil {
		return nil, err
	}

	if location.IsEmpty() {
		return nil, ErrInvalidLocation
	}

	return &Courier{
		ID:        uuid.New(),
		Name:      name,
		Transport: transport,
		Location:  location,
		Status:    StatusFree,
	}, nil
}

func MustNewCourier(name string, transportName string, transportSpeed int, location kernel.Location) *Courier {
	t, err := NewCourier(name, transportName, transportSpeed, location)
	if err != nil {
		panic(err)
	}
	return t
}

func (c *Courier) SetBusy() error {
	if c.IsBusy() {
		return ErrCourierAlreadyBusy
	}

	c.Status = StatusBusy
	return nil
}

func (c *Courier) SetFree() error {
	if c.IsFree() {
		return ErrCourierAlreadyFree
	}

	c.Status = StatusFree
	return nil
}

func (c *Courier) StepsToOrder(orderLocation kernel.Location) (steps int, _ error) {
	if orderLocation.IsEmpty() {
		return steps, errs.NewValueIsRequiredError("orderLocation")
	}

	for c.Location != orderLocation {
		err := c.Move(orderLocation)
		if err != nil {
			return steps, err
		}
		steps++
	}
	return steps, nil
}

func (c *Courier) Move(target kernel.Location) error {
	newLocation, err := c.Transport.Move(c.Location, target)
	if err != nil {
		return err
	}
	c.Location = newLocation
	return nil
}

func (c *Courier) IsFree() bool {
	return c.Status == StatusFree
}

func (c *Courier) IsBusy() bool {
	return c.Status == StatusBusy
}
