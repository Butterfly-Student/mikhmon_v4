package ws

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"github.com/irhabi89/mikhmon/pkg/pubsub"
	"go.uber.org/zap"
)

// LogMonitorHandler handles WebSocket log streaming for all, hotspot, or PPP logs.
// The topics field is set at construction time:
//
//	"" = all logs
//	"hotspot,info" = hotspot logs
//	"ppp,pppoe,info" = PPP logs
type LogMonitorHandler struct {
	upgrader      websocket.Upgrader
	routerRepo    repository.RouterRepository
	mikrotikSvc   *mikrotik.Manager
	ps            *pubsub.PubSub // nil if Redis unavailable
	topics        string
	internalWSKey string
	log           *zap.Logger
}

func newLogMonitorHandler(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
	ps *pubsub.PubSub,
	topics string,
	internalWSKey string,
	log *zap.Logger,
	name string,
) *LogMonitorHandler {
	return &LogMonitorHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin:     func(r *http.Request) bool { return true },
			ReadBufferSize:  1024,
			WriteBufferSize: 4096,
		},
		routerRepo:    routerRepo,
		mikrotikSvc:   mikrotikSvc,
		ps:            ps,
		topics:        topics,
		internalWSKey: internalWSKey,
		log:           log.Named(name),
	}
}

// NewLogMonitorHandler streams all logs (topics="").
func NewLogMonitorHandler(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
	ps *pubsub.PubSub,
	internalWSKey string,
	log *zap.Logger,
) *LogMonitorHandler {
	return newLogMonitorHandler(routerRepo, mikrotikSvc, ps, "", internalWSKey, log, "ws-log-monitor")
}

// NewHotspotLogMonitorHandler streams hotspot logs (topics="hotspot,info").
func NewHotspotLogMonitorHandler(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
	ps *pubsub.PubSub,
	internalWSKey string,
	log *zap.Logger,
) *LogMonitorHandler {
	return newLogMonitorHandler(routerRepo, mikrotikSvc, ps, "hotspot,info", internalWSKey, log, "ws-hotspot-log-monitor")
}

// NewPPPLogMonitorHandler streams PPP logs (topics="ppp,pppoe,info").
func NewPPPLogMonitorHandler(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
	ps *pubsub.PubSub,
	internalWSKey string,
	log *zap.Logger,
) *LogMonitorHandler {
	return newLogMonitorHandler(routerRepo, mikrotikSvc, ps, "ppp,pppoe,info", internalWSKey, log, "ws-ppp-log-monitor")
}

// Handle handles WebSocket connections for log streaming.
func (h *LogMonitorHandler) Handle(c *gin.Context) {
	if c.Query("key") != h.internalWSKey {
		h.log.Warn("WebSocket auth failed", zap.String("key", c.Query("key")))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	routerID, err := strconv.ParseUint(c.Param("router_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid router_id"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("WebSocket upgrade failed", zap.Uint64("routerID", routerID), zap.Error(err))
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Keepalive: ping every 20s, expect pong within 60s.
	conn.SetPongHandler(func(string) error {
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
				if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second)); err != nil {
					cancel()
					return
				}
			}
		}
	}()

	router, err := h.routerRepo.GetByID(ctx, uint(routerID))
	if err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "router not found"}) //nolint:errcheck
		return
	}

	cfg := mikrotik.Config{
		Host:     router.Host,
		Port:     router.Port,
		Username: router.Username,
		Password: router.Password,
		UseTLS:   router.UseSSL,
		Timeout:  time.Duration(router.Timeout) * time.Second,
	}
	client, err := h.mikrotikSvc.GetOrConnect(ctx, router.Name, cfg)
	if err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "router not connected: " + err.Error()}) //nolint:errcheck
		return
	}

	// 1. Send history snapshot (up to 1000 entries).
	history, err := client.GetLogs(ctx, h.topics, 1000)
	if err != nil {
		h.log.Warn("Failed to fetch log history", zap.Uint64("routerID", routerID), zap.Error(err))
	} else {
		if err := conn.WriteJSON(gin.H{"type": "history", "data": history}); err != nil {
			h.log.Warn("Failed to write history", zap.Error(err))
			return
		}
	}

	// 2. Start follow-only stream.
	resultChan := make(chan *dto.LogEntry, 64)
	cancelFn, err := client.ListenLogs(ctx, h.topics, resultChan)
	if err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "failed to start log stream: " + err.Error()}) //nolint:errcheck
		return
	}
	defer cancelFn() //nolint:errcheck

	conn.WriteJSON(gin.H{"type": "status", "status": "started"}) //nolint:errcheck

	// Forward goroutine.
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case entry, ok := <-resultChan:
				if !ok {
					cancel()
					return
				}
				if err := conn.WriteJSON(gin.H{"type": "entry", "data": entry}); err != nil {
					h.log.Warn("Failed to write log entry", zap.Error(err))
					cancel()
					return
				}
				if h.ps != nil {
					ch := pubsub.LogsChannel(uint(routerID), h.topics)
					if err := h.ps.Publish(ctx, ch, entry); err != nil {
						h.log.Debug("pubsub publish failed", zap.Error(err))
					}
				}
			}
		}
	}()

	// Read loop: detect client disconnect.
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.log.Info("Log monitor WS closed", zap.Uint64("routerID", routerID), zap.Error(err))
			return
		}
	}
}
