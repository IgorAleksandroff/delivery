package commands

import (
	"context"
	"errors"
	"log"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/services"
	"github.com/IgorAleksandroff/delivery/internal/core/ports"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
	"github.com/IgorAleksandroff/delivery/internal/pkg/uow"
)

var (
	NotAvailableOrders   = errors.New("not available orders")
	NotAvailableCouriers = errors.New("not available couriers")
)

type AssignOrdersCommandHandler struct {
	unitOfWork        uow.UnitOfWork
	orderRepository   ports.OrderRepository
	courierRepository ports.CourierRepository
	orderDispatcher   *services.Dispatcher
}

func NewAssignOrdersCommandHandler(
	unitOfWork uow.UnitOfWork,
	orderRepository ports.OrderRepository,
	courierRepository ports.CourierRepository,
	orderDispatcher *services.Dispatcher,
) (*AssignOrdersCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}
	if orderRepository == nil {
		return nil, errs.NewValueIsRequiredError("orderRepository")
	}
	if courierRepository == nil {
		return nil, errs.NewValueIsRequiredError("courierRepository")
	}
	if orderDispatcher == nil {
		return nil, errs.NewValueIsRequiredError("orderDispatcher")
	}

	return &AssignOrdersCommandHandler{
		unitOfWork:        unitOfWork,
		orderRepository:   orderRepository,
		courierRepository: courierRepository,
		orderDispatcher:   orderDispatcher}, nil
}

func (ch *AssignOrdersCommandHandler) Handle(ctx context.Context, command AssignOrdersCommand) (err error) {
	if command.isEmpty() {
		return errs.NewValueIsRequiredError("add address command")
	}

	// Восстановили
	orderAggregate, err := ch.orderRepository.GetFirstInCreatedStatus(ctx)
	if err != nil {
		return err
	}
	if orderAggregate == nil {
		return NotAvailableOrders
	}

	couriers, err := ch.courierRepository.GetAllInFreeStatus(ctx)
	if err != nil {
		return err
	}
	if len(couriers) == 0 {
		return NotAvailableCouriers
	}

	// Изменили
	courier, err := ch.orderDispatcher.Dispatch(orderAggregate, couriers)
	if err != nil {
		return err
	}

	// Сохранили
	ctx = ch.unitOfWork.Begin(ctx)
	defer func() {
		if err != nil {
			errRollback := ch.unitOfWork.Rollback(ctx)
			if errRollback != nil {
				log.Println("AssignOrdersCommandHandler Rollback error:", errRollback)
			}
		}
	}()

	err = ch.orderRepository.Update(ctx, orderAggregate)
	if err != nil {
		return err
	}
	err = ch.courierRepository.Update(ctx, courier)
	if err != nil {
		return err
	}

	err = ch.unitOfWork.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

type AssignOrdersCommand struct {
	isSet bool
}

func NewAssignOrdersCommand() (AssignOrdersCommand, error) {
	return AssignOrdersCommand{isSet: true}, nil
}
func (c AssignOrdersCommand) isEmpty() bool {
	return !c.isSet
}
