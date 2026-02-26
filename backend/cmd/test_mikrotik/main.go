package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
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
	fmt.Println("MikroTik Client Direct Testing")
	fmt.Println("==========================================")
	fmt.Printf("Router: %s:%d\n", RouterHost, RouterPort)
	fmt.Println()

	testHotspotUsers(ctx, client, router)
	fmt.Println()

	testHotspotProfiles(ctx, client, router)
	fmt.Println()

	testHotspotActive(ctx, client, router)
	fmt.Println()

	testHotspotHosts(ctx, client, router)
	fmt.Println()

	testHotspotServers(ctx, client, router)
	fmt.Println()

	testSystem(ctx, client, router)
	fmt.Println()

	testInterfaces(ctx, client, router)
	fmt.Println()

	testNAT(ctx, client, router)
	fmt.Println()

	testQueues(ctx, client, router)
	fmt.Println()

	testPools(ctx, client, router)
	fmt.Println()

	testLogs(ctx, client, router)
	fmt.Println()

	testReports(ctx, client, router)
	fmt.Println()

	testPing(ctx, client, router)
	fmt.Println()

	testExpireMonitor(ctx, client, router)
	fmt.Println()

	testVoucherGenerator(ctx)
	fmt.Println()

	testOnLoginGenerator(ctx)
	fmt.Println()

	client.CloseAll()
}

func testHotspotUsers(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- Hotspot Users ---")

	users, err := client.GetHotspotUsers(ctx, router, "")
	if err != nil {
		log.Printf("GetHotspotUsers Error: %v\n", err)
		return
	}
	log.Printf("Total Users: %d\n", len(users))
	if len(users) > 0 {
		log.Printf("First User: %+v\n", users[0])
	}

	count, err := client.GetHotspotUsersCount(ctx, router)
	if err != nil {
		log.Printf("GetHotspotUsersCount Error: %v\n", err)
	} else {
		log.Printf("Users Count: %d\n", count)
	}
}

func testHotspotProfiles(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- Hotspot Profiles ---")

	profiles, err := client.GetUserProfiles(ctx, router)
	if err != nil {
		log.Printf("GetUserProfiles Error: %v\n", err)
		return
	}
	log.Printf("Total Profiles: %d\n", len(profiles))
	if len(profiles) > 0 {
		log.Printf("First Profile: Name=%s, ExpireMode=%s, Price=%.0f\n",
			profiles[0].Name, profiles[0].ExpireMode, profiles[0].Price)
	}

	if len(profiles) > 0 {
		profile, err := client.GetUserProfileByID(ctx, router, profiles[0].ID)
		if err != nil {
			log.Printf("GetUserProfileByID Error: %v\n", err)
		} else {
			log.Printf("GetProfileByID: Name=%s\n", profile.Name)
		}
	}

	profile, err := client.GetUserProfileByName(ctx, router, profiles[0].Name)
	if err != nil {
		log.Printf("GetUserProfileByName Error: %v\n", err)
	} else {
		log.Printf("GetProfileByName: Name=%s\n", profile.Name)
	}
}

func testHotspotActive(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- Hotspot Active Sessions ---")

	active, err := client.GetHotspotActive(ctx, router)
	if err != nil {
		log.Printf("GetHotspotActive Error: %v\n", err)
		return
	}
	log.Printf("Active Sessions: %d\n", len(active))
	if len(active) > 0 {
		log.Printf("First Active: User=%s, Address=%s, Uptime=%s\n",
			active[0].User, active[0].Address, active[0].Uptime)
	}

	count, err := client.GetHotspotActiveCount(ctx, router)
	if err != nil {
		log.Printf("GetHotspotActiveCount Error: %v\n", err)
	} else {
		log.Printf("Active Count: %d\n", count)
	}
}

func testHotspotHosts(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- Hotspot Hosts ---")

	hosts, err := client.GetHotspotHosts(ctx, router)
	if err != nil {
		log.Printf("GetHotspotHosts Error: %v\n", err)
		return
	}
	log.Printf("Total Hosts: %d\n", len(hosts))
	if len(hosts) > 0 {
		log.Printf("First Host: MAC=%s, Address=%s\n", hosts[0].MACAddress, hosts[0].Address)
	}
}

