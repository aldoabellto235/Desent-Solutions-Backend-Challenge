package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type EchoHandler struct{}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (h *EchoHandler) Echo(c echo.Context) error {
	var body interface{}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON body"})
	}
	return c.JSON(http.StatusOK, body)
}
