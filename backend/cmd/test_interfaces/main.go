package main

import (
	"context"
	"fmt"
	"log"

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
	fmt.Println("MikroTik Interface List Test")
	fmt.Println("==========================================")
	fmt.Printf("Router: %s:%d\n", RouterHost, RouterPort)
	fmt.Println()

	interfaces, err := client.GetInterfaces(ctx, router)
	if err != nil {
		log.Printf("GetInterfaces Error: %v\n", err)
		return
	}

	fmt.Printf("Total Interfaces: %d\n\n", len(interfaces))

	fmt.Printf("%-5s %-20s %-15s %-10s %-8s %-10s %-20s\n", "ID", "NAME", "TYPE", "MTU", "RUNNING", "DISABLED", "MAC-ADDRESS")
	fmt.Println("----------------------------------------------------------------------------------------------------")
	
	for _, iface := range interfaces {
		running := "No"
		if iface.Running {
			running = "Yes"
		}
		disabled := "No"
		if iface.Disabled {
			disabled = "Yes"
		}
		fmt.Printf("%-5s %-20s %-15s %-10d %-8s %-10s %-20s\n",
			iface.ID, iface.Name, iface.Type, iface.MTU, running, disabled, iface.MacAddress)
	}
}
