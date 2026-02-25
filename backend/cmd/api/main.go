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
	"github.com/irhabi89/mikhmon/internal/infrastructure/logger"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"github.com/irhabi89/mikhmon/internal/infrastructure/repository/postgres"
	"github.com/irhabi89/mikhmon/internal/usecase"
	"go.uber.org/zap"
)

func main() {
	// --- 1. Init Logger (must be first) ---
	logger.FromEnv()
	log := logger.Log
	defer log.Sync() //nolint:errcheck

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
	mikrotikClient := mikrotik.NewClient(log)
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL)

	// --- 7. Initialize infrastructure services ---
	hotspotService := mikrotik.NewHotspotService(mikrotikClient, routerRepo)

	// --- 8. Initialize use cases ---
	authUC := usecase.NewAuthUseCase(adminRepo, jwtService)
	routerUC := usecase.NewRouterUseCase(routerRepo, mikrotikClient)
	hotspotUC := usecase.NewHotspotUseCase(hotspotService)
	voucherUC := usecase.NewVoucherUseCase(routerRepo, hotspotService, nil)
	reportUC := usecase.NewReportUseCase(routerRepo, mikrotikClient)
	dashboardUC := usecase.NewDashboardUseCase(routerRepo, mikrotikClient, log)

	// --- 9. Initialize handlers (inject logger) ---
	authHandler := handler.NewAuthHandler(authUC, log)
	routerHandler := handler.NewRouterHandler(routerUC, log)
	hotspotHandler := handler.NewHotspotHandler(hotspotUC, log)
	voucherHandler := handler.NewVoucherHandler(voucherUC, log)
	reportHandler := handler.NewReportHandler(reportUC, log)
	dashboardHandler := handler.NewDashboardHandler(dashboardUC, log)
	pingWSHandler := handler.NewPingWebSocketHandler(routerRepo, mikrotikClient, cfg.InternalWSKey, log)

	// --- 10. Setup HTTP router ---
	router := httpInfra.NewRouter(
		authHandler,
		routerHandler,
		hotspotHandler,
		voucherHandler,
		reportHandler,
		dashboardHandler,
		pingWSHandler,
		jwtService,
		log,
	)

	// --- 11. Create HTTP server ---
	// ReadTimeout/WriteTimeout di-set 0 agar WebSocket tidak di-terminasi.
	// WebSocket keepalive dikelola sendiri oleh gorilla/websocket.
	srv := &http.Server{
		Addr:        fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:     router.GetEngine(),
		ReadTimeout: 0, // Disable: required for WebSocket
		WriteTimeout: 0, // Disable: required for WebSocket
		IdleTimeout: cfg.Server.IdleTimeout,
	}

	// --- 12. Start server ---
	go func() {
		log.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// --- 13. Graceful shutdown ---
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
