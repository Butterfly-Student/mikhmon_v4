package entity

import (
	"time"
)

// Setting represents an application setting
type Setting struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Key       string    `json:"key" gorm:"type:varchar(50);not null;uniqueIndex"`
	Value     string    `json:"value" gorm:"type:text"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SettingKey represents predefined setting keys
const (
	SettingCompanyName    = "company_name"
	SettingCompanyAddress = "company_address"
	SettingCompanyPhone   = "company_phone"
	SettingTheme          = "theme" // light, dark, blue, green, pink
	SettingSessionTimeout = "session_timeout"
	SettingCurrency       = "currency"
	SettingDateFormat     = "date_format"
	SettingTimeFormat     = "time_format"
	SettingHotspotAddress = "hotspot_address"
	SettingHotspotPhone   = "hotspot_phone"
)

// PrintTemplate represents a voucher print template
type PrintTemplate struct {
	ID          string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"type:varchar(50);not null"`
	IsDefault   bool      `json:"isDefault" gorm:"default:false"`
	Description string    `json:"description" gorm:"type:varchar(255)"`
	Content     string    `json:"content" gorm:"type:text;not null"` // HTML template
	CSS         string    `json:"css" gorm:"type:text"`
	Width       int       `json:"width" gorm:"default:80"`  // mm
	Height      int       `json:"height" gorm:"default:0"`  // 0 = auto
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateTemplateRequest represents a create request
type CreateTemplateRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description,omitempty" validate:"max=255"`
	Content     string `json:"content" validate:"required"`
	CSS         string `json:"css,omitempty"`
	Width       int    `json:"width,omitempty" validate:"min=50,max=200"`
	Height      int    `json:"height,omitempty" validate:"min=0,max=300"`
}

// UpdateTemplateRequest represents an update request
type UpdateTemplateRequest struct {
	Name        string `json:"name,omitempty" validate:"max=50"`
	Description string `json:"description,omitempty" validate:"max=255"`
	Content     string `json:"content,omitempty"`
	CSS         string `json:"css,omitempty"`
	Width       int    `json:"width,omitempty" validate:"min=50,max=200"`
	Height      int    `json:"height,omitempty" validate:"min=0,max=300"`
	IsDefault   bool   `json:"isDefault,omitempty"`
}

// DefaultTemplateContent returns the default voucher template
func DefaultTemplateContent() string {
	return `<div class="voucher">
  <div class="header">
    <h2>{{company_name}}</h2>
    <p>{{hotspot_address}}</p>
  </div>
  <div class="content">
    <table>
      <tr><td>Username</td><td>: <b>{{username}}</b></td></tr>
      <tr><td>Password</td><td>: <b>{{password}}</b></td></tr>
      <tr><td>Profile</td><td>: {{profile}}</td></tr>
      {{#if time_limit}}<tr><td>Time Limit</td><td>: {{time_limit}}</td></tr>{{/if}}
      {{#if data_limit}}<tr><td>Data Limit</td><td>: {{data_limit}}</td></tr>{{/if}}
      {{#if price}}<tr><td>Price</td><td>: {{currency}}{{price}}</td></tr>{{/if}}
    </table>
  </div>
  <div class="footer">
    <p>{{footer_text}}</p>
  </div>
</div>`
}

// DefaultTemplateCSS returns the default CSS
func DefaultTemplateCSS() string {
	return `.voucher {
  width: 80mm;
  padding: 5mm;
  font-family: Arial, sans-serif;
  font-size: 12px;
  border: 1px dashed #000;
  margin-bottom: 5mm;
}
.voucher .header {
  text-align: center;
  border-bottom: 1px solid #000;
  padding-bottom: 3mm;
}
.voucher .header h2 {
  margin: 0;
  font-size: 16px;
}
.voucher .content {
  margin: 3mm 0;
}
.voucher table {
  width: 100%;
}
.voucher td {
  padding: 1mm 0;
}
.voucher .footer {
  text-align: center;
  border-top: 1px solid #000;
  padding-top: 3mm;
  font-size: 10px;
}`
}
