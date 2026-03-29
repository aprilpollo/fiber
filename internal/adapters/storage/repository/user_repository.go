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

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) output.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(ctx context.Context, opts query.QueryOptions) ([]domain.User, int64, error) {
	var rows []models.UserModel
	var total int64

	base := r.db.WithContext(ctx).Model(&models.UserModel{})

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.ApplyToGorm(r.db.WithContext(ctx).Model(&models.UserModel{}), opts).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	users := make([]domain.User, len(rows))
	for i, row := range rows {
		users[i] = *row.ToDomain()
	}

	return users, total, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var row models.UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return row.ToDomain(), nil
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
	row := models.FromUserDomain(user)
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return err
	}
	// sync generated fields (id, created_at, updated_at) back to domain
	*user = *row.ToDomain()
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	row := models.FromUserDomain(user)
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return err
	}
	*user = *row.ToDomain()
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.UserModel{}).Error
}
