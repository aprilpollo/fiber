package input

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

// BookService is the input port — defines what the application can do with books.
type BookService interface {
	List(ctx context.Context, opts query.QueryOptions) ([]domain.Book, int64, error)
	GetByID(ctx context.Context, id string) (*domain.Book, error)
	Create(ctx context.Context, book *domain.Book) error
	Update(ctx context.Context, book *domain.Book) error
	Delete(ctx context.Context, id string) error
}
