package health

import (
	"net/http"

	"github.com/MartyHub/mail-api/db"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Get(c echo.Context) error
}

func NewHandler(repo db.Repository) Handler {
	return handler{service: newService(repo)}
}

type handler struct {
	service Service
}

func (h handler) Get(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, h.service.Health(c.Request().Context()), "  ")
}
