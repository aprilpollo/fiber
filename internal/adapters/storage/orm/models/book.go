package models

import (
	"time"

	"aprilpollo/internal/core/domain"

	"github.com/google/uuid"
)

type BookModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title     string    `gorm:"not null"`
	Author    string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (BookModel) TableName() string { return "books" }

func (m *BookModel) ToDomain() *domain.Book {
	return &domain.Book{
		ID:        m.ID,
		Title:     m.Title,
		Author:    m.Author,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromBookDomain(b *domain.Book) *BookModel {
	return &BookModel{
		ID:        b.ID,
		Title:     b.Title,
		Author:    b.Author,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
