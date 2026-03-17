//go:build modern

package entity

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;size:100;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	Role         string `gorm:"size:50;not null;default:'admin'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
