package services

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
)

type bookService struct {
	repo output.BookRepository
}

// NewBookService returns an input.BookService backed by the given repository.
func NewBookService(repo output.BookRepository) input.BookService {
	return &bookService{repo: repo}
}

func (s *bookService) List(ctx context.Context, opts query.QueryOptions) ([]domain.Book, int64, error) {
	return s.repo.FindAll(ctx, opts)
}

func (s *bookService) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *bookService) Create(ctx context.Context, book *domain.Book) error {
	return s.repo.Save(ctx, book)
}

func (s *bookService) Update(ctx context.Context, book *domain.Book) error {
	return s.repo.Update(ctx, book)
}

func (s *bookService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
