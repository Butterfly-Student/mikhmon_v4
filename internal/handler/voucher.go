package handler

import (
	"embed"
	"encoding/base64"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	qrcode "github.com/skip2/go-qrcode"

	"mikhmon_v4/config"
	ros "mikhmon_v4/internal/routeros"
	"mikhmon_v4/internal/util"
)

// VoucherHandler handles voucher generation, caching, and printing.
type VoucherHandler struct {
	store     sessions.Store
	pool      *ros.Pool
	templates embed.FS
}

func NewVoucherHandler(store sessions.Store, pool *ros.Pool, templates embed.FS) *VoucherHandler {
	return &VoucherHandler{store: store, pool: pool, templates: templates}
}

// GeneratePage renders the voucher generation form.
func (h *VoucherHandler) GeneratePage(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.Redirect(http.StatusFound, "/admin/settings")
		return
	}
	c.HTML(http.StatusOK, "web/templates/generate.html", gin.H{
		"Title":   "Generate Vouchers",
		"Session": session,
		"Router":  router,
	})
}

// GenerateRequest is the JSON body for voucher generation.
type GenerateRequest struct {
	Qty       int    `json:"qty"`
	Server    string `json:"server"`
	User      string `json:"user"` // "up" or "vc"
	UserLen   int    `json:"userl"`
	Prefix    string `json:"prefix"`
	Char      string `json:"char"`
	Profile   string `json:"profile"`
	TimeLimit string `json:"timelimit"`
	DataLimit string `json:"datalimit"` // e.g. "100m" or "1g"
	GComment  string `json:"gcomment"`
	GenCode   string `json:"gencode"`
	SessName  string `json:"sessname"`
}

// GenerateVouchers creates N hotspot users in bulk on the router.
func (h *VoucherHandler) GenerateVouchers(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "data": gin.H{"error": err.Error()}})
		return
	}

	var body GenerateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "data": gin.H{"error": err.Error()}})
		return
	}

	// Calculate data limit in bytes.
	dataBytes := parseDataLimit(body.DataLimit)

	// Build comment for this batch.
	commt := fmt.Sprintf("%s-%s-%s", body.User, body.GenCode, body.GComment)

	// Generate username/password pairs.
	type creds struct{ u, p string }
	var pairs []creds

	qty := body.Qty
	userl := body.UserLen
	if qty <= 0 {
		qty = 1
	}

	for i := 0; i < qty; i++ {
		var u, p string
		if body.User == "up" {
			u = genUsername(body.Char, userl)
			p = genPinByLen(userl)
			u = body.Prefix + u
		} else { // vc
			base, pin := genVoucherCode(body.Char, userl)
			u = body.Prefix + base + pin
			p = u
			// Override for numeric/mix chars.
			switch body.Char {
			case "num":
				pin2 := genPinByLen(userl)
				u = body.Prefix + pin2
				p = u
			case "mix":
				u = body.Prefix + util.RandNLC(userl)
				p = u
			case "mix1":
				u = body.Prefix + util.RandNUC(userl)
				p = u
			case "mix2":
				u = body.Prefix + util.RandNULC(userl)
				p = u
			}
		}
		pairs = append(pairs, creds{u: u, p: p})
	}

	// Add users to routers.
	for _, pair := range pairs {
		args := []string{
			"/ip/hotspot/user/add",
			"=server=" + body.Server,
			"=name=" + pair.u,
			"=password=" + pair.p,
			"=profile=" + body.Profile,
			"=limit-uptime=" + body.TimeLimit,
			"=comment=" + commt,
		}
		if dataBytes > 0 {
			args = append(args, fmt.Sprintf("=limit-bytes-total=%d", dataBytes))
		}
		ros.RunArgs(client, args...)
	}

	// Retrieve created users by comment.
	users, err := ros.RunArgs(client, "/ip/hotspot/user/print", "?comment="+commt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "data": gin.H{"error": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"count":   len(users),
			"comment": commt,
			"profile": body.Profile,
		},
	})
}

