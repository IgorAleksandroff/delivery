package courier

import (
	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
)

func RestoreCourier(ID uuid.UUID, name string, transport *Transport, location kernel.Location, status Status) *Courier {
	return &Courier{
		id:        ID,
		name:      name,
		transport: transport,
		location:  location,
		status:    status,
	}
}

func RestoreTransport(ID uuid.UUID, name string, speed int) *Transport {
	return &Transport{
		id:    ID,
		name:  name,
		speed: speed,
	}
}
