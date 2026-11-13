package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Book is a pure domain entity — no infrastructure or transport concerns.
type Book struct {
	ID        uuid.UUID
	Title     string
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewBook creates a valid book with a fresh identity.
func NewBook(title, author string) (*Book, error) {
	b := &Book{ID: uuid.New(), Title: strings.TrimSpace(title), Author: strings.TrimSpace(author)}
	if err := b.Validate(); err != nil {
		return nil, err
	}
	return b, nil
}

// Validate reports the first invariant the book breaks.
func (b *Book) Validate() error {
	if err := requireText("title", b.Title); err != nil {
		return err
	}
	return requireText("author", b.Author)
}
