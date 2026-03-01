// Package main provides a test program for MikroTik system operations.
//
// This program directly tests the functions from internal/infrastructure/mikrotik/system.go
// by connecting to a real MikroTik router and executing the operations.
//
// Usage:
//
//	go run ./cmd/test/ -host=192.168.88.1 -user=admin -pass=admin
//
// Or using environment variables:
//
//	set MIKROTIK_HOST=192.168.88.1
//	set MIKROTIK_USER=admin
//	set MIKROTIK_PASS=admin
//	go run ./cmd/test/
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
//   -monitor duration
//		Run resource monitor for specified duration (e.g., 5s, 10s)
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
	host    string
	port    int
	user    string
	pass    string
	useTLS  bool
	monitor time.Duration
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
	flag.DurationVar(&monitor, "monitor", 0, "Run resource monitor for specified duration (e.g., 5s, 10s)")
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
	fmt.Println("║          MikroTik System Test Program                      ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Validate connection parameters
	if host == "" {
		fmt.Println("❌ Error: MikroTik host is required")
		fmt.Println("   Use -host flag or set MIKROTIK_HOST environment variable")
		os.Exit(1)
	}

	fmt.Printf("📡 Connecting to MikroTik router...\n")
	fmt.Printf("   Host: %s:%d\n", host, port)
	fmt.Printf("   User: %s\n", user)
	fmt.Printf("   TLS:  %v\n", useTLS)
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

	// Test 1: GetSystemResource
	testCount++
	if testGetSystemResource(client) {
		passCount++
	}

	// Test 2: GetSystemHealth
	testCount++
	if testGetSystemHealth(client) {
		passCount++
	}

	// Test 3: GetSystemIdentity
	testCount++
	if testGetSystemIdentity(client) {
		passCount++
	}

	// Test 4: GetSystemClock
	testCount++
	if testGetSystemClock(client) {
		passCount++
	}

	// Test 5: GetRouterBoardInfo
	testCount++
	if testGetRouterBoardInfo(client) {
		passCount++
	}

	// Test 6: StartSystemResourceMonitorListen (if requested)
	if monitor > 0 {
		testCount++
		if testSystemResourceMonitor(client, monitor) {
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

func testGetSystemResource(client *mikrotik.Client) bool {
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Test: GetSystemResource                                    │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resource, err := client.GetSystemResource(ctx)
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
		return false
	}

	fmt.Println("✅ SUCCESS")
	fmt.Printf("   Uptime:        %s\n", resource.Uptime)
	fmt.Printf("   Version:       %s\n", resource.Version)
	fmt.Printf("   Build Time:    %s\n", resource.BuildTime)
	fmt.Printf("   Memory:        %s / %s\n", formatBytes(resource.FreeMemory), formatBytes(resource.TotalMemory))
	fmt.Printf("   HDD Space:     %s / %s\n", formatBytes(resource.FreeHddSpace), formatBytes(resource.TotalHddSpace))
	fmt.Printf("   CPU:           %s (%d cores @ %dMHz)\n", resource.Cpu, resource.CpuCount, resource.CpuFrequency)
	fmt.Printf("   CPU Load:      %d%%\n", resource.CpuLoad)
	fmt.Printf("   Architecture:  %s\n", resource.ArchitectureName)
	fmt.Printf("   Board:         %s\n", resource.BoardName)
	fmt.Printf("   Platform:      %s\n", resource.Platform)
	fmt.Printf("   Bad Blocks:    %.2f%%\n", resource.BadBlocks)
	fmt.Println()
	return true
}

func testGetSystemHealth(client *mikrotik.Client) bool {
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Test: GetSystemHealth                                      │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	health, err := client.GetSystemHealth(ctx)
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
		return false
	}

	fmt.Println("✅ SUCCESS")
	if health.Voltage != "" {
		fmt.Printf("   Voltage:     %s\n", health.Voltage)
	}
	if health.Temperature != "" {
		fmt.Printf("   Temperature: %s\n", health.Temperature)
	}
	if health.FanSpeed != "" {
		fmt.Printf("   Fan Speed:   %s\n", health.FanSpeed)
	}
	if health.Voltage == "" && health.Temperature == "" && health.FanSpeed == "" {
		fmt.Println("   (No health data available on this router)")
	}
	fmt.Println()
	return true
}

func testGetSystemIdentity(client *mikrotik.Client) bool {
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Test: GetSystemIdentity                                    │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	identity, err := client.GetSystemIdentity(ctx)
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
		return false
	}

	fmt.Println("✅ SUCCESS")
	fmt.Printf("   Identity Name: %s\n", identity.Name)
	fmt.Println()
	return true
}

