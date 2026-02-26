package dto

// VoucherGenerateRequest represents a request to generate vouchers
type VoucherGenerateRequest struct {
	RouterID     uint   `json:"routerId" validate:"required"`
	Profile      string `json:"profile" validate:"required,max=50"`
	Quantity     int    `json:"quantity" validate:"required,min=1,max=500"`
	Server       string `json:"server,omitempty" validate:"max=50"`
	Mode         string `json:"mode" validate:"required,oneof=vc up"` // vc = voucher, up = username/password
	Gencode      string `json:"gencode,omitempty" validate:"max=10"`
	NameLength   int    `json:"nameLength" validate:"required,min=3,max=12"`
	Prefix       string `json:"prefix,omitempty" validate:"max=20"`
	CharacterSet string `json:"characterSet" validate:"required,oneof=lower upper upplow lower1 upper1 upplow1 mix mix1 mix2 num"`
	TimeLimit    string `json:"timeLimit,omitempty" validate:"max=20"`
	DataLimit    string `json:"dataLimit,omitempty" validate:"max=20"`
	Comment      string `json:"comment,omitempty" validate:"max=100"`
}

// Voucher represents a generated voucher
type Voucher struct {
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	Profile   string `json:"profile"`
	Server    string `json:"server,omitempty"`
	TimeLimit string `json:"timeLimit,omitempty"`
	DataLimit string `json:"dataLimit,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

// PrintVoucherRequest represents a request to print vouchers
type PrintVoucherRequest struct {
	RouterID   string   `json:"routerId" validate:"required"`
	TemplateID string   `json:"templateId,omitempty"`
	IDs        []string `json:"ids,omitempty"`     // Print by user IDs
	Comment    string   `json:"comment,omitempty"` // Print by comment (batch)
	Profile    string   `json:"profile,omitempty"` // Print by profile
}

// PrintVoucherResponse represents a print voucher response
type PrintVoucherResponse struct {
	HTML     string            `json:"html"`
	CSS      string            `json:"css,omitempty"`
	Settings map[string]string `json:"settings,omitempty"`
}

// VoucherFilter represents filter for getting vouchers
type VoucherFilter struct {
	Profile string `json:"profile,omitempty"`
	Comment string `json:"comment,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

// VoucherBatchResult represents the result of voucher generation
type VoucherBatchResult struct {
	Count    int       `json:"count"`
	Comment  string    `json:"comment"`
	Vouchers []Voucher `json:"vouchers"`
}
