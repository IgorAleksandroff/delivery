package orderrepo

import (
	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/orders"
)

type OrderDTO struct {
	ID        uuid.UUID     `gorm:"type:uuid;primaryKey"`
	CourierID *uuid.UUID    `gorm:"type:uuid;index"`
	Location  LocationDTO   `gorm:"embedded;embeddedPrefix:location_"`
	Status    orders.Status `gorm:"type:varchar(20)"`
}

type LocationDTO struct {
	X int
	Y int
}

func (OrderDTO) TableName() string {
	return "orders"
}

func DomainToDTO(aggregate *orders.Order) OrderDTO {
	var orderDTO OrderDTO
	orderDTO.ID = aggregate.ID()
	orderDTO.CourierID = aggregate.AssignedCourier()
	orderDTO.Location = LocationDTO{
		X: aggregate.Location().X(),
		Y: aggregate.Location().Y(),
	}
	orderDTO.Status = aggregate.Status()
	return orderDTO
}

func DtoToDomain(dto OrderDTO) *orders.Order {
	var aggregate *orders.Order
	location, _ := kernel.NewLocation(dto.Location.X, dto.Location.Y)
	aggregate = orders.RestoreOrder(dto.ID, dto.CourierID, location, dto.Status)
	return aggregate
}
