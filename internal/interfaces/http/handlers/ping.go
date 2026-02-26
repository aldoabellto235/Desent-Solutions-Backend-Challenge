package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}
