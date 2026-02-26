package mikrotik

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler"
	"github.com/irhabi89/mikhmon/internal/usecase/mikrotik"
	"go.uber.org/zap"
)

// QueueHandler handles queue management
type QueueHandler struct {
	handler.BaseHandler
	queueUC *mikrotik.QueueUseCase
}

// NewQueueHandler creates a new queue handler
func NewQueueHandler(queueUC *mikrotik.QueueUseCase, log *zap.Logger) *QueueHandler {
	return &QueueHandler{
		BaseHandler: handler.BaseHandler{Log: log.Named("queue")},
		queueUC:     queueUC,
	}
}

// GetAllQueues retrieves all queues
func (h *QueueHandler) GetAllQueues(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	queues, err := h.queueUC.GetAllQueues(c.Request.Context(), uint(routerID))
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, queues)
}

// GetParentQueues retrieves parent queues
func (h *QueueHandler) GetParentQueues(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	queues, err := h.queueUC.GetParentQueues(c.Request.Context(), uint(routerID))
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, queues)
}

// RegisterRoutes registers queue routes
func (h *QueueHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/queues", h.GetAllQueues)
	r.GET("/queues/parent", h.GetParentQueues)
}
