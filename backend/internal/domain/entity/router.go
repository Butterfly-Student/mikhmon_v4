package entity

import (
	"time"
)

// Router represents a MikroTik router configuration
type Router struct {
	ID            uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string     `json:"name" gorm:"type:varchar(50);not null;uniqueIndex"`
	Host          string     `json:"host" gorm:"type:varchar(255);not null"`
	Port          int        `json:"port" gorm:"default:8728"`
	Username      string     `json:"-" gorm:"type:varchar(50);not null"` // Hide password
	Password      string     `json:"-" gorm:"type:varchar(255);not null"` // Encrypted MikroTik password
	UseSSL        bool       `json:"useSsl" gorm:"default:false"`
	Timeout       int        `json:"timeout" gorm:"default:3"`
	IsActive      bool       `json:"isActive" gorm:"default:true"`
	Description   string     `json:"description" gorm:"type:text"`
	LastConnected *time.Time `json:"lastConnected,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// RouterResponse represents router in API responses (without sensitive data)
type RouterResponse struct {
	ID            uint       `json:"id"`
	Name          string     `json:"name"`
	Host          string     `json:"host"`
	Port          int        `json:"port"`
	UseSSL        bool       `json:"useSsl"`
	Timeout       int        `json:"timeout"`
	IsActive      bool       `json:"isActive"`
	Description   string     `json:"description"`
	LastConnected *time.Time `json:"lastConnected,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// ToResponse converts Router to RouterResponse
func (r *Router) ToResponse() *RouterResponse {
	return &RouterResponse{
		ID:            r.ID,
		Name:          r.Name,
		Host:          r.Host,
		Port:          r.Port,
		UseSSL:        r.UseSSL,
		Timeout:       r.Timeout,
		IsActive:      r.IsActive,
		Description:   r.Description,
		LastConnected: r.LastConnected,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

// RouterConnectionRequest represents a connection test request
type RouterConnectionRequest struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"min=1,max=65535"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	UseSSL   bool   `json:"useSsl"`
	Timeout  int    `json:"timeout" validate:"min=1,max=60"`
}

// RouterCreateRequest represents a create request
type RouterCreateRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Host        string `json:"host" validate:"required"`
	Port        int    `json:"port" validate:"min=1,max=65535"`
	Username    string `json:"username" validate:"required,max=50"`
	Password    string `json:"password" validate:"required"`
	UseSSL      bool   `json:"useSsl"`
	Timeout     int    `json:"timeout" validate:"min=1,max=60"`
	Description string `json:"description" validate:"max=500"`
}

// RouterUpdateRequest represents an update request
type RouterUpdateRequest struct {
	Name        string `json:"name,omitempty" validate:"max=50"`
	Host        string `json:"host,omitempty"`
	Port        int    `json:"port,omitempty" validate:"min=1,max=65535"`
	Username    string `json:"username,omitempty" validate:"max=50"`
	Password    string `json:"password,omitempty"`
	UseSSL      bool   `json:"useSsl"`
	Timeout     int    `json:"timeout,omitempty" validate:"min=1,max=60"`
	IsActive    bool   `json:"isActive"`
	Description string `json:"description,omitempty" validate:"max=500"`
}
