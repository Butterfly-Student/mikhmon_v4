package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/entity"
	"github.com/irhabi89/mikhmon/internal/infrastructure/logger"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
)

const (
	RouterHost     = "192.168.233.1"
	RouterPort     = 8728
	RouterUser     = "admin"
	RouterPassword = "r00t"
)

func main() {
	logger.FromEnv()

	router := &entity.Router{
		Name:     "Test Router",
		Host:     RouterHost,
		Port:     RouterPort,
		Username: RouterUser,
		Password: RouterPassword,
		UseSSL:   false,
		Timeout:  10,
	}

	client := mikrotik.NewClient(logger.Log)
	ctx := context.Background()

	fmt.Println("==========================================")
	fmt.Println("MikroTik Traffic Monitor Test")
	fmt.Println("==========================================")
	fmt.Printf("Router: %s:%d\n", RouterHost, RouterPort)
	fmt.Println()

	// Monitor ether1
	ifaceName := "ether1"
	fmt.Printf("Monitoring interface: %s (5 seconds)...\n\n", ifaceName)

	resultChan := make(chan mikrotik.TrafficMonitorStats, 100)
	cancel, err := client.StartTrafficMonitorListen(ctx, router, ifaceName, resultChan)
	if err != nil {
		log.Printf("StartTrafficMonitorListen Error: %v\n", err)
		return
	}
	defer cancel()

	for i := 0; i < 5; i++ {
		select {
		case stats := <-resultChan:
			fmt.Printf("[%d] Interface: %s\n", i+1, stats.Name)
			fmt.Printf("  RX: %s (%d pps)\n", formatRate(stats.RxBitsPerSecond), stats.RxPacketsPerSecond)
			fmt.Printf("  TX: %s (%d pps)\n", formatRate(stats.TxBitsPerSecond), stats.TxPacketsPerSecond)
			fmt.Printf("  FP-RX: %s (%d pps)\n", formatRate(stats.FpRxBitsPerSecond), stats.FpRxPacketsPerSecond)
			fmt.Printf("  FP-TX: %s (%d pps)\n", formatRate(stats.FpTxBitsPerSecond), stats.FpTxPacketsPerSecond)
			fmt.Printf("  Drops: RX=%d/s, TX=%d/s (queue=%d/s)\n", stats.RxDropsPerSecond, stats.TxDropsPerSecond, stats.TxQueueDropsPerSecond)
			fmt.Printf("  Errors: RX=%d/s, TX=%d/s\n", stats.RxErrorsPerSecond, stats.TxErrorsPerSecond)
			fmt.Println()
		case <-time.After(10 * time.Second):
			fmt.Printf("[%d] Timeout - No data received\n", i+1)
		}
	}

	fmt.Println("--- Monitoring Stopped ---")
}

func formatRate(bps int64) string {
	if bps < 1000 {
		return fmt.Sprintf("%d bps", bps)
	} else if bps < 1000000 {
		return fmt.Sprintf("%.2f Kbps", float64(bps)/1000)
	} else if bps < 1000000000 {
		return fmt.Sprintf("%.2f Mbps", float64(bps)/1000000)
	} else {
		return fmt.Sprintf("%.2f Gbps", float64(bps)/1000000000)
	}
}
