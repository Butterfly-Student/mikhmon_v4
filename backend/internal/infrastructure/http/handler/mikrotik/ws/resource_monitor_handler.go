package ws

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/irhabi89/mikhmon/internal/domain/entity"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"go.uber.org/zap"
)

// ResourceMonitorResult is sent to frontend
type ResourceMonitorResult struct {
	CPUUsed    int   `json:"cpuUsed"`
	FreeMemory int64 `json:"freeMemory"`
	Timestamp  int64 `json:"timestamp"`
}

// ResourceMonitorHandler handles WebSocket for /system/resource/monitor
type ResourceMonitorHandler struct {
	upgrader      websocket.Upgrader
	routerRepo    repository.RouterRepository
	mikrotikSvc   *mikrotik.Client
	internalWSKey string
	log           *zap.Logger
}

// NewResourceMonitorHandler creates a new WebSocket handler for resource monitoring
func NewResourceMonitorHandler(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Client,
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

	var forwardingCancel context.CancelFunc
	stopForwarding := func() {
		if forwardingCancel != nil {
			forwardingCancel()
		}
	}
	defer stopForwarding()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			h.log.Info("WebSocket connection closed",
				zap.Uint64("routerID", routerID),
				zap.Error(err),
			)
			stopForwarding()
			return
		}

		var cmd struct {
			Action string `json:"action"`
		}
		if err := conn.ReadJSON(message); err != nil {
			conn.WriteJSON(gin.H{"type": "error", "message": "invalid command"})
			continue
		}

		switch cmd.Action {
		case "start":
			stopForwarding()
			forwardingCancel = h.startResourceMonitor(ctx, conn, router)
		case "stop":
			stopForwarding()
			conn.WriteJSON(gin.H{"type": "status", "status": "stopped"})
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

		resultChan := make(chan mikrotik.SystemResourceMonitorStats, 10)
		cancelFn, err := h.mikrotikSvc.StartSystemResourceMonitorListen(ctx, router, resultChan)
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
					CPUUsed:    result.CPUUsed,
					FreeMemory: result.FreeMemory,
					Timestamp:  result.Timestamp.UnixMilli(),
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
