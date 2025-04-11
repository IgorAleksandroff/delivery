package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/orders"
)

type OrderRepository interface {
	Add(ctx context.Context, aggregate *orders.Order) error
	Update(ctx context.Context, aggregate *orders.Order) error
	Get(ctx context.Context, ID uuid.UUID) (*orders.Order, error)
	GetFirstInCreatedStatus(ctx context.Context) (*orders.Order, error)
	GetAllInAssignedStatus(ctx context.Context) ([]*orders.Order, error)
}
