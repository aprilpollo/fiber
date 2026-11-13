package output

import (
	"context"

	"aprilpollo/internal/core/domain"

	"github.com/google/uuid"
)

// BookRepository is the output port — defines how the core communicates with storage.
//
// FindByID, Update and Delete return domain.ErrNotFound for a missing book.
type BookRepository interface {
	FindAll(ctx context.Context, q domain.ListQuery) ([]domain.Book, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Book, error)
	Save(ctx context.Context, book *domain.Book) error
	Update(ctx context.Context, book *domain.Book) error
	Delete(ctx context.Context, id uuid.UUID) error
}
