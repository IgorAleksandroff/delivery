package http

import (
	"errors"
	"github.com/IgorAleksandroff/delivery/internal/api/http/problems"
	"github.com/IgorAleksandroff/delivery/internal/core/application/usecases/queries"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
	servers "github.com/IgorAleksandroff/delivery/pkg/servers"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) GetOrders(c echo.Context) error {
	query, err := queries.NewGetNotCompletedOrdersQuery()
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	response, err := s.getNotCompletedOrdersQueryHandler.Handle(query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, problems.NewNotFound(err.Error()))
		}
	}

	var orders []servers.Order
	for _, courier := range response.Orders {
		location := servers.Location{
			X: courier.Location.X,
			Y: courier.Location.Y,
		}

		var courier = servers.Order{
			Id:       courier.ID,
			Location: location,
		}
		orders = append(orders, courier)
	}
	return c.JSON(http.StatusOK, orders)
}
