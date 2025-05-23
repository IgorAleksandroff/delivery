package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/IgorAleksandroff/delivery/internal/adapters/in/http/problems"
	"github.com/IgorAleksandroff/delivery/internal/core/application/usecases/commands"
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
