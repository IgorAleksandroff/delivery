package commands

import (
	"context"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/courier"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
)

type GeoClient interface {
	GetGeolocation(ctx context.Context, street string) (kernel.Location, error)
}

type CourierRepository interface {
	Add(ctx context.Context, aggregate *courier.Courier) error
	Update(ctx context.Context, aggregate *courier.Courier) error
	Get(ctx context.Context, ID uuid.UUID) (*courier.Courier, error)
	GetAllInFreeStatus(ctx context.Context) ([]*courier.Courier, error)
}

type OrderRepository interface {
	Add(ctx context.Context, aggregate *order.Order) error
	Update(ctx context.Context, aggregate *order.Order) error
	Get(ctx context.Context, ID uuid.UUID) (*order.Order, error)
	GetFirstInCreatedStatus(ctx context.Context) (*order.Order, error)
	GetAllInAssignedStatus(ctx context.Context) ([]*order.Order, error)
}
