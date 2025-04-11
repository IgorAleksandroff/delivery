package commands

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/orders"
	"github.com/IgorAleksandroff/delivery/internal/core/ports"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

var OrderAlreadyExists = errors.New("order already exists")

type CreateOrderCommandHandler struct {
	orderRepository ports.OrderRepository
	geoClient       ports.GeoClient
}

func NewCreateOrderCommandHandler(
	orderRepository ports.OrderRepository, geoClient ports.GeoClient) (*CreateOrderCommandHandler, error) {
	if orderRepository == nil {
		return nil, errs.NewValueIsRequiredError("orderRepository")
	}

	return &CreateOrderCommandHandler{
		orderRepository: orderRepository,
		geoClient:       geoClient}, nil
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

	// Получили геопозицию из Geo.
	location, err := ch.geoClient.GetGeolocation(ctx, command.Street())
	if err != nil {
		return err
	}

	// Изменили
	orderAggregate, err = orders.NewOrder(command.orderID, location)
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

func (c CreateOrderCommand) Street() string {
	return c.street
}

func (c CreateOrderCommand) isEmpty() bool {
	return !c.isSet
}
