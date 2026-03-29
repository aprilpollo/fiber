package repository

import (
	"context"
	"errors"

	"aprilpollo/internal/adapters/storage/orm/models"
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"

	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) output.BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) FindAll(ctx context.Context, opts query.QueryOptions) ([]domain.Book, int64, error) {
	var rows []models.BookModel
	var total int64

	base := r.db.WithContext(ctx).Model(&models.BookModel{})

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.ApplyToGorm(r.db.WithContext(ctx).Model(&models.BookModel{}), opts).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	books := make([]domain.Book, len(rows))
	for i, row := range rows {
		books[i] = *row.ToDomain()
	}

	return books, total, nil
}

func (r *bookRepository) FindByID(ctx context.Context, id string) (*domain.Book, error) {
	var row models.BookModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return row.ToDomain(), nil
}

func (r *bookRepository) Save(ctx context.Context, book *domain.Book) error {
	row := models.FromBookDomain(book)
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return err
	}
	// sync generated fields (id, created_at, updated_at) back to domain
	*book = *row.ToDomain()
	return nil
}

func (r *bookRepository) Update(ctx context.Context, book *domain.Book) error {
	row := models.FromBookDomain(book)
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return err
	}
	*book = *row.ToDomain()
	return nil
}

func (r *bookRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.BookModel{}).Error
}
