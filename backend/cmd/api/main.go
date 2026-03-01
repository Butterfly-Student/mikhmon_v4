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
	cache "github.com/irhabi89/mikhmon/internal/infrastructure/cache"
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
	"github.com/irhabi89/mikhmon/pkg/pubsub"
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

	// --- 5. Initialize Redis (optional — pub-sub disabled if unavailable) ---
	var ps *pubsub.PubSub
	if redisClient, err := cache.NewRedis(cfg.Redis); err != nil {
		log.Warn("Redis unavailable, pub-sub disabled", zap.Error(err))
	} else {
		ps = pubsub.New(redisClient.Client())
	}

	// --- 6. Initialize repositories ---
	adminRepo := postgres.NewAdminUserRepository(db)
	routerRepo := postgres.NewRouterRepository(db)

	// --- 6. Initialize services ---
	mikrotikManager := mikrotikInfra.NewManager(log)
	hotspotService := mikrotikInfra.NewHotspotService(mikrotikManager, routerRepo)
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL)

	// --- 7. Initialize use cases ---
	authUC := usecase.NewAuthUseCase(adminRepo, jwtService)
	routerUC := usecase.NewRouterUseCase(routerRepo, mikrotikManager)

	// MikroTik use cases
	hotspotUC := mikrotikUsecase.NewHotspotUseCase(routerRepo, hotspotService, log)
	voucherUC := mikrotikUsecase.NewVoucherUseCase(routerRepo, hotspotService, nil, log)
	reportUC := mikrotikUsecase.NewReportUseCase(routerRepo, mikrotikManager, log)
	interfaceUC := mikrotikUsecase.NewInterfaceUseCase(routerRepo, mikrotikManager, log)
	systemUC := mikrotikUsecase.NewSystemUseCase(routerRepo, mikrotikManager, log)
	natUC := mikrotikUsecase.NewNATUseCase(routerRepo, mikrotikManager, log)
	queueUC := mikrotikUsecase.NewQueueUseCase(routerRepo, mikrotikManager, log)
	logUC := mikrotikUsecase.NewLogUseCase(routerRepo, mikrotikManager, log)
	poolUC := mikrotikUsecase.NewPoolUseCase(routerRepo, mikrotikManager, log)
	pppActiveUC := mikrotikUsecase.NewPPPActiveUseCase(routerRepo, mikrotikManager, log)
	pppProfileUC := mikrotikUsecase.NewPPPProfileUseCase(routerRepo, mikrotikManager, log)
	pppSecretUC := mikrotikUsecase.NewPPPSecretUseCase(routerRepo, mikrotikManager, log)

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
	pppActiveHandler := mikrotik.NewPPPActiveHandler(pppActiveUC, log)
	pppProfileHandler := mikrotik.NewPPPProfileHandler(pppProfileUC, log)
	pppSecretHandler := mikrotik.NewPPPSecretHandler(pppSecretUC, log)

	// WebSocket handlers
	resourceMonitorHandler := ws.NewResourceMonitorHandler(routerRepo, mikrotikManager, cfg.InternalWSKey, log)
	trafficMonitorHandler := ws.NewTrafficMonitorHandler(routerRepo, mikrotikManager, cfg.InternalWSKey, log)
	queueMonitorHandler := ws.NewQueueMonitorHandler(routerRepo, mikrotikManager, cfg.InternalWSKey, log)
	pingHandler := ws.NewPingHandler(routerRepo, mikrotikManager, cfg.InternalWSKey, log)
	logMonitorHandler := ws.NewLogMonitorHandler(routerRepo, mikrotikManager, ps, cfg.InternalWSKey, log)
	hotspotLogMonitorHandler := ws.NewHotspotLogMonitorHandler(routerRepo, mikrotikManager, ps, cfg.InternalWSKey, log)
	pppLogMonitorHandler := ws.NewPPPLogMonitorHandler(routerRepo, mikrotikManager, ps, cfg.InternalWSKey, log)
	pppActiveMonitorHandler := ws.NewPPPActiveMonitorHandler(routerRepo, mikrotikManager, cfg.InternalWSKey, log)
	pppInactiveMonitorHandler := ws.NewPPPInactiveMonitorHandler(routerRepo, mikrotikManager, cfg.InternalWSKey, log)
	hotspotActiveMonitorHandler := ws.NewHotspotActiveMonitorHandler(routerRepo, mikrotikManager, cfg.InternalWSKey, log)
	hotspotInactiveMonitorHandler := ws.NewHotspotInactiveMonitorHandler(routerRepo, mikrotikManager, cfg.InternalWSKey, log)

	// --- 9. Setup HTTP router ---
	router := httpInfra.NewRouter(
		authHandler,
		routerHandler,
		&httpInfra.MikrotikHandlers{
			Hotspot:    hotspotHandler,
			Voucher:    voucherHandler,
			Report:     reportHandler,
			Interface:  interfaceHandler,
			System:     systemHandler,
			NAT:        natHandler,
			Queue:      queueHandler,
			Log:        logHandler,
			Pool:       poolHandler,
			PPPActive:  pppActiveHandler,
			PPPProfile: pppProfileHandler,
			PPPSecret:  pppSecretHandler,
		},
		&httpInfra.WSHandlers{
			ResourceMonitor:   resourceMonitorHandler,
			TrafficMonitor:    trafficMonitorHandler,
			QueueMonitor:      queueMonitorHandler,
			Ping:              pingHandler,
			LogMonitor:        logMonitorHandler,
			HotspotLogMonitor: hotspotLogMonitorHandler,
			PPPLogMonitor:     pppLogMonitorHandler,
			PPPActiveMonitor:       pppActiveMonitorHandler,
			PPPInactiveMonitor:     pppInactiveMonitorHandler,
			HotspotActiveMonitor:   hotspotActiveMonitorHandler,
			HotspotInactiveMonitor: hotspotInactiveMonitorHandler,
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
	mikrotikManager.CloseAll()

	// Close database connection
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Info("Server exited cleanly")
}
