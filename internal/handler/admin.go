package handler

import (
	"encoding/json"
	"fmt"
	htmltemplate "html/template"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"mikhmon_v4/config"
	ros "mikhmon_v4/internal/routeros"
)

// AdminHandler handles router config CRUD and settings.
type AdminHandler struct {
	store sessions.Store
}

func NewAdminHandler(store sessions.Store) *AdminHandler {
	return &AdminHandler{store: store}
}

// SettingsPage renders the admin settings page.
func (h *AdminHandler) SettingsPage(c *gin.Context) {
	cfg := config.Get()
	routers := cfg.Routers
	if routers == nil {
		routers = []config.RouterConfig{}
	}
	slog.Info("SettingsPage rendered", "routers_count", len(routers))
	jsonBytes, _ := json.Marshal(routers)
	c.HTML(http.StatusOK, "web/templates/settings.html", gin.H{
		"Title":       "Mikhmon — Settings",
		"Routers":     routers,
		"RoutersJSON": htmltemplate.JS(string(jsonBytes)),
		"AdminUser":   cfg.Admin.Username,
		"Active":      "settings",
	})
}

// AddRouter creates a new empty router config entry and returns the session name.
func (h *AdminHandler) AddRouter(c *gin.Context) {
	cfg := config.Get()

	base := fmt.Sprintf("session%d", rand.Intn(900)+100)
	// Ensure uniqueness.
	for config.FindRouter(base) != nil {
		base = fmt.Sprintf("session%d", rand.Intn(900)+100)
	}

	cfg.Routers = append(cfg.Routers, config.RouterConfig{
		SessionName: base,
		IdleTimeout: "30",
		ReportFlag:  "disable",
		Currency:    "Rp",
	})
	if err := config.Save(cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving config: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success", "sesname": base})
}

// SaveRouter updates an existing router config.
func (h *AdminHandler) SaveRouter(c *gin.Context) {
	oldName := c.Param("name")
	cfg := config.Get()

	idx := -1
	for i, r := range cfg.Routers {
		if r.SessionName == oldName {
			idx = i
			break
		}
	}
	if idx < 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Router not found"})
		return
	}

	var body struct {
		Session     string `json:"session"`
		IpMik       string `json:"ipmik"`
		UserMik     string `json:"usermik"`
		PassMik     string `json:"passmik"`
		HotspotName string `json:"hotspotname"`
		DNSName     string `json:"dnsname"`
		Currency    string `json:"currency"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		InfoLP      string `json:"infolp"`
		IdleTo      string `json:"idleto"`
		Report      string `json:"report"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	// Sanitize new session name.
	newName := sanitizeSessionName(body.Session)
	if newName == "" {
		newName = oldName
	}

	// Check if newName already exists (different router).
	if newName != oldName && config.FindRouter(newName) != nil {
		suffix := rand.Intn(900) + 100
		newName = fmt.Sprintf("%s%d", newName, suffix)
	}

	cfg.Routers[idx] = config.RouterConfig{
		SessionName: newName,
		Host:        strings.TrimSpace(body.IpMik),
		Username:    body.UserMik,
		Password:    body.PassMik,
		HotspotName: body.HotspotName,
		DNSName:     body.DNSName,
		Currency:    body.Currency,
		Phone:       body.Phone,
		Email:       body.Email,
		InfoLP:      body.InfoLP,
		IdleTimeout: body.IdleTo,
		ReportFlag:  body.Report,
	}

	if err := config.Save(cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving config: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success", "sess": newName})
}

// RemoveRouter deletes a router config entry.
func (h *AdminHandler) RemoveRouter(c *gin.Context) {
	name := c.Param("name")
	cfg := config.Get()

	filtered := cfg.Routers[:0]
	for _, r := range cfg.Routers {
		if r.SessionName != name {
			filtered = append(filtered, r)
		}
	}
	cfg.Routers = filtered

	if err := config.Save(cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving config: " + err.Error()})
		return
	}

	// Remove logo if exists.
	os.Remove("web/static/img/logo-" + name + ".png")

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// SaveTemplate writes a voucher template file to disk.
func (h *AdminHandler) SaveTemplate(c *gin.Context) {
	var body struct {
		Section  string `json:"section"`  // header | row | footer
		Template string `json:"template"` // default | small | thermal
		Content  string `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	filename := fmt.Sprintf("voucher_templates/%s.%s.txt", body.Section, body.Template)
	if err := os.WriteFile(filename, []byte(body.Content), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error writing template: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// EditorPage renders the voucher template editor.
// Template data is pre-encoded to JSON to avoid using {{range}} inside a
// <script> block, which html/template's context-aware parser does not allow.
func (h *AdminHandler) EditorPage(c *gin.Context) {
	tmplMap := map[string]map[string]string{}
	sections := []string{"header", "row", "footer"}
	sizes := []string{"default", "small", "thermal"}
	for _, sz := range sizes {
		tmplMap[sz] = map[string]string{}
		for _, sec := range sections {
			fname := fmt.Sprintf("voucher_templates/%s.%s.txt", sec, sz)
			data, _ := os.ReadFile(fname)
			tmplMap[sz][sec] = string(data)
		}
	}
	jsonBytes, _ := json.Marshal(tmplMap)
	c.HTML(http.StatusOK, "web/templates/editor.html", gin.H{
		"Title":         "Template Editor",
		"TemplatesJSON": htmltemplate.JS(jsonBytes),
		"Active":        "editor",
	})
}

// GetTemplate returns a specific voucher template content as JSON.
func (h *AdminHandler) GetTemplate(c *gin.Context) {
	size := c.DefaultQuery("size", "default")
	section := c.DefaultQuery("section", "header")
	fname := fmt.Sprintf("voucher_templates/%s.%s.txt", section, size)
	data, err := os.ReadFile(fname)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"content": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": string(data)})
}

// TestConnection tests a RouterOS connection without saving it.
func (h *AdminHandler) TestConnection(c *gin.Context) {
	var body struct {
		IpMik   string `json:"ipmik"`
		UserMik string `json:"usermik"`
		PassMik string `json:"passmik"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	client, err := ros.Connect(body.IpMik, body.UserMik, body.PassMik)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "Connection failed: " + err.Error()})
		return
	}
	client.Close()
	c.JSON(http.StatusOK, gin.H{"message": "Connection successful"})
}

// SaveAdmin updates the admin username and password.
func (h *AdminHandler) SaveAdmin(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	cfg := config.Get()
	cfg.Admin.Username = strings.ReplaceAll(body.Username, "'", "")
	if body.Password != "" {
		hash, err := config.HashPassword(body.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
			return
		}
		cfg.Admin.PasswordHash = hash
	}
	if err := config.Save(cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving config"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

var nonAlnum = regexp.MustCompile(`[^a-zA-Z0-9]`)

func sanitizeSessionName(s string) string {
	return nonAlnum.ReplaceAllString(strings.ReplaceAll(s, " ", ""), "")
}
