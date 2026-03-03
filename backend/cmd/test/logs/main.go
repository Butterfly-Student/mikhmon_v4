// Package main provides a test program for MikroTik Log streaming operations.
//
// This program tests ListenAllLogs, ListenHotspotLogs, or ListenPPPLogs.
// Run one test at a time.
//
// Usage:
//
//	go run ./cmd/test/logs -type=all -duration=10s
//	go run ./cmd/test/logs -type=hotspot -duration=10s
//	go run ./cmd/test/logs -type=ppp -duration=10s
//
// Environment variables:
//	MIKROTIK_HOST (default: 192.168.233.1)
//	MIKROTIK_PORT (default: 8728)
//	MIKROTIK_USER (default: admin)
//	MIKROTIK_PASS (default: r00t)
//	MIKROTIK_TLS (default: false)
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
	host     string
	port     int
	user     string
	pass     string
	useTLS   bool
	logType  string
	duration time.Duration
)

func init() {
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
	flag.StringVar(&logType, "type", "all", "Log type to test: all, hotspot, ppp")
	flag.DurationVar(&duration, "duration", 10*time.Second, "Listen duration (e.g., 10s, 30s, 1m)")
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
	fmt.Println("║          MikroTik Logs Streaming Test Program              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	if host == "" {
		fmt.Println("❌ Error: MikroTik host is required")
		os.Exit(1)
	}

	fmt.Printf("📡 Configuration:\n")
	fmt.Printf("   Host:     %s:%d\n", host, port)
	fmt.Printf("   User:     %s\n", user)
	fmt.Printf("   Type:     %s\n", logType)
	fmt.Printf("   Duration: %s\n", duration)
	fmt.Println()

	// Create client
	logger := zap.NewNop()
	cfg := mikrotik.Config{
		Host:     host,
		Port:     port,
		Username: user,
		Password: pass,
		UseTLS:   useTLS,
		Timeout:  10 * time.Second,
	}

	client := mikrotik.NewClient(cfg, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	if err := client.Connect(ctx); err != nil {
		cancel()
		fmt.Printf("❌ Failed to connect: %v\n", err)
		os.Exit(1)
	}
	cancel()
	defer client.Close()

	fmt.Println("✅ Connected successfully!")
	fmt.Println()

	// Run test based on type
	var success bool
	switch logType {
	case "all":
		success = testListenAllLogs(client, duration)
	case "hotspot":
		success = testListenHotspotLogs(client, duration)
	case "ppp":
		success = testListenPPPLogs(client, duration)
	default:
		fmt.Printf("❌ Unknown log type: %s (use: all, hotspot, ppp)\n", logType)
		os.Exit(1)
	}

	fmt.Println()
	if success {
		fmt.Println("╔════════════════════════════════════════════════════════════╗")
		fmt.Println("║  ✅ Test PASSED                                            ║")
		fmt.Println("╚════════════════════════════════════════════════════════════╝")
		os.Exit(0)
	} else {
		fmt.Println("╔════════════════════════════════════════════════════════════╗")
		fmt.Println("║  ❌ Test FAILED                                            ║")
		fmt.Println("╚════════════════════════════════════════════════════════════╝")
		os.Exit(1)
	}
}

func testListenAllLogs(client *mikrotik.Client, d time.Duration) bool {
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Test: ListenAllLogs                                        │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")

	ctx, cancel := context.WithTimeout(context.Background(), d)

	resultChan := make(chan *dto.LogEntry, 20)
	stopFunc, err := client.ListenAllLogs(ctx, resultChan)
	if err != nil {
		fmt.Printf("❌ FAILED to start listener: %v\n", err)
		return false
	}
	fmt.Printf("✅ Listener started (will run for %s)\n", d)
	fmt.Println("   Collecting all system logs...")
	fmt.Println()

	logCount := 0
	maxLogs := 1000
	done := false

	for !done {
		select {
		case log, ok := <-resultChan:
			if !ok {
				fmt.Println("   (Channel closed)")
				done = true
				break
			}
			logCount++
			fmt.Printf("   [%d] [%s] %s | %s\n",
				logCount,
				log.Time,
				log.Topics,
				truncateString(log.Message, 50),
			)
			if logCount >= maxLogs {
				fmt.Println("   (Reached max logs)")
				done = true
			}

		case <-ctx.Done():
			fmt.Println("   (Timeout reached)")
			done = true
		}
	}

	fmt.Printf("\n✅ Collected %d log(s) successfully\n", logCount)
	
	// Clean up: stop listener and cancel context
	// Use goroutine to avoid blocking on stopFunc
	stopDone := make(chan struct{})
	go func() {
		stopFunc()
		close(stopDone)
	}()
	
	select {
	case <-stopDone:
		// stopped successfully
	case <-time.After(2 * time.Second):
		fmt.Println("   (Stop timeout - forcing cancel)")
	}
	cancel()
	
	return true
}

func testListenHotspotLogs(client *mikrotik.Client, d time.Duration) bool {
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Test: ListenHotspotLogs                                    │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")

	ctx, cancel := context.WithTimeout(context.Background(), d)

	resultChan := make(chan *dto.LogEntry, 20)
	stopFunc, err := client.ListenHotspotLogs(ctx, resultChan)
	if err != nil {
		fmt.Printf("❌ FAILED to start listener: %v\n", err)
		return false
	}

	fmt.Printf("✅ Listener started (will run for %s)\n", d)
	fmt.Println("   Collecting hotspot logs...")
	fmt.Println("   Tip: Generate hotspot activity to see logs")
	fmt.Println()

	logCount := 0
	maxLogs := 20

	for {
		select {
		case log, ok := <-resultChan:
			if !ok {
				fmt.Println("   (Channel closed)")
				goto done
			}
			logCount++
			fmt.Printf("   [%d] [%s] %s | %s\n",
				logCount,
				log.Time,
				log.Topics,
				truncateString(log.Message, 50),
			)
			if logCount >= maxLogs {
				fmt.Println("   (Reached max logs)")
				goto done
			}

		case <-ctx.Done():
			fmt.Println("   (Timeout reached)")
			goto done
		}
	}
done:
	fmt.Printf("\n✅ Collected %d hotspot log(s) successfully\n", logCount)
	
	stopDone := make(chan struct{})
	go func() {
		stopFunc()
		close(stopDone)
	}()
	
	select {
	case <-stopDone:
	case <-time.After(2 * time.Second):
		fmt.Println("   (Stop timeout - forcing cancel)")
	}
	cancel()
	
	return true
}

func testListenPPPLogs(client *mikrotik.Client, d time.Duration) bool {
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Test: ListenPPPLogs                                        │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")

	ctx, cancel := context.WithTimeout(context.Background(), d)

	resultChan := make(chan *dto.LogEntry, 20)
	stopFunc, err := client.ListenPPPLogs(ctx, resultChan)
	if err != nil {
		fmt.Printf("❌ FAILED to start listener: %v\n", err)
		return false
	}

	fmt.Printf("✅ Listener started (will run for %s)\n", d)
	fmt.Println("   Collecting PPP logs...")
	fmt.Println("   Tip: Connect/disconnect a PPP user to see logs")
	fmt.Println()

	logCount := 0
	maxLogs := 0

	for {
		select {
		case log, ok := <-resultChan:
			if !ok {
				fmt.Println("   (Channel closed)")
				goto done
			}
			logCount++
			fmt.Printf("   [%d] [%s] %s | %s\n",
				logCount,
				log.Time,
				log.Topics,
				truncateString(log.Message, 50),
			)
			if logCount >= maxLogs {
				fmt.Println("   (Reached max logs)")
				goto done
			}

		case <-ctx.Done():
			fmt.Println("   (Timeout reached)")
			goto done
		}
	}
done:
	fmt.Printf("\n✅ Collected %d PPP log(s) successfully\n", logCount)
	
	stopDone := make(chan struct{})
	go func() {
		stopFunc()
		close(stopDone)
	}()
	
	select {
	case <-stopDone:
	case <-time.After(2 * time.Second):
		fmt.Println("   (Stop timeout - forcing cancel)")
	}
	cancel()
	
	return true
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
