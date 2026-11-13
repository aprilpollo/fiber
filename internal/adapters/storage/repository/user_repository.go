package repository

import (
	"context"

	"aprilpollo/internal/adapters/storage/orm/models"
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/output"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var userColumns = columns{
	"id":         "id",
	"name":       "name",
	"email":      "email",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) output.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(ctx context.Context, q domain.ListQuery) ([]domain.User, int64, error) {
	where, page, err := listScopes(q, userColumns)
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := r.db.WithContext(ctx).Model(&models.UserModel{}).Scopes(where).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []models.UserModel
	if err := r.db.WithContext(ctx).Scopes(where, page).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	users := make([]domain.User, len(rows))
	for i := range rows {
		users[i] = *rows[i].ToDomain()
	}
	return users, total, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var row models.UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		return nil, translate(err)
	}
	return row.ToDomain(), nil
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
	row := models.FromUserDomain(user)
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return translate(err)
	}
	// sync the timestamps storage set back to the domain
	*user = *row.ToDomain()
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	row := models.FromUserDomain(user)
	// Updates, not Save: Save turns a missing row into an INSERT.
	res := r.db.WithContext(ctx).Model(row).Select("*").Omit("id", "created_at").Updates(row)
	if res.Error != nil {
		return translate(res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	*user = *row.ToDomain()
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.UserModel{})
	if res.Error != nil {
		return translate(res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
