package mail

import (
	"net/http"

	"github.com/MartyHub/mail-api/db"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Create(c echo.Context) error
	Get(c echo.Context) error
}

func NewHandler(repo db.Repository) Handler {
	return handler{
		service: newService(repo),
	}
}

type handler struct {
	service Service
}

func (h handler) Create(c echo.Context) error {
	input := CreateInput{}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "")
	}

	if err := c.Validate(input); err != nil {
		return echo.ErrBadRequest
	}

	entity, err := h.service.Create(c.Request().Context(), input)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, entity)
}

func (h handler) Get(c echo.Context) error {
	input := GetInput{}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "")
	}

	if err := c.Validate(input); err != nil {
		return echo.ErrBadRequest
	}

	entity, err := h.service.Get(c.Request().Context(), input)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, entity)
}
