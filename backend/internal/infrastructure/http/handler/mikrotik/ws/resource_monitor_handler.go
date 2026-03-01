package ws

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/entity"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"go.uber.org/zap"
)

// ResourceMonitorResult is sent to frontend.
// Semua field diperbarui setiap detik dari /system/resource/print interval=1s.
type ResourceMonitorResult struct {
	Uptime               string  `json:"uptime"`
	Version              string  `json:"version"`
	BuildTime            string  `json:"buildTime"`
	FreeMemory           int64   `json:"freeMemory"`
	TotalMemory          int64   `json:"totalMemory"`
	CPU                  string  `json:"cpu"`
	CPUCount             int     `json:"cpuCount"`
	CPUFrequency         int     `json:"cpuFrequency"`
	CPULoad              int     `json:"cpuLoad"`
	FreeHddSpace         int64   `json:"freeHddSpace"`
	TotalHddSpace        int64   `json:"totalHddSpace"`
	WriteSectSinceReboot int64   `json:"writeSectSinceReboot"`
	WriteSectTotal       int64   `json:"writeSectTotal"`
	BadBlocks            float64 `json:"badBlocks"`
	ArchitectureName     string  `json:"architectureName"`
	BoardName            string  `json:"boardName"`
	Platform             string  `json:"platform"`
}

// ResourceMonitorHandler handles WebSocket for /system/resource/monitor
type ResourceMonitorHandler struct {
	upgrader      websocket.Upgrader
	routerRepo    repository.RouterRepository
	mikrotikSvc   *mikrotik.Manager
	internalWSKey string
	log           *zap.Logger
}

// NewResourceMonitorHandler creates a new WebSocket handler for resource monitoring
func NewResourceMonitorHandler(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
	internalWSKey string,
	log *zap.Logger,
) *ResourceMonitorHandler {
	return &ResourceMonitorHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		routerRepo:    routerRepo,
		mikrotikSvc:   mikrotikSvc,
		internalWSKey: internalWSKey,
		log:           log.Named("ws-resource-monitor"),
	}
}

// Handle handles WebSocket connections
func (h *ResourceMonitorHandler) Handle(c *gin.Context) {
	providedKey := c.Query("key")

	if providedKey != h.internalWSKey {
		h.log.Warn("WebSocket auth failed", zap.String("providedKey", providedKey))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	routerID, err := strconv.ParseUint(c.Param("router_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid router_id"})
		return
	}

	h.log.Info("Resource monitor WebSocket connection attempt",
		zap.Uint64("routerID", routerID),
	)

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("WebSocket upgrade failed",
			zap.Uint64("routerID", routerID),
			zap.Error(err),
		)
		return
	}
	defer conn.Close()

	h.log.Info("Resource monitor WebSocket upgraded",
		zap.Uint64("routerID", routerID),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn.SetPongHandler(func(data string) error {
		return conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := conn.WriteControl(
					websocket.PingMessage,
					[]byte{},
					time.Now().Add(5*time.Second),
				); err != nil {
					h.log.Warn("WebSocket keepalive ping failed",
						zap.Uint64("routerID", routerID),
						zap.Error(err),
					)
					cancel()
					return
				}
			}
		}
	}()

	router, err := h.routerRepo.GetByID(ctx, uint(routerID))
	if err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "router not found"})
		return
	}

	// Auto-start resource monitor immediately — no command needed
	h.log.Info("Resource monitor WebSocket ready, starting monitor", zap.Uint64("routerID", routerID))

	forwardingCancel := h.startResourceMonitor(ctx, conn, router)
	defer forwardingCancel()

	// Read loop: only used to detect client disconnect
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			h.log.Info("Resource monitor WebSocket connection closed",
				zap.Uint64("routerID", routerID),
				zap.Error(err),
			)
			return
		}
	}
}

func (h *ResourceMonitorHandler) startResourceMonitor(
	parentCtx context.Context,
	conn *websocket.Conn,
	router *entity.Router,
) context.CancelFunc {
	ctx, cancel := context.WithCancel(parentCtx)

	go func() {
		defer cancel()

		h.log.Info("Starting resource monitor",
			zap.String("host", router.Host),
		)

		cfg := mikrotik.Config{
			Host:     router.Host,
			Port:     router.Port,
			Username: router.Username,
			Password: router.Password,
			UseTLS:   router.UseSSL,
			Timeout:  time.Duration(router.Timeout) * time.Second,
		}
		routerClient, err := h.mikrotikSvc.GetOrConnect(ctx, router.Name, cfg)
		if err != nil {
			h.log.Error("Router not connected", zap.String("name", router.Name), zap.Error(err))
			conn.WriteJSON(gin.H{"type": "error", "message": "router not connected: " + err.Error()})
			return
		}

		resultChan := make(chan dto.SystemResourceMonitorStats, 10)
		cancelFn, err := routerClient.StartSystemResourceMonitorListen(ctx, resultChan)
		if err != nil {
			h.log.Error("Failed to start resource monitor", zap.Error(err))
			conn.WriteJSON(gin.H{"type": "error", "message": "failed to start monitor: " + err.Error()})
			return
		}
		defer cancelFn()

		conn.WriteJSON(gin.H{"type": "status", "status": "started", "monitor": "resource"})

		for {
			select {
			case <-ctx.Done():
				return
			case result, ok := <-resultChan:
				if !ok {
					return
				}
				dto := ResourceMonitorResult{
					Uptime:               result.Uptime,
					Version:              result.Version,
					BuildTime:            result.BuildTime,
					FreeMemory:           result.FreeMemory,
					TotalMemory:          result.TotalMemory,
					CPU:                  result.CPU,
					CPUCount:             result.CPUCount,
					CPUFrequency:         result.CPUFrequency,
					CPULoad:              result.CPULoad,
					FreeHddSpace:         result.FreeHddSpace,
					TotalHddSpace:        result.TotalHddSpace,
					WriteSectSinceReboot: result.WriteSectSinceReboot,
					WriteSectTotal:       result.WriteSectTotal,
					BadBlocks:            result.BadBlocks,
					ArchitectureName:     result.ArchitectureName,
					BoardName:            result.BoardName,
					Platform:             result.Platform,
				}
				if err := conn.WriteJSON(dto); err != nil {
					h.log.Warn("Failed to write resource monitor result", zap.Error(err))
					return
				}
			}
		}
	}()

	return cancel
}
