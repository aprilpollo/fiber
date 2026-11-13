package services

import (
	"context"
	"strings"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"

	"github.com/google/uuid"
)

type bookService struct {
	repo output.BookRepository
}

// NewBookService returns an input.BookService backed by the given repository.
func NewBookService(repo output.BookRepository) input.BookService {
	return &bookService{repo: repo}
}

func (s *bookService) List(ctx context.Context, q domain.ListQuery) ([]domain.Book, int64, error) {
	return s.repo.FindAll(ctx, q)
}

func (s *bookService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *bookService) Create(ctx context.Context, in input.CreateBook) (*domain.Book, error) {
	book, err := domain.NewBook(in.Title, in.Author)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, book); err != nil {
		return nil, err
	}
	return book, nil
}

func (s *bookService) Update(ctx context.Context, id uuid.UUID, in input.UpdateBook) (*domain.Book, error) {
	book, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Title != nil {
		book.Title = strings.TrimSpace(*in.Title)
	}
	if in.Author != nil {
		book.Author = strings.TrimSpace(*in.Author)
	}
	if err := book.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, book); err != nil {
		return nil, err
	}
	return book, nil
}

func (s *bookService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
