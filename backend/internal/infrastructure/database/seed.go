package database

import (
	"fmt"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed creates default data if not exists
func Seed(db *gorm.DB) error {
	// Seed default admin user
	if err := seedAdminUser(db); err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}

	// Seed default router
	if err := seedDefaultRouter(db); err != nil {
		return fmt.Errorf("failed to seed default router: %w", err)
	}

	return nil
}

// seedAdminUser creates default admin user if no users exist
func seedAdminUser(db *gorm.DB) error {
	var count int64
	if err := db.Model(&entity.AdminUser{}).Count(&count).Error; err != nil {
		return err
	}

	// Only create default user if no users exist
	if count == 0 {
		// Hash default password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		now := time.Now()
		admin := &entity.AdminUser{
			Username:     "admin",
			Password:     string(hashedPassword),
			Email:        "admin@mikhmon.local",
			IsActive:     true,
			LastLogin:    nil,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		if err := db.Create(admin).Error; err != nil {
			return fmt.Errorf("failed to create default admin: %w", err)
		}

		fmt.Println("✓ Default admin user created: admin / admin123")
	}

	return nil
}

// seedDefaultRouter creates a default router if no routers exist
func seedDefaultRouter(db *gorm.DB) error {
	var count int64
	if err := db.Model(&entity.Router{}).Count(&count).Error; err != nil {
		return err
	}

	// Only create default router if no routers exist
	if count == 0 {
		now := time.Now()
		router := &entity.Router{
			Name:        "MikroTik-1",
			Host:        "192.168.233.1",
			Port:        8728,
			Username:    "admin",
			Password:    "r00t", // Default password, should be changed
			UseSSL:      false,
			Timeout:     3,
			IsActive:    true,
			Description: "Default router - please update with your MikroTik credentials",
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if err := db.Create(router).Error; err != nil {
			return fmt.Errorf("failed to create default router: %w", err)
		}

		fmt.Println("✓ Default router created: MikroTik-1 (192.168.88.1)")
		fmt.Println("  Please update the router credentials in the settings")
	}

	return nil
}
