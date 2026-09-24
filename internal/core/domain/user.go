package domain

import (
	"time"

	"github.com/google/uuid"
)

// User is a pure domain entity — no infrastructure dependencies.
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name" validate:"required,max=255"`
	Email     string    `json:"email" validate:"required,email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
