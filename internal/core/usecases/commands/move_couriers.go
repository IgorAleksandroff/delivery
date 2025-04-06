package commands

import (
	"context"
	"errors"
	"log"

	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
	"github.com/IgorAleksandroff/delivery/internal/pkg/uow"
)

type MoveCouriersCommandHandler struct {
	unitOfWork        uow.UnitOfWork
	orderRepository   OrderRepository
	courierRepository CourierRepository
}

func NewMoveCouriersCommandHandler(
	unitOfWork uow.UnitOfWork,
	orderRepository OrderRepository,
	courierRepository CourierRepository,
) (*MoveCouriersCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}
	if orderRepository == nil {
		return nil, errs.NewValueIsRequiredError("orderRepository")
	}
	if courierRepository == nil {
		return nil, errs.NewValueIsRequiredError("courierRepository")
	}

	return &MoveCouriersCommandHandler{
		unitOfWork:        unitOfWork,
		orderRepository:   orderRepository,
		courierRepository: courierRepository}, nil
}

func (ch *MoveCouriersCommandHandler) Handle(ctx context.Context, command MoveCouriersCommand) (err error) {
	if command.isEmpty() {
		return errs.NewValueIsRequiredError("add address command")
	}

	// Восстановили
	assignedOrders, err := ch.orderRepository.GetAllInAssignedStatus(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return nil
		}
		return err
	}

	// Изменили и сохранили
	ctx = ch.unitOfWork.Begin(ctx)
	defer func() {
		if err != nil {
			errRollback := ch.unitOfWork.Rollback(ctx)
			if errRollback != nil {
				if err != nil {
					log.Println("MoveCouriersCommandHandler Rollback error:", err)
				}
			}
		}
	}()

	for _, assignedOrder := range assignedOrders {
		courier, err := ch.courierRepository.Get(ctx, *assignedOrder.AssignedCourier())
		if err != nil {
			if errors.Is(err, errs.ErrObjectNotFound) {
				log.Printf("AssignedCourier %v for order %v not found", *assignedOrder.AssignedCourier(), assignedOrder.ID())
			}
			return err
		}

		err = courier.Move(assignedOrder.Location())
		if err != nil {
			return err
		}

		if courier.Location().Equals(assignedOrder.Location()) {
			err := assignedOrder.Complete()
			if err != nil {
				return err
			}
			err = courier.SetFree()
			if err != nil {
				return err
			}
		}

		err = ch.orderRepository.Update(ctx, assignedOrder)
		if err != nil {
			return err
		}
		err = ch.courierRepository.Update(ctx, courier)
		if err != nil {
			return err
		}
	}

	err = ch.unitOfWork.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

type MoveCouriersCommand struct {
	isSet bool
}

func NewMoveCouriersCommand() (MoveCouriersCommand, error) {
	return MoveCouriersCommand{isSet: true}, nil
}
func (c MoveCouriersCommand) isEmpty() bool {
	return !c.isSet
}
