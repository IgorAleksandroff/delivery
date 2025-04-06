package http

import (
	"github.com/IgorAleksandroff/delivery/internal/api/http/problems"
	"github.com/IgorAleksandroff/delivery/internal/core/usecases/commands"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) CreateOrder(c echo.Context) error {
	createOrderCommand, err := commands.NewCreateOrderCommand(uuid.New(), "Бажная")
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	err = s.createOrderCommandHandler.Handle(c.Request().Context(), createOrderCommand)
	if err != nil {
		return problems.NewConflict(err.Error(), "/")
	}

	return c.JSON(http.StatusOK, nil)
}
