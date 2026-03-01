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

// PPPInactiveMonitorHandler handles WebSocket streaming of inactive PPP secrets
type PPPInactiveMonitorHandler struct {
	upgrader      websocket.Upgrader
	routerRepo    repository.RouterRepository
	mikrotikSvc   *mikrotik.Manager
	internalWSKey string
	log           *zap.Logger
}

// NewPPPInactiveMonitorHandler creates a new WebSocket handler for PPP inactive monitoring
func NewPPPInactiveMonitorHandler(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
	internalWSKey string,
	log *zap.Logger,
) *PPPInactiveMonitorHandler {
	return &PPPInactiveMonitorHandler{
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
		log:           log.Named("ws-ppp-inactive-monitor"),
	}
}

// Handle handles WebSocket connections for PPP inactive streaming
func (h *PPPInactiveMonitorHandler) Handle(c *gin.Context) {
	if c.Query("key") != h.internalWSKey {
		h.log.Warn("WebSocket auth failed", zap.String("providedKey", c.Query("key")))
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

	h.log.Info("PPP inactive monitor WebSocket upgraded", zap.Uint64("routerID", routerID))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
					h.log.Warn("WebSocket ping failed", zap.Uint64("routerID", routerID), zap.Error(err))
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

	forwardingCancel := h.startMonitor(ctx, conn, router)
	defer forwardingCancel()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			h.log.Info("PPP inactive monitor connection closed", zap.Uint64("routerID", routerID), zap.Error(err))
			return
		}
	}
}

func (h *PPPInactiveMonitorHandler) startMonitor(
	parentCtx context.Context,
	conn *websocket.Conn,
	router *entity.Router,
) context.CancelFunc {
	ctx, cancel := context.WithCancel(parentCtx)

	go func() {
		defer cancel()

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

		resultChan := make(chan []*dto.PPPSecret, 10)
		cancelFn, err := routerClient.ListenPPPInactive(ctx, resultChan)
		if err != nil {
			h.log.Error("Failed to start PPP inactive monitor", zap.Error(err))
			conn.WriteJSON(gin.H{"type": "error", "message": "failed to start monitor: " + err.Error()})
			return
		}
		defer cancelFn()

		conn.WriteJSON(gin.H{"type": "status", "status": "started", "monitor": "ppp-inactive"})

		for {
			select {
			case <-ctx.Done():
				return
			case result, ok := <-resultChan:
				if !ok {
					return
				}
				if err := conn.WriteJSON(result); err != nil {
					h.log.Warn("Failed to write PPP inactive result", zap.Error(err))
					return
				}
			}
		}
	}()

	return cancel
}
