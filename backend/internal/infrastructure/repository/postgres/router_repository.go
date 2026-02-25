package postgres

import (
	"context"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/entity"
	"gorm.io/gorm"
)

// RouterRepository implements repository.RouterRepository
type RouterRepository struct {
	db *gorm.DB
}

// NewRouterRepository creates a new router repository
func NewRouterRepository(db *gorm.DB) *RouterRepository {
	return &RouterRepository{db: db}
}

// Create creates a new router
func (r *RouterRepository) Create(ctx context.Context, router *entity.Router) error {
	return r.db.WithContext(ctx).Create(router).Error
}

// GetByID retrieves a router by ID
func (r *RouterRepository) GetByID(ctx context.Context, id uint) (*entity.Router, error) {
	var router entity.Router
	if err := r.db.WithContext(ctx).First(&router, id).Error; err != nil {
		return nil, err
	}
	return &router, nil
}

// GetByName retrieves a router by name
func (r *RouterRepository) GetByName(ctx context.Context, name string) (*entity.Router, error) {
	var router entity.Router
	if err := r.db.WithContext(ctx).First(&router, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &router, nil
}

// GetAll retrieves all routers
func (r *RouterRepository) GetAll(ctx context.Context) ([]*entity.Router, error) {
	var routers []*entity.Router
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&routers).Error; err != nil {
		return nil, err
	}
	return routers, nil
}

// Update updates a router
func (r *RouterRepository) Update(ctx context.Context, router *entity.Router) error {
	return r.db.WithContext(ctx).Save(router).Error
}

// Delete soft deletes a router
func (r *RouterRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Router{}, id).Error
}

// UpdateLastConnected updates the last connected timestamp
func (r *RouterRepository) UpdateLastConnected(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.Router{}).
		Where("id = ?", id).
		Update("last_connected", now).Error
}
