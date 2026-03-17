package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

const (
	ConfigFile = "config/mikhmon.json" // server setting only
	DBFile     = "config/mikhmon.db"   // users + routers
)

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

type serverFile struct {
	Server ServerConfig `json:"server"`
}

var (
	mu  sync.RWMutex
	cfg *AppConfig
)

func Load() (*AppConfig, error) {
	mu.Lock()
	defer mu.Unlock()

	srv, err := loadServerConfig()
	if err != nil {
		return nil, err
	}

	if err := initDB(); err != nil {
		return nil, err
	}
	if err := ensureDefaultAdmin(); err != nil {
		return nil, err
	}

	admin, routers, err := readAdminRouters()
	if err != nil {
		return nil, err
	}

	cfg = &AppConfig{Admin: admin, Routers: routers, Server: srv}
	slog.Info("config loaded", "db", DBFile, "routers_count", len(cfg.Routers))
	return cfg, nil
}

func Get() *AppConfig {
	mu.Lock()
	defer mu.Unlock()
	if cfg == nil {
		return nil
	}
	admin, routers, err := readAdminRouters()
	if err == nil {
		cfg.Admin = admin
		cfg.Routers = routers
	}
	copyCfg := *cfg
	return &copyCfg
}

func Save(c *AppConfig) error {
	mu.Lock()
	defer mu.Unlock()

	if err := saveAdminRouters(c.Admin, c.Routers); err != nil {
		return err
	}
	if err := saveServerConfig(c.Server); err != nil {
		return err
	}
	cfg = c
	return nil
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func FindRouter(sessionName string) *RouterConfig {
	mu.Lock()
	defer mu.Unlock()
	if cfg == nil {
		return nil
	}
	_, routers, err := readAdminRouters()
	if err != nil {
		return nil
	}
	cfg.Routers = routers
	for i := range cfg.Routers {
		if cfg.Routers[i].SessionName == sessionName {
			r := cfg.Routers[i]
			return &r
		}
	}
	return nil
}

func loadServerConfig() (ServerConfig, error) {
	b, err := os.ReadFile(ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			s := ServerConfig{Port: 8080, SessionSecret: randomHex(32)}
			return s, saveServerConfig(s)
		}
		return ServerConfig{}, err
	}

	var sf serverFile
	if err := json.Unmarshal(b, &sf); err != nil {
		return ServerConfig{}, err
	}
	if sf.Server.Port == 0 {
		sf.Server.Port = 8080
	}
	if sf.Server.SessionSecret == "" {
		sf.Server.SessionSecret = randomHex(32)
	}
	return sf.Server, nil
}

func saveServerConfig(s ServerConfig) error {
	data, err := json.MarshalIndent(serverFile{Server: s}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile, data, 0600)
}

func initDB() error {
	schema := `
CREATE TABLE IF NOT EXISTS users (
	username TEXT PRIMARY KEY,
	password_hash TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS routers (
	session_name TEXT PRIMARY KEY,
	host TEXT,
	username TEXT,
	password TEXT,
	hotspot_name TEXT,
	dns_name TEXT,
	currency TEXT,
	phone TEXT,
	email TEXT,
	info_lp TEXT,
	idle_timeout TEXT,
	report_flag TEXT,
	token TEXT
);`
	_, err := runSQLite(schema)
	return err
}

func ensureDefaultAdmin() error {
	out, err := runSQLite("SELECT COUNT(*) FROM users;")
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != "0" {
		return nil
	}
	h, err := HashPassword("admin")
	if err != nil {
		return err
	}
	_, err = runSQLite("INSERT INTO users(username,password_hash) VALUES ('admin', " + q(h) + ");")
	return err
}

func readAdminRouters() (AdminConfig, []RouterConfig, error) {
	adminRows, err := runSQLite("SELECT username,password_hash FROM users ORDER BY username LIMIT 1;")
	if err != nil {
		return AdminConfig{}, nil, err
	}
	admin := AdminConfig{Username: "admin"}
	if strings.TrimSpace(adminRows) != "" {
		parts := strings.Split(strings.TrimSpace(adminRows), "\x1f")
		if len(parts) >= 2 {
			admin.Username = parts[0]
			admin.PasswordHash = parts[1]
		}
	}

	rRows, err := runSQLite("SELECT session_name,host,username,password,hotspot_name,dns_name,currency,phone,email,info_lp,idle_timeout,report_flag,token FROM routers ORDER BY session_name;")
	if err != nil {
		return AdminConfig{}, nil, err
	}
	routers := make([]RouterConfig, 0)
	for _, line := range strings.Split(strings.TrimSpace(rRows), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		p := strings.Split(line, "\x1f")
		for len(p) < 13 {
			p = append(p, "")
		}
		routers = append(routers, RouterConfig{
			SessionName: p[0], Host: p[1], Username: p[2], Password: p[3], HotspotName: p[4], DNSName: p[5], Currency: p[6], Phone: p[7], Email: p[8], InfoLP: p[9], IdleTimeout: p[10], ReportFlag: p[11], Token: p[12],
		})
	}
	return admin, routers, nil
}

func saveAdminRouters(admin AdminConfig, routers []RouterConfig) error {
	var sb strings.Builder
	sb.WriteString("BEGIN;\n")
	sb.WriteString("DELETE FROM users;\n")
	sb.WriteString("INSERT INTO users(username,password_hash) VALUES (")
	sb.WriteString(q(admin.Username))
	sb.WriteString(",")
	sb.WriteString(q(admin.PasswordHash))
	sb.WriteString(");\n")
	sb.WriteString("DELETE FROM routers;\n")
	for _, r := range routers {
		sb.WriteString("INSERT INTO routers(session_name,host,username,password,hotspot_name,dns_name,currency,phone,email,info_lp,idle_timeout,report_flag,token) VALUES (")
		sb.WriteString(strings.Join([]string{q(r.SessionName), q(r.Host), q(r.Username), q(r.Password), q(r.HotspotName), q(r.DNSName), q(r.Currency), q(r.Phone), q(r.Email), q(r.InfoLP), q(r.IdleTimeout), q(r.ReportFlag), q(r.Token)}, ","))
		sb.WriteString(");\n")
	}
	sb.WriteString("COMMIT;")
	_, err := runSQLite(sb.String())
	return err
}

func runSQLite(sql string) (string, error) {
	cmd := exec.Command("sqlite3", "-separator", "\x1f", DBFile, sql)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("sqlite error: %v (%s)", err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func q(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