func testGetSystemClock(client *mikrotik.Client) bool {
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Test: GetSystemClock                                       │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clock, err := client.GetSystemClock(ctx)
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
		return false
	}

	fmt.Println("✅ SUCCESS")
	fmt.Printf("   Time:         %s\n", clock.Time)
	fmt.Printf("   Date:         %s\n", clock.Date)
	fmt.Printf("   Time Zone:    %s\n", clock.TimeZoneName)
	fmt.Printf("   DST Active:   %s\n", clock.DSTActive)
	fmt.Println()
	return true
}

func testGetRouterBoardInfo(client *mikrotik.Client) bool {
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Test: GetRouterBoardInfo                                   │")
	fmt.Println("└────────────────────────────────────────────────────────────┘")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	info, err := client.GetRouterBoardInfo(ctx)
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
		return false
	}

	fmt.Println("✅ SUCCESS")
	fmt.Printf("   RouterBoard:     %s\n", info.RouterBoard)
	fmt.Printf("   Model:           %s\n", info.Model)
	fmt.Printf("   Serial Number:   %s\n", info.SerialNumber)
	fmt.Printf("   Firmware Type:   %s\n", info.FirmwareType)
	fmt.Printf("   Current Firmware: %s\n", info.CurrentFirmware)
	fmt.Printf("   Factory Firmware: %s\n", info.FactoryFirmware)
	if info.UpgradeFirmware != "" {
		fmt.Printf("   Upgrade Firmware: %s\n", info.UpgradeFirmware)
	}
	fmt.Println()
	return true
}

func testSystemResourceMonitor(client *mikrotik.Client, duration time.Duration) bool {
	fmt.Println("┌────────────────────────────────────────────────────────────┐")
	fmt.Printf("│ Test: StartSystemResourceMonitorListen (%s)              │\n", duration)
	fmt.Println("└────────────────────────────────────────────────────────────┘")

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	resultChan := make(chan dto.SystemResourceMonitorStats, 10)

	stopFunc, err := client.StartSystemResourceMonitorListen(ctx, resultChan)
	if err != nil {
		fmt.Printf("❌ FAILED to start monitor: %v\n", err)
		return false
	}
	defer stopFunc()

	fmt.Println("✅ Monitor started successfully")
	fmt.Println("   Collecting samples...")
	fmt.Println()

	sampleCount := 0
	maxSamples := 5

	for {
		select {
		case stats, ok := <-resultChan:
			if !ok {
				fmt.Println("   Channel closed")
				goto done
			}
			sampleCount++
			fmt.Printf("   Sample %d: CPU=%d%% | Memory=%s/%s | Uptime=%s\n",
				sampleCount,
				stats.CPULoad,
				formatBytes(stats.FreeMemory),
				formatBytes(stats.TotalMemory),
				stats.Uptime,
			)
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
	fmt.Printf("\n✅ Collected %d samples successfully\n", sampleCount)
	fmt.Println()
	return sampleCount > 0
}

// formatBytes converts bytes to human-readable format
func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	if bytes >= GB {
		return fmt.Sprintf("%.2f GiB", float64(bytes)/float64(GB))
	}
	if bytes >= MB {
		return fmt.Sprintf("%.2f MiB", float64(bytes)/float64(MB))
	}
	if bytes >= KB {
		return fmt.Sprintf("%.2f KiB", float64(bytes)/float64(KB))
	}
	return fmt.Sprintf("%d B", bytes)
}
