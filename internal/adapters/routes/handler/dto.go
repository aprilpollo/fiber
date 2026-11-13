package handler

import (
	"time"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"

	"github.com/google/uuid"
)

// Request and response shapes belong to the HTTP adapter, not the domain: the
// client can only send the fields listed here (never id or timestamps), and a
// change to the entity does not silently change the API.

type createUserRequest struct {
	Name  string `json:"name" validate:"required,max=255"`
	Email string `json:"email" validate:"required,email"`
}

func (r createUserRequest) toInput() input.CreateUser {
	return input.CreateUser{Name: r.Name, Email: r.Email}
}

// updateUserRequest uses pointers so an omitted field is left unchanged.
type updateUserRequest struct {
	Name  *string `json:"name" validate:"omitempty,max=255"`
	Email *string `json:"email" validate:"omitempty,email"`
}

func (r updateUserRequest) toInput() input.UpdateUser {
	return input.UpdateUser{Name: r.Name, Email: r.Email}
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toUserResponse(u *domain.User) userResponse {
	return userResponse{ID: u.ID, Name: u.Name, Email: u.Email, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt}
}

type createBookRequest struct {
	Title  string `json:"title" validate:"required,max=255"`
	Author string `json:"author" validate:"required,max=255"`
}

func (r createBookRequest) toInput() input.CreateBook {
	return input.CreateBook{Title: r.Title, Author: r.Author}
}

// updateBookRequest uses pointers so an omitted field is left unchanged.
type updateBookRequest struct {
	Title  *string `json:"title" validate:"omitempty,max=255"`
	Author *string `json:"author" validate:"omitempty,max=255"`
}

func (r updateBookRequest) toInput() input.UpdateBook {
	return input.UpdateBook{Title: r.Title, Author: r.Author}
}

type bookResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toBookResponse(b *domain.Book) bookResponse {
	return bookResponse{ID: b.ID, Title: b.Title, Author: b.Author, CreatedAt: b.CreatedAt, UpdatedAt: b.UpdatedAt}
}

func mapSlice[T, R any](in []T, f func(*T) R) []R {
	out := make([]R, len(in))
	for i := range in {
		out[i] = f(&in[i])
	}
	return out
}
