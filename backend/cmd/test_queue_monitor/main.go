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
	fmt.Println("MikroTik Queue Monitor Test")
	fmt.Println("==========================================")
	fmt.Printf("Router: %s:%d\n", RouterHost, RouterPort)
	fmt.Println()

	testGetAllParentQueues(ctx, client, router)
	fmt.Println()
	testStartQueueStatsListen(ctx, client, router)

	// client.CloseAll()
}

func testGetAllParentQueues(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- Get All Parent Queues ---")

	queues, err := client.GetAllParentQueues(ctx, router)
	if err != nil {
		log.Printf("GetAllParentQueues Error: %v\n", err)
		return
	}

	log.Printf("Total Parent Queues: %d\n", len(queues))
	fmt.Printf("Parent Queues: %v\n", queues)

	if len(queues) > 0 {
		log.Printf("First Queue: %s\n", queues[0])
	}
}

func testStartQueueStatsListen(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- Start Queue Stats Monitor (5 seconds) ---")

	queues, err := client.GetAllParentQueues(ctx, router)
	if err != nil {
		log.Printf("GetAllParentQueues Error: %v\n", err)
		return
	}
	if len(queues) == 0 {
		log.Println("No parent queues found")
		return
	}

	queueName := queues[0]
	fmt.Printf("Monitoring queue: %s\n", queueName)

	cfg := mikrotik.DefaultQueueStatsConfig(queueName)
	resultChan := make(chan mikrotik.QueueStats, 100)

	cancel, err := client.StartQueueStatsListen(ctx, router, cfg, resultChan)
	if err != nil {
		log.Printf("StartQueueStatsListen Error: %v\n", err)
		return
	}
	defer cancel()

	fmt.Println("Monitoring started (5 seconds)...")
	fmt.Println()

	for i := 0; i < 5; i++ {
		select {
		case stats := <-resultChan:
			fmt.Printf("[%d] Queue: %s\n", i+1, stats.Name)
			fmt.Printf("  Bytes: In=%s, Out=%s\n", formatBytes(stats.BytesIn), formatBytes(stats.BytesOut))
			fmt.Printf("  Rate: In=%s, Out=%s\n", formatRate(stats.RateIn), formatRate(stats.RateOut))
			fmt.Printf("  Packets: In=%d, Out=%d\n", stats.PacketsIn, stats.PacketsOut)
			fmt.Printf("  Packet Rate: In=%d pps, Out=%d pps\n", stats.PacketRateIn, stats.PacketRateOut)
		case <-time.After(100 * time.Second):
			fmt.Printf("[%.2fs] No data received\n", float64(i+1))
		}
	}

	fmt.Println()
	fmt.Println("--- Monitoring Stopped ---")
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < unit*unit {
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(unit))
	} else if bytes < unit*unit*unit {
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(unit*unit))
	} else {
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(unit*unit*unit))
	}
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
