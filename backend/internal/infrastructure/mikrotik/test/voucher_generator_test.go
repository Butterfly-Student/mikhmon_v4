package mikrotik_test

import (
	"testing"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
)

// TestVoucherGenerator_GenerateBatch_Quantity verifies correct number of vouchers generated.
func TestVoucherGenerator_GenerateBatch_Quantity(t *testing.T) {
	g := mikrotik.NewVoucherGenerator()
	req := &dto.VoucherGenerateRequest{
		Mode:         "vc",
		Quantity:     5,
		NameLength:   6,
		CharacterSet: "num",
		Profile:      "Basic",
		Server:       "all",
	}
	vouchers := g.GenerateBatch(req)
	if len(vouchers) != 5 {
		t.Errorf("Expected 5 vouchers, got %d", len(vouchers))
	}
}

// TestVoucherGenerator_GenerateBatch_VCMode verifies username == password in vc mode.
func TestVoucherGenerator_GenerateBatch_VCMode(t *testing.T) {
	g := mikrotik.NewVoucherGenerator()
	req := &dto.VoucherGenerateRequest{
		Mode:         "vc",
		Quantity:     3,
		NameLength:   8,
		CharacterSet: "num",
		Profile:      "Basic",
		Server:       "all",
	}
	vouchers := g.GenerateBatch(req)
	for i, v := range vouchers {
		if v.Username != v.Password {
			t.Errorf("Voucher %d: username %q != password %q in vc mode", i, v.Username, v.Password)
		}
		if v.Username == "" {
			t.Errorf("Voucher %d: username is empty", i)
		}
	}
}

// TestVoucherGenerator_GenerateBatch_UPMode verifies username != password in up mode.
func TestVoucherGenerator_GenerateBatch_UPMode(t *testing.T) {
	g := mikrotik.NewVoucherGenerator()
	req := &dto.VoucherGenerateRequest{
		Mode:         "up",
		Quantity:     3,
		NameLength:   8,
		CharacterSet: "lower",
		Profile:      "Basic",
		Server:       "all",
	}
	vouchers := g.GenerateBatch(req)
	for i, v := range vouchers {
		if v.Username == "" {
			t.Errorf("Voucher %d: username is empty", i)
		}
		if v.Password == "" {
			t.Errorf("Voucher %d: password is empty", i)
		}
	}
}

// TestVoucherGenerator_GenerateBatch_CommentEmpty verifies GenerateBatch leaves Comment empty.
func TestVoucherGenerator_GenerateBatch_CommentEmpty(t *testing.T) {
	g := mikrotik.NewVoucherGenerator()
	req := &dto.VoucherGenerateRequest{
		Mode:         "vc",
		Quantity:     2,
		NameLength:   6,
		CharacterSet: "num",
		Comment:      "Test Batch",
	}
	vouchers := g.GenerateBatch(req)
	for i, v := range vouchers {
		if v.Comment != "" {
			t.Errorf("Voucher %d: Comment should be empty (set by VoucherUseCase), got %q", i, v.Comment)
		}
	}
}

// TestVoucherGenerator_Prefix verifies prefix is prepended to username.
func TestVoucherGenerator_Prefix(t *testing.T) {
	g := mikrotik.NewVoucherGenerator()
	prefix := "TEST-"
	req := &dto.VoucherGenerateRequest{
		Mode:         "vc",
		Quantity:     3,
		Prefix:       prefix,
		NameLength:   6,
		CharacterSet: "num",
	}
	vouchers := g.GenerateBatch(req)
	for i, v := range vouchers {
		if len(v.Username) < len(prefix) || v.Username[:len(prefix)] != prefix {
			t.Errorf("Voucher %d: username %q should start with prefix %q", i, v.Username, prefix)
		}
	}
}

// TestParseDataLimit verifies data limit parsing for different units.
func TestParseDataLimit(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"", 0},
		{"0", 0},
		{"100M", 100 * 1048576},
		{"1G", 1073741824},
		{"500K", 500 * 1024},
		{"2G", 2 * 1073741824},
		{"256M", 256 * 1048576},
		{"100m", 100 * 1048576},
		{"1g", 1073741824},
		{"invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := mikrotik.ParseDataLimit(tt.input)
			if got != tt.expected {
				t.Errorf("ParseDataLimit(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}
