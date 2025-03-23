package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/order"
)

type OrderRepository interface {
	Add(ctx context.Context, aggregate *order.Order) error
	Update(ctx context.Context, aggregate *order.Order) error
	Get(ctx context.Context, ID uuid.UUID) (*order.Order, error)
	GetFirstInCreatedStatus(ctx context.Context) (*order.Order, error)
	GetAllInAssignedStatus(ctx context.Context) ([]*order.Order, error)
}
