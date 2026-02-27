package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	bookapp "api-quest/internal/application/book"
	"api-quest/internal/domain/book"

	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	service *bookapp.Service
}

func NewBookHandler(service *bookapp.Service) *BookHandler {
	return &BookHandler{service: service}
}

type bookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
}

type bookResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	ISBN      string    `json:"isbn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toResponse(b *book.Book) bookResponse {
	return bookResponse{
		ID:        b.ID,
		Title:     b.Title,
		Author:    b.Author,
		ISBN:      b.ISBN,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

func (h *BookHandler) Create(c echo.Context) error {
	var req bookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	b := &book.Book{
		Title:  req.Title,
		Author: req.Author,
		ISBN:   req.ISBN,
	}

	created, err := h.service.Create(c.Request().Context(), b)
	if errors.Is(err, book.ErrTitleRequired) || errors.Is(err, book.ErrAuthorRequired) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create book"})
	}

	return c.JSON(http.StatusCreated, toResponse(created))
}

func (h *BookHandler) List(c echo.Context) error {
	author := c.QueryParam("author")
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page := 1
	limit := 10

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	filter := book.Filter{Author: author, Page: page, Limit: limit}
	books, _, err := h.service.List(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list books"})
	}

	resp := make([]bookResponse, len(books))
	for i, b := range books {
		resp[i] = toResponse(b)
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *BookHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	b, err := h.service.GetByID(c.Request().Context(), id)
	if errors.Is(err, book.ErrNotFound) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Book not found"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get book"})
	}
	return c.JSON(http.StatusOK, toResponse(b))
}

func (h *BookHandler) Update(c echo.Context) error {
	id := c.Param("id")

	var req bookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	b := &book.Book{
		Title:  req.Title,
		Author: req.Author,
		ISBN:   req.ISBN,
	}

	updated, err := h.service.Update(c.Request().Context(), id, b)
	if errors.Is(err, book.ErrNotFound) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Book not found"})
	}
	if errors.Is(err, book.ErrTitleRequired) || errors.Is(err, book.ErrAuthorRequired) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update book"})
	}

	return c.JSON(http.StatusOK, toResponse(updated))
}

func (h *BookHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	err := h.service.Delete(c.Request().Context(), id)
	if errors.Is(err, book.ErrNotFound) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Book not found"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete book"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Book deleted"})
}
