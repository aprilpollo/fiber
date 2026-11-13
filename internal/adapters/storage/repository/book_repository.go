package repository

import (
	"context"

	"aprilpollo/internal/adapters/storage/orm/models"
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/output"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var bookColumns = columns{
	"id":         "id",
	"title":      "title",
	"author":     "author",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) output.BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) FindAll(ctx context.Context, q domain.ListQuery) ([]domain.Book, int64, error) {
	where, page, err := listScopes(q, bookColumns)
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := r.db.WithContext(ctx).Model(&models.BookModel{}).Scopes(where).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []models.BookModel
	if err := r.db.WithContext(ctx).Scopes(where, page).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	books := make([]domain.Book, len(rows))
	for i := range rows {
		books[i] = *rows[i].ToDomain()
	}
	return books, total, nil
}

func (r *bookRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	var row models.BookModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		return nil, translate(err)
	}
	return row.ToDomain(), nil
}

func (r *bookRepository) Save(ctx context.Context, book *domain.Book) error {
	row := models.FromBookDomain(book)
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return translate(err)
	}
	// sync the timestamps storage set back to the domain
	*book = *row.ToDomain()
	return nil
}

func (r *bookRepository) Update(ctx context.Context, book *domain.Book) error {
	row := models.FromBookDomain(book)
	// Updates, not Save: Save turns a missing row into an INSERT.
	res := r.db.WithContext(ctx).Model(row).Select("*").Omit("id", "created_at").Updates(row)
	if res.Error != nil {
		return translate(res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	*book = *row.ToDomain()
	return nil
}

func (r *bookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.BookModel{})
	if res.Error != nil {
		return translate(res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
