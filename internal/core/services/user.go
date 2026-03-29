package services

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
)

type userService struct {
	repo output.UserRepository
}

// NewUserService returns an input.UserService backed by the given repository.
func NewUserService(repo output.UserRepository) input.UserService {
	return &userService{repo: repo}
}

func (s *userService) List(ctx context.Context, opts query.QueryOptions) ([]domain.User, int64, error) {
	return s.repo.FindAll(ctx, opts)
}

func (s *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) Create(ctx context.Context, user *domain.User) error {
	return s.repo.Save(ctx, user)
}

func (s *userService) Update(ctx context.Context, user *domain.User) error {
	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
