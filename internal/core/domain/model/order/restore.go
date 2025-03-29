package order

import (
	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
)

func RestoreOrder(ID uuid.UUID, courierID *uuid.UUID, location kernel.Location, status Status) *Order {
	return &Order{
		id:        ID,
		courierID: courierID,
		location:  location,
		status:    status,
	}
}
