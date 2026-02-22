package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"mikhmon_v4/config"
)

// migrate reads legacy PHP config/config.php and emits config/mikhmon.json.
//
// The PHP config.php has two formats:
//
//  1. Admin record:
//     $data['mikhmon'] = array('1' => 'mikhmon<|<USERNAME', 'mikhmon>|>BASE64PASS');
//
//  2. Router record (one or more):
//     $data['SESS'] = array(
//     'host' => 'IP:PORT', 'user' => 'USER', 'pass' => 'PASS',
//     'hotspotname' => '...', 'dnsname' => '...', 'currency' => '...',
//     'phone'=>'...','email'=>'...','infolp'=>'...','idleto'=>'...','report'=>'yes'
//     );
//
// Usage: go run ./cmd/migrate [path/to/config.php]
func main() {
	phpPath := "config/config.php"
	if len(os.Args) > 1 {
		phpPath = os.Args[1]
	}

	data, err := os.ReadFile(phpPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", phpPath, err)
		os.Exit(1)
	}
	content := string(data)

	// ── Parse admin block ──────────────────────────────────────────────────────
	// Pattern: $data['mikhmon'] = array('1' => 'mikhmon<|<USERNAME', 'mikhmon>|>BASE64PASS');
	adminRe := regexp.MustCompile(
		`\$data\['mikhmon'\]\s*=\s*array\([^)]*'mikhmon<\|<([^']+)'[^)]*'mikhmon>\|>([^']+)'`)
	adminUser := "mikhmon"
	adminPass := ""
	if m := adminRe.FindStringSubmatch(content); m != nil {
		adminUser = m[1]
		// Password is base64 encoded in the PHP config.
		if decoded, err := base64.StdEncoding.DecodeString(m[2]); err == nil {
			adminPass = string(decoded)
		} else {
			// Not base64 — use raw value.
			adminPass = m[2]
		}
	}

	// ── Parse router blocks ────────────────────────────────────────────────────
	// Pattern: $data['SESS'] = array( ... key => 'value' ... );
	// Skip the mikhmon admin block.
	routerBlockRe := regexp.MustCompile(
		`\$data\['([^']+)'\]\s*=\s*array\(([^)]+)\)`)
	kvRe := regexp.MustCompile(`'([^']+)'\s*=>\s*'([^']*)'`)

	var routers []config.RouterConfig
	for _, block := range routerBlockRe.FindAllStringSubmatch(content, -1) {
		sess := block[1]
		if sess == "mikhmon" {
			continue // already handled as admin
		}
		body := block[2]
		kv := make(map[string]string)
		for _, pair := range kvRe.FindAllStringSubmatch(body, -1) {
			kv[pair[1]] = pair[2]
		}
		r := config.RouterConfig{
			SessionName: sanitize(sess),
			Host:        kv["host"],
			Username:    kv["user"],
			Password:    kv["pass"],
			HotspotName: kv["hotspotname"],
			DNSName:     kv["dnsname"],
			Currency:    kv["currency"],
			Phone:       kv["phone"],
			Email:       kv["email"],
			InfoLP:      kv["infolp"],
			IdleTimeout: kv["idleto"],
			ReportFlag:  kv["report"],
		}
		if r.SessionName != "" {
			routers = append(routers, r)
		}
	}

	// ── Build and write config ─────────────────────────────────────────────────
	adminHash, err := config.HashPassword(adminPass)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not hash admin password: %v\n", err)
		adminHash = ""
	}

	cfg := &config.AppConfig{
		Admin: config.AdminConfig{
			Username:     adminUser,
			PasswordHash: adminHash,
		},
		Routers: routers,
		Server: config.ServerConfig{
			Port:          8080,
			SessionSecret: "change-me-in-production",
		},
	}

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal error: %v\n", err)
		os.Exit(1)
	}

	outPath := "config/mikhmon.json"
	if err := os.WriteFile(outPath, out, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "error writing %s: %v\n", outPath, err)
		os.Exit(1)
	}

	fmt.Printf("✓ Migrated %d router(s) to %s\n", len(routers), outPath)
	fmt.Printf("  Admin: %s / %s (original plain-text)\n", adminUser, adminPass)
	fmt.Println("  ⚠  Remember to update Server.SessionSecret in", outPath)
}

func sanitize(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return re.ReplaceAllString(strings.ReplaceAll(s, " ", ""), "")
}
