package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"go.uber.org/zap"
)

// PingCommand from frontend
type PingCommand struct {
	Action   string  `json:"action"`   // "start", "stop", "ping"
	Address  string  `json:"address"`  // target address to ping
	Interval float64 `json:"interval"` // interval in seconds (default: 1)
	Count    int     `json:"count"`    // number of pings: 0 = infinite (default: 0)
	Size     int     `json:"size"`     // packet size in bytes (default: 64)
}

// PingWebSocketHandler handles WebSocket connections for real-time ping
type PingWebSocketHandler struct {
	upgrader      websocket.Upgrader
	routerRepo    repository.RouterRepository
	mikrotikSvc   *mikrotik.Client
	internalWSKey string // Internal key for simple auth
	log           *zap.Logger
}

// NewPingWebSocketHandler creates a new WebSocket handler for ping
func NewPingWebSocketHandler(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Client,
	internalWSKey string,
	log *zap.Logger,
) *PingWebSocketHandler {
	return &PingWebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow any origin
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		routerRepo:    routerRepo,
		mikrotikSvc:   mikrotikSvc,
		internalWSKey: internalWSKey,
		log:           log.Named("ping-ws"),
	}
}

// PingResultDTO is the data sent to frontend
type PingResultDTO struct {
	TimeMs   float64 `json:"timeMs"`
	Received bool    `json:"received"`
	Address  string  `json:"address"`
}

// HandleWebSocket upgrades HTTP to WebSocket and handles bidirectional communication
func (h *PingWebSocketHandler) HandleWebSocket(c *gin.Context) {
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
	h.log.Info("WebSocket connection attempt", zap.Uint64("routerID", routerID))

	// Upgrade to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("WebSocket upgrade failed",
			zap.Uint64("routerID", routerID),
			zap.Error(err),
		)
		return
	}
	defer conn.Close()
	h.log.Info("WebSocket upgraded", zap.Uint64("routerID", routerID))

	// Setup context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// WebSocket keepalive
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
		h.log.Error("Router not found", zap.Uint64("routerID", routerID), zap.Error(err))
		conn.WriteJSON(gin.H{"type": "error", "message": "router not found"})
		return
	}
	h.log.Info("Router found for ping",
		zap.Uint64("routerID", routerID),
		zap.String("name", router.Name),
		zap.String("host", router.Host),
	)

	var resultChan chan mikrotik.PingResult
	var pingCancel func() error
	// Initialize forwarding context with a no-op cancel so go vet is satisfied.
	// Real cancel is replaced each time a "start" command is received.
	forwardingCtx, forwardingCancel := context.WithCancel(ctx)
	_ = forwardingCtx

	stopForwarding := func() {
		if forwardingCancel != nil {
			forwardingCancel()
		}
		if pingCancel != nil {
			pingCancel()
		}
	}
	defer stopForwarding() // always clean up on function exit

	h.log.Info("WebSocket ready, waiting for commands", zap.Uint64("routerID", routerID))
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

		var cmd PingCommand
		if err := json.Unmarshal(message, &cmd); err != nil {
			conn.WriteJSON(gin.H{"type": "error", "message": "invalid command"})
			continue
		}
		h.log.Debug("WebSocket command received",
			zap.Uint64("routerID", routerID),
			zap.String("action", cmd.Action),
			zap.String("address", cmd.Address),
		)

		switch cmd.Action {
		case "start", "ping":
			stopForwarding()

			if cmd.Address == "" {
				conn.WriteJSON(gin.H{"type": "error", "message": "address required"})
				continue
			}

			h.log.Info("Starting ping",
				zap.Uint64("routerID", routerID),
				zap.String("address", cmd.Address),
				zap.Float64("interval", cmd.Interval),
				zap.Int("count", cmd.Count),
				zap.Int("size", cmd.Size),
			)
			resultChan = make(chan mikrotik.PingResult, 10)

			pingCfg := mikrotik.PingConfig{
				Address: cmd.Address,
				Count:   cmd.Count,
				Size:    cmd.Size,
			}
			if cmd.Interval > 0 {
				pingCfg.Interval = time.Duration(cmd.Interval * float64(time.Second))
			}

			cancelFn, err := h.mikrotikSvc.StartPingListen(ctx, router, pingCfg, resultChan)
			if err != nil {
				h.log.Error("Failed to start ping",
					zap.Uint64("routerID", routerID),
					zap.String("address", cmd.Address),
					zap.Error(err),
				)
				conn.WriteJSON(gin.H{"type": "error", "message": "failed to start ping: " + err.Error()})
				continue
			}
			pingCancel = cancelFn

			forwardingCtx, forwardingCancel = context.WithCancel(ctx)
			go func(pingCtx context.Context, ch chan mikrotik.PingResult) {
				h.log.Debug("Ping forwarding goroutine started",
					zap.Uint64("routerID", routerID),
					zap.String("address", cmd.Address),
				)
				for {
					select {
					case <-pingCtx.Done():
						return
					case result, ok := <-ch:
						if !ok {
							return
						}
						dto := PingResultDTO{
							TimeMs:   result.TimeMs,
							Received: result.Received,
							Address:  result.Address,
						}
						if err := conn.WriteJSON(dto); err != nil {
							h.log.Warn("Failed to write ping result to WebSocket",
								zap.Uint64("routerID", routerID),
								zap.Error(err),
							)
							return
						}
					}
				}
			}(forwardingCtx, resultChan)

			conn.WriteJSON(gin.H{"type": "status", "status": "started", "address": cmd.Address})

		case "stop":
			stopForwarding()
			resultChan = nil
			pingCancel = nil
			h.log.Info("Ping stopped", zap.Uint64("routerID", routerID))
			conn.WriteJSON(gin.H{"type": "status", "status": "stopped"})
		}
	}
}
