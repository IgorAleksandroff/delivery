package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/courier"
)

type CourierRepository interface {
	Add(ctx context.Context, aggregate *courier.Courier) error
	Update(ctx context.Context, aggregate *courier.Courier) error
	Get(ctx context.Context, ID uuid.UUID) (*courier.Courier, error)
	GetAllInFreeStatus(ctx context.Context) ([]*courier.Courier, error)
}
