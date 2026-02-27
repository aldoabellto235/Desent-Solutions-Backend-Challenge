package handlers

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EchoHandler struct{}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (h *EchoHandler) Echo(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to read body"})
	}
	return c.Blob(http.StatusOK, "application/json", body)
}
