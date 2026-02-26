package book

import (
	"errors"
	"time"
)

var (
	ErrNotFound       = errors.New("book not found")
	ErrTitleRequired  = errors.New("title is required")
	ErrAuthorRequired = errors.New("author is required")
)

type Book struct {
	ID        string
	Title     string
	Author    string
	ISBN      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Book) Validate() error {
	if b.Title == "" {
		return ErrTitleRequired
	}
	if b.Author == "" {
		return ErrAuthorRequired
	}
	return nil
}
