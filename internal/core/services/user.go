package services

import (
	"context"
	"strings"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"

	"github.com/google/uuid"
)

type userService struct {
	repo output.UserRepository
}

// NewUserService returns an input.UserService backed by the given repository.
func NewUserService(repo output.UserRepository) input.UserService {
	return &userService{repo: repo}
}

func (s *userService) List(ctx context.Context, q domain.ListQuery) ([]domain.User, int64, error) {
	return s.repo.FindAll(ctx, q)
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) Create(ctx context.Context, in input.CreateUser) (*domain.User, error) {
	user, err := domain.NewUser(in.Name, in.Email)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, in input.UpdateUser) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Name != nil {
		user.Name = strings.TrimSpace(*in.Name)
	}
	if in.Email != nil {
		user.Email = strings.TrimSpace(*in.Email)
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
