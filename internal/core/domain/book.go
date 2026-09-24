package domain

import (
	"time"

	"github.com/google/uuid"
)

// Book is a pure domain entity — no infrastructure dependencies.
type Book struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title" validate:"required,max=255"`
	Author    string    `json:"author" validate:"required,max=255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
