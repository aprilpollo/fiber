package routes

import (
	"aprilpollo/internal/core/port"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService port.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService port.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}