func testHotspotServers(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- Hotspot Servers ---")

	servers, err := client.GetHotspotServers(ctx, router)
	if err != nil {
		log.Printf("GetHotspotServers Error: %v\n", err)
		return
	}
	log.Printf("Servers: %v\n", servers)
}

func testSystem(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- System Info ---")

	resource, err := client.GetSystemResource(ctx, router)
	if err != nil {
		log.Printf("GetSystemResource Error: %v\n", err)
	} else {
		log.Printf("Resource: CPU=%d%%, FreeMemory=%dMB\n", resource.CpuLoad, resource.FreeMemory/1024/1024)
	}

	health, err := client.GetSystemHealth(ctx, router)
	if err != nil {
		log.Printf("GetSystemHealth Error: %v\n", err)
	} else {
		log.Printf("Health: Voltage=%s, Temp=%s\n", health.Voltage, health.Temperature)
	}

	identity, err := client.GetSystemIdentity(ctx, router)
	if err != nil {
		log.Printf("GetSystemIdentity Error: %v\n", err)
	} else {
		log.Printf("Identity: %s\n", identity.Name)
	}

	rb, err := client.GetRouterBoardInfo(ctx, router)
	if err != nil {
		log.Printf("GetRouterBoardInfo Error: %v\n", err)
	} else {
		log.Printf("RouterBoard: Model=%s, Serial=%s\n", rb.Model, rb.SerialNumber)
	}

	clock, err := client.GetSystemClock(ctx, router)
	if err != nil {
		log.Printf("GetSystemClock Error: %v\n", err)
	} else {
		log.Printf("Clock: %s %s\n", clock.Time, clock.Date)
	}
}

func testInterfaces(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- Network Interfaces ---")

	interfaces, err := client.GetInterfaces(ctx, router)
	if err != nil {
		log.Printf("GetInterfaces Error: %v\n", err)
		return
	}
	log.Printf("Total Interfaces: %d\n", len(interfaces))
	if len(interfaces) > 0 {
		log.Printf("First Interface: Name=%s, Type=%s, Running=%v\n",
			interfaces[0].Name, interfaces[0].Type, interfaces[0].Running)
	}

	if len(interfaces) > 0 {
		resultChan := make(chan mikrotik.TrafficMonitorStats, 1)
		cancel, err := client.StartTrafficMonitorListen(ctx, router, interfaces[0].Name, resultChan)
		if err != nil {
			log.Printf("StartTrafficMonitorListen Error: %v\n", err)
		} else {
			select {
			case stats := <-resultChan:
				log.Printf("Traffic: TX=%d bps, RX=%d bps\n", stats.TxBitsPerSecond, stats.RxBitsPerSecond)
			case <-ctx.Done():
				return
			case <-time.After(3 * time.Second):
				return
			}
			cancel()
		}
	}
}

func testNAT(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- NAT Rules ---")

	natRules, err := client.GetNATRules(ctx, router)
	if err != nil {
		log.Printf("GetNATRules Error: %v\n", err)
		return
	}
	log.Printf("Total NAT Rules: %d\n", len(natRules))
	if len(natRules) > 0 {
		log.Printf("First NAT Rule: Chain=%s, Action=%s\n", natRules[0].Chain, natRules[0].Action)
	}
}

func testQueues(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- Simple Queues ---")

	queues, err := client.GetAllQueues(ctx, router)
	if err != nil {
		log.Printf("GetAllQueues Error: %v\n", err)
		return
	}
	log.Printf("Total Queues: %d\n", len(queues))
	if len(queues) > 0 {
		log.Printf("First 3 Queues: %v\n", min(queues, 3))
	}

	parentQueues, err := client.GetAllParentQueues(ctx, router)
	if err != nil {
		log.Printf("GetAllParentQueues Error: %v\n", err)
	} else {
		log.Printf("Parent Queues: %v\n", parentQueues)
	}
}

func testPools(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- IP Address Pools ---")

	pools, err := client.GetAddressPools(ctx, router)
	if err != nil {
		log.Printf("GetAddressPools Error: %v\n", err)
		return
	}
	log.Printf("Total Pools: %d\n", len(pools))
	if len(pools) > 0 {
		log.Printf("First Pool: %s\n", pools[0])
	}
}

