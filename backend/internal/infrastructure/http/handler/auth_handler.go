package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/usecase"
	"go.uber.org/zap"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	BaseHandler
	authUC *usecase.AuthUseCase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUC *usecase.AuthUseCase, log *zap.Logger) *AuthHandler {
	return &AuthHandler{
		BaseHandler: BaseHandler{Log: log.Named("auth")},
		authUC:      authUC,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	resp, err := h.authUC.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.Log.Warn("Login failed", zap.String("username", req.Username), zap.Error(err))
		h.ErrorWithCode(c, http.StatusUnauthorized, ErrCodeUnauthorized, err.Error())
		return
	}
	h.Log.Info("Login success", zap.String("username", req.Username))

	h.Success(c, resp)
}

// GetMe retrieves current user info
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := h.GetUserID(c)
	if userID == "" {
		h.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.authUC.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.Success(c, user)
}

// Logout handles logout for stateless JWT mode.
// Token invalidation/blacklist is not enabled; client should remove token locally.
func (h *AuthHandler) Logout(c *gin.Context) {
	h.SuccessWithMessage(c, "Logged out successfully", nil)
}

// RegisterRoutes registers auth routes
func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.GET("/me", h.GetMe)
		auth.POST("/logout", h.Logout)
	}
}
