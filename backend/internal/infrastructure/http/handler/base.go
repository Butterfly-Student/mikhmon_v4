package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorCode represents a machine-readable API error category.
type ErrorCode string

const (
	ErrCodeOK                  ErrorCode = "OK"
	ErrCodeValidation          ErrorCode = "VALIDATION_ERROR"
	ErrCodeUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden           ErrorCode = "FORBIDDEN"
	ErrCodeNotFound            ErrorCode = "NOT_FOUND"
	ErrCodeMikrotikConnection  ErrorCode = "MIKROTIK_CONNECTION"
	ErrCodeMikrotikTimeout     ErrorCode = "MIKROTIK_TIMEOUT"
	ErrCodeMikrotikAuth        ErrorCode = "MIKROTIK_AUTH"
	ErrCodeInternal            ErrorCode = "INTERNAL_ERROR"
)

// Response represents a standard API response
type Response struct {
	Success   bool        `json:"success"`
	Code      ErrorCode   `json:"code"`
	RequestID string      `json:"requestId,omitempty"`
	Timestamp string      `json:"timestamp"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
}

// PaginatedResponse wraps a list response with pagination metadata.
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
}

// BaseHandler provides common handler functionality
type BaseHandler struct {
	Log *zap.Logger
}

// newResponse builds a base response with common fields filled in.
func newResponse(c *gin.Context) Response {
	return Response{
		RequestID: c.GetHeader("X-Request-ID"),
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// Success returns a 200 success response
func (h *BaseHandler) Success(c *gin.Context, data interface{}) {
	r := newResponse(c)
	r.Success = true
	r.Code = ErrCodeOK
	r.Data = data
	c.JSON(http.StatusOK, r)
}

// SuccessWithMessage returns a 200 success response with a message
func (h *BaseHandler) SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	r := newResponse(c)
	r.Success = true
	r.Code = ErrCodeOK
	r.Message = message
	r.Data = data
	c.JSON(http.StatusOK, r)
}

// Created returns a 201 created response
func (h *BaseHandler) Created(c *gin.Context, data interface{}) {
	r := newResponse(c)
	r.Success = true
	r.Code = ErrCodeOK
	r.Data = data
	c.JSON(http.StatusCreated, r)
}

// Paginated returns a 200 response with pagination metadata
func (h *BaseHandler) Paginated(c *gin.Context, items interface{}, total, page, pageSize int) {
	totalPages := 0
	if pageSize > 0 && total > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}
	r := newResponse(c)
	r.Success = true
	r.Code = ErrCodeOK
	r.Data = PaginatedResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
	c.JSON(http.StatusOK, r)
}

// Error returns an error response with the given HTTP status, error code, and message.
func (h *BaseHandler) Error(c *gin.Context, status int, message string) {
	h.ErrorWithCode(c, status, inferErrorCode(status), message)
}

// ErrorWithCode returns an error response with an explicit error code.
func (h *BaseHandler) ErrorWithCode(c *gin.Context, status int, code ErrorCode, message string) {
	r := newResponse(c)
	r.Success = false
	r.Code = code
	r.Error = message
	c.JSON(status, r)
}

// inferErrorCode maps HTTP status codes to default ErrorCode values.
func inferErrorCode(status int) ErrorCode {
	switch status {
	case http.StatusBadRequest:
		return ErrCodeValidation
	case http.StatusUnauthorized:
		return ErrCodeUnauthorized
	case http.StatusForbidden:
		return ErrCodeForbidden
	case http.StatusNotFound:
		return ErrCodeNotFound
	case http.StatusBadGateway:
		return ErrCodeMikrotikConnection
	case http.StatusGatewayTimeout:
		return ErrCodeMikrotikTimeout
	default:
		return ErrCodeInternal
	}
}

// GetUserID gets the user ID from the context
func (h *BaseHandler) GetUserID(c *gin.Context) string {
	userID, _ := c.Get("user_id")
	if id, ok := userID.(string); ok {
		return id
	}
	return ""
}
