package output

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

// BookRepository is the output port — defines how the core communicates with storage.
type BookRepository interface {
	FindAll(ctx context.Context, opts query.QueryOptions) ([]domain.Book, int64, error)
	FindByID(ctx context.Context, id string) (*domain.Book, error)
	Save(ctx context.Context, book *domain.Book) error
	Update(ctx context.Context, book *domain.Book) error
	Delete(ctx context.Context, id string) error
}
