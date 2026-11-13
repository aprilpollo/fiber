package input

import (
	"context"

	"aprilpollo/internal/core/domain"

	"github.com/google/uuid"
)

// CreateUser is the data needed to register a user.
type CreateUser struct {
	Name  string
	Email string
}

// UpdateUser changes only the fields that are set.
type UpdateUser struct {
	Name  *string
	Email *string
}

// UserService is the input port — defines what the application can do with users.
type UserService interface {
	List(ctx context.Context, q domain.ListQuery) ([]domain.User, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Create(ctx context.Context, in CreateUser) (*domain.User, error)
	Update(ctx context.Context, id uuid.UUID, in UpdateUser) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
