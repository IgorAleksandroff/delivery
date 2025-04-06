package ports

import (
	"context"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/order"
)

type OrderProducer interface {
	Publish(ctx context.Context, domainEvent order.CompletedDomainEvent) error
	Close() error
}
