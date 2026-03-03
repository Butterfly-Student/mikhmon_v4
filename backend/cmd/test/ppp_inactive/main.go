// Package main provides a test program for MikroTik PPP Inactive operations.
//
// This program directly tests the functions from:
// - internal/infrastructure/mikrotik/ppp_active.go (ListenPPPActive, ListenPPPInactive)
// - internal/infrastructure/mikrotik/ppp_secret.go (GetPPPSecrets)
// by connecting to a real MikroTik router and executing the operations.
//
// Usage:
//
//	go run ./cmd/test/ppp_inactive.go -host=192.168.88.1 -user=admin -pass=admin
//
// Or using environment variables:
//
//	set MIKROTIK_HOST=192.168.88.1
//	set MIKROTIK_USER=admin
//	set MIKROTIK_PASS=admin
//	go run ./cmd/test/ppp_inactive.go
//
// Flags:
//   -host string
//		MikroTik host (default "192.168.88.1")
//   -port int
//		MikroTik API port (default 8728)
//   -user string
//		MikroTik username (default "admin")
//   -pass string
//		MikroTik password
//   -tls
//		Use TLS connection
//   -profile string
//		Filter PPP secrets by profile (optional)
//   -listen duration
//		Run PPP listeners for specified duration (e.g., 10s, 30s, 1m)
//   -listen-active
//		Enable ListenPPPActive test
//   -listen-inactive
//		Enable ListenPPPInactive test
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"go.uber.org/zap"
)

var (
	host           string
	port           int
	user           string
	pass           string
	useTLS         bool
	profileFilter  string
	listenDuration time.Duration
	listenActive   bool
	listenInactive bool
)

