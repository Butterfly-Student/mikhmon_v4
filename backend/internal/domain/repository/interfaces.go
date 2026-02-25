package repository

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/entity"
)

// AdminUserRepository defines the interface for admin user data access
type AdminUserRepository interface {
	Create(ctx context.Context, user *entity.AdminUser) error
	GetByID(ctx context.Context, id string) (*entity.AdminUser, error)
	GetByUsername(ctx context.Context, username string) (*entity.AdminUser, error)
	GetAll(ctx context.Context) ([]*entity.AdminUser, error)
	Update(ctx context.Context, user *entity.AdminUser) error
	UpdateLastLogin(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

// RouterRepository defines the interface for router configuration data access
type RouterRepository interface {
	Create(ctx context.Context, router *entity.Router) error
	GetByID(ctx context.Context, id uint) (*entity.Router, error)
	GetByName(ctx context.Context, name string) (*entity.Router, error)
	GetAll(ctx context.Context) ([]*entity.Router, error)
	Update(ctx context.Context, router *entity.Router) error
	Delete(ctx context.Context, id uint) error
	UpdateLastConnected(ctx context.Context, id uint) error
}

// SettingRepository defines the interface for application settings
type SettingRepository interface {
	GetByKey(ctx context.Context, key string) (*entity.Setting, error)
	GetAll(ctx context.Context) ([]*entity.Setting, error)
	Set(ctx context.Context, setting *entity.Setting) error
	Delete(ctx context.Context, key string) error
}

// PrintTemplateRepository defines the interface for print templates
type PrintTemplateRepository interface {
	Create(ctx context.Context, template *entity.PrintTemplate) error
	GetByID(ctx context.Context, id string) (*entity.PrintTemplate, error)
	GetDefault(ctx context.Context) (*entity.PrintTemplate, error)
	GetAll(ctx context.Context) ([]*entity.PrintTemplate, error)
	Update(ctx context.Context, template *entity.PrintTemplate) error
	UpdateSetDefault(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}
