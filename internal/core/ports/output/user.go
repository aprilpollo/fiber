package output

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

// UserRepository is the output port — defines how the core communicates with storage.
type UserRepository interface {
	FindAll(ctx context.Context, opts query.QueryOptions) ([]domain.User, int64, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}
