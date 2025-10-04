package service

import (
	"aprilpollo/internal/core/port"
)

type AuthService struct {
	authRepo port.AuthRepository
}

func NewAuthService(authRepo port.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}