func init() {
	// Default values from environment or hardcoded defaults
	defaultHost := getEnv("MIKROTIK_HOST", "192.168.233.1")
	defaultPort, _ := strconv.Atoi(getEnv("MIKROTIK_PORT", "8728"))
	defaultUser := getEnv("MIKROTIK_USER", "admin")
	defaultPass := getEnv("MIKROTIK_PASS", "r00t")
	defaultTLS := getEnv("MIKROTIK_TLS", "false") == "true"

	flag.StringVar(&host, "host", defaultHost, "MikroTik host")
	flag.IntVar(&port, "port", defaultPort, "MikroTik API port")
	flag.StringVar(&user, "user", defaultUser, "MikroTik username")
	flag.StringVar(&pass, "pass", defaultPass, "MikroTik password")
	flag.BoolVar(&useTLS, "tls", defaultTLS, "Use TLS connection")
	flag.StringVar(&profileFilter, "profile", "", "Filter PPP secrets by profile (optional)")
	flag.DurationVar(&listenDuration, "listen", 0, "Run PPP listeners for specified duration (e.g., 10s, 30s, 1m)")
	flag.BoolVar(&listenActive, "listen-active", false, "Enable ListenPPPActive test")
	flag.BoolVar(&listenInactive, "listen-inactive", false, "Enable ListenPPPInactive test")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	flag.Parse()

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║          MikroTik PPP Inactive Test Program                ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Validate connection parameters
	if host == "" {
		fmt.Println("❌ Error: MikroTik host is required")
		fmt.Println("   Use -host flag or set MIKROTIK_HOST environment variable")
		os.Exit(1)
	}

	fmt.Printf("📡 Connecting to MikroTik router...\n")
	fmt.Printf("   Host:    %s:%d\n", host, port)
	fmt.Printf("   User:    %s\n", user)
	fmt.Printf("   TLS:     %v\n", useTLS)
	if profileFilter != "" {
		fmt.Printf("   Profile: %s\n", profileFilter)
	}
	fmt.Println()

	// Create logger
	logger := zap.NewNop()

	// Create client config
	cfg := mikrotik.Config{
		Host:     host,
		Port:     port,
		Username: user,
		Password: pass,
		UseTLS:   useTLS,
		Timeout:  10 * time.Second,
	}

	// Create and connect client
	client := mikrotik.NewClient(cfg, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		fmt.Printf("❌ Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Println("✅ Connected successfully!")
	fmt.Println()

	// Run tests
	testCount := 0
	passCount := 0

	// Test 1: GetPPPSecrets
	testCount++
	if testGetPPPSecrets(client) {
		passCount++
	}

	// // Test 2: ListenPPPActive (if requested)
	// if listenActive || listenDuration > 0 {
	// 	testCount++
	// 	if testListenPPPActive(client, listenDuration) {
	// 		passCount++
	// 	}
	// }

	// Test 3: ListenPPPInactive (if requested)
	if listenInactive || listenDuration > 0 {
		testCount++
		if testListenPPPInactive(client, listenDuration) {
			passCount++
		}
	}

	// Summary
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Printf("║  Test Results: %d/%d passed                                  ║\n", passCount, testCount)
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	if passCount == testCount {
		fmt.Println("✅ All tests passed!")
		os.Exit(0)
	} else {
		fmt.Printf("❌ %d test(s) failed\n", testCount-passCount)
		os.Exit(1)
	}
}

func testGetPPPSecrets(client *mikrotik.Client) bool {
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Test: GetPPPSecrets                                        │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	secrets, err := client.GetPPPSecrets(ctx, profileFilter)
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
		return false
	}

	fmt.Printf("✅ SUCCESS - Found %d PPP secret(s)\n", len(secrets))
	if len(secrets) > 0 {
		fmt.Println()
		fmt.Println("   PPP Secrets:")
		for i, s := range secrets {
			if i >= 10 {
				fmt.Printf("   ... and %d more\n", len(secrets)-10)
				break
			}
			status := "enabled"
			if s.Disabled {
				status = "disabled"
			}
			fmt.Printf("   - %-20s | Profile: %-15s | %s\n", s.Name, s.Profile, status)
		}
	}
	fmt.Println()
	return true
}

func testListenPPPActive(client *mikrotik.Client, duration time.Duration) bool {
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Test: ListenPPPActive                                      │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")

	if duration == 0 {
		duration = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	resultChan := make(chan []*dto.PPPActive, 10)

	stopFunc, err := client.ListenPPPActive(ctx, resultChan)
	if err != nil {
		fmt.Printf("❌ FAILED to start listener: %v\n", err)
		return false
	}
	defer stopFunc()

	fmt.Printf("✅ Listener started (duration: %s)\n", duration)
	fmt.Println("   Collecting active PPP sessions...")
	fmt.Println()

	sampleCount := 0
	maxSamples := 5
	uniqueSessions := make(map[string]bool)

	for {
		select {
		case batch, ok := <-resultChan:
			if !ok {
				fmt.Println("   Channel closed")
				goto done
			}
			sampleCount++
			for _, session := range batch {
				uniqueSessions[session.Name] = true
				fmt.Printf("   [%d] User: %-20s | Service: %-10s | Uptime: %s\n",
					sampleCount,
					session.Name,
					session.Service,
					session.Uptime,
				)
			}
			if sampleCount >= maxSamples {
				fmt.Println("   (Reached max samples)")
				goto done
			}

		case <-ctx.Done():
			fmt.Println("   (Timeout reached)")
			goto done
		}
	}

done:
	fmt.Printf("\n✅ Collected %d update(s), %d unique session(s)\n", sampleCount, len(uniqueSessions))
	fmt.Println()
	return true
}

func testListenPPPInactive(client *mikrotik.Client, duration time.Duration) bool {
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Test: ListenPPPInactive                                    │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")

	if duration == 0 {
		duration = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	resultChan := make(chan []*dto.PPPSecret, 10)

	stopFunc, err := client.ListenPPPInactive(ctx, resultChan)
	if err != nil {
		fmt.Printf("❌ FAILED to start listener: %v\n", err)
		return false
	}
	defer stopFunc()

	fmt.Printf("✅ Listener started (duration: %s)\n", duration)
	fmt.Println("   Collecting inactive PPP secrets...")
	fmt.Println()

	sampleCount := 0
	maxSamples := 5
	var lastInactive []*dto.PPPSecret

	for {
		select {
		case inactive, ok := <-resultChan:
			if !ok {
				fmt.Println("   Channel closed")
				goto doneInactive
			}
			sampleCount++
			lastInactive = inactive
			fmt.Printf("   [%d] Inactive secrets count: %d\n", sampleCount, len(inactive))
			for i, s := range inactive {
				if i >= 5 {
					fmt.Printf("       ... and %d more\n", len(inactive)-5)
					break
				}
				fmt.Printf("       - %s (Profile: %s)\n", s.Name, s.Profile)
			}
			if sampleCount >= maxSamples {
				fmt.Println("   (Reached max samples)")
				goto doneInactive
			}

		case <-ctx.Done():
			fmt.Println("   (Timeout reached)")
			goto doneInactive
		}
	}
doneInactive:
	fmt.Printf("\n✅ Collected %d update(s), final count: %d inactive secret(s)\n",
		sampleCount, len(lastInactive))
	fmt.Println()
	return true
}
