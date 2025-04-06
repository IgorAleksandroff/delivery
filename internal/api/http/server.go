package http

import (
	"github.com/IgorAleksandroff/delivery/internal/core/application/usecases/commands"
	"github.com/IgorAleksandroff/delivery/internal/core/application/usecases/queries"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

type Server struct {
	createOrderCommandHandler *commands.CreateOrderCommandHandler

	getAllCouriersQueryHandler        *queries.GetAllCouriersQueryHandler
	getNotCompletedOrdersQueryHandler *queries.GetNotCompletedOrdersQueryHandler
}

func NewServer(
	createOrderCommandHandler *commands.CreateOrderCommandHandler,

	getAllCouriersQueryHandler *queries.GetAllCouriersQueryHandler,
	getNotCompletedOrdersQueryHandler *queries.GetNotCompletedOrdersQueryHandler,
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
