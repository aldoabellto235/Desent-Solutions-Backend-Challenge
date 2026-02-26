package book

import "context"

type Filter struct {
	Author string
	Page   int
	Limit  int
}

type Repository interface {
	Create(ctx context.Context, book *Book) (*Book, error)
	FindAll(ctx context.Context, filter Filter) ([]*Book, int64, error)
	FindByID(ctx context.Context, id string) (*Book, error)
	Update(ctx context.Context, id string, book *Book) (*Book, error)
	Delete(ctx context.Context, id string) error
}
