package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/infrastructure/auth"
	"github.com/irhabi89/mikhmon/internal/infrastructure/config"
	"github.com/irhabi89/mikhmon/internal/infrastructure/database"
	httpInfra "github.com/irhabi89/mikhmon/internal/infrastructure/http"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler/mikrotik"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler/mikrotik/ws"
	"github.com/irhabi89/mikhmon/internal/infrastructure/logger"
	mikrotikInfra "github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"github.com/irhabi89/mikhmon/internal/infrastructure/repository/postgres"
	"github.com/irhabi89/mikhmon/internal/usecase"
	mikrotikUsecase "github.com/irhabi89/mikhmon/internal/usecase/mikrotik"
	"go.uber.org/zap"
)

func main() {
	// --- 1. Init Logger (must be first) ---
	logger.FromEnv()
	log := logger.Log
	defer log.Sync()

	// --- 2. Load configuration ---
	cfg, err := config.Load("")
	if err != nil {
		log.Fatal("Failed to load config", zap.Error(err))
	}

	// --- 3. Set gin mode ---
	if cfg.Server.Port == "80" || cfg.Server.Port == "443" {
		gin.SetMode(gin.ReleaseMode)
	}

	// --- 4. Initialize database ---
	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database", zap.Error(err))
	}

	if err := database.Seed(db); err != nil {
		log.Fatal("Failed to seed database", zap.Error(err))
	}

	// --- 5. Initialize repositories ---
	adminRepo := postgres.NewAdminUserRepository(db)
	routerRepo := postgres.NewRouterRepository(db)

	// --- 6. Initialize services ---
	mikrotikClient := mikrotikInfra.NewClient(log)
	hotspotService := mikrotikInfra.NewHotspotService(mikrotikClient, routerRepo)
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL)

	// --- 7. Initialize use cases ---
	authUC := usecase.NewAuthUseCase(adminRepo, jwtService)
	routerUC := usecase.NewRouterUseCase(routerRepo, mikrotikClient)

	// MikroTik use cases
	hotspotUC := mikrotikUsecase.NewHotspotUseCase(routerRepo, hotspotService, mikrotikClient, log)
	voucherUC := mikrotikUsecase.NewVoucherUseCase(routerRepo, hotspotService, nil, log)
	reportUC := mikrotikUsecase.NewReportUseCase(routerRepo, mikrotikClient, log)
	interfaceUC := mikrotikUsecase.NewInterfaceUseCase(routerRepo, mikrotikClient, log)
	systemUC := mikrotikUsecase.NewSystemUseCase(routerRepo, mikrotikClient, log)
	natUC := mikrotikUsecase.NewNATUseCase(routerRepo, mikrotikClient, log)
	queueUC := mikrotikUsecase.NewQueueUseCase(routerRepo, mikrotikClient, log)
	logUC := mikrotikUsecase.NewLogUseCase(routerRepo, mikrotikClient, log)
	poolUC := mikrotikUsecase.NewPoolUseCase(routerRepo, mikrotikClient, log)

	// --- 8. Initialize handlers (inject logger) ---
	authHandler := handler.NewAuthHandler(authUC, log)
	routerHandler := handler.NewRouterHandler(routerUC, log)

	// MikroTik handlers
	hotspotHandler := mikrotik.NewHotspotHandler(hotspotUC, log)
	voucherHandler := mikrotik.NewVoucherHandler(voucherUC, log)
	reportHandler := mikrotik.NewReportHandler(reportUC, log)
	interfaceHandler := mikrotik.NewInterfaceHandler(interfaceUC, log)
	systemHandler := mikrotik.NewSystemHandler(systemUC, log)
	natHandler := mikrotik.NewNATHandler(natUC, log)
	queueHandler := mikrotik.NewQueueHandler(queueUC, log)
	logHandler := mikrotik.NewLogHandler(logUC, log)
	poolHandler := mikrotik.NewPoolHandler(poolUC, log)

	// WebSocket handlers
	resourceMonitorHandler := ws.NewResourceMonitorHandler(routerRepo, mikrotikClient, cfg.InternalWSKey, log)
	trafficMonitorHandler := ws.NewTrafficMonitorHandler(routerRepo, mikrotikClient, cfg.InternalWSKey, log)
	queueMonitorHandler := ws.NewQueueMonitorHandler(routerRepo, mikrotikClient, cfg.InternalWSKey, log)
	pingHandler := ws.NewPingHandler(routerRepo, mikrotikClient, cfg.InternalWSKey, log)

	// --- 9. Setup HTTP router ---
	router := httpInfra.NewRouter(
		authHandler,
		routerHandler,
		&httpInfra.MikrotikHandlers{
			Hotspot:   hotspotHandler,
			Voucher:   voucherHandler,
			Report:    reportHandler,
			Interface: interfaceHandler,
			System:    systemHandler,
			NAT:       natHandler,
			Queue:     queueHandler,
			Log:       logHandler,
			Pool:      poolHandler,
		},
		&httpInfra.WSHandlers{
			ResourceMonitor: resourceMonitorHandler,
			TrafficMonitor:  trafficMonitorHandler,
			QueueMonitor:    queueMonitorHandler,
			Ping:            pingHandler,
		},
		jwtService,
		log,
	)

	// --- 10. Create HTTP server ---
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router.GetEngine(),
		ReadTimeout:  0,
		WriteTimeout: 0,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// --- 11. Start server ---
	go func() {
		log.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// --- 12. Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	// Close MikroTik connections
	mikrotikClient.CloseAll()

	// Close database connection
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Info("Server exited cleanly")
}
