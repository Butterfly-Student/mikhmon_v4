package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

const ConfigFile = "config/mikhmon.json"

// AdminConfig holds admin credentials.
type AdminConfig struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

// RouterConfig mirrors a MikroTik router session entry.
type RouterConfig struct {
	SessionName string `json:"session_name"`
	Host        string `json:"host"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	HotspotName string `json:"hotspot_name"`
	DNSName     string `json:"dns_name"`
	Currency    string `json:"currency"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	InfoLP      string `json:"info_lp"`
	IdleTimeout string `json:"idle_timeout"`
	ReportFlag  string `json:"report_flag"`
	Token       string `json:"token"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port          int    `json:"port"`
	SessionSecret string `json:"session_secret"`
}

// AppConfig is the root configuration structure.
type AppConfig struct {
	Admin   AdminConfig    `json:"admin"`
	Routers []RouterConfig `json:"routers"`
	Server  ServerConfig   `json:"server"`
}

var (
	mu  sync.RWMutex
	cfg *AppConfig
)

// Load reads config from disk. Creates default config if not found.
func Load() (*AppConfig, error) {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			c := defaultConfig()
			cfg = c
			return c, saveUnlocked(c)
		}
		return nil, err
	}

	var c AppConfig
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	cfg = &c
	slog.Info("config loaded", "path", ConfigFile, "routers_count", len(cfg.Routers))
	return &c, nil
}

// Get returns the cached config (must call Load first).
func Get() *AppConfig {
	mu.RLock()
	defer mu.RUnlock()
	return cfg
}

// Save writes the config to disk atomically.
func Save(c *AppConfig) error {
	mu.Lock()
	defer mu.Unlock()
	cfg = c
	return saveUnlocked(c)
}

func saveUnlocked(c *AppConfig) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile, data, 0600)
}

// HashPassword creates a bcrypt hash of the given password.
func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPassword compares a plain password against a bcrypt hash.
func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// FindRouter returns the RouterConfig for a given session name (nil if not found).
func FindRouter(sessionName string) *RouterConfig {
	mu.RLock()
	defer mu.RUnlock()
	if cfg == nil {
		return nil
	}
	for i := range cfg.Routers {
		if cfg.Routers[i].SessionName == sessionName {
			return &cfg.Routers[i]
		}
	}
	return nil
}

func defaultConfig() *AppConfig {
	hash, _ := HashPassword("admin")
	secret := randomHex(32)
	return &AppConfig{
		Admin: AdminConfig{
			Username:     "admin",
			PasswordHash: hash,
		},
		Routers: []RouterConfig{},
		Server: ServerConfig{
			Port:          8080,
			SessionSecret: secret,
		},
	}
}

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
