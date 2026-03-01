package http

import (
	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/infrastructure/auth"
	httpHandler "github.com/irhabi89/mikhmon/internal/infrastructure/http/handler"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler/mikrotik"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler/mikrotik/ws"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/middleware"
	"go.uber.org/zap"
)

// Router holds all HTTP handlers
type Router struct {
	engine           *gin.Engine
	authHandler      *httpHandler.AuthHandler
	routerHandler    *httpHandler.RouterHandler
	mikrotikHandlers *MikrotikHandlers
	wsHandlers       *WSHandlers
	jwtService       *auth.JWTService
}

// MikrotikHandlers holds all MikroTik handlers
type MikrotikHandlers struct {
	Hotspot    *mikrotik.HotspotHandler
	Voucher    *mikrotik.VoucherHandler
	Report     *mikrotik.ReportHandler
	Interface  *mikrotik.InterfaceHandler
	System     *mikrotik.SystemHandler
	NAT        *mikrotik.NATHandler
	Queue      *mikrotik.QueueHandler
	Log        *mikrotik.LogHandler
	Pool       *mikrotik.PoolHandler
	PPPActive  *mikrotik.PPPActiveHandler
	PPPProfile *mikrotik.PPPProfileHandler
	PPPSecret  *mikrotik.PPPSecretHandler
}

// WSHandlers holds all WebSocket handlers
type WSHandlers struct {
	ResourceMonitor   *ws.ResourceMonitorHandler
	TrafficMonitor    *ws.TrafficMonitorHandler
	QueueMonitor      *ws.QueueMonitorHandler
	Ping              *ws.PingHandler
	LogMonitor        *ws.LogMonitorHandler
	HotspotLogMonitor *ws.LogMonitorHandler
	PPPLogMonitor     *ws.LogMonitorHandler
	PPPActiveMonitor        *ws.PPPActiveMonitorHandler
	PPPInactiveMonitor      *ws.PPPInactiveMonitorHandler
	HotspotActiveMonitor    *ws.HotspotActiveMonitorHandler
	HotspotInactiveMonitor  *ws.HotspotInactiveMonitorHandler
}

// NewRouter creates a new HTTP router
func NewRouter(
	authHandler *httpHandler.AuthHandler,
	routerHandler *httpHandler.RouterHandler,
	mikrotikHandlers *MikrotikHandlers,
	wsHandlers *WSHandlers,
	jwtService *auth.JWTService,
	log *zap.Logger,
) *Router {
	r := &Router{
		engine:           gin.New(),
		authHandler:      authHandler,
		routerHandler:    routerHandler,
		mikrotikHandlers: mikrotikHandlers,
		wsHandlers:       wsHandlers,
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
			// Router management
			r.routerHandler.RegisterRoutes(protected)

			// MikroTik routes
			mikrotik := protected.Group("/mikrotik/:router_id")
			{
				r.mikrotikHandlers.Hotspot.RegisterRoutes(mikrotik)
				r.mikrotikHandlers.Voucher.RegisterRoutes(mikrotik)
				r.mikrotikHandlers.Report.RegisterRoutes(mikrotik)
				r.mikrotikHandlers.Interface.RegisterRoutes(mikrotik)
				r.mikrotikHandlers.System.RegisterRoutes(mikrotik)
				r.mikrotikHandlers.NAT.RegisterRoutes(mikrotik)
				r.mikrotikHandlers.Queue.RegisterRoutes(mikrotik)
				r.mikrotikHandlers.Log.RegisterRoutes(mikrotik)
				r.mikrotikHandlers.Pool.RegisterRoutes(mikrotik)
				r.mikrotikHandlers.PPPActive.RegisterRoutes(mikrotik)
				r.mikrotikHandlers.PPPProfile.RegisterRoutes(mikrotik)
				r.mikrotikHandlers.PPPSecret.RegisterRoutes(mikrotik)
			}
		}
	}

	// WebSocket routes
	ws := r.engine.Group("/api/v1/ws/mikrotik/monitor")
	{
		ws.GET("/resource/:router_id",     r.wsHandlers.ResourceMonitor.Handle)
		ws.GET("/interface/:router_id",    r.wsHandlers.TrafficMonitor.Handle)
		ws.GET("/queue/:router_id",        r.wsHandlers.QueueMonitor.Handle)
		ws.GET("/ping/:router_id",         r.wsHandlers.Ping.Handle)
		ws.GET("/logs/:router_id",         r.wsHandlers.LogMonitor.Handle)
		ws.GET("/hotspot-logs/:router_id", r.wsHandlers.HotspotLogMonitor.Handle)
		ws.GET("/ppp-logs/:router_id",     r.wsHandlers.PPPLogMonitor.Handle)
		ws.GET("/ppp-active/:router_id",       r.wsHandlers.PPPActiveMonitor.Handle)
		ws.GET("/ppp-inactive/:router_id",     r.wsHandlers.PPPInactiveMonitor.Handle)
		ws.GET("/hotspot-active/:router_id",   r.wsHandlers.HotspotActiveMonitor.Handle)
		ws.GET("/hotspot-inactive/:router_id", r.wsHandlers.HotspotInactiveMonitor.Handle)
	}
}

// GetEngine returns the gin engine
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
