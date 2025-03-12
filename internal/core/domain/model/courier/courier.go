package courier

import (
	"errors"

	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
)

type Status string

const (
	StatusFree Status = "free"
	StatusBusy Status = "busy"
)

type Courier struct {
	ID        uuid.UUID
	Name      string
	Transport Transport
	Location  kernel.Location
	Status    Status
}

var (
	ErrCourierAlreadyBusy = errors.New("courier is already busy")
	ErrCourierAlreadyFree = errors.New("courier is already free")
)

func NewCourier(name string, transport Transport, location kernel.Location) *Courier {
	return &Courier{
		ID:        uuid.New(),
		Name:      name,
		Transport: transport,
		Location:  location,
		Status:    StatusFree,
	}
}

func (c *Courier) SetBusy() error {
	if c.Status == StatusBusy {
		return ErrCourierAlreadyBusy
	}

	c.Status = StatusBusy
	return nil
}

func (c *Courier) SetFree() error {
	if c.Status == StatusFree {
		return ErrCourierAlreadyFree
	}

	c.Status = StatusFree
	return nil
}

func (c *Courier) StepsToOrder(orderLocation kernel.Location) (steps int, _ error) {
	for c.Location != orderLocation {
		err := c.Move(orderLocation)
		if err != nil {
			return steps, err
		}
		steps++
	}
	return steps, nil
}

func (c *Courier) Move(location kernel.Location) error {
	newLocation, err := c.Transport.Move(c.Location, location)
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
