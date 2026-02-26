package mikrotik

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
)

// VoucherGenerator generates voucher codes
type VoucherGenerator struct{}

// NewVoucherGenerator creates a new voucher generator
func NewVoucherGenerator() *VoucherGenerator {
	return &VoucherGenerator{}
}

// GenerateBatch generates a batch of vouchers without comment (caller is responsible for setting comment).
// The caller (VoucherUseCase.GenerateVouchers) must set the correct comment format: "mode-gencode-date-comment".
func (g *VoucherGenerator) GenerateBatch(req *dto.VoucherGenerateRequest) []*dto.Voucher {
	vouchers := make([]*dto.Voucher, 0, req.Quantity)

	for i := 0; i < req.Quantity; i++ {
		var username, password string

		switch req.Mode {
		case "vc":
			// Voucher mode: username = password
			username = g.generateVoucherCode(req.Prefix, req.NameLength, req.CharacterSet)
			password = username
		case "up":
			// Username/Password mode: separate username and password
			username = g.generateUsername(req.Prefix, req.NameLength, req.CharacterSet)
			password = g.generatePassword(req.NameLength)
		}

		vouchers = append(vouchers, &dto.Voucher{
			Username:  username,
			Password:  password,
			Profile:   req.Profile,
			Server:    req.Server,
			TimeLimit: req.TimeLimit,
			DataLimit: req.DataLimit,
			// Comment is intentionally empty here; VoucherUseCase sets the correct format:
			// "mode-gencode-date-comment"
			Comment: "",
		})
	}

	return vouchers
}

// generateVoucherCode generates a voucher code
func (g *VoucherGenerator) generateVoucherCode(prefix string, length int, charset string) string {
	chars := g.getCharset(charset)
	
	var sb strings.Builder
	sb.Grow(length + len(prefix))
	sb.WriteString(prefix)
	
	// For voucher mode, generate the code
	switch charset {
	case "num":
		// Pure numeric
		sb.WriteString(g.randomString(chars, length))
	case "lower1", "upper1", "upplow1":
		// Letters + numbers: letters first, then numbers
		letterLen := length - 2
		if letterLen < 1 {
			letterLen = 2
		}
		sb.WriteString(g.randomString(g.getLetters(charset), letterLen))
		sb.WriteString(g.randomString("0123456789", length-letterLen))
	default:
		// Mix: use full charset
		sb.WriteString(g.randomString(chars, length))
	}
	
	return sb.String()
}

// generateUsername generates a username
func (g *VoucherGenerator) generateUsername(prefix string, length int, charset string) string {
	chars := g.getCharset(charset)
	
	var sb strings.Builder
	sb.Grow(length + len(prefix))
	sb.WriteString(prefix)
	sb.WriteString(g.randomString(chars, length))
	
	return sb.String()
}

// generatePassword generates a password (numeric)
func (g *VoucherGenerator) generatePassword(length int) string {
	return g.randomString("0123456789", length)
}

// randomString generates a random string from given characters
func (g *VoucherGenerator) randomString(chars string, length int) string {
	result := make([]byte, length)
	charLen := big.NewInt(int64(len(chars)))
	
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, charLen)
		result[i] = chars[n.Int64()]
	}
	
	return string(result)
}

// getCharset returns the character set based on charset name
func (g *VoucherGenerator) getCharset(charset string) string {
	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	
	switch charset {
	case "lower":
		return lower
	case "upper":
		return upper
	case "upplow":
		return lower + upper
	case "lower1":
		return lower + digits
	case "upper1":
		return upper + digits
	case "upplow1":
		return lower + upper + digits
	case "mix":
		return lower + digits
	case "mix1":
		return upper + digits
	case "mix2":
		return lower + upper + digits
	case "num":
		return digits
	default:
		return lower + digits
	}
}

// getLetters returns only letters from charset
func (g *VoucherGenerator) getLetters(charset string) string {
	switch charset {
	case "lower1":
		return "abcdefghijklmnopqrstuvwxyz"
	case "upper1":
		return "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "upplow1":
		return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	default:
		return "abcdefghijklmnopqrstuvwxyz"
	}
}

// ParseDataLimit parses data limit string to bytes
// Supports formats: 100M, 1G, 500K
func ParseDataLimit(limit string) int64 {
	if limit == "" {
		return 0
	}
	
	limit = strings.ToUpper(limit)
	var multiplier int64 = 1
	
	// Extract unit
	if strings.HasSuffix(limit, "G") {
		multiplier = 1073741824 // 1 GB
		limit = strings.TrimSuffix(limit, "G")
	} else if strings.HasSuffix(limit, "M") {
		multiplier = 1048576 // 1 MB
		limit = strings.TrimSuffix(limit, "M")
	} else if strings.HasSuffix(limit, "K") {
		multiplier = 1024 // 1 KB
		limit = strings.TrimSuffix(limit, "K")
	}
	
	// Parse number
	var value int64
	fmt.Sscanf(limit, "%d", &value)
	
	return value * multiplier
}
