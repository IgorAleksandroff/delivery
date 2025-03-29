package courierrepo

import (
	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/courier"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
)

type CourierDTO struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string
	Transport TransportDTO   `gorm:"foreignKey:CourierID;constraint:OnDelete:CASCADE;"`
	Location  LocationDTO    `gorm:"embedded;embeddedPrefix:location_"`
	Status    courier.Status `gorm:"type:varchar(20)"`
}

type TransportDTO struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Name      string
	Speed     int
	CourierID uuid.UUID `gorm:"type:uuid;index"`
}

type LocationDTO struct {
	X int
	Y int
}

func (CourierDTO) TableName() string {
	return "transports"
}

func DomainToDTO(aggregate *courier.Courier) CourierDTO {
	var courierDTO CourierDTO
	courierDTO.ID = aggregate.ID()
	courierDTO.Name = aggregate.Name()
	courierDTO.Transport = TransportDTO{
		ID:        aggregate.Transport().ID(),
		Name:      aggregate.Transport().Name(),
		Speed:     aggregate.Transport().Speed(),
		CourierID: aggregate.ID(),
	}
	courierDTO.Location = LocationDTO{
		X: aggregate.Location().X(),
		Y: aggregate.Location().Y(),
	}
	courierDTO.Status = aggregate.Status()
	return courierDTO
}

func DtoToDomain(dto CourierDTO) *courier.Courier {
	var aggregate *courier.Courier
	transport := courier.RestoreTransport(dto.Transport.ID, dto.Transport.Name, dto.Transport.Speed)
	location, _ := kernel.NewLocation(dto.Location.X, dto.Location.Y)
	aggregate = courier.RestoreCourier(dto.ID, dto.Name, transport, location, dto.Status)
	return aggregate
}
