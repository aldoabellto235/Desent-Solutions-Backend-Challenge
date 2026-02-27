package routes

import (
	"api-quest/internal/interfaces/http/handlers"
"github.com/labstack/echo/v4"
)

func Register(
	e *echo.Echo,
	ping *handlers.PingHandler,
	echoH *handlers.EchoHandler,
	book *handlers.BookHandler,
	auth *handlers.AuthHandler,
	jwtSecret string,
) {
	// Level 1
	e.GET("/ping", ping.Ping)

	// Level 2
	e.POST("/echo", echoH.Echo)

	// Level 5 — auth token endpoint
	e.POST("/auth/token", auth.GenerateToken)

	// Books group — Level 3, 4, 6, 7
	b := e.Group("/books")
	b.POST("", book.Create)
	b.GET("", book.List)
	b.GET("/:id", book.GetByID)
	b.PUT("/:id", book.Update)
	b.DELETE("/:id", book.Delete)
}
