package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// User is a pure domain entity — no infrastructure or transport concerns.
type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a valid user with a fresh identity.
func NewUser(name, email string) (*User, error) {
	u := &User{ID: uuid.New(), Name: strings.TrimSpace(name), Email: strings.TrimSpace(email)}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	return u, nil
}

// Validate reports the first invariant the user breaks.
func (u *User) Validate() error {
	if err := requireText("name", u.Name); err != nil {
		return err
	}
	return requireEmail("email", u.Email)
}
