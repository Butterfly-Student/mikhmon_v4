package ws

import (
	"context"
	"encoding/json"
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

// QueueMonitorResult is sent to frontend
type QueueMonitorResult struct {
	Name      string `json:"name"`
	BytesIn   int64  `json:"bytesIn"`
	BytesOut  int64  `json:"bytesOut"`
	RateIn    int64  `json:"rateIn"`
	RateOut   int64  `json:"rateOut"`
	Timestamp int64  `json:"timestamp"`
}

// QueueMonitorHandler handles WebSocket for /queue/simple/print stats
type QueueMonitorHandler struct {
	upgrader      websocket.Upgrader
	routerRepo    repository.RouterRepository
	mikrotikSvc   *mikrotik.Client
	internalWSKey string
	log           *zap.Logger
}

// NewQueueMonitorHandler creates a new WebSocket handler for queue monitoring
func NewQueueMonitorHandler(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Client,
	internalWSKey string,
	log *zap.Logger,
) *QueueMonitorHandler {
	return &QueueMonitorHandler{
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
		log:           log.Named("ws-queue-monitor"),
	}
}

// Handle handles WebSocket connections
func (h *QueueMonitorHandler) Handle(c *gin.Context) {
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

	h.log.Info("Queue monitor WebSocket connection attempt",
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

	h.log.Info("Queue monitor WebSocket upgraded",
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

		var cmd MonitorCommand
		if err := json.Unmarshal(message, &cmd); err != nil {
			conn.WriteJSON(gin.H{"type": "error", "message": "invalid command"})
			continue
		}

		h.log.Debug("Monitor command received",
			zap.Uint64("routerID", routerID),
			zap.String("action", cmd.Action),
		)

		switch cmd.Action {
		case "start":
			stopForwarding()
			forwardingCancel = h.startQueueMonitor(ctx, conn, router, cmd)
		case "stop":
			stopForwarding()
			conn.WriteJSON(gin.H{"type": "status", "status": "stopped"})
		}
	}
}

func (h *QueueMonitorHandler) startQueueMonitor(
	parentCtx context.Context,
	conn *websocket.Conn,
	router *entity.Router,
	cmd MonitorCommand,
) context.CancelFunc {
	ctx, cancel := context.WithCancel(parentCtx)

	if cmd.Name == "" {
		conn.WriteJSON(gin.H{"type": "error", "message": "queue name is required"})
		return func() {}
	}

	go func() {
		defer cancel()

		h.log.Info("Starting queue monitor",
			zap.String("host", router.Host),
			zap.String("queue", cmd.Name),
		)

		cfg := mikrotik.DefaultQueueStatsConfig(cmd.Name)
		resultChan := make(chan mikrotik.QueueStats, 10)
		cancelFn, err := h.mikrotikSvc.StartQueueStatsListen(ctx, router, cfg, resultChan)
		if err != nil {
			h.log.Error("Failed to start queue monitor", zap.Error(err))
			conn.WriteJSON(gin.H{"type": "error", "message": "failed to start monitor: " + err.Error()})
			return
		}
		defer cancelFn()

		conn.WriteJSON(gin.H{"type": "status", "status": "started", "monitor": "queue", "name": cmd.Name})

		for {
			select {
			case <-ctx.Done():
				return
			case result, ok := <-resultChan:
				if !ok {
					return
				}
				dto := QueueMonitorResult{
					Name:      result.Name,
					BytesIn:   result.BytesIn,
					BytesOut:  result.BytesOut,
					RateIn:    result.RateIn,
					RateOut:   result.RateOut,
					Timestamp: result.Timestamp.UnixMilli(),
				}
				if err := conn.WriteJSON(dto); err != nil {
					h.log.Warn("Failed to write queue monitor result", zap.Error(err))
					return
				}
			}
		}
	}()

	return cancel
}