func testLogs(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- Hotspot Logs ---")

	logs, err := client.GetHotspotLogs(ctx, router, 5)
	if err != nil {
		log.Printf("GetHotspotLogs Error: %v\n", err)
		return
	}
	log.Printf("Total Logs: %d\n", len(logs))
	if len(logs) > 0 {
		log.Printf("First Log: %s - %s\n", logs[0].Time, logs[0].Message)
	}
}

func testReports(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- Sales Reports ---")

	reports, err := client.GetSalesReports(ctx, router, "")
	if err != nil {
		log.Printf("GetSalesReports Error: %v\n", err)
		return
	}
	log.Printf("Total Reports: %d\n", len(reports))

	reportsByDay, err := client.GetSalesReportsByDay(ctx, router, "jan/26/2026")
	if err != nil {
		log.Printf("GetSalesReportsByDay Error: %v\n", err)
	} else {
		log.Printf("Reports by Day: %d\n", len(reportsByDay))
	}
}

func testPing(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- Ping (Streaming) ---")

	cfg := mikrotik.PingConfig{
		Address:  "8.8.8.8",
		Count:    3,
		Size:     64,
		Interval: time.Second,
	}

	resultChan := make(chan mikrotik.PingResult, 10)
	cancel, err := client.StartPingListen(ctx, router, cfg, resultChan)
	if err != nil {
		log.Printf("StartPingListen Error: %v\n", err)
		return
	}

	ctxPing, cancelPing := context.WithTimeout(ctx, 10*time.Second)
	defer cancelPing()
	defer cancel()

	pingCount := 0
	for pingCount < 3 {
		select {
		case <-ctxPing.Done():
			return
		case result, ok := <-resultChan:
			if !ok {
				return
			}
			log.Printf("Ping #%d: %s -> %.2fms, Received: %v\n",
				pingCount+1, result.Address, result.TimeMs, result.Received)
			pingCount++
		}
	}
}

func testExpireMonitor(ctx context.Context, client *mikrotik.Client, router *entity.Router) {
	fmt.Println("--- Expire Monitor ---")

	script := mikrotik.NewOnLoginGenerator().GenerateExpireMonitorScript()
	log.Printf("Expire Monitor Script length: %d\n", len(script))

	status, err := client.EnsureExpireMonitor(ctx, router, script)
	if err != nil {
		log.Printf("EnsureExpireMonitor Error: %v\n", err)
	} else {
		log.Printf("Expire Monitor Status: %s\n", status)
	}
}

func testVoucherGenerator(ctx context.Context) {
	fmt.Println("--- Voucher Generator ---")

	generator := mikrotik.NewVoucherGenerator()

	req := &dto.VoucherGenerateRequest{
		Profile:      "default",
		Quantity:     5,
		Server:       "all",
		Mode:         "vc",
		NameLength:   8,
		CharacterSet: "lower",
		Prefix:       "V",
		TimeLimit:    "1h",
		DataLimit:    "500M",
		Comment:      "Test",
	}

	vouchers := generator.GenerateBatch(req)
	log.Printf("Generated %d vouchers\n", len(vouchers))

	if len(vouchers) > 0 {
		log.Printf("First Voucher: %s / %s\n", vouchers[0].Username, vouchers[0].Password)
	}
}

func testOnLoginGenerator(ctx context.Context) {
	fmt.Println("--- On-Login Script Generator ---")

	generator := mikrotik.NewOnLoginGenerator()

	req := &dto.ProfileRequest{
		Name:         "Premium",
		ExpireMode:   "remc",
		Validity:     "30d",
		Price:        5000,
		SellingPrice: 5500,
		LockUser:     "Enable",
		LockServer:   "Disable",
	}

	script := generator.Generate(req)
	log.Printf("On-Login Script generated, length: %d\n", len(script))

	parsed := generator.Parse(script)
	log.Printf("Parsed: ExpireMode=%s, Price=%.0f, Validity=%s\n",
		parsed.ExpireMode, parsed.Price, parsed.Validity)
}

func min(slice []string, n int) []string {
	if n > len(slice) {
		n = len(slice)
	}
	return slice[:n]
}
