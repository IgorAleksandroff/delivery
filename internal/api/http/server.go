package http

import (
	"github.com/IgorAleksandroff/delivery/internal/core/usecases/commands"
	queries2 "github.com/IgorAleksandroff/delivery/internal/core/usecases/queries"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

type Server struct {
	createOrderCommandHandler *commands.CreateOrderCommandHandler

	getAllCouriersQueryHandler        *queries2.GetAllCouriersQueryHandler
	getNotCompletedOrdersQueryHandler *queries2.GetNotCompletedOrdersQueryHandler
}

func NewServer(
	createOrderCommandHandler *commands.CreateOrderCommandHandler,

	getAllCouriersQueryHandler *queries2.GetAllCouriersQueryHandler,
	getNotCompletedOrdersQueryHandler *queries2.GetNotCompletedOrdersQueryHandler,
) (*Server, error) {
	if createOrderCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("createOrderCommandHandler")
	}
	if getAllCouriersQueryHandler == nil {
		return nil, errs.NewValueIsRequiredError("getAllCouriersQueryHandler")
	}
	if getNotCompletedOrdersQueryHandler == nil {
		return nil, errs.NewValueIsRequiredError("getNotCompletedOrdersQueryHandler")
	}
	return &Server{
		createOrderCommandHandler: createOrderCommandHandler,

		getAllCouriersQueryHandler:        getAllCouriersQueryHandler,
		getNotCompletedOrdersQueryHandler: getNotCompletedOrdersQueryHandler,
	}, nil
}
