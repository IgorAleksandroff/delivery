package courier

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

const (
	SPEED_MIN = 1
	SPEED_MAX = 3
)

type Transport struct {
	id    uuid.UUID
	name  string
	speed int
}

func NewTransport(name string, speed int) (*Transport, error) {
	if name == "" {
		return nil, errs.NewValueIsRequiredError("transport name cannot be empty")
	}

	if speed < SPEED_MIN || speed > SPEED_MAX {
		return nil, errs.NewValueIsOutOfRangeError("speed", speed, SPEED_MIN, SPEED_MAX)
	}

	return &Transport{
		id:    uuid.New(),
		name:  name,
		speed: speed,
	}, nil
}

func MustNewTransport(name string, speed int) *Transport {
	t, err := NewTransport(name, speed)
	if err != nil {
		panic(err)
	}
	return t
}

func (t Transport) ID() uuid.UUID {
	return t.id
}

func (t Transport) Name() string {
	return t.name
}

func (t Transport) Speed() int {
	return t.speed
}

func (t Transport) Equals(other Transport) bool {
	return t.id == other.id
}

func (t Transport) MoveTowards(current, target kernel.Location) (kernel.Location, error) {
	if current.Equals(target) {
		return current, nil
	}

	dx := -1
	if target.X() > current.X() {
		dx = 1
	}

	dy := -1
	if target.Y() > current.Y() {
		dy = 1
	}

	absStepsX := dx * (target.X() - current.X())
	absStepsY := dy * (target.Y() - current.Y())

	stepsX := min(absStepsX, t.speed)
	stepsY := min(absStepsY, max(t.speed-stepsX, 0))

	newX := current.X() + dx*stepsX
	newY := current.Y() + dy*stepsY

	return kernel.NewLocation(newX, newY)
}

func (t Transport) String() string {
	return fmt.Sprintf("Transport{id=%s, name=%s, speed=%d}", t.id, t.name, t.speed)
}
