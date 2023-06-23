package delivery

import (
	"net/http"

	"yafa/internal/repository"

	"github.com/labstack/echo/v4"
)

type ServiceHandler interface {
	Status(ctx echo.Context) error
	Clear(ctx echo.Context) error
}

type serviceHandler struct {
	serviceRepo repository.ServiceRepo
}

func NewServiceHandler(s repository.ServiceRepo) ServiceHandler {
	return &serviceHandler{serviceRepo: s}
}

func (h *serviceHandler) Status(ctx echo.Context) error {
	status, err := h.serviceRepo.Status()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, status)
}

func (h *serviceHandler) Clear(ctx echo.Context) error {
	err := h.serviceRepo.Clear()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}
