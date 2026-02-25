package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/domain/entity"
	"github.com/irhabi89/mikhmon/internal/usecase"
	"go.uber.org/zap"
)

// RouterHandler handles router management
type RouterHandler struct {
	BaseHandler
	routerUC *usecase.RouterUseCase
}

// NewRouterHandler creates a new router handler
func NewRouterHandler(routerUC *usecase.RouterUseCase, log *zap.Logger) *RouterHandler {
	return &RouterHandler{
		BaseHandler: BaseHandler{Log: log.Named("router")},
		routerUC:    routerUC,
	}
}

// CreateRouter creates a new router
func (h *RouterHandler) CreateRouter(c *gin.Context) {
	var req entity.RouterCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	router, err := h.routerUC.Create(c.Request.Context(), &req)
	if err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.Created(c, router)
}

// ListRouters lists all routers
func (h *RouterHandler) ListRouters(c *gin.Context) {
	routers, err := h.routerUC.GetAll(c.Request.Context())
	if err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.Success(c, routers)
}

// GetRouter gets a router by ID
func (h *RouterHandler) GetRouter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	router, err := h.routerUC.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.Error(c, http.StatusNotFound, err.Error())
		return
	}

	h.Success(c, router)
}

// UpdateRouter updates a router
func (h *RouterHandler) UpdateRouter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req entity.RouterUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	router, err := h.routerUC.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.Success(c, router)
}

// DeleteRouter deletes a router
func (h *RouterHandler) DeleteRouter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.routerUC.Delete(c.Request.Context(), uint(id)); err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.SuccessWithMessage(c, "Router deleted successfully", nil)
}

// TestConnection tests router connection
func (h *RouterHandler) TestConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.routerUC.TestConnection(c.Request.Context(), uint(id)); err != nil {
		h.Error(c, http.StatusBadGateway, err.Error())
		return
	}

	h.SuccessWithMessage(c, "Connection successful", nil)
}

// RegisterRoutes registers router routes
func (h *RouterHandler) RegisterRoutes(r *gin.RouterGroup) {
	router := r.Group("/routers")
	{
		router.POST("", h.CreateRouter)
		router.GET("", h.ListRouters)
		router.GET("/:id", h.GetRouter)
		router.PUT("/:id", h.UpdateRouter)
		router.DELETE("/:id", h.DeleteRouter)
		router.POST("/:id/test", h.TestConnection)
	}
}
