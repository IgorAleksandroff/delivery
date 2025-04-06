package eventhandlers

import (
	"context"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/order"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

type OrderCompleted struct {
	orderProducer OrderProducer
}

func NewOrderCompleted(orderProducer OrderProducer) (*OrderCompleted, error) {
	if orderProducer == nil {
		return nil, errs.NewValueIsRequiredError("orderProducer")
	}

	return &OrderCompleted{orderProducer: orderProducer}, nil
}

func (eh *OrderCompleted) Handle(ctx context.Context, domainEvent *order.CompletedDomainEvent) error {
	err := eh.orderProducer.Publish(ctx, domainEvent)
	if err != nil {
		return err
	}
	return nil
}
