package input

import (
	"context"

	"aprilpollo/internal/core/domain"

	"github.com/google/uuid"
)

// CreateBook is the data needed to add a book.
type CreateBook struct {
	Title  string
	Author string
}

// UpdateBook changes only the fields that are set.
type UpdateBook struct {
	Title  *string
	Author *string
}

// BookService is the input port — defines what the application can do with books.
type BookService interface {
	List(ctx context.Context, q domain.ListQuery) ([]domain.Book, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error)
	Create(ctx context.Context, in CreateBook) (*domain.Book, error)
	Update(ctx context.Context, id uuid.UUID, in UpdateBook) (*domain.Book, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
