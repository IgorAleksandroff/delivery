package kernel

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

type Location struct {
	x int
	y int
}

func NewLocation(x, y int) (Location, error) {
	if x < 1 || x > 10 || y < 1 || y > 10 {
		return Location{}, errs.NewValueIsInvalidError("coordinates must be between 1 and 10 inclusive")
	}
	return Location{x: x, y: y}, nil
}

func MustNewLocation(x, y int) Location {
	loc, err := NewLocation(x, y)
	if err != nil {
		panic(err)
	}
	return loc
}

func (l Location) X() int {
	return l.x
}

func (l Location) Y() int {
	return l.y
}

func (l Location) Equals(other Location) bool {
	return l.x == other.x && l.y == other.y
}

func (l Location) Distance(other Location) int {
	dx := abs(l.x - other.x)
	dy := abs(l.y - other.y)
	return dx + dy
}

func (l Location) String() string {
	return fmt.Sprintf("(%d,%d)", l.x, l.y)
}

func CreateRandomLocation() Location {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	x := r.Intn(10) + 1
	y := r.Intn(10) + 1

	location, err := NewLocation(x, y)
	if err != nil {
		fmt.Println(err)
	}

	return location
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
