package entity

import (
	"time"
)

// AdminUser represents an administrator user for Mikhmon login
type AdminUser struct {
	ID        string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username  string     `json:"username" gorm:"type:varchar(50);not null;uniqueIndex"`
	Password  string     `json:"-" gorm:"type:varchar(255);not null"` // Hide in JSON
	Email     string     `json:"email" gorm:"type:varchar(100)"`
	IsActive  bool       `json:"isActive" gorm:"default:true"`
	LastLogin *time.Time `json:"lastLogin,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// AdminUserResponse represents admin user in API responses
type AdminUserResponse struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email,omitempty"`
	IsActive  bool       `json:"isActive"`
	LastLogin *time.Time `json:"lastLogin,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

// ToResponse converts AdminUser to AdminUserResponse
func (u *AdminUser) ToResponse() *AdminUserResponse {
	return &AdminUserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		IsActive:  u.IsActive,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
	}
}
