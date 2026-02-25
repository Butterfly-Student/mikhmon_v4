package http

import (
	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/infrastructure/auth"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/middleware"
	"go.uber.org/zap"
)

// Router holds all HTTP handlers
type Router struct {
	engine           *gin.Engine
	authHandler      *handler.AuthHandler
	routerHandler    *handler.RouterHandler
	hotspotHandler   *handler.HotspotHandler
	voucherHandler   *handler.VoucherHandler
	reportHandler    *handler.ReportHandler
	dashboardHandler *handler.DashboardHandler
	pingWSHandler    *handler.PingWebSocketHandler
	jwtService       *auth.JWTService
}

// NewRouter creates a new HTTP router
func NewRouter(
	authHandler *handler.AuthHandler,
	routerHandler *handler.RouterHandler,
	hotspotHandler *handler.HotspotHandler,
	voucherHandler *handler.VoucherHandler,
	reportHandler *handler.ReportHandler,
	dashboardHandler *handler.DashboardHandler,
	pingWSHandler *handler.PingWebSocketHandler,
	jwtService *auth.JWTService,
	log *zap.Logger,
) *Router {
	r := &Router{
		engine:           gin.New(),
		authHandler:      authHandler,
		routerHandler:    routerHandler,
		hotspotHandler:   hotspotHandler,
		voucherHandler:   voucherHandler,
		reportHandler:    reportHandler,
		dashboardHandler: dashboardHandler,
		pingWSHandler:    pingWSHandler,
		jwtService:       jwtService,
	}

	r.setupMiddleware(log)
	r.setupRoutes()

	return r
}

// setupMiddleware configures middleware
func (r *Router) setupMiddleware(log *zap.Logger) {
	// CORS middleware — must be first
	r.engine.Use(middleware.CORS())

	// Recovery middleware
	r.engine.Use(gin.Recovery())

	// Zap structured request logger — replaces gin.Logger()
	r.engine.Use(middleware.ZapLogger(log.Named("http")))
}

// setupRoutes configures routes
func (r *Router) setupRoutes() {
	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success":   true,
			"code":      "OK",
			"timestamp": c.Request.Header.Get("X-Request-ID"),
			"data":      gin.H{"status": "ok"},
		})
	})

	// API v1
	v1 := r.engine.Group("/api/v1")
	{
		// Public routes
		r.authHandler.RegisterRoutes(v1)

		// Protected routes (HTTP API)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(r.jwtService))
		{
			r.routerHandler.RegisterRoutes(protected)
			r.hotspotHandler.RegisterRoutes(protected)
			r.voucherHandler.RegisterRoutes(protected)
			r.reportHandler.RegisterRoutes(protected)
			r.dashboardHandler.RegisterRoutes(protected)
		}
	}

	// WebSocket routes — handler has its own internal key auth
	r.engine.GET("/api/v1/ws/ping/:router_id", r.pingWSHandler.HandleWebSocket)
}

// GetEngine returns the gin engine
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
