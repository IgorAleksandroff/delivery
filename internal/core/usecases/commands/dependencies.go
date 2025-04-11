package commands

import (
	"context"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/courier"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/orders"
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
	Add(ctx context.Context, aggregate *orders.Order) error
	Update(ctx context.Context, aggregate *orders.Order) error
	Get(ctx context.Context, ID uuid.UUID) (*orders.Order, error)
	GetFirstInCreatedStatus(ctx context.Context) (*orders.Order, error)
	GetAllInAssignedStatus(ctx context.Context) ([]*orders.Order, error)
}
