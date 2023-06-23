package delivery

import (
	"net/http"

	"yafa/internal/repository"

	"github.com/labstack/echo/v4"
)

type ServiceHandlersI interface {
	Status(ctx echo.Context) error
	Clear(ctx echo.Context) error
}

type serviceH struct {
	serviceRepo repository.ServiceRepo
}

func NewServiceHandler(s repository.ServiceRepo) ServiceHandlersI {
	return &serviceH{serviceRepo: s}
}

func (h *serviceH) Status(ctx echo.Context) error {
	status, err := h.serviceRepo.Status()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, status)
}

func (h *serviceH) Clear(ctx echo.Context) error {
	err := h.serviceRepo.Clear()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}
