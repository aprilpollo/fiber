package output

import (
	"context"

	"aprilpollo/internal/core/domain"

	"github.com/google/uuid"
)

// UserRepository is the output port — defines how the core communicates with storage.
//
// FindByID, Update and Delete return domain.ErrNotFound for a missing user;
// Save and Update return domain.ErrConflict when the email is already taken.
type UserRepository interface {
	FindAll(ctx context.Context, q domain.ListQuery) ([]domain.User, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