// CacheVouchers fetches users by comment from the router.
func (h *VoucherHandler) CacheVouchers(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "data": gin.H{"error": err.Error()}})
		return
	}

	var body struct {
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "data": gin.H{"error": err.Error()}})
		return
	}

	users, err := ros.RunArgs(client, "/ip/hotspot/user/print", "?comment="+body.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "data": gin.H{"error": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": users})
}

// Voucher holds per-voucher data for the print template.
type Voucher struct {
	Username    string
	Password    string
	Profile     string
	Validity    string
	LimitUptime string
	LimitBytes  string
	Price       string
	DNSName     string
	HotspotName string
	QRCode      string // base64 PNG data URI
	Num         int
}

// voucherData is the template context for print-voucher.
type voucherData struct {
	Header   string
	Row      string
	Footer   string
	Vouchers []Voucher
	Size     string
}

// PrintVouchers renders the voucher print page with QR codes.
// Query params: comment=..., profile=..., size=default|small|thermal
func (h *VoucherHandler) PrintVouchers(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.Redirect(http.StatusFound, "/admin/settings")
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "web/templates/print_voucher.html", gin.H{"Error": err.Error()})
		return
	}

	comment := c.Query("comment")
	size := c.DefaultQuery("size", "default")

	// Load templates from disk first, then fall back to embedded.
	header := readTemplate(h.templates, "header", size)
	row := readTemplate(h.templates, "row", size)
	footer := readTemplate(h.templates, "footer", size)

	// Fetch users.
	users, _ := ros.RunArgs(client, "/ip/hotspot/user/print", "?comment="+comment)

	var vouchers []Voucher
	for i, u := range users {
		loginURL := fmt.Sprintf("http://%s", router.DNSName)
		qr, _ := generateQR(loginURL, 60)

		// Parse on-login script to get price/validity from profile.
		vouchers = append(vouchers, Voucher{
			Username:    u["name"],
			Password:    u["password"],
			Profile:     u["profile"],
			LimitUptime: u["limit-uptime"],
			LimitBytes:  u["limit-bytes-total"],
			DNSName:     router.DNSName,
			HotspotName: router.HotspotName,
			QRCode:      qr,
			Num:         i + 1,
		})
	}

	c.HTML(http.StatusOK, "web/templates/print_voucher.html", gin.H{
		"Title":    "Print Vouchers",
		"Session":  session,
		"Router":   router,
		"Header":   header,
		"Row":      row,
		"Footer":   footer,
		"Vouchers": vouchers,
		"Size":     size,
	})
}

// ── Helpers ───────────────────────────────────────────────────────────────

func parseDataLimit(s string) int64 {
	if s == "" {
		return 0
	}
	s = strings.ToLower(strings.TrimSpace(s))
	var multiplier int64 = 1048576 // MB default
	var numStr string

	if strings.HasSuffix(s, "g") {
		multiplier = 1073741824
		numStr = s[:len(s)-1]
	} else if strings.HasSuffix(s, "m") {
		numStr = s[:len(s)-1]
	} else {
		numStr = s
		multiplier = 1
	}

	var n int64
	fmt.Sscan(numStr, &n)
	return n * multiplier
}

func genUsername(char string, length int) string {
	switch char {
	case "lower":
		return util.RandLC(length)
	case "upper":
		return util.RandUC(length)
	case "upplow":
		return util.RandULC(length)
	case "mix":
		return util.RandNLC(length)
	case "mix1":
		return util.RandNUC(length)
	case "mix2":
		return util.RandNULC(length)
	default:
		return util.RandNULC(length)
	}
}

func genPinByLen(length int) string {
	if length >= 4 && length <= 8 {
		return util.RandN(length)
	}
	return util.RandN(4)
}

func genVoucherCode(char string, length int) (base, pin string) {
	var baseLen, pinLen int
	switch {
	case length >= 4 && length <= 5:
		pinLen = 2
	case length >= 6 && length <= 7:
		pinLen = 3
	case length == 8:
		pinLen = 4
	default:
		pinLen = 2
	}
	baseLen = length - pinLen

	switch char {
	case "lower1":
		base = util.RandLC(baseLen)
	case "upper1":
		base = util.RandUC(baseLen)
	case "upplow1":
		base = util.RandULC(baseLen)
	default:
		base = util.RandNULC(baseLen)
	}
	pin = util.RandN(pinLen)
	return
}

func generateQR(content string, size int) (string, error) {
	png, err := qrcode.Encode(content, qrcode.Medium, size)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png), nil
}

func readTemplate(embedded embed.FS, section, size string) string {
	// Try disk first.
	fname := fmt.Sprintf("voucher_templates/%s.%s.txt", section, size)
	if data, err := fs.ReadFile(embedded, fname); err == nil {
		return string(data)
	}
	return ""
}
