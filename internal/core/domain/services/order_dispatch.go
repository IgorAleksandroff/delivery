package services

import (
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/courier"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/orders"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

type Dispatcher struct{}

func NewOrderDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (p *Dispatcher) Dispatch(order *orders.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if order == nil {
		return nil, errs.NewValueIsRequiredError("order")
	}

	if len(couriers) == 0 {
		return nil, errs.NewValueIsRequiredError("couriers")
	}

	bestCourier := couriers[0]
	minSteps, err := couriers[0].StepsToOrder(order.Location())
	if err != nil {
		return nil, err
	}
	for idx := range len(couriers) - 1 {
		stepsToOrder, err := couriers[idx+1].StepsToOrder(order.Location())
		if err != nil {
			return nil, err
		}

		if stepsToOrder < minSteps {
			minSteps = stepsToOrder
			bestCourier = couriers[idx+1]
		}
	}

	err = order.AssignToCourier(bestCourier.ID())
	if err != nil {
		return nil, err
	}

	err = bestCourier.SetBusy()
	if err != nil {
		return nil, err
	}

	return bestCourier, nil
}
