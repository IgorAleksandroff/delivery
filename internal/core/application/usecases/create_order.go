package usecases

import (
	"context"
	"errors"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
	"strings"

	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/order"
	"github.com/IgorAleksandroff/delivery/internal/core/ports"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

var OrderAlreadyExists = errors.New("order already exists")

type CreateOrderCommandHandler struct {
	orderRepository ports.OrderRepository
}

func NewCreateOrderCommandHandler(
	orderRepository ports.OrderRepository) (*CreateOrderCommandHandler, error) {
	if orderRepository == nil {
		return nil, errs.NewValueIsRequiredError("orderRepository")
	}

	return &CreateOrderCommandHandler{
		orderRepository: orderRepository}, nil
}

func (ch *CreateOrderCommandHandler) Handle(ctx context.Context, command CreateOrderCommand) error {
	if command.isEmpty() {
		return errs.NewValueIsRequiredError("add address command")
	}

	// Проверяем нет ли уже такого заказа
	orderAggregate, err := ch.orderRepository.Get(ctx, command.orderID)
	if err != nil {
		return err
	}
	if orderAggregate != nil {
		return OrderAlreadyExists
	}

	// Получили геопозицию из Geo. Пока ставим рандом значение
	location := kernel.CreateRandomLocation()

	// Изменили
	orderAggregate, err = order.NewOrder(command.orderID, location)
	if err != nil {
		return err
	}

	// Сохранили
	err = ch.orderRepository.Add(ctx, orderAggregate)
	if err != nil {
		return err
	}

	return nil
}

type CreateOrderCommand struct {
	orderID uuid.UUID
	street  string

	isSet bool
}

func NewCreateOrderCommand(orderID uuid.UUID, street string) (CreateOrderCommand, error) {
	if orderID == uuid.Nil {
		return CreateOrderCommand{}, errs.NewValueIsInvalidError("basketID")
	}
	if strings.TrimSpace(street) == "" {
		return CreateOrderCommand{}, errs.NewValueIsRequiredError("street")
	}
	return CreateOrderCommand{orderID: orderID, street: street, isSet: true}, nil
}

func (c CreateOrderCommand) isEmpty() bool {
	return !c.isSet
}
