package dto

// SalesReport represents a sales report entry from MikroTik /system/script
// Format name: $date-|-$time-|-$user-|-$price-|-$address-|-$mac-|-$validity-|-$profile-|-$comment
type SalesReport struct {
	ID         string  `json:".id,omitempty"`
	Name       string  `json:"name,omitempty"`    // Raw script name
	Owner      string  `json:"owner,omitempty"`   // month+year (jan2024)
	Source     string  `json:"source,omitempty"`  // day (jan/25/2024)
	Comment    string  `json:"comment,omitempty"` // mikhmon
	DontReq    string  `json:"dont-require-permissions,omitempty"`
	RunCount   string  `json:"run-count,omitempty"`
	CopyOf     string  `json:"copy-of,omitempty"`
	
	// Parsed fields from name
	Date       string  `json:"date"`
	Time       string  `json:"time"`
	Username   string  `json:"username"`
	Price      float64 `json:"price"`
	IPAddress  string  `json:"ipAddress"`
	MACAddress string  `json:"macAddress"`
	Validity   string  `json:"validity"`
	Profile    string  `json:"profile"`
	VoucherComment string `json:"voucherComment,omitempty"`
}

// ReportFilter represents filter parameters for reports
type ReportFilter struct {
	RouterID string `json:"routerId" validate:"required"`
	Day      string `json:"day,omitempty"`      // Format: jan/25/2024
	Month    string `json:"month,omitempty"`    // Format: jan2024 (used as owner)
	Year     string `json:"year,omitempty"`     // Format: 2024
	Profile  string `json:"profile,omitempty"`  // Filter by profile
}

// ReportSummary represents a summary of sales
type ReportSummary struct {
	TotalVouchers int                 `json:"totalVouchers"`
	TotalAmount   float64             `json:"totalAmount"`
	ByProfile     map[string]ProfileSummary `json:"byProfile,omitempty"`
}

// ProfileSummary represents summary by profile
type ProfileSummary struct {
	Count int     `json:"count"`
	Total float64 `json:"total"`
}

// ReportResponse represents report API response
type ReportResponse struct {
	Data    []*SalesReport `json:"data"`
	Summary *ReportSummary `json:"summary,omitempty"`
	Filter  *ReportFilter  `json:"filter,omitempty"`
}

// LiveReportRequest represents a live report request
type LiveReportRequest struct {
	RouterID string `json:"routerId" validate:"required"`
	Month    string `json:"month,omitempty"` // jan2024
	Day      string `json:"day,omitempty"`   // jan/25/2024
}
