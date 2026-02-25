package postgres

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/entity"
	"gorm.io/gorm"
)

// SettingRepository implements repository.SettingRepository
type SettingRepository struct {
	db *gorm.DB
}

// NewSettingRepository creates a new setting repository
func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

// GetByKey retrieves a setting by key
func (r *SettingRepository) GetByKey(ctx context.Context, key string) (*entity.Setting, error) {
	var setting entity.Setting
	if err := r.db.WithContext(ctx).First(&setting, "key = ?", key).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// GetAll retrieves all settings
func (r *SettingRepository) GetAll(ctx context.Context) ([]*entity.Setting, error) {
	var settings []*entity.Setting
	if err := r.db.WithContext(ctx).Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

// Set creates or updates a setting
func (r *SettingRepository) Set(ctx context.Context, setting *entity.Setting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}

// Delete deletes a setting
func (r *SettingRepository) Delete(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).Delete(&entity.Setting{}, "key = ?", key).Error
}

// PrintTemplateRepository implements repository.PrintTemplateRepository
type PrintTemplateRepository struct {
	db *gorm.DB
}

// NewPrintTemplateRepository creates a new print template repository
func NewPrintTemplateRepository(db *gorm.DB) *PrintTemplateRepository {
	return &PrintTemplateRepository{db: db}
}

// Create creates a new print template
func (r *PrintTemplateRepository) Create(ctx context.Context, template *entity.PrintTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

// GetByID retrieves a print template by ID
func (r *PrintTemplateRepository) GetByID(ctx context.Context, id string) (*entity.PrintTemplate, error) {
	var template entity.PrintTemplate
	if err := r.db.WithContext(ctx).First(&template, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

// GetDefault retrieves the default print template
func (r *PrintTemplateRepository) GetDefault(ctx context.Context) (*entity.PrintTemplate, error) {
	var template entity.PrintTemplate
	if err := r.db.WithContext(ctx).Where("is_default = ?", true).First(&template).Error; err != nil {
		// Return nil if no default template
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

// GetAll retrieves all print templates
func (r *PrintTemplateRepository) GetAll(ctx context.Context) ([]*entity.PrintTemplate, error) {
	var templates []*entity.PrintTemplate
	if err := r.db.WithContext(ctx).
		Order("is_default DESC, created_at DESC").
		Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// Update updates a print template
func (r *PrintTemplateRepository) Update(ctx context.Context, template *entity.PrintTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

// UpdateSetDefault sets a template as default and unsets others
func (r *PrintTemplateRepository) UpdateSetDefault(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset all defaults
		if err := tx.Model(&entity.PrintTemplate{}).
			Where("is_default = ?", true).
			Update("is_default", false).Error; err != nil {
			return err
		}

		// Set new default
		return tx.Model(&entity.PrintTemplate{}).
			Where("id = ?", id).
			Update("is_default", true).Error
	})
}

// Delete deletes a print template
func (r *PrintTemplateRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.PrintTemplate{}, "id = ?", id).Error
}
