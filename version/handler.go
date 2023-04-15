package version

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	Get(c echo.Context) error
}

func NewHandler() Handler {
	return handler{
		service: newService(),
	}
}

type handler struct {
	service Service
}

func (h handler) Get(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, h.service.Version(), "  ")
}
