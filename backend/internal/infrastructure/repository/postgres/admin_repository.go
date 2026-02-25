package postgres

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/entity"
	"gorm.io/gorm"
)

// AdminUserRepository implements repository.AdminUserRepository
type AdminUserRepository struct {
	db *gorm.DB
}

// NewAdminUserRepository creates a new admin user repository
func NewAdminUserRepository(db *gorm.DB) *AdminUserRepository {
	return &AdminUserRepository{db: db}
}

// Create creates a new admin user
func (r *AdminUserRepository) Create(ctx context.Context, user *entity.AdminUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves an admin user by ID
func (r *AdminUserRepository) GetByID(ctx context.Context, id string) (*entity.AdminUser, error) {
	var user entity.AdminUser
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves an admin user by username
func (r *AdminUserRepository) GetByUsername(ctx context.Context, username string) (*entity.AdminUser, error) {
	var user entity.AdminUser
	if err := r.db.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll retrieves all admin users
func (r *AdminUserRepository) GetAll(ctx context.Context) ([]*entity.AdminUser, error) {
	var users []*entity.AdminUser
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Update updates an admin user
func (r *AdminUserRepository) Update(ctx context.Context, user *entity.AdminUser) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// UpdateLastLogin updates the last login timestamp
func (r *AdminUserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&entity.AdminUser{}).
		Where("id = ?", id).
		Update("last_login", gorm.Expr("NOW()")).Error
}

// Delete soft deletes an admin user
func (r *AdminUserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.AdminUser{}, "id = ?", id).Error
}
