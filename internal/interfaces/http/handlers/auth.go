package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	jwtSecret    string
	authUsername string
	authPassword string
}

func NewAuthHandler(jwtSecret, username, password string) *AuthHandler {
	return &AuthHandler{
		jwtSecret:    jwtSecret,
		authUsername: username,
		authPassword: password,
	}
}

type tokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) GenerateToken(c echo.Context) error {
	var req tokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Username != h.authUsername || req.Password != h.authPassword {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": req.Username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": signed})
}
