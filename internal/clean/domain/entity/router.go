//go:build modern

package entity

import "time"

type Router struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;size:100;not null"`
	Host      string `gorm:"size:200;not null"`
	Username  string `gorm:"size:100;not null"`
	Password  string `gorm:"size:200;not null"`
	Port      int    `gorm:"not null;default:8728"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
