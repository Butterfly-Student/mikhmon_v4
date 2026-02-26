package service

import (
	"context"

	"github.com/go-routeros/routeros/v3"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/entity"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
)

// MikrotikService defines the interface for MikroTik API operations
// Semua data diambil real-time dari RouterOS
type MikrotikService interface {
	// Connection Management
	Connect(ctx context.Context, router *entity.Router) (*routeros.Client, error)
	Disconnect(client *routeros.Client)
	TestConnection(ctx context.Context, router *entity.Router) error

	// System Information
	GetSystemResource(ctx context.Context, client *routeros.Client) (*dto.SystemResource, error)
	GetSystemHealth(ctx context.Context, client *routeros.Client) (*dto.SystemHealth, error)
	GetSystemIdentity(ctx context.Context, client *routeros.Client) (*dto.SystemIdentity, error)
	GetSystemClock(ctx context.Context, client *routeros.Client) (*dto.SystemClock, error)
	GetRouterBoardInfo(ctx context.Context, client *routeros.Client) (*dto.RouterBoardInfo, error)

	// Interfaces & Traffic
	GetInterfaces(ctx context.Context, client *routeros.Client) ([]*dto.Interface, error)
	MonitorTraffic(ctx context.Context, client *routeros.Client, iface string) (*dto.TrafficStats, error)

	// Real-time Streaming Monitoring ( menggunakan ListenArgsContext )
	StartQueueStatsListen(ctx context.Context, router *entity.Router, cfg mikrotik.QueueStatsConfig, resultChan chan<- mikrotik.QueueStats) (func() error, error)
	StartTrafficMonitorListen(ctx context.Context, router *entity.Router, interfaceName string, resultChan chan<- mikrotik.TrafficMonitorStats) (func() error, error)
	StartSystemResourceMonitorListen(ctx context.Context, router *entity.Router, resultChan chan<- mikrotik.SystemResourceMonitorStats) (func() error, error)

	// Hotspot Users
	GetHotspotUsers(ctx context.Context, client *routeros.Client, profile string) ([]*dto.HotspotUser, error)
	GetHotspotUsersByComment(ctx context.Context, client *routeros.Client, comment string) ([]*dto.HotspotUser, error)
	GetHotspotUserByID(ctx context.Context, client *routeros.Client, id string) (*dto.HotspotUser, error)
	GetHotspotUserByName(ctx context.Context, client *routeros.Client, name string) (*dto.HotspotUser, error)
	AddHotspotUser(ctx context.Context, client *routeros.Client, user *dto.HotspotUser) (string, error)
	UpdateHotspotUser(ctx context.Context, client *routeros.Client, id string, user *dto.HotspotUser) error
	RemoveHotspotUser(ctx context.Context, client *routeros.Client, id string) error
	RemoveHotspotUsersByComment(ctx context.Context, client *routeros.Client, comment string) error
	ResetHotspotUserCounters(ctx context.Context, client *routeros.Client, id string) error
	GetHotspotUsersCount(ctx context.Context, client *routeros.Client) (int, error)

	// Hotspot Active
	GetHotspotActive(ctx context.Context, client *routeros.Client) ([]*dto.HotspotActive, error)
	GetHotspotActiveCount(ctx context.Context, client *routeros.Client) (int, error)
	RemoveHotspotActive(ctx context.Context, client *routeros.Client, id string) error

	// Hotspot Hosts
	GetHotspotHosts(ctx context.Context, client *routeros.Client) ([]*dto.HotspotHost, error)
	RemoveHotspotHost(ctx context.Context, client *routeros.Client, id string) error

	// User Profiles (dengan on-login script!)
	GetUserProfiles(ctx context.Context, client *routeros.Client) ([]*dto.UserProfile, error)
	GetUserProfileByID(ctx context.Context, client *routeros.Client, id string) (*dto.UserProfile, error)
	GetUserProfileByName(ctx context.Context, client *routeros.Client, name string) (*dto.UserProfile, error)
	AddUserProfile(ctx context.Context, client *routeros.Client, profile *dto.UserProfile) (string, error)
	UpdateUserProfile(ctx context.Context, client *routeros.Client, id string, profile *dto.UserProfile) error
	RemoveUserProfile(ctx context.Context, client *routeros.Client, id string) error

	// Logging
	GetHotspotLogs(ctx context.Context, client *routeros.Client, limit int) ([]*dto.LogEntry, error)
	EnableHotspotLogging(ctx context.Context, client *routeros.Client) error

	// Sales Reports (via /system/script)
	GetSalesReports(ctx context.Context, client *routeros.Client, owner string) ([]*dto.SalesReport, error)
	GetSalesReportsByDay(ctx context.Context, client *routeros.Client, day string) ([]*dto.SalesReport, error)
	AddSalesReport(ctx context.Context, client *routeros.Client, report *dto.SalesReport) error

	// Helpers
	GetAddressPools(ctx context.Context, client *routeros.Client) ([]string, error)
	GetParentQueues(ctx context.Context, client *routeros.Client) ([]string, error)
	GetHotspotServers(ctx context.Context, client *routeros.Client) ([]string, error)
}

// OnLoginScriptGenerator defines the interface for generating on-login scripts
// INI ADALAH FITUR PALING PENTING DARI MIKHMON!
type OnLoginScriptGenerator interface {
	// Generate creates an on-login script for user profile
	// Script ini akan dieksekusi setiap kali user login
	Generate(req *dto.ProfileRequest) string
	
	// Parse extracts Mikhmon metadata dari existing on-login script
	Parse(script string) *dto.ProfileRequest
}
