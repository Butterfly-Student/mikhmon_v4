package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ZapLogger returns a Gin middleware that logs requests with Zap structured logging.
// It logs method, path, status code, latency, client IP, and the X-Request-ID header.
func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = c.GetString("request_id")
		}

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
			zap.String("requestId", requestID),
		}

		if query != "" {
			fields = append(fields, zap.String("query", query))
		}

		errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()
		if errMsg != "" {
			fields = append(fields, zap.String("errors", errMsg))
		}

		switch {
		case status >= 500:
			logger.Error("Request", fields...)
		case status >= 400:
			logger.Warn("Request", fields...)
		default:
			logger.Info("Request", fields...)
		}
	}
}
