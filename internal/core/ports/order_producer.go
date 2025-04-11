package ports

import (
	"context"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/orders"
)

type OrderProducer interface {
	Publish(ctx context.Context, domainEvent *orders.CompletedDomainEvent) error
	Close() error
}
