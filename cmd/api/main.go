package main

import (
	"context"
	"log"
	"net/http"

	"api-quest/config"
	bookapp "api-quest/internal/application/book"
	"api-quest/internal/infrastructure/mongodb"
	"api-quest/internal/interfaces/http/handlers"
	"api-quest/internal/interfaces/http/routes"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	// Infrastructure: MongoDB
	client, err := mongodb.NewClient(cfg.MongoURI)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(cfg.DBName)
	bookRepo := mongodb.NewBookRepository(db)

	// Application layer
	bookService := bookapp.NewService(bookRepo)

	// Handlers (interfaces layer)
	pingH := handlers.NewPingHandler()
	echoH := handlers.NewEchoHandler()
	bookH := handlers.NewBookHandler(bookService)
	authH := handlers.NewAuthHandler(cfg.JWTSecret, cfg.AuthUsername, cfg.AuthPassword)

	// Echo server
	e := echo.New()
	e.HideBanner = true
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	// Global error handler for consistent JSON errors
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		msg := "Internal server error"
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			if m, ok := he.Message.(string); ok {
				msg = m
			}
		}
		_ = c.JSON(code, map[string]string{"error": msg})
	}

	routes.Register(e, pingH, echoH, bookH, authH, cfg.JWTSecret)

	log.Printf("Server starting on :%s", cfg.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
