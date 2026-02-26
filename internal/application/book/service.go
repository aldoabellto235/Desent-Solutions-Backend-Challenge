package bookapp

import (
	"context"
	"time"

	"api-quest/internal/domain/book"
)

type Service struct {
	repo book.Repository
}

func NewService(repo book.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, b *book.Book) (*book.Book, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return s.repo.Create(ctx, b)
}

func (s *Service) List(ctx context.Context, filter book.Filter) ([]*book.Book, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	return s.repo.FindAll(ctx, filter)
}

func (s *Service) GetByID(ctx context.Context, id string) (*book.Book, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id string, b *book.Book) (*book.Book, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}
	b.UpdatedAt = time.Now()
	return s.repo.Update(ctx, id, b)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
